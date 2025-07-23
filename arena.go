// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
)

type SegmentID uint32

type Segment struct {
	b []byte
}

// uncheckedOpenSlice returns a slice starting at the provided offset up to the
// end of the buffer without checking bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedTailSlice(offset WordOffset) []byte {
	return ms.b[offset*WordSize:]
}

func (ms *Segment) intLen() int {
	return len(ms.b)
}

// uncheckedGetWord returns the word at the given offset without checking for
// valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64(ms.uncheckedTailSlice(offset)))
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
	return pointer(binary.LittleEndian.Uint64(ms.uncheckedTailSlice(offset)))
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
	if endOffset > ms.intLen() {
		return ErrInvalidMemOffset{AvailableLen: ms.intLen(), Offset: endOffset}
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
	if byteOffset >= ms.intLen() {
		return 0, ErrInvalidMemOffset{AvailableLen: ms.intLen(), Offset: byteOffset}
	}

	n := copy(b, ms.b[byteOffset:])
	return n, nil
}

func (ms *Segment) CheckBounds(offset WordOffset, size WordCount) error {
	byteOffset := int(offset) * WordSize // TODO: check if 32 bits arch?
	byteSize := int(size) * WordSize
	if byteOffset < 0 || byteOffset+byteSize > ms.intLen() {
		return ErrObjectOutOfBounds{Offset: offset, Size: size, Len: ms.intLen()}
	}
	return nil
}

type Arena struct {
	// fb is the full, framed data for the arena (includes header and arena
	// size framing when != nil).
	fb []byte

	// s is the first segment. It is the only segment in single-segment
	// arenas.
	s Segment

	// segs are the additional segments in multi-segment arenas. The segment
	// at index 0 is the segment with id 1, and so on.
	segs *[]*Segment

	rl *ReadLimiter
}

func (arena *Arena) ReadLimiter() *ReadLimiter {
	return arena.rl
}

func (arena *Arena) Segment(id SegmentID) (*Segment, error) {
	if arena == nil {
		return nil, errArenaNotInitialized
	}

	if id == 0 {
		return &arena.s, nil
	}

	index := int(id - 1)
	segs := *arena.segs
	if index >= len(segs) {
		return nil, ErrUnknownSegment(id)
	}

	return segs[index], nil
}

// DecodeSingleSegment decodes the given buffer as a single segment arena.
func (arena *Arena) DecodeSingleSegment(fb []byte) error {
	b, err := decodeSingleSegmentStream(fb)
	if err != nil {
		return err
	}
	arena.Reset(b, false)
	arena.fb = fb
	return nil
}

func (arena *Arena) Reset(b []byte, writable bool) {
	arena.s.b = b
	arena.fb = nil
	arena.rl.Reset()
}

func NewSingleSegmentArena(b []byte, writable bool, rl *ReadLimiter) *Arena {
	var arena Arena
	arena.rl = rl
	arena.Reset(b, writable)
	return &arena
}
