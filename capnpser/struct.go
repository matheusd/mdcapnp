// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"errors"
	"fmt"
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

type StructType = struct {
	seg   *Segment
	arena *Arena
	dl    depthLimit
	ptr   structPointer
}

type Struct StructType

/*
type Struct struct {
	seg   *Segment
	arena *Arena
	dl    depthLimit
	ptr   structPointer
}
*/

// structSize returns the (shallow) size of this object. Guaranteed to be a
// valid count.
func (s *Struct) structSize() WordCount {
	return WordCount(s.ptr.dataSectionSize) + WordCount(s.ptr.pointerSectionSize)
}

func (s *Struct) AsAnyPointer() AnyPointer {
	return AnyPointer{
		seg:   s.seg,
		arena: s.arena,
		dl:    s.dl,
		ptr:   s.ptr.toPointer(),
	}
}

// HasData returns true if the specified word in the data section is set in this
// struct.
func (s *Struct) HasData(dataIndex DataFieldIndex) bool {
	return s.ptr.dataOffset > 0 && dataIndex < DataFieldIndex(s.ptr.dataSectionSize)
}

func (s *Struct) Uint64(dataIndex DataFieldIndex) (res uint64) {
	if s.HasData(dataIndex) {
		res = uint64(s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset)))
	}
	return
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

type Uint16DataFieldShift int

const (
	Uint16FieldShift0 Uint16DataFieldShift = 0
	Uint16FieldShift1 Uint16DataFieldShift = 16
	Uint16FieldShift2 Uint16DataFieldShift = 32
	Uint16FieldShift3 Uint16DataFieldShift = 48
)

type Uint16DataFieldSetMask Word

const (
	Uint16FieldSetMask0 Uint16DataFieldSetMask = 0x000000000000ffff
	Uint16FieldSetMask1 Uint16DataFieldSetMask = 0x00000000ffff0000
	Uint16FieldSetMask2 Uint16DataFieldSetMask = 0x0000ffff0000000
	Uint16FieldSetMask3 Uint16DataFieldSetMask = 0xffff00000000000
)

func (s *Struct) Uint16(dataIndex DataFieldIndex, shift Uint16DataFieldShift) (res uint16) {
	if s.HasData(dataIndex) {
		data := s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = uint16(data >> shift)
	}
	return
}

type Uint32DataFieldShift int

const (
	Int32FieldLo  Uint32DataFieldShift = 0
	Uint32FieldLo Uint32DataFieldShift = 0
	Int32FieldHi  Uint32DataFieldShift = 32
	Uint32FieldHi Uint32DataFieldShift = 32
)

type Uint32DataFieldSetMask Word

const (
	Uint32FieldSetMaskLo Uint32DataFieldSetMask = 0x00000000ffffffff
	Uint32FieldSetMaskHi Uint32DataFieldSetMask = 0xffffffff0000000
)

// Uint32 returns a data field as an uint32. Given that an uint32 field occupies
// either the low or high end of data word, the second parameter disambiguates
// between the two.
func (s *Struct) Uint32(dataIndex DataFieldIndex, shift Uint32DataFieldShift) (res uint32) {
	if s.HasData(dataIndex) {
		data := s.seg.uncheckedGetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = uint32(data >> shift)
	}
	return
}

// Int32 returns a data field as an int32. Given that an int32 field occupies
// either the low or high end of data word, the second parameter disambiguates
// between the two.
//
// TODO: review if this is the way to go.
func (s *Struct) Int32(dataIndex DataFieldIndex, shift Uint32DataFieldShift) (res int32) {
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

func (s *Struct) readFieldPtr(ptrIndex PointerFieldIndex) (seg *Segment, ptrType pointerType,
	ptr pointer, dl depthLimit, pointerOffset WordOffset, err error) {

	// Check if we can descend further into the struct (to fetch the first
	// list pointer).
	var ok bool
	dl, ok = s.dl.dec()
	if !ok {
		err = errDepthLimitExceeded
		return
	}

	// Check if this pointer is set within the pointer section.
	if ptrIndex >= PointerFieldIndex(s.ptr.pointerSectionSize) {
		// TODO: return default if it exists? Or handle this at a higher
		// level?
		err = errStructDoesNotContainPointerField(ptrIndex)
		return
	}

	// Determine the offset of the pointer word (given its index) and fetch
	// it.
	//
	// This does not need an overflow check because the entire struct
	// (including this pointer offset which is <= pointerSectionSize) is
	// known to be in bounds.
	pointerOffset = s.ptr.dataOffset + WordOffset(s.ptr.dataSectionSize) + WordOffset(ptrIndex)
	ptr = s.seg.uncheckedGetWordAsPointer(pointerOffset)

	// De-ref far pointers into the concrete list segment and near pointer.
	ptrType = ptr.pointerType()
	if ptrType == pointerTypeFarPointer /*ptr.isFarPointer()*/ {
		seg, ptr, dl, err = derefFarPointer(s.arena, dl, ptr)
		if err != nil {
			return
		}
		ptrType = ptr.pointerType()
	} else {
		seg = s.seg
	}

	return
}

func (s *Struct) readListPtr(ptrIndex PointerFieldIndex) (seg *Segment, lp listPointer, listDL depthLimit, err error) {

	var ptr pointer
	var ptrType pointerType
	var pointerOffset WordOffset
	seg, ptrType, ptr, listDL, pointerOffset, err = s.readFieldPtr(ptrIndex)

	if ptr.isNullPointer() {
		// A null pointer is an either an empty list or the default
		// value list for a field.
		return
	}

	// Check if the final pointer (after potential deref) is a list pointer.
	if ptrType != pointerTypeList /*!ptr.isListPointer() */ {
		err = errNotListPointer
		return
	}
	lp = ptr.toListPointer()

	// Determine concrete offset into segment of where the list actually
	// starts.
	var ok bool
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
	ls.arena = s.arena
	ls.ptr = lp
	ls.dl = listDL
	return nil
}

func (s *Struct) ReadStructList(ptrIndex PointerFieldIndex, sls *StructList) error {
	seg, lp, listDL, err := s.readListPtr(ptrIndex)
	if err != nil {
		return err
	}

	if lp.elSize != listElSizeComposite {
		// A null pointer means this is an empty list or the default
		// field value.
		if lp.toPointer().isNullPointer() {
			*sls = StructList{
				l: List{
					seg: seg,
					ptr: lp,
					dl:  listDL,
				},
				itemSize: StructSize{},
				listLen:  0,
			}
			return nil
		}

		// TODO: support re-interpreting native lists as composite.
		return errors.New("not a composite list")
	}

	// Read the tag word, which contains the per-item information. The list
	// has already been verified to be entirely in-bounds (by word count),
	// therefore the tag word is in-bounds.
	tagWord := s.seg.uncheckedGetWordAsPointer(lp.startOffset)
	if !tagWord.isStructPointer() {
		return errors.New("composite list tag word is not a struct pointer")
	}
	listLen := listSize(tagWord.dataOffset())
	itemSize := StructSize{DataSectionSize: tagWord.dataSectionSize(), PointerSectionSize: tagWord.pointerSectionSize()}

	// Double check the total list word count size is correct, when
	// calculated as itemSize * len (readListPtr verified by by total word
	// count only).
	gotListWordCount := uint64(itemSize.TotalSize())*uint64(listLen) + 1 // +1 tag word
	if gotListWordCount != uint64(listWordCount(lp.elSize, lp.listSize)) {
		return fmt.Errorf("incongruent list sizes: in list pointer %d, in tag word: %d",
			listWordCount(lp.elSize, lp.listSize), gotListWordCount)
	}

	*sls = StructList{
		l: List{
			seg: seg,
			ptr: lp,
			dl:  listDL,
		},
		itemSize: itemSize,
		listLen:  listLen,
	}

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

func (s *Struct) ReadStruct(ptrIndex PointerFieldIndex, res *Struct) (err error) {
	seg, ptrType, ptr, structDL, pointerOffset, err := s.readFieldPtr(ptrIndex)

	// Check if the final pointer (after potential deref) is a struct
	// pointer.
	if ptrType != pointerTypeStruct {
		err = errNotStructPointer
		return
	}
	sp := ptr.toStructPointer()

	// Determine concrete offset into segment of where the struct actually
	// starts.
	var ok bool
	if sp.dataOffset, ok = addWordOffsetsWithCarry(pointerOffset, sp.dataOffset, 1); !ok {
		err = errWordOffsetSumOverflows{sp.dataOffset, pointerOffset}
		return
	}

	// Check if entire struct is readable.
	fullSize := WordCount(sp.dataSectionSize) + WordCount(sp.pointerSectionSize)
	if err = seg.checkBounds(sp.dataOffset, fullSize); err != nil {
		return
	}
	if err = s.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return
	}

	*res = Struct{
		seg:   seg,
		arena: s.arena,
		dl:    structDL,
		ptr:   sp,
	}
	return
}

func (s *Struct) ReadAnyPointer(ptrIndex PointerFieldIndex, res *AnyPointer) (err error) {
	seg, ptrType, ptr, dl, pointerOffset, err := s.readFieldPtr(ptrIndex)

	// Determine concrete offset into segment of where the object actually
	// starts.
	if ptrType == pointerTypeList || ptrType == pointerTypeStruct {
		dataOffset, ok := addWordOffsetsWithCarry(pointerOffset, ptr.dataOffset(), 1)
		if !ok {
			return errWordOffsetSumOverflows{ptr.dataOffset(), pointerOffset}
		}
		ptr = ptr.withDataOffset(dataOffset)
	}

	if err == nil {
		*res = AnyPointer{
			seg:   seg,
			arena: s.arena,
			dl:    dl,
			ptr:   ptr,
			//pointerOffset: pointerOffset,
			//parentOffset:  s.ptr.dataOffset,
		}
	}

	return
}

func StructToAnyPointer[T ~StructType](v T) AnyPointer {
	s := Struct(v)
	return AnyPointer{
		seg:   s.seg,
		arena: s.arena,
		dl:    s.dl,
		ptr:   s.ptr.toPointer(),
	}
}
