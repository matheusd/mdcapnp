// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
)

type Word uint64

const WordSize = 8

type WordCount uint64

type wordCount16 uint16

// PointerFieldIndex is the index of a pointer field in a struct. The first
// pointer of a struct's pointer section has index 0, the second one has index
// 1, and so on.
//
// Given that 1 pointer == 1 word, and a struct is limited to 16 bits worth of
// words in its pointer section, it can have up to 2^16 pointer fields.
type PointerFieldIndex uint16

// DataFieldIndex is the index of a data field in a struct (in words). When the
// field is smaller than one word, then further indexing may be necessary to
// extract its value.
type DataFieldIndex uint16

type WordOffset uint64

type SignedWordOffset int64

type listOrStructOffset int32 // Up to 30 bits usable.

type ByteOffset uint64

type ByteCount uint64

type SegmentID uint64

type Arena interface {
	Segment(id SegmentID) (*Segment, error)
}

type Segment struct {
	b  []byte
	rl *ReadLimiter
}

func (ms *Segment) GetWord(offset WordOffset) (res Word, err error) {
	if err = ms.rl.CanRead(1); err != nil {
	} else if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.LittleEndian.Uint64(ms.b[byteOffset:]))

		// copy((*[8]byte)(unsafe.Pointer(&res))[:], ms.b[byteOffset:])

		// Assumes a big endian version is written. Note: this is
		// counterintuitive, double check.
		// res = Word(binary.BigEndian.Uint64((*[8]byte)(unsafe.Pointer(&res))[:]))
	}
	return
}

// checkSliceBounds checks whether a subsequent call to [uncheckedSlice] with
// the same arguments will fail. If this function returns true, immediately
// calling [uncheckedSlice] will generate a valid slice.
func (ms *Segment) checkSliceBounds(offset WordOffset, size ByteCount) error {
	if err := ms.rl.CanRead(WordCount(size) / WordSize); err != nil {
		return err
	}

	startOffset := int(offset * WordSize) // FIXME: check for overflows in 32bit archs
	endOffset := startOffset + int(size)
	if endOffset > len(ms.b) {
		return ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: endOffset}
	}

	return nil
}

// uncheckedSlice returns a slice without checking for bounds. Bounds MUST be
// checked first by calling checkSliceBounds, otherwise this may panic.
//
// These functions are split to allow uncheckedSlice to be trivially inlineable.
func (ms *Segment) uncheckedSlice(offset WordOffset, size ByteCount) []byte {
	startOffset := int(offset * WordSize)
	return ms.b[startOffset : startOffset+int(size)]
}

func (ms *Segment) Read(offset WordOffset, b []byte) (int, error) {
	if err := ms.rl.CanRead(1); err != nil {
		return 0, err
	}

	byteOffset := int(offset * WordSize)
	if byteOffset >= len(ms.b) {
		return 0, ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: byteOffset}
	}

	n := copy(b, ms.b[byteOffset:])
	return n, nil
}

func (ms *Segment) CheckBounds(offset WordOffset, size WordCount) error {
	byteOffset := int(offset * WordSize)
	byteSize := int(size * WordSize)
	if byteOffset < 0 || byteOffset+byteSize > len(ms.b) {
		return ErrObjectOutOfBounds{Offset: offset, Size: size, Len: len(ms.b)}
	}
	return nil
}

type SingleSegmentArena struct {
	s Segment
}

func (arena *SingleSegmentArena) Segment(id SegmentID) (*Segment, error) {
	if id != 0 {
		return nil, ErrUnknownSegment(id)
	}
	if arena == nil {
		return nil, errSegmentNotInitialized
	}

	return &arena.s, nil
}

func (arena *SingleSegmentArena) Reset(b []byte, writable bool) {
	arena.s.b = b
	arena.s.rl.Reset()
}

func MakeSingleSegmentArena(b []byte, writable bool, rl *ReadLimiter) SingleSegmentArena {
	var arena SingleSegmentArena
	arena.s.rl = rl
	arena.Reset(b, writable)
	return arena
}
