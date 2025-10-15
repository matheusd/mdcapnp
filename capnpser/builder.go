// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type AnyPointerBuilder struct {
	mb *MessageBuilder

	off WordOffset // Concrete offset into segment where pointer data begins (if this is a struct/list).
	ptr pointer
	sid SegmentID
}

func CapPointerAsAnyPointerBuilder(index uint32) AnyPointerBuilder {
	return AnyPointerBuilder{
		ptr: buildRawCapPointer(index),
	}
}

func ZeroStructAsPointerBuilder() AnyPointerBuilder {
	return AnyPointerBuilder{
		ptr: zeroStructPointer,
	}
}

type StructListBuilder struct {
	off WordOffset // Concrete offset into segment where struct data begins.

	itemSize StructSize
	listLen  listSize // Cannot be > listCap
	listCap  listSize

	mb  *MessageBuilder
	urb UnsafeRawBuilder
	sid SegmentID
}

// Len is the current length of the list.
func (slb *StructListBuilder) Len() int {
	return int(slb.listLen)
}

// Cap is the capacity of the list.
func (slb *StructListBuilder) Cap() int {
	return int(slb.listCap)
}

func (slb *StructListBuilder) writeTagWord() {
	pointer := buildRawStructPointer(WordOffset(slb.listLen), slb.itemSize)
	slb.urb.SetWord(0, Word(pointer))
}

// SetLen modifies the length of the list. It can increase or decrease the
// length without modifiying the internal structures.
func (slb *StructListBuilder) SetLen(newLen int) {
	if newLen < 0 {
		panic("new len cannot be < 0")
	}
	if newLen > int(slb.listCap) {
		panic("new len cannot be > cap")
	}
	slb.listLen = listSize(newLen)
	slb.writeTagWord()
}

// At returns the builder for the i'th element of the list. Does NOT perform
// bounds check.
func (slb *StructListBuilder) at(i int) (res StructBuilder) {
	// Determine offset of this structure. No need to check for bounds
	// because the list is assumed initialized at least up to len i.
	relOff := WordOffset(1 + WordOffset(WordCount(i)*slb.itemSize.TotalSize()))

	res = StructBuilder{
		off: slb.off + relOff,
		sz:  slb.itemSize,
		mb:  slb.mb,
		sid: slb.sid,
		// urb: slb.urb.Child(relOff),
	}
	return
}

// At returns the builder for the i'th element of the list. Panics if the item
// is out of bounds.
func (slb *StructListBuilder) At(i int) (res StructBuilder) {
	if i < 0 {
		panic("i is < 0 (out of bounds)")
	}
	if i > slb.Len() {
		panic("i is > len (out of bounds)")
	}
	return slb.at(i)
}

// Add an item to the list.
func (slb *StructListBuilder) Add() (res StructBuilder) {
	if slb.listLen == slb.listCap {
		panic(errors.New("list is full"))
	}

	i := int(slb.listLen)
	slb.listLen++

	// Write the new tag word with added element.
	slb.writeTagWord()

	return slb.at(i)
}

type GenericStructListBuilder[T ~StructBuilderType] struct {
	slb StructListBuilder
}

func (gslb *GenericStructListBuilder[T]) SetLen(i int) { gslb.slb.SetLen(i) }
func (gslb *GenericStructListBuilder[T]) Len() int     { return gslb.slb.Len() }
func (gslb *GenericStructListBuilder[T]) Cap() int     { return gslb.slb.Cap() }
func (gslb *GenericStructListBuilder[T]) At(i int) T   { return T(gslb.slb.At(i)) }
func (gslb *GenericStructListBuilder[T]) Add(i int) T  { return T(gslb.slb.Add()) }

func NewGenericStructListBuilder[T ~StructBuilderType](mb *MessageBuilder,
	itemSize StructSize, listLen, listCap int) (res GenericStructListBuilder[T], err error) {

	res.slb, err = mb.NewStructList(itemSize, listLen, listCap)
	return
}

func NewGenericStructListBuilderField[T ~StructBuilderType](s *StructBuilder,
	ptrIndex PointerFieldIndex, itemSize StructSize, listLen, listCap int) (
	res GenericStructListBuilder[T], err error) {

	res.slb, err = s.NewStructListField(ptrIndex, itemSize, listLen, listCap)
	return
}

type StructBuilderType = struct {
	off WordOffset // Concrete offset into segment where struct data begins.
	sz  StructSize
	mb  *MessageBuilder
	urb SegmentBuilder
	// urb UnsafeRawBuilder
	sid SegmentID
}
type StructBuilder StructBuilderType

func (sb *StructBuilder) hasData(dataIndex DataFieldIndex) bool {
	return dataIndex < DataFieldIndex(sb.sz.DataSectionSize)
}

func (sb *StructBuilder) hasPointer(ptrIndex PointerFieldIndex) bool {
	return ptrIndex < PointerFieldIndex(sb.sz.PointerSectionSize)
}

func (sb *StructBuilder) SetUint64(dataIndex DataFieldIndex, v uint64) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// finalOff := dataIndex.uncheckedWordOffset(sb.off)
		// sb.seg.uncheckedSetWord(finalOff, Word(v))
		sb.urb.SetWord(WordOffset(dataIndex), Word(v))
	}
	return
}
func (sb *StructBuilder) SetInt64(dataIndex DataFieldIndex, v int64) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// finalOff := dataIndex.uncheckedWordOffset(sb.off)
		// sb.seg.uncheckedSetWord(finalOff, Word(v))
		sb.urb.SetWord(WordOffset(dataIndex), Word(v))
	}
	return
}

func (sb *StructBuilder) SetUint16(dataIndex DataFieldIndex, shift Uint16DataFieldShift, v uint16) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), Word(mask), Word(v))
		sb.urb.maskAndMergeWord(WordOffset(dataIndex), Word(0xffff)<<shift, Word(v)<<shift)
	}
	return
}

func (sb *StructBuilder) SetUint32(dataIndex DataFieldIndex, shift Uint32DataFieldShift, v uint32) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), Word(mask), Word(v))
		sb.urb.maskAndMergeWord(WordOffset(dataIndex), Word(0xffffffff)<<shift, Word(v)<<shift)
	}
	return
}

func (sb *StructBuilder) SetInt32(dataIndex DataFieldIndex, shift Uint32DataFieldShift, v int32) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), Word(mask), Word(v))
		sb.urb.maskAndMergeWord(WordOffset(dataIndex), Word(0xffffffff)<<shift, Word(v)<<shift)
	}
	return
}

func (sb *StructBuilder) SetBool(dataIndex DataFieldIndex, bit byte, v bool) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), ^(1 << bit), boolToWord(v)<<bit)
		sb.urb.maskAndMergeWord(WordOffset(dataIndex), (1 << bit), boolToWord(v)<<bit)
	}
	return
}

func (sb *StructBuilder) SetFloat64(dataIndex DataFieldIndex, v float64) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		// sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(math.Float64bits(v)))
		sb.urb.SetWord(WordOffset(dataIndex), Word(math.Float64bits(v)))
	}
	return
}

func (sb *StructBuilder) SetString(ptrIndex PointerFieldIndex, v string) (err error) {
	if !sb.hasPointer(ptrIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		return errStructBuilderDoesNotContainPointerField(ptrIndex)
	}

	// Allocate the string as a new object (list of bytes).
	// segb, lsPtr, err := sb.seg.mb.newText(sb.seg.id, v)
	// segb, lsPtr, err := sb.mb.newText(sb.sid, v)
	sid, lsPtr, err := sb.mb.newText(sb.sid, v)
	if err != nil {
		return err
	}

	// TODO: handle allocs in new segments.
	// if segb.id != sb.seg.id {
	// if segb.id != sb.sid {
	if sid != sb.sid {
		return errors.New("needs handling")
	}

	// Offset of the pointer field that will reference the new list. This
	// is relative to the start of this struct (sb).
	ptrOff := ptrIndex.uncheckedWordOffset(WordOffset(sb.sz.DataSectionSize))

	// Determine concrete pointer offset inside struct. This doesn't need
	// overflow checks because the entire struct has been allocated, thus
	// this pointer offset is known to be in bounds.
	concretePtrOff := ptrOff + sb.off

	// Determine the relative offset from the field pointer offset to the
	// actual data. This finishes the construction of the list pointer.
	lsPtr.startOffset = lsPtr.startOffset - concretePtrOff - 1

	// Structure already fully allocated, no need to check for
	// bounds.
	// sb.seg.uncheckedSetWord(ptrOff, Word(lsPtr.toPointer()))
	sb.urb.SetWord(ptrOff, Word(lsPtr.toPointer()))
	return nil
}

func (sb *StructBuilder) NewStructField(ptrIndex PointerFieldIndex, size StructSize) (nsb StructBuilder, err error) {
	if !sb.hasPointer(ptrIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainPointerField(ptrIndex)
		return
	}

	nsb, err = sb.mb.NewStruct(size)
	if err != nil {
		return
	}

	// TODO: handle inter-segment pointer.
	// if sb.seg.id != 0 {
	if nsb.sid != 0 {
		panic("needs handling")
	}

	// Offset of the pointer field that will reference the new struct. This
	// is relative to the start of this struct (sb).
	ptrOff := ptrIndex.uncheckedWordOffset(WordOffset(sb.sz.DataSectionSize))

	// Determine concrete pointer offset of ptrOff inside struct. This
	// doesn't need overflow checks because the entire struct has been
	// allocated, thus this pointer offset is known to be in bounds.
	concretePtrOff := ptrOff + sb.off

	// Determine the relative offset from the field pointer offset to the
	// actual data. This finishes the construction of the struct pointer.
	sp := structPointer{
		dataOffset:         nsb.off - concretePtrOff - 1,
		dataSectionSize:    size.DataSectionSize,
		pointerSectionSize: size.PointerSectionSize,
	}

	// Structure already fully allocated, no need to check for
	// bounds.
	sb.urb.SetWord(ptrOff, Word(sp.toPointer()))

	return
}

func (sb *StructBuilder) NewStructAsUnionValue(ptrIndex PointerFieldIndex,
	size StructSize, unionField DataFieldIndex, unionFieldShift Uint16DataFieldShift,
	unionFieldValue uint16) (nsb StructBuilder, err error) {

	nsb, err = sb.NewStructField(ptrIndex, size)
	if err != nil {
		return
	}

	err = sb.SetUint16(unionField, unionFieldShift, unionFieldValue)
	return
}

func (sb *StructBuilder) SetAnyPointer(ptrIndex PointerFieldIndex, v AnyPointerBuilder) error {
	if !sb.hasPointer(ptrIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		return errStructBuilderDoesNotContainPointerField(ptrIndex)
	}

	// Offset of the pointer field inside sb. No need to check for overflow
	// because the struct has been validated to contain this pointer.
	ptrOff := ptrIndex.uncheckedWordOffset(WordOffset(sb.sz.DataSectionSize))

	// Cap pointers and null pointers are written directly. Cap pointers are
	// simple indices, and null pointers don't point anywhere. Neither
	// requires actual data to link to.
	if v.ptr.isCapPointer() || v.ptr.isNullPointer() || v.ptr.isZeroStruct() {
		sb.urb.SetWord(ptrOff, Word(v.ptr))
		return nil
	}

	ptrType := v.ptr.pointerType()
	if ptrType == pointerTypeFarPointer {
		// TODO: support pointing to ther segments.
		return errors.New("far pointers unsupported")
	}

	// The only cases left are non-null list and struct pointers. Sanity
	// check.
	if ptrType != pointerTypeList && ptrType != pointerTypeStruct {
		return errors.New("unhandled pointer type case in StructBuilder.SetAnyPointer")
	}

	if sb.mb != v.mb {
		// TODO: Copy data from v.mb?
		return errors.New("v not in same messageBuilder as sb")
	}

	if sb.sid != v.sid {
		// TODO: support pointing to other segments.
		return errors.New("v not in the same segment as sb")
	}

	// Determine concrete pointer offset of ptrOff inside struct. This
	// doesn't need overflow checks because the entire struct has been
	// allocated, thus this pointer offset is known to be in bounds.
	concretePtrOff := ptrOff + sb.off

	// Create the new pointer, with the relative offset from sb's pointer
	// field to the target AnyPointer offset.
	//
	// FIXME: needs overflow and bounds check?
	ptr := v.ptr.withDataOffset(v.off - concretePtrOff - 1)

	// Write the pointer field.
	sb.urb.SetWord(ptrOff, Word(ptr))
	return nil
}

func (sb *StructBuilder) NewStructListField(ptrIndex PointerFieldIndex, itemSize StructSize, listLen, listCap int) (res StructListBuilder, err error) {
	if !sb.hasPointer(ptrIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainPointerField(ptrIndex)
		return
	}

	res, err = sb.mb.NewStructList(itemSize, listLen, listCap)
	if err != nil {
		return
	}

	if res.sid != sb.sid {
		// TODO: add support.
		err = errors.New("unhandled list in different segment")
		return
	}

	// Offset of the pointer field inside sb. No need to check for overflow
	// because the struct has been validated to contain this pointer.
	ptrOff := ptrIndex.uncheckedWordOffset(WordOffset(sb.sz.DataSectionSize))

	// Determine concrete offset from structure start until ptrOff.
	concretePtrOff := sb.off + ptrOff

	// Build the final list pointer.
	lsPtr := listPointer{
		startOffset: res.off - concretePtrOff - 1,
		elSize:      listElSizeComposite,
		listSize:    listSize(listWordCount(listElSizeComposite, res.listLen)),
	}

	// Write the pointer field.
	sb.urb.SetWord(ptrOff, Word(lsPtr.toPointer()))
	return
}

type SegmentBuilder struct {
	mb *MessageBuilder
	b  *[]byte
	id SegmentID
}

func (sb *SegmentBuilder) ID() SegmentID {
	return sb.id
}

// SetWord sets the word at the given offset in the segment. This must
// only be called when the caller is sure the given word is already allocated in
// the segment.
func (sb *SegmentBuilder) SetWord(offset WordOffset, value Word) {
	// binary.LittleEndian.PutUint64(sb.as.uncheckedSegSlice(sb.id, offset, 1), uint64(value))
	binary.LittleEndian.PutUint64((*sb.b)[offset*WordSize:], uint64(value))
	// binary.LittleEndian.PutUint64(sb.b[offset*WordSize:(offset+1)*WordSize], uint64(value))
	// *(*Word)(unsafe.Add(sb.ptr, offset*WordSize)) = value
}

func (sb *SegmentBuilder) maskAndMergeWord(offset WordOffset, mask, value Word) {
	old := binary.LittleEndian.Uint64((*sb.b)[offset*WordSize:])
	binary.LittleEndian.PutUint64((*sb.b)[offset*WordSize:], old&^uint64(mask)|uint64(value))

	/*
		ptr := (*Word)(unsafe.Add(sb.ptr, offset*WordSize))
		*ptr = *ptr&mask | value
	*/

	/*
		old := *(*Word)(unsafe.Add(sb.ptr, offset*WordSize))
		*(*Word)(unsafe.Add(sb.ptr, offset*WordSize)) = old&mask | value
	*/

	/*
		buf := (*sb.b)[offset*WordSize:]
		old := binary.LittleEndian.Uint64(buf)
		binary.LittleEndian.PutUint64(buf, old&uint64(mask)|uint64(value))
	*/

	/*
		ptr := (*uint64)(unsafe.Add(sb.ub, offset*8))
		old := *ptr
		*ptr = old&uint64(mask) | uint64(value)
	*/
}

func (sb *SegmentBuilder) GetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64((*sb.b)[offset*WordSize:]))
	// return *(*Word)(unsafe.Add(sb.ptr, offset*WordSize))
}

/*
// uncheckedSegSlice slices part of a segment (aligned to a word) without
// checking for valid bounds.
//
// This must only be called when the assumption holds that the bounds have
// already been validated.
func (sb *SegmentBuilder) uncheckedSegSlice(offset WordOffset, size WordCount) []byte {
	return (*sb.b)[offset*WordSize : (offset+WordOffset(size))*WordSize]
}
*/

// copyStringTo copies s into the segment, starting at the given offset.
func (sb *SegmentBuilder) copyStringTo(offset WordOffset, s string) {
	copy((*sb.b)[offset*WordSize:], s)
	// copy(unsafe.Slice((*byte)(unsafe.Add(sb.ptr, offset*WordSize)), len(s)), s)
}

type Allocator interface {
	Init(state *AllocState) (err error)
	Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error)
	Reset(state *AllocState) (err error)
}

type AllocState struct {
	HeaderBuf []byte
	FirstSeg  []byte
	Segs      [][]byte
	Extra     any

	// firstSegPtr unsafe.Pointer
}

func (as *AllocState) GetHeader() []byte {
	return as.HeaderBuf
}

func (as *AllocState) GetSeg0() []byte {
	return as.FirstSeg
}

func (as *AllocState) SetSeg0(b []byte) {
	as.FirstSeg = b
	// as.firstSegPtr = unsafe.Pointer(unsafe.SliceData(as.FirstSeg))
}

func (as *AllocState) SetHeaderAndSeg0(buf []byte, expectedSegCount SegmentCount) {
	headerSize := alignToWord(4 + Word(expectedSegCount)*4)
	as.HeaderBuf = buf[:headerSize]
	as.FirstSeg = buf[headerSize:]
	// as.firstSegPtr = unsafe.Pointer(unsafe.SliceData(as.FirstSeg))
}

func (as *AllocState) GetSeg(id SegmentID) []byte {
	if id == 0 {
		return as.FirstSeg
	} else {
		return as.Segs[id]
	}
}

// ValidAfterInitReset checks if the AllocState is valid after having been
// (re-)initialized by an [Allocator] Init() or Reset() call.
func (as *AllocState) ValidAfterInitReset() (err error) {
	if len(as.FirstSeg) < WordSize {
		err = errAllocNoRootWord
	}
	return nil
}

/*
// uncheckedSegTailSlice returns the tail slice of a segment (aligned to a word)
// without checking for valid bounds.
//
// This must only be called when the assumption holds that the bounds have
// already been validated.
func (as *AllocState) uncheckedSegTailSlice(seg SegmentID, offset WordOffset) []byte {
	return as.Segs[seg][offset*WordSize:]
}

// uncheckedSegSlice slices part of a segment (aligned to a word) without
// checking for valid bounds.
//
// This must only be called when the assumption holds that the bounds have
// already been validated.
func (as *AllocState) uncheckedSegSlice(seg SegmentID, offset WordOffset, size WordCount) []byte {
	return as.Segs[seg][offset*WordSize : (offset+WordOffset(size))*WordSize]
}
*/

// headerBufPrefixesSeg0Buf returns true if the underlying array in HeaderBuf
// exactly prefixes the underlying array of Segs[0].
//
// This is used to detect if HeaderBuf and Seg[0] have been allocated in such a
// way as to be the same underlying array.
//
// Note: this can be called only after checking that HeaderBuf != nil, that
// len(Segs) > 0 and that len(Segs[0]) > 0, otherwise it panics.
func (as *AllocState) headerBufPrefixesSeg0Buf() bool {
	// The two buffers are contiguous (i.e. seg0 is aliased on the same
	// underlying array as headerBuf) if the data for seg0 starts
	// immediately after the data for the header (i.e. the pointer address
	// for seg0 data is exactly len(headerBuf) bytes after the pointer
	// address for headerBuf) and the headerBuf slice could be extended
	// towards the seg0 buffer.
	headerBufPtr := unsafe.Pointer(unsafe.SliceData(as.HeaderBuf))
	seg0BufPtr := unsafe.Pointer(unsafe.SliceData(as.FirstSeg))
	return cap(as.HeaderBuf) >= len(as.HeaderBuf)+len(as.FirstSeg) &&
		seg0BufPtr == unsafe.Add(headerBufPtr, len(as.HeaderBuf))
}

// putSingleSegHeaderInBuf writes the framing header in headerBuf for the case
// where a single segment is used.
//
// This must only be called in the single segment case, after ensuring the
// header buf has enough room for the header.
func (as *AllocState) putSingleSegHeaderInBuf() {
	seg0size := uint64(len(as.FirstSeg))
	if seg0size > maxValidBytes {
		// This should never happen for correctly implemented
		// allocators.
		panic("allocator allocated single segment too large")
	}

	// clear(as.HeaderBuf[:4]) // Segment count is all zeroes (==1 segment).
	// binary.LittleEndian.PutUint32(as.HeaderBuf[4:], uint32(len(as.FirstSeg)/WordSize))

	// The capnp spec interprets the first (Q)word as two little-endian
	// D(words). So we shift the target segment count 32 bits to the left,
	// so that when converted to little-endian, it ends up in the correct
	// location. The shift naturally zeroes the LSBs, which clears any
	// leftover data in segment count DWORD (signifying one segment).
	binary.LittleEndian.PutUint64(as.HeaderBuf, (seg0size/WordSize)<<32)
	// *(*uint64)(as.firstSegPtr) = (seg0size / WordSize) << 32
}

type MessageBuilder struct {
	state       AllocState
	alloc       Allocator
	segsCap     int
	readerArena Arena
}

func NewMessageBuilder(alloc Allocator) (mb *MessageBuilder, err error) {
	mb = &MessageBuilder{alloc: alloc}
	if err := alloc.Init(&mb.state); err != nil {
		return nil, err
	}
	if err := mb.state.ValidAfterInitReset(); err != nil {
		return nil, err
	}
	mb.segsCap = cap(mb.state.Segs)
	return mb, nil
}

func (mb *MessageBuilder) Reset() error {
	if err := mb.alloc.Reset(&mb.state); err != nil {
		return err
	}
	if err := mb.state.ValidAfterInitReset(); err != nil {
		return err
	}
	mb.segsCap = cap(mb.state.Segs)
	return nil
}

// MessageReader returns a [Message] to read data already written to the
// builder.
//
// NOTE: the returned message and its structures are only guaranteed to be valid
// until the next modification of the message being built. It should be
// discarded and re-read if any modifications are made, otherwise orphaned data
// may be read.
func (mb *MessageBuilder) MessageReader() Message {
	mb.readerArena.notResetable = true
	mb.readerArena.s.b = mb.state.FirstSeg
	mb.readerArena.ReadLimiter().InitNoLimit()
	if len(mb.state.Segs) > 0 {
		panic("TODO: handle multiseg")
	}
	return Message{arena: &mb.readerArena, dl: maxDepthLimit}
}

// allocate allocates size words, preferably (but not necessarily) on the
// preferred segment.
//
// If size has already been validated to be a valid word count, use
// allocateValidSize().
func (mb *MessageBuilder) allocate(preferred SegmentID, size WordCount) (segb SegmentBuilder, offset WordOffset, err error) {
	if size > MaxValidWordCount {
		return SegmentBuilder{}, 0, errAllocOverMaxWordCount
	}
	return mb.allocateValidSize(preferred, size)
}

// allocateValidSize allocates size words, preferably (but not necessarily) on
// the preferred segment.
//
// This does NOT validate that size is a valid word count.
func (mb *MessageBuilder) allocateValidSize(preferred SegmentID, size WordCount) (segb SegmentBuilder, offset WordOffset, err error) {
	// Ask the allocator to allocate.
	segb.id, offset, err = mb.alloc.Allocate(&mb.state, preferred, size)
	if err != nil {
		return
	}

	// This assertion is necessary because SegmentBuilders track the
	// segment buffers by pointers into mb.state.Segs. Changing the
	// capacity (but _not_ the length) would invalidate such pointers
	// (because of the reallocation of the Segs slice). Thus we impose this
	// restriction on allocators, that they must define at init time the
	// max number of segments they are likely to use (while actual usage is
	// still dynamic, given by the length of Segs).
	if cap(mb.state.Segs) != mb.segsCap {
		return SegmentBuilder{}, 0, errCannotChangeSegsCap
	}

	if segb.id == 0 {
		segb.b = &mb.state.FirstSeg
		// segb.ptr = mb.state.firstSegPtr
	} else {
		segb.b = &mb.state.Segs[segb.id-1]
	}

	// Sanity check allocator didn't do something silly.
	lenB := len(*segb.b)
	if lenB > maxValidBytes {
		return SegmentBuilder{}, 0, errAllocatedTooLargeSeg
	}
	if !isWordAligned(lenB) {
		return SegmentBuilder{}, 0, errAllocatedUnalignedSeg
	}
	if endOff, ok := addWordOffsets(offset, WordOffset(size)); !ok || endOff > WordOffset(lenB/WordSize) {
		return SegmentBuilder{}, 0, errAllocatedOutOfRange
	}

	// All good.
	segb.mb = mb
	return
}

func (mb *MessageBuilder) allocateValidSizeXXX(preferred SegmentID, size WordCount) (sid SegmentID, b []byte, offset WordOffset, err error) {
	// Ask the allocator to allocate.
	sid, offset, err = mb.alloc.Allocate(&mb.state, preferred, size)
	if err != nil {
		return
	}

	// This assertion is necessary because SegmentBuilders track the
	// segment buffers by pointers into mb.state.Segs. Changing the
	// capacity (but _not_ the length) would invalidate such pointers
	// (because of the reallocation of the Segs slice). Thus we impose this
	// restriction on allocators, that they must define at init time the
	// max number of segments they are likely to use (while actual usage is
	// still dynamic, given by the length of Segs).
	if cap(mb.state.Segs) != mb.segsCap {
		err = errCannotChangeSegsCap
		return
	}

	if sid == 0 {
		b = mb.state.FirstSeg
		// segb.ptr = mb.state.firstSegPtr
	} else {
		b = mb.state.Segs[sid-1]
	}

	// Sanity check allocator didn't do something silly.
	lenB := len(b)
	if lenB > maxValidBytes {
		err = errAllocatedTooLargeSeg
		return
	}
	if !isWordAligned(lenB) {
		err = errAllocatedUnalignedSeg
		return
	}
	if endOff, ok := addWordOffsets(offset, WordOffset(size)); !ok || endOff > WordOffset(lenB/WordSize) {
		err = errAllocatedOutOfRange
		return
	}

	// All good.
	return
}
func (mb *MessageBuilder) SetRoot(sb *StructBuilder) error {
	// NewMessageBuilder() ensures the allocator returns at least one
	// segment with at least enough room for the root pointer.
	if mb == nil || mb.state.FirstSeg == nil || len(mb.state.FirstSeg) < WordSize {
		return errAllocStateNoRootWord
	}
	// if sb.seg.mb != mb {
	if sb.mb != mb {
		return fmt.Errorf("sb.mb vs mb: %w", errDifferentMsgBuilders)
	}

	// TODO: handle inter-segment pointer.
	// if sb.seg.id != 0 {
	if sb.sid != 0 {
		panic("needs handling")
	}

	// Write the first word of seg0.
	sp := structPointer{
		// The wire offset is relative to the end of first word.
		dataOffset:         sb.off - 1,
		dataSectionSize:    sb.sz.DataSectionSize,
		pointerSectionSize: sb.sz.PointerSectionSize,
	}
	binary.LittleEndian.PutUint64(mb.state.FirstSeg, uint64(sp.toPointer()))

	return nil
}

func (mb *MessageBuilder) NewStruct(size StructSize) (sb StructBuilder, err error) {
	// TotalSize() is necessarily a valid size because it is only up to
	// 2^17-2 words.
	/*
		sid, b, off, err := mb.allocateValidSizeXXX(0, size.TotalSize())
		if err != nil {
			return StructBuilder{}, err
		}
		urb := UnsafeRawBuilder{ptr: unsafe.Pointer(&(b[off*WordSize]))}
	*/

	urb, off, err := mb.allocateValidSize(0, size.TotalSize())
	if err != nil {
		return StructBuilder{}, err
	}

	return StructBuilder{
		// seg: seg,
		mb:  mb,
		sid: urb.id, // sid, // seg.id,
		urb: urb,
		off: off,
		sz:  size,
	}, nil
}

// NewRootStruct initializes a new struct and sets it as the root struct on this
// message.
func (mb *MessageBuilder) NewRootStruct(size StructSize) (sb StructBuilder, err error) {
	sb, err = mb.NewStruct(size)
	if err == nil {
		err = mb.SetRoot(&sb)
	}
	return
}

func (mb *MessageBuilder) NewStructList(itemSize StructSize, listLen, listCap int) (res StructListBuilder, err error) {
	// Sanity checks.
	if listCap < listLen {
		err = errors.New("listCap cannot be < listLen")
		return
	}
	if listCap > MaxListSize {
		err = errors.New("list capacity cannot be larger than MaxListSize")
		return
	}
	totalWordCount := uint64(itemSize.TotalSize())*uint64(listCap) + 1 // + 1 tag word
	if totalWordCount > maxWordOffset {
		err = fmt.Errorf("total struct size %d overflows max word offset %d", totalWordCount, maxWordOffset)
		return
	}

	var b []byte

	res.sid, b, res.off, err = mb.allocateValidSizeXXX(0, WordCount(totalWordCount))
	if err != nil {
		return
	}

	res.urb = UnsafeRawBuilder{ptr: unsafe.Pointer(&(b[res.off*WordSize]))}
	res.listCap = listSize(listCap)
	res.listLen = listSize(listLen)
	res.mb = mb
	res.itemSize = itemSize

	// Write the initial tag word.
	res.writeTagWord()
	return
}

// newText allocates and places s as a new text in the meesage. The text is
// preferably (but not necessarily) put into segment preferSeg.
func (mb *MessageBuilder) newText(preferSeg SegmentID, s string) ( /*segb SegmentBuilder, ptr listPointer*/ sid SegmentID, ptr listPointer, err error) {
	// Length of texts (strings) in capnp is +1 due to null at the end.
	textLen := uint(len(s) + 1)
	if textLen > MaxListSize {
		// return SegmentBuilder{}, listPointer{}, errStringTooLarge
		return 0, listPointer{}, errStringTooLarge
	}

	words := WordCount(uintBytesToWordAligned(textLen))
	/*
		segb, off, err := mb.allocateValidSize(preferSeg, words)
		if err != nil {
			return SegmentBuilder{}, listPointer{}, err
		}
		segb.copyStringTo(off, s)
	*/
	sid, b, off, err := mb.allocateValidSizeXXX(preferSeg, words)
	if err != nil {
		// return SegmentBuilder{}, listPointer{}, err
		return 0, listPointer{}, err
	}
	copy(b[off*WordSize:], s)

	return sid,
		listPointer{
			startOffset: off,
			elSize:      listElSizeByte,
			listSize:    listSize(textLen),
		}, nil
}

// Allocate a number of words into a segment.
func (mb *MessageBuilder) Allocate(size WordCount) (segb SegmentBuilder, offset WordOffset, err error) {
	return mb.allocate(0, size)
}

/*
func (mb *MessageBuilder) AllocateRawBuilder(size WordCount) (rb RawBuilder, err error) {
	// Ask the allocator to allocate.
	var segID SegmentID
	var off WordOffset
	segID, off, err = mb.alloc.Allocate(&mb.state, 0, size)
	if err != nil {
		return
	}

	var b *[]byte
	if segID == 0 {
		b = &mb.state.FirstSeg
	} else {
		b = &mb.state.Segs[segID-1]
	}

	rb.b = unsafe.Slice((*Word)(unsafe.Pointer(&(*b)[off*WordSize])), size)
	return
}
*/

// AllocateRootRawBuilder allocates a raw builder for the root struct. The
// returned builder is aliased over the entire allocated size, including the
// root struct pointer of the message.
//
// This should only be called on empty messages, otherwise it errors.
func (mb *MessageBuilder) AllocateRootRawBuilder(size WordCount) (rb RawBuilder, err error) {
	// Ask the allocator to allocate.
	var segID SegmentID
	var off WordOffset
	segID, off, err = mb.alloc.Allocate(&mb.state, 0, size)
	if err != nil {
		return
	}

	if segID != 0 {
		err = errRootNotOnSeg0
		return
	}
	if off != 1 {
		err = errRootNotAtOffset1
		return
	}

	b := &mb.state.FirstSeg
	rb.b = unsafe.Slice((*Word)(unsafe.Pointer(&(*b)[0])), size+1)
	return
}

func (mb *MessageBuilder) AllocateUnsafeRootRawBuilder(size WordCount) (rb UnsafeRawBuilder, err error) {
	// Ask the allocator to allocate.
	var segID SegmentID
	var off WordOffset
	segID, off, err = mb.alloc.Allocate(&mb.state, 0, size)
	if err != nil {
		return
	}

	if segID != 0 {
		err = errRootNotOnSeg0
		return
	}
	if off != 1 {
		err = errRootNotAtOffset1
		return
	}

	b := &mb.state.FirstSeg
	rb.ptr = unsafe.Pointer(&(*b)[off*WordSize])
	return
}

func (mb *MessageBuilder) Serialize() ([]byte, error) {
	if len(mb.state.FirstSeg) == 0 {
		return nil, errMsgBuilderNoSegData
	}

	// Special case where the allocator allocated both the header buffer
	// and the single segment data in the same contiguous buffer. In this
	// case, the data is already fully framed and serialized.
	//
	// A single segment header is 4 bytes segment count (== 0) and 4 bytes
	// segment size (in words).
	const singleSegmentHeaderSize = 8
	if len(mb.state.HeaderBuf) == singleSegmentHeaderSize &&
		len(mb.state.Segs) == 0 &&
		mb.state.headerBufPrefixesSeg0Buf() {
		mb.state.putSingleSegHeaderInBuf() // Write single segment header
		return mb.state.HeaderBuf[:len(mb.state.HeaderBuf)+len(mb.state.FirstSeg)], nil
	}

	// TODO: proceed to standard framing and serialization.
	panic("boo")
}

func (mb *MessageBuilder) SerializeXXX() ([]byte, error) {
	mb.state.putSingleSegHeaderInBuf() // Write single segment header
	return mb.state.HeaderBuf[:len(mb.state.HeaderBuf)+len(mb.state.FirstSeg)], nil
}
