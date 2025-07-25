// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

type Word uint64

const WordSize = 8

type wordCount16 uint16

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

// ByteCount returns the number of bytes that correspond to this number of
// words.
func (wc WordCount) ByteCount() ByteCount {
	return ByteCount(wc) * WordSize
}

func addWordOffsetAndCount(off WordOffset, c WordCount) (r WordOffset, ok bool) {
	return addWordOffsets(off, WordOffset(c))
}

const MaxValidWordCount = 1<<30 - 1

type ByteCount uint64

type StructSize struct {
	DataSectionSize    wordCount16
	PointerSectionSize wordCount16
}

func (ss StructSize) TotalSize() WordCount {
	return WordCount(ss.DataSectionSize) + WordCount(ss.PointerSectionSize)
}
