// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
)

type Word uint64

const WordSize = 8

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

// uncheckedWordOffset returns the final offset of a data word of a struct
// inside a segment, without checking for the validity of this operation. This
// is used in cases where it is assumed that the index has already been
// determined to be valid inside the struct (because the entire struct size has
// been bounds checked already).
func (i DataFieldIndex) uncheckedWordOffset(base WordOffset) WordOffset {
	return base + WordOffset(i)
}

// WordOffset is a signed offset into a segment. Segments can have up to 2^29
// words.
//
// Segment offsets are signed for two reasons: they may validly point to objects
// that have been written to the segment before their pointer was written (e.g.
// orphans or relocated objects) or because empty structs point to their own
// pointer as an offset.
type WordOffset int32

const (
	maxWordOffset = (1 << 29) - 1
	minWordOffset = -(1 << 29)
)

// Valid determines if the value of this offset is valid.
func (w WordOffset) Valid() bool {
	return w >= minWordOffset && w <= maxWordOffset
}

// addWordOffsets adds two offsets, detecting whether the resulting offset
// remains valid.
//
// Returns true if the sum was valid.
func addWordOffsets(a, b WordOffset) (c WordOffset, ok bool) {
	sum64 := int64(a) + int64(b)
	ok = sum64 >= minWordOffset && sum64 <= maxWordOffset
	c = WordOffset(sum64)
	return
}

// WordCount is a count of addressable words within a segment. Only up to 2^29
// words are addressable within a segment, therefore a count of words can only
// go up to that amount. Counts cannot be negative.
type WordCount uint32

// Valid returns true if this is a valid word count.
func (wc WordCount) Valid() bool {
	// Valid counts cannot be negative and cannot have more than 29 bits
	// set. Thus to test for validity, check if any of the highest bits in
	// the most significant nibble are set.
	const invalidBitsMask = 0b1111 << 28
	return wc&invalidBitsMask == 0
}

func addWordOffsetAndCount(off WordOffset, c WordCount) (r WordOffset, ok bool) {
	return addWordOffsets(off, WordOffset(c))
}

const MaxValidWordCount = 1<<30 - 1

type SignedWordOffset int64

type listOrStructOffset int32 // Up to 30 bits usable.

type ByteOffset uint64

type ByteCount uint64

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
