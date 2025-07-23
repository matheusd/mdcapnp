// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
)

type SegmentID uint64

type Arena interface {
	Segment(id SegmentID) (*Segment, error)
	ReadLimiter() *ReadLimiter
}

type Segment struct {
	b []byte
}

// uncheckedGetWord returns the word at the given offset without checking for
// valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64(ms.b[offset*WordSize:]))
}

func (ms *Segment) GetWord(offset WordOffset) (res Word, err error) {
	if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.LittleEndian.Uint64(ms.b[byteOffset:]))
	}
	return
}

// uncheckedGetWord returns the word at the given offset as a pointer without
// checking for valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWordAsPointer(offset WordOffset) pointer {
	return pointer(binary.LittleEndian.Uint64(ms.b[offset*WordSize:]))
}

func (ms *Segment) getWordAsPointer(offset WordOffset) (pointer, error) {
	w, err := ms.GetWord(offset)
	return pointer(w), err
}

// checkSliceBounds checks whether a subsequent call to [uncheckedSlice] with
// the same arguments will fail. If this function returns true, immediately
// calling [uncheckedSlice] will generate a valid slice.
func (ms *Segment) checkSliceBounds(offset WordOffset, size ByteCount) error {
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
	byteOffset := int(offset * WordSize)
	if byteOffset >= len(ms.b) {
		return 0, ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: byteOffset}
	}

	n := copy(b, ms.b[byteOffset:])
	return n, nil
}

func (ms *Segment) CheckBounds(offset WordOffset, size WordCount) error {
	byteOffset := int(offset) * WordSize // TODO: check if 32 bits arch?
	byteSize := int(size) * WordSize
	if byteOffset < 0 || byteOffset+byteSize > len(ms.b) {
		return ErrObjectOutOfBounds{Offset: offset, Size: size, Len: len(ms.b)}
	}
	return nil
}

type SingleSegmentArena struct {
	s  Segment
	rl *ReadLimiter
}

func (arena *SingleSegmentArena) ReadLimiter() *ReadLimiter {
	return arena.rl
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
	arena.rl.Reset()
}

func NewSingleSegmentArena(b []byte, writable bool, rl *ReadLimiter) *SingleSegmentArena {
	var arena SingleSegmentArena
	arena.rl = rl
	arena.Reset(b, writable)
	return &arena
}
