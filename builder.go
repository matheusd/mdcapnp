// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"errors"
	"unsafe"
)

type StructBuilder struct {
	mb  *MessageBuilder
	seg SegmentBuilder
	off WordOffset // Concrete offset into segment where struct data begins.
	sz  StructSize
}

func (sb *StructBuilder) hasData(dataIndex DataFieldIndex) bool {
	return dataIndex < DataFieldIndex(sb.sz.DataSectionSize)
}

func (sb *StructBuilder) SetInt64(dataIndex DataFieldIndex, v int64) error {
	if !sb.hasData(dataIndex) {
		// TODO: allocate new struct, copy over old fields to new fields
		// or error out?
		return errors.New("cannot resize struct")
	} else {
		// Structure already fully allocated, no need to check for
		// bounds.
		sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(v))
		return nil
	}
}

type SegmentBuilder struct {
	mb *MessageBuilder
	id SegmentID
}

// uncheckedSetWord sets the word at the given offset in the segment. This must
// only be called when the caller is sure the given word is already allocated in
// the segment.
func (sb *SegmentBuilder) uncheckedSetWord(offset WordOffset, value Word) {
	byteOffset := int(offset * WordSize)
	binary.LittleEndian.PutUint64(sb.mb.state.Segs[sb.id][byteOffset:], uint64(value))
}

type Allocator interface {
	Init() (initState AllocState, err error)
	Allocate(prevState AllocState, preferred SegmentID, size WordCount) (nextState AllocState, seg SegmentID, off WordOffset, err error)
	Reset(lastState AllocState) (blankState AllocState, err error)
}

type AllocState struct {
	HeaderBuf []byte
	Segs      [][]byte
	Extra     any
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
	mb = new(MessageBuilder)
	if mb.state, err = alloc.Init(); err != nil {
		return nil, err
	}
	if len(mb.state.Segs) < 1 {
		return nil, errAllocInitNoFirstSeg
	}
	if len(mb.state.Segs[0]) < WordSize {
		return nil, errAllocInitNoRootWord
	}
	mb.alloc = alloc
	return mb, nil
}

func (mb *MessageBuilder) Reset() error {
	newState, err := mb.alloc.Reset(mb.state)
	if err != nil {
		return err
	}
	if len(newState.Segs) < 1 {
		return errAllocInitNoFirstSeg
	}
	if len(newState.Segs[0]) < WordSize {
		return errAllocInitNoRootWord
	}
	mb.state = newState
	return nil
}

func (mb *MessageBuilder) allocate(preferred SegmentID, size WordCount) (SegmentBuilder, WordOffset, error) {
	if size > MaxValidWordCount {
		return SegmentBuilder{}, 0, errAllocOverMaxWordCount
	}
	newState, segID, offset, err := mb.alloc.Allocate(mb.state, preferred, size)
	if err != nil {
		return SegmentBuilder{}, 0, err
	}
	mb.state = newState
	return SegmentBuilder{id: segID, mb: mb}, offset, nil
}

func (mb *MessageBuilder) SetRoot(sb *StructBuilder) error {
	// NewMessageBuilder() ensure the allocator returns at least one segment
	// with at least enough room for the root pointer.
	if mb == nil || len(mb.state.Segs) == 0 || len(mb.state.Segs[0]) < WordSize {
		return errAllocStateNoRootWord
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
	binary.LittleEndian.PutUint64(sb.mb.state.Segs[0], uint64(sp.toPointer()))

	return nil
}

func (mb *MessageBuilder) NewStruct(size StructSize) (sb StructBuilder, err error) {
	seg, off, err := mb.allocate(0, size.TotalSize())
	if err != nil {
		return StructBuilder{}, err
	}

	return StructBuilder{
		mb:  mb,
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
