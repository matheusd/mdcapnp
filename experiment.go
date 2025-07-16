// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "encoding/binary"

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

// WordOffset is a signed offset into a segment. Segments can have up to 2^29
// words.
//
// Segment offsets are signed for two reasons: they may validly point to objects
// that have been written to the segment before their pointer was written (e.g.
// orphans or relocated objects) or because empty structs point to their own
// pointer as an offset.
type WordOffset int32

// Valid determines if the value within this offset is valid.
func (w WordOffset) Valid() bool {
	// Valid word offsets have up to 29 bits set, optionally with the sign
	// bit set. This means the invalid bits are the first three bits of the
	// most significant nibble of the value. We test directly whether any of
	// these bits are set here, to determine if the value is valid.
	const invalidBitsMask = 0b0111 << 28
	return w&invalidBitsMask == 0
}

// AddWordOffsets adds two offsets, setting the resulting argument to the sum if
// the sum generates a still valid offset.
//
// Returns true if the sum was valid.
func AddWordOffsets(a, b WordOffset, r *WordOffset) (ok bool) {
	// Could this use bits.Add64??
	c := a + b
	ok = ((c > a) == (b > 0)) && c.Valid()
	if ok {
		*r = c
	}
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

func AddWordOffsetAndCount(off WordOffset, c WordCount, r *WordOffset) (ok bool) {
	return AddWordOffsets(off, WordOffset(c), r)
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

func (ms *Segment) GetWord(offset WordOffset) (res Word, err error) {
	if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.LittleEndian.Uint64(ms.b[byteOffset:]))
	}
	return
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
