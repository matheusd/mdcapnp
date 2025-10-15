// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import "errors"

type pointerType pointer

const (
	pointerTypeStruct     pointerType = 0x00
	pointerTypeList       pointerType = 0x01
	pointerTypeFarPointer pointerType = 0x02
	pointerTypeOther      pointerType = 0x03

	zeroStructPointer pointer = 0x00000000fffffffc
)

type pointer Word

func (ptr pointer) dataOffset() WordOffset {
	return WordOffset(ptr&0xfffffffc) >> 2
}

// withDataOffset returns a new pointer with the given data offset (for lists
// and structs).
func (ptr pointer) withDataOffset(off WordOffset) pointer {
	return pointer(ptr&0xfffffffc | pointer(off<<2))
}

func (ptr pointer) dataSectionSize() wordCount16 {
	return wordCount16(ptr & 0xffff00000000 >> 32)
}

func (ptr pointer) pointerSectionSize() wordCount16 {
	return wordCount16(ptr >> 48)
}

func (ptr pointer) elSize() listElementSize {
	return listElementSize((ptr >> 32) & 0b111)
}

func (ptr pointer) listSize() listSize {
	return listSize(ptr & 0xfffffff800000000 >> 35)
}

func (ptr pointer) capPointerIndex() uint32 {
	return uint32(ptr >> 32)
}

func (ptr pointer) isStructPointer() bool {
	return (ptr & 0x03) == 0
}

func (ptr pointer) isListPointer() bool {
	return (ptr & 0x03) == 1
}

func (ptr pointer) isFarPointer() bool {
	return (ptr & 0x03) == 2
}

func (ptr pointer) isCapPointer() bool {
	return (ptr & 0xffffffff) == 3
}

func (ptr pointer) pointerType() pointerType {
	return pointerType(ptr & 0x03)
}

// isNullPointer returns true if this is a "null" pointer. A null pointer is all
// zeros (except for the first two bits which may denote the type of pointer).
func (ptr pointer) isNullPointer() bool {
	return (ptr & 0xfffffffffffffffc) == 0
}

// isZeroStruct returns true if this is a zero-struct pointer. A zero-struct
// pointer has zero size and offset == -1.
func (ptr pointer) isZeroStruct() bool {
	return ptr == 0x00000000fffffffc
}

func (ptr pointer) toStructPointer() (sp structPointer) {
	sp.dataOffset = ptr.dataOffset()
	sp.dataSectionSize = ptr.dataSectionSize()
	sp.pointerSectionSize = ptr.pointerSectionSize()
	return
}

func (ptr pointer) toListPointer() (lp listPointer) {
	lp.startOffset = ptr.dataOffset()
	lp.elSize = ptr.elSize()
	lp.listSize = ptr.listSize()
	return
}

type structPointer struct {
	dataOffset         WordOffset
	dataSectionSize    wordCount16
	pointerSectionSize wordCount16
}

func (sp structPointer) toPointer() pointer {
	return pointer(pointerTypeStruct) |
		pointer(uint32(sp.dataOffset<<2)) |
		pointer(sp.dataSectionSize)<<32 |
		pointer(sp.pointerSectionSize)<<48
}

func buildRawStructPointer(off WordOffset, sz StructSize) pointer {
	return pointer(pointerTypeStruct) |
		pointer(uint32(off<<2)) |
		pointer(sz.DataSectionSize)<<32 |
		pointer(sz.PointerSectionSize)<<48
}

type listPointer struct {
	startOffset WordOffset
	elSize      listElementSize
	listSize    listSize
}

func (lp listPointer) toPointer() pointer {
	return pointer(pointerTypeList) |
		pointer(uint32(lp.startOffset<<2)) |
		pointer(lp.elSize)<<32 |
		pointer(lp.listSize)<<35
}

func buildRawListPointer(startOffset WordOffset, elSize listElementSize, lsSize listSize) pointer {
	return pointer(pointerTypeList) |
		pointer(uint32(startOffset<<2)) |
		pointer(elSize)<<32 |
		pointer(lsSize)<<35
}

// derefFarPointer de-references a far pointer into a concrete segment pointer.
// It follows pointers (up to the depth limit) until a non-far pointer is found.
//
// Returns the resulting pointer and the remaining depth limit.
//
//go:noinline
func derefFarPointer(arena *Arena, dl depthLimit, ptr pointer) (*Segment, pointer, depthLimit, error) {
	// TODO: implement.
	return nil, 0, 0, errors.New("not implemented")
}

type CapPointer struct {
	index uint32
}

func (cp *CapPointer) Index() uint32 {
	return cp.index
}

func buildRawCapPointer(index uint32) pointer {
	return pointer(pointerTypeOther) |
		pointer(index)<<32
}

type AnyPointer struct {
	seg           *Segment
	arena         *Arena
	dl            depthLimit
	ptr           pointer
	pointerOffset WordOffset
	parentOffset  WordOffset
}

// IsZeroStruct returns true if the pointer represents a zero-struct (a struct
// with zero size).
func (ap *AnyPointer) IsZeroStruct() bool {
	return ap.ptr.isZeroStruct()
}

func (ap *AnyPointer) IsStruct() bool {
	return ap.ptr.isStructPointer()
}

func (ap *AnyPointer) IsCapPointer() bool {
	return ap.ptr.isCapPointer()
}

func (ap *AnyPointer) AsStruct() Struct {
	return Struct{
		seg:   ap.seg,
		arena: ap.arena,
		dl:    ap.dl,
		ptr:   ap.ptr.toStructPointer(),
	}
}

func (ap *AnyPointer) AsCapPointer() CapPointer {
	return CapPointer{
		index: ap.ptr.capPointerIndex(),
	}
}
