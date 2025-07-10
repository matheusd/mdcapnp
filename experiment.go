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

type WordOffset uint64

type SignedWordOffset int64

type ByteOffset uint64

type SegmentID uint64

type Arena interface {
	Segment(id SegmentID) (Segment, error)
}

type Segment interface {
	GetWord(offset WordOffset) (res Word, err error)
	Read(offset WordOffset, b []byte) (int, error)
	CheckBounds(offset WordOffset, size WordCount) error
}

type MemSegment struct {
	b  []byte
	rl *ReadLimiter
}

func (ms *MemSegment) GetWord(offset WordOffset) (res Word, err error) {
	if err = ms.rl.CanRead(1); err != nil {
	} else if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.BigEndian.Uint64(ms.b[byteOffset:]))

		// copy((*[8]byte)(unsafe.Pointer(&res))[:], ms.b[byteOffset:])

		// Assumes a big endian version is written. Note: this is
		// counterintuitive, double check.
		// res = Word(binary.BigEndian.Uint64((*[8]byte)(unsafe.Pointer(&res))[:]))
	}
	return
}

func (ms *MemSegment) Read(offset WordOffset, b []byte) (int, error) {
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

func (ms *MemSegment) CheckBounds(offset WordOffset, size WordCount) error {
	byteOffset := int(offset * WordSize)
	byteSize := int(size * WordSize)
	if byteOffset < 0 || byteOffset+byteSize >= len(ms.b) {
		return ErrObjectOutOfBounds{Offset: offset, Size: size, Len: len(ms.b)}
	}
	return nil
}

type SingleSegmentMemArena struct {
	s MemSegment
}

func (arena *SingleSegmentMemArena) Segment(id SegmentID) (Segment, error) {
	if id != 0 {
		return nil, ErrUnknownSegment(id)
	}
	if arena == nil {
		return nil, errSegmentNotInitialized
	}

	return &arena.s, nil
}

func (arena *SingleSegmentMemArena) Reset(b []byte, writable bool) {
	arena.s.b = b
}

func MakeSingleSegmentMemArena(b []byte, writable bool) SingleSegmentMemArena {
	var arena SingleSegmentMemArena
	arena.Reset(b, writable)
	return arena
}
