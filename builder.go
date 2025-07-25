// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type StructBuilder struct {
	off WordOffset // Concrete offset into segment where struct data begins.
	sz  StructSize
	seg SegmentBuilder
}

func (sb *StructBuilder) hasData(dataIndex DataFieldIndex) bool {
	return dataIndex < DataFieldIndex(sb.sz.DataSectionSize)
}

func (sb *StructBuilder) SetInt64(dataIndex DataFieldIndex, v int64) (err error) {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		err = errStructBuilderDoesNotContainDataField(dataIndex)
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(v))
	}
	return
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

func (sb *SegmentBuilder) uncheckedGetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64((*sb.b)[offset*WordSize:]))
}

type Allocator interface {
	Init(state *AllocState) (err error)
	Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error)
	Reset(state *AllocState) (err error)
}

type AllocState struct {
	HeaderBuf []byte
	Segs      [][]byte
	Extra     any
}

// ValidAfterInitReset checks if the AllocState is valid after having been
// (re-)initialized by an [Allocator] Init() or Reset() call.
func (as *AllocState) ValidAfterInitReset() error {
	if len(as.Segs) < 1 {
		return errAllocNoFirstSeg
	}
	if len(as.Segs[0]) < WordSize {
		return errAllocNoRootWord
	}
	return nil
}

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
	seg0BufPtr := unsafe.Pointer(unsafe.SliceData(as.Segs[0]))
	return cap(as.HeaderBuf) >= len(as.HeaderBuf)+len(as.Segs[0]) &&
		seg0BufPtr == unsafe.Add(headerBufPtr, len(as.HeaderBuf))
}

// putSingleSegHeaderInBuf writes the framing header in headerBuf for the case
// where a single segment is used.
//
// This must only be called in the single segment case, after ensuring the
// header buf has enough room for the header.
func (as *AllocState) putSingleSegHeaderInBuf() {
	seg0size := uint64(len(as.Segs[0]))
	if seg0size > MaxValidWordCount*WordSize {
		// This should never happen for correctly implemented
		// allocators.
		panic("allocator allocated single segment too large")
	}
	clear(as.HeaderBuf[:4]) // Segment count is all zeroes.
	binary.LittleEndian.PutUint32(as.HeaderBuf[4:], uint32(len(as.Segs[0])/WordSize))
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

func (mb *MessageBuilder) allocate(preferred SegmentID, size WordCount) (SegmentBuilder, WordOffset, error) {
	if size > MaxValidWordCount {
		return SegmentBuilder{}, 0, errAllocOverMaxWordCount
	}
	oldSegsCap := cap(mb.state.Segs)

	// Ask the allocator to allocate.
	segID, offset, err := mb.alloc.Allocate(&mb.state, preferred, size)
	if err != nil {
		return SegmentBuilder{}, 0, err
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

	// All good.
	return SegmentBuilder{
		id: segID,
		mb: mb,
		b:  &mb.state.Segs[segID],
	}, offset, nil
}

func (mb *MessageBuilder) SetRoot(sb *StructBuilder) error {
	// NewMessageBuilder() ensure the allocator returns at least one segment
	// with at least enough room for the root pointer.
	if mb == nil || len(mb.state.Segs) == 0 || len(mb.state.Segs[0]) < WordSize {
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
	binary.LittleEndian.PutUint64(mb.state.Segs[0], uint64(sp.toPointer()))

	return nil
}

func (mb *MessageBuilder) NewStruct(size StructSize) (sb StructBuilder, err error) {
	seg, off, err := mb.allocate(0, size.TotalSize())
	if err != nil {
		return StructBuilder{}, err
	}

	return StructBuilder{
		seg: seg,
		off: off,
		sz:  size,
	}, nil
}

func (mb *MessageBuilder) Serialize() ([]byte, error) {
	if len(mb.state.Segs) == 0 {
		return nil, errMsgBuilderNoSegments
	}
	if len(mb.state.Segs[0]) == 0 {
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
		len(mb.state.Segs) == 1 &&
		mb.state.headerBufPrefixesSeg0Buf() {
		mb.state.putSingleSegHeaderInBuf() // Write single segment header
		return mb.state.HeaderBuf[:len(mb.state.HeaderBuf)+len(mb.state.Segs[0])], nil
	}

	// TODO: proceed to standard framing and serialization.
	panic("boo")
}
