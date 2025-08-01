// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "math/bits"

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
	maxWordOffset                = (1 << 29) - 1
	minWordOffset                = -(1 << 29)
	minWordOffsetAsUint64 uint64 = 0xffffffffe0000000 // This MUST be equal to minWordOffset.
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
	sum64 := uint64(a) + uint64(b)
	c = WordOffset(sum64)
	// ok = sum64 >= minWordOffset && sum64 <= maxWordOffset
	ok = (sum64 - minWordOffsetAsUint64) <= (maxWordOffset - minWordOffset)
	return
}

// add3WordOffsets adds 3 word offsets, determining whether the resulting offset
// is valid.
func add3WordOffsets(a, b, c WordOffset) (d WordOffset, ok bool) {
	sum64 := uint64(a) + uint64(b) + uint64(c)
	d = WordOffset(sum64)
	// ok = sum64 >= minWordOffset && sum64 <= maxWordOffset
	ok = (sum64 - minWordOffsetAsUint64) <= (maxWordOffset - minWordOffset)
	return
}

// addWordOffsetsWithCarry returns a + b + c, with c being a carry (a value of
// either 0 or 1). If c is not 0 or 1, the results are undefined.
//
// It returns ok if the resulting value is a valid word offset.
func addWordOffsetsWithCarry(a, b WordOffset, c uint64) (d WordOffset, ok bool) {
	sum64, _ := bits.Add64(uint64(int64(a)), uint64(int64(b)), c)
	d = WordOffset(sum64)
	// ok = sum64 >= minWordOffset && sum64 <= maxWordOffset
	ok = (sum64 - minWordOffsetAsUint64) <= (maxWordOffset - minWordOffset)
	return
}

// WordCount is a count of addressable words within a segment. Only up to 2^29
// words are addressable within a segment, therefore a count of words can only
// go up to that amount. Counts cannot be negative.
type WordCount uint32

// Valid returns true if this is a valid word count.
func (wc WordCount) Valid() bool {
	return wc <= MaxValidWordCount
}

// ByteCount returns the number of bytes that correspond to this number of
// words.
func (wc WordCount) ByteCount() ByteCount {
	return ByteCount(wc) * WordSize
}

func addWordOffsetAndCount(off WordOffset, c WordCount) (r WordOffset, ok bool) {
	return addWordOffsets(off, WordOffset(c))
}

// MaxValidWordCount is the maximum number of words a segment may have.
const MaxValidWordCount = 1<<29 - 1

// maxValidBytes is the maximum number of bytes a segment may have.
const maxValidBytes = MaxValidWordCount * WordSize

type ByteCount uint64

type StructSize struct {
	DataSectionSize    wordCount16
	PointerSectionSize wordCount16
}

func (ss StructSize) TotalSize() WordCount {
	return WordCount(ss.DataSectionSize) + WordCount(ss.PointerSectionSize)
}

type ListSize struct {
	elSize   listElementSize
	listSize listSize
}

func isWordAligned(i int) bool {
	const alignMask = WordSize - 1 // Only works because WordSize is a power of 2.
	return i&alignMask == 0
}
