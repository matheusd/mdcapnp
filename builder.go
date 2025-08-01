// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type ListBuilder struct {
	seg SegmentBuilder

	elSize   listElementSize
	listSize listSize

	// sz  ListSize
	off WordOffset
}

type StructBuilder struct {
	off WordOffset // Concrete offset into segment where struct data begins.
	sz  StructSize
	seg SegmentBuilder
}

func (sb *StructBuilder) hasData(dataIndex DataFieldIndex) bool {
	return dataIndex < DataFieldIndex(sb.sz.DataSectionSize)
}

func (sb *StructBuilder) hasPointer(ptrIndex PointerFieldIndex) bool {
	return ptrIndex < PointerFieldIndex(sb.sz.PointerSectionSize)
}

func (sb *StructBuilder) SetInt64(dataIndex DataFieldIndex, v int64) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		finalOff := dataIndex.uncheckedWordOffset(sb.off)
		sb.seg.uncheckedSetWord(finalOff, Word(v))
	}
	return
}

func (sb *StructBuilder) SetInt32(dataIndex DataFieldIndex, mask Int32DataFieldSetMask, v int32) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), Word(mask), Word(v))
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
		sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), ^(1 << bit), boolToWord(v)<<bit)
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
		sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(math.Float64bits(v)))
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
	segb, lsPtr, err := sb.seg.mb.newText(sb.seg.id, v)
	if err != nil {
		return err
	}

	// TODO: handle allocs in new segments.
	if segb.id != sb.seg.id {
		return errors.New("needs handling")
	}

	// Determine concrete pointer offset inside struct. This doesn't need
	// overflow checks because the entire struct has been allocated, thus
	// this pointer offset is known to be in bounds.
	ptrOff := ptrIndex.uncheckedWordOffset(sb.off + WordOffset(sb.sz.DataSectionSize))

	// Determine the relative offset from the field pointer offset to the
	// actual data. This finishes the construction of the list pointer.
	// lsPtr.startOffset = lsPtr.startOffset - ptrOff - 1
	lsPtr.startOffset = lsPtr.startOffset - ptrOff - 1

	// Structure already fully allocated, no need to check for
	// bounds.
	sb.seg.uncheckedSetWord(ptrOff, Word(lsPtr.toPointer()))
	return nil
}

type SegmentBuilder struct {
	mb *MessageBuilder
	b  *[]byte
	id SegmentID
}

// uncheckedSetWord sets the word at the given offset in the segment. This must
// only be called when the caller is sure the given word is already allocated in
// the segment.
func (sb *SegmentBuilder) uncheckedSetWord(offset WordOffset, value Word) {
	// binary.LittleEndian.PutUint64(sb.as.uncheckedSegSlice(sb.id, offset, 1), uint64(value))
	binary.LittleEndian.PutUint64((*sb.b)[offset*WordSize:], uint64(value))
	// binary.LittleEndian.PutUint64(sb.b[offset*WordSize:(offset+1)*WordSize], uint64(value))
}

func (sb *SegmentBuilder) uncheckedMaskAndMergeWord(offset WordOffset, mask, value Word) {
	old := binary.LittleEndian.Uint64((*sb.b)[offset*WordSize:])
	binary.LittleEndian.PutUint64((*sb.b)[offset*WordSize:], old&uint64(mask)|uint64(value))
}

func (sb *SegmentBuilder) uncheckedGetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64((*sb.b)[offset*WordSize:]))
}

// uncheckedSegSlice slices part of a segment (aligned to a word) without
// checking for valid bounds.
//
// This must only be called when the assumption holds that the bounds have
// already been validated.
func (sb *SegmentBuilder) uncheckedSegSlice(offset WordOffset, size WordCount) []byte {
	return (*sb.b)[offset*WordSize : (offset+WordOffset(size))*WordSize]
}

// copyStringTo copies s into the segment, starting at the given offset.
func (sb *SegmentBuilder) copyStringTo(offset WordOffset, s string) {
	copy((*sb.b)[offset*WordSize:], s)
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
}

// ValidAfterInitReset checks if the AllocState is valid after having been
// (re-)initialized by an [Allocator] Init() or Reset() call.
func (as *AllocState) ValidAfterInitReset() error {
	if len(as.FirstSeg) == 0 {
		return errAllocNoFirstSeg
	}
	if len(as.FirstSeg) < WordSize {
		return errAllocNoRootWord
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
	if seg0size > MaxValidWordCount*WordSize {
		// This should never happen for correctly implemented
		// allocators.
		panic("allocator allocated single segment too large")
	}
	clear(as.HeaderBuf[:4]) // Segment count is all zeroes.
	binary.LittleEndian.PutUint32(as.HeaderBuf[4:], uint32(len(as.FirstSeg)/WordSize))
}

type MessageBuilder struct {
	state AllocState
	alloc Allocator
}

func NewMessageBuilder(alloc Allocator) (mb *MessageBuilder, err error) {
	mb = &MessageBuilder{alloc: alloc}
	if err := alloc.Init(&mb.state); err != nil {
		return nil, err
	}
	if err := mb.state.ValidAfterInitReset(); err != nil {
		return nil, err
	}
	return mb, nil
}

func (mb *MessageBuilder) Reset() error {
	if err := mb.alloc.Reset(&mb.state); err != nil {
		return err
	}
	if err := mb.state.ValidAfterInitReset(); err != nil {
		return err
	}
	return nil
}

// allocate allocates size words, preferably (but not necessarily) on the
// preferred segment.
//
// If size has already been validated to be a valid word count, use
// allocateValidSize().
func (mb *MessageBuilder) allocate(preferred SegmentID, size WordCount) (segb SegmentBuilder, offset WordOffset, err error) {
	if size > MaxValidWordCount {
		err = errAllocOverMaxWordCount
		return
	}
	return mb.allocateValidSize(preferred, size)
}

// allocateValidSize allocates size words, preferably (but not necessarily) on
// the preferred segment.
//
// This does NOT validate that size is a valid word count.
func (mb *MessageBuilder) allocateValidSize(preferred SegmentID, size WordCount) (segb SegmentBuilder, offset WordOffset, err error) {
	oldSegsCap := cap(mb.state.Segs)

	// Ask the allocator to allocate.
	segb.id, offset, err = mb.alloc.Allocate(&mb.state, preferred, size)
	if err != nil {
		return
	}

	// This assertion is necessary because SegmentBuilders track the segment
	// buffers by pointers into mb.state.Segs. Changing the capacity (but
	// _not_ the length) would invalidate such pointers (because of the
	// reallocation of the Segs slice). Thus we impose this restriction on
	// allocators, that they must define at init time the max number of
	// segments they are likely to use (while actual usage is still dynamic,
	// given by the length of Segs).
	if cap(mb.state.Segs) != oldSegsCap {
		return SegmentBuilder{}, 0, errCannotChangeSegsCap
	}

	if segb.id == 0 {
		segb.b = &mb.state.FirstSeg
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

func (mb *MessageBuilder) SetRoot(sb *StructBuilder) error {
	// NewMessageBuilder() ensure the allocator returns at least one segment
	// with at least enough room for the root pointer.
	if mb == nil || mb.state.FirstSeg == nil || len(mb.state.FirstSeg) < WordSize {
		return errAllocStateNoRootWord
	}
	if sb.seg.mb != mb {
		return fmt.Errorf("sb.mb vs mb: %w", errDifferentMsgBuilders)
	}

	// TODO: handle inter-segment pointer.
	if sb.seg.id != 0 {
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
	seg, off, err := mb.allocateValidSize(0, size.TotalSize())
	if err != nil {
		return StructBuilder{}, err
	}

	return StructBuilder{
		seg: seg,
		off: off,
		sz:  size,
	}, nil
}

func (mb *MessageBuilder) newText(preferSeg SegmentID, s string) (segb SegmentBuilder, ptr listPointer, err error) {
	// Length of texts (strings) in capnp is +1 due to null at the end.
	textLen := uint(len(s) + 1)
	if textLen > MaxListSize {
		return SegmentBuilder{}, listPointer{}, errStringTooLarge
	}

	words := WordCount(uintBytesToWordAligned(textLen))
	segb, off, err := mb.allocateValidSize(preferSeg, words)
	if err != nil {
		return SegmentBuilder{}, listPointer{}, err
	}
	segb.copyStringTo(off, s)

	return segb,
		listPointer{
			startOffset: off,
			elSize:      listElSizeByte,
			listSize:    listSize(textLen),
		}, nil
}

func (mb *MessageBuilder) Serialize() ([]byte, error) {
	if mb.state.FirstSeg == nil {
		return nil, errMsgBuilderNoSegments
	}
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
