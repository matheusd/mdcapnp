// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"math"
	"unsafe"
)

// PointerFieldIndex is the index of a pointer field in a struct. The first
// pointer of a struct's pointer section has index 0, the second one has index
// 1, and so on.
//
// Given that 1 pointer == 1 word, and a struct is limited to 16 bits worth of
// words in its pointer section, it can have up to 2^16 pointer fields.
type PointerFieldIndex uint16

func (i PointerFieldIndex) uncheckedWordOffset(base WordOffset) WordOffset {
	return base + WordOffset(i)
}

// DataFieldIndex is the index of a data field in a struct (in words). When the
// field is smaller than one word, then further indexing may be necessary to
// extract its value.
type DataFieldIndex uint16

// uncheckedWordOffset returns the final offset of a data word of a struct
// inside a segment, without checking for the validity of this operation. This
// is used in cases where it is assumed that the index has already been
// determined to be valid inside the struct (because the entire struct size has
// been bounds checked already).
func (i DataFieldIndex) uncheckedWordOffset(base WordOffset) WordOffset {
	return base + WordOffset(i)
}

type Struct struct {
	seg   *Segment
	arena *Arena
	dl    depthLimit
	ptr   structPointer
}

// HasData returns true if the specified word in the data section is set in this
// struct.
func (s *Struct) HasData(dataIndex DataFieldIndex) bool {
	return s.ptr.dataOffset > 0 && dataIndex < DataFieldIndex(s.ptr.dataSectionSize)
}

func (s *Struct) Int64(dataIndex DataFieldIndex) (res int64) {
	if s.HasData(dataIndex) {
		res = int64(s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset)))
	}
	return
}

func (s *Struct) Float64(dataIndex DataFieldIndex) (res float64) {
	if s.HasData(dataIndex) {
		res = math.Float64frombits(uint64(s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))))
	}
	return
}

type Int32DataFieldShift int

const (
	Int32FieldLo Int32DataFieldShift = 0
	Int32FieldHi Int32DataFieldShift = 32
)

type Int32DataFieldSetMask Word

const (
	Int32FieldSetMaskLo Int32DataFieldSetMask = 0x00000000ffffffff
	Int32FieldSetMaskHi Int32DataFieldSetMask = 0xffffffff0000000
)

// Int32 returns a data field as an int32. Given that an int32 field occupies
// either the low or high end of data word, the second parameter disambiguates
// between the two.
//
// TODO: review if this is the way to go.
func (s *Struct) Int32(dataIndex DataFieldIndex, shift Int32DataFieldShift) (res int32) {
	if s.HasData(dataIndex) {
		data := s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = int32(data >> shift)
	}
	return
}

// Bool returns a data field as a bool. fieldIndex points to the data word
// within the struct, while bit determines which bit (within the word)
// corresponds to the target field.
func (s *Struct) Bool(dataIndex DataFieldIndex, bit byte) (res bool) {
	if s.HasData(dataIndex) {
		data := s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = data&(1<<bit) != 0
	}
	return res
}

func (s *Struct) readListPtr(ptrIndex PointerFieldIndex) (seg *Segment, lp listPointer, listDL depthLimit, err error) {
	// Check if we can descend further into the struct (to fetch the first
	// list pointer).
	var ok bool
	listDL, ok = s.dl.dec()
	if !ok {
		err = errDepthLimitExceeded
		return
	}

	// Check if this pointer is set within the pointer section.
	if ptrIndex >= PointerFieldIndex(s.ptr.pointerSectionSize) {
		// TODO: return default if it exists? Or handle this at a higher
		// level?
		err = errStructBuilderDoesNotContainPointerField(s.ptr.pointerSectionSize)
		return
	}

	// Determine the offset of the pointer word (given its index) and fetch
	// it.
	//
	// This does not need an overflow check because the entire struct
	// (including this pointer offset which is <= pointerSectionSize) is
	// known to be in bounds.
	pointerOffset := s.ptr.dataOffset + WordOffset(s.ptr.dataSectionSize) + WordOffset(ptrIndex)
	ptr := s.seg.uncheckedGetWordAsPointer(pointerOffset)

	// De-ref far pointers into the concrete list segment and near pointer.
	ptrType := ptr.pointerType()
	if ptrType == pointerTypeFarPointer /*ptr.isFarPointer()*/ {
		seg, ptr, listDL, err = derefFarPointer(s.arena, listDL, ptr)
		if err != nil {
			return
		}
		ptrType = ptr.pointerType()
	} else {
		seg = s.seg
	}

	// Check if the final pointer (after potential deref) is a list pointer.
	if ptrType != pointerTypeList /*!ptr.isListPointer() */ {
		err = errNotListPointer
		return
	}
	lp = ptr.toListPointer()

	// Determine concrete offset into segment of where the list actually
	// starts.
	if lp.startOffset, ok = addWordOffsetsWithCarry(pointerOffset, lp.startOffset, 1); !ok {
		err = errWordOffsetSumOverflows{lp.startOffset, pointerOffset}
		return
	}

	// Check if entire list is readable.
	fullSize := listWordCount(lp.elSize, lp.listSize)
	if err = seg.checkBounds(lp.startOffset, fullSize); err != nil {
		return
	}

	// If list elements have zero size, count them as one byte per element,
	// to avoid vulnerabilities where large simulated lists are iterated.
	if lp.elSize == listElSizeVoid {
		fullSize = WordCount(lp.listSize / WordSize)
	}
	if err = s.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return
	}

	return
}

func (s *Struct) ReadList(ptrIndex PointerFieldIndex, ls *List) error {
	seg, lp, listDL, err := s.readListPtr(ptrIndex)
	if err != nil {
		return err
	}

	// All good.
	ls.seg = seg
	ls.ptr = lp
	ls.dl = listDL
	return nil
}

func (s *Struct) UnsafeString(ptrIndex PointerFieldIndex) string {
	seg, lp, _, err := s.readListPtr(ptrIndex)
	if err != nil || lp.elSize != listElSizeByte || lp.listSize == 0 {
		return ""
	}

	return seg.uncheckedUnsafeString(lp.startOffset, ByteCount(lp.listSize)-1)
}

func (s *Struct) String(ptrIndex PointerFieldIndex) string {
	seg, lp, _, err := s.readListPtr(ptrIndex)
	if err != nil || lp.elSize != listElSizeByte || lp.listSize == 0 {
		return ""
	}

	buf := seg.uncheckedSlice(lp.startOffset, ByteCount(lp.listSize)-1) // -1 to remove last null
	return string(buf)
}

func (s *Struct) UnsafeStringXXX(dataIndex DataFieldIndex, size WordCount) string {
	// XXX HasData()
	hasData := s.ptr.dataOffset > 0 && (dataIndex+DataFieldIndex(size)) <= DataFieldIndex(s.ptr.dataSectionSize)
	if !hasData {
		return ""
	}

	buf := s.seg.uncheckedSlice(dataIndex.uncheckedWordOffset(s.ptr.dataOffset), ByteCount(size*WordSize))
	strLen := buf[len(buf)-1]

	buf = buf[:strLen-1]
	return *(*string)(unsafe.Pointer(&buf))
}
