// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

type pointerType int

const (
	pointerTypeStruct     pointerType = 0x00
	pointerTypeList       pointerType = 0x01
	pointerTypeFarPointer pointerType = 0x02
)

type pointer Word

func (ptr pointer) dataOffset() WordOffset {
	return WordOffset(ptr&0xfffffffc) >> 2
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

func (ptr pointer) isStructPointer() bool {
	return (ptr & 0x03) == 0
}

func (ptr pointer) isListPointer() bool {
	return (ptr & 0x03) == 1
}

func (ptr pointer) isFarPointer() bool {
	return (ptr & 0x03) == 2
}

func (ptr pointer) pointerType() pointerType {
	return pointerType(ptr & 0x03)
}

// isNullPointer returns true if this is a "null" pointer. A null pointer is all
// zeros (except for the first two bits which may denote the type of pointer).
func (ptr pointer) isNullPointer() bool {
	return (ptr & 0xfffffffffffffffc) == 0
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

// derefFarPointer de-references a far pointer into a concrete segment pointer.
// It follows pointers (up to the depth limit) until a non-far pointer is found.
//
// Returns the resulting pointer and the remaining depth limit.
func derefFarPointer(arena *Arena, dl depthLimit, ptr pointer) (*Segment, pointer, depthLimit, error) {
	// TODO: implement.
	return nil, 0, 0, nil
}
