// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// appendWords appends little endian encoded words into a slice.
func appendWords(b []byte, words ...Word) []byte {
	for _, w := range words {
		b = binary.LittleEndian.AppendUint64(b, uint64(w))
	}
	return b
}

// newSegmentData creates a raw segment data with little-endian encoded words.
func newSegmentData(words ...Word) []byte {
	return appendWords([]byte(nil), words...)
}

// TestValidWordCounts tests validity of various word counts.
func TestValidWordCounts(t *testing.T) {
	tests := []struct {
		wc WordCount
		v  bool
	}{
		{wc: 0, v: true},
		{wc: 1, v: true},
		{wc: MaxValidWordCount - 1, v: true},
		{wc: MaxValidWordCount, v: true},
		{wc: MaxValidWordCount + 1, v: false},
		{wc: 0xffffffff, v: false}, // -1
		{wc: 0x80000000, v: false}, // min word offset
		{wc: 0x40000000, v: false}, // other invalid bit set
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%08x", tc.wc), func(t *testing.T) {
			got := tc.wc.Valid()
			require.EqualValues(t, tc.v, got)
		})
	}
}

// BenchmarkWordCountValid benchmarks the Valid() call of WordCount.
func BenchmarkWordCountValid(b *testing.B) {
	var wc WordCount
	var N = WordCount(b.N)
	var res bool

	b.ResetTimer()
	for wc < N {
		wc += 1
		res = wc.Valid()
	}

	require.EqualValues(b, N.Valid(), res)
}

// TestAdd3WordOffsets tests that the add3WordOffsets works as expected.
func TestAdd3WordOffsets(t *testing.T) {
	maxInt32 := WordOffset(math.MaxInt32)
	minInt32 := WordOffset(math.MinInt32)

	tests := []struct {
		a, b, c WordOffset
		r       WordOffset
		ok      bool
	}{
		{a: 0, b: 0, c: 0, r: 0, ok: true},
		{a: 1, b: 2, c: 3, r: 6, ok: true},
		{a: -1, b: -2, c: -3, r: -6, ok: true},
		{a: maxWordOffset, b: 0, c: 0, r: maxWordOffset, ok: true},
		{a: maxWordOffset, b: +1, c: -1, r: maxWordOffset, ok: true},
		{a: maxWordOffset, b: +1, c: 0, r: maxWordOffset + 1, ok: false},
		{a: maxWordOffset, b: maxWordOffset, c: maxWordOffset, r: maxWordOffset * 3, ok: false},
		{a: maxWordOffset, b: 0x60000000, c: 0, r: maxInt32, ok: false},
		{a: maxWordOffset, b: 0x60000000, c: 1, r: minInt32, ok: false}, // Overflow.
		{a: maxWordOffset, b: -maxWordOffset, c: maxWordOffset, r: maxWordOffset, ok: true},
		{a: maxInt32, b: 0, c: 0, r: maxInt32, ok: false},
		{a: maxInt32, b: 1, c: 0, r: minInt32, ok: false}, // Overflow.
		{a: maxInt32, b: maxInt32, c: maxInt32, r: maxInt32 * 3, ok: false},
		{a: maxInt32, b: -maxInt32, c: 0, r: 0, ok: true},
		{a: minWordOffset, b: 0, c: 0, r: minWordOffset, ok: true},
		{a: minWordOffset, b: -1, c: 0, r: minWordOffset - 1, ok: false},
		{a: minWordOffset, b: -minWordOffset, c: 0, r: 0, ok: true},
		{a: minWordOffset, b: -minWordOffset, c: minWordOffset, r: minWordOffset, ok: true},
		{a: minInt32, b: 0, c: 0, r: minInt32, ok: false},
		{a: minInt32, b: maxInt32, c: 0, r: -1, ok: true},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("%08x+%08x+%08x", tc.a, tc.b, tc.c)
		t.Run(name, func(t *testing.T) {
			got, gotOk := add3WordOffsets(tc.a, tc.b, tc.c)
			require.EqualValues(t, tc.r, got)
			require.EqualValues(t, tc.ok, gotOk)
		})
	}
}

// BenchmarkAdd3WordCounts benchmarks the add3WordCounts function.
func BenchmarkAdd3WordCounts(b *testing.B) {
	var r WordOffset
	var ok bool
	var seed = rand.Uint64()

	for i := range b.N {
		var a, b, c WordOffset
		u := uint64(i)

		// Make the values unpredictable (to prevent compiler using
		// constants).
		a, seed = WordOffset(u^seed), bits.RotateLeft64(seed, i%64)
		b, seed = WordOffset(u^seed), bits.RotateLeft64(seed, i%64)
		c, seed = WordOffset(u^seed), bits.RotateLeft64(seed, i%64)

		r, ok = add3WordOffsets(a, b, c)
	}

	if r == 666 {
		b.Logf("Result: r=%016x ok=%v", r, ok)
	}
}

// TestAddWordOffsetsWithCarry tests that the addWordOffsetsWithCarry works as
// expected.
func TestAddWordOffsetsWithCarry(t *testing.T) {
	maxInt32 := WordOffset(math.MaxInt32)
	minInt32 := WordOffset(math.MinInt32)

	tests := []struct {
		a, b WordOffset
		c    uint64
		r    WordOffset
		ok   bool
	}{
		{a: 0, b: 0, c: 0, r: 0, ok: true},
		{a: 0, b: 0, c: 1, r: 1, ok: true},
		{a: 2, b: 3, c: 1, r: 6, ok: true},
		{a: maxWordOffset, b: 0, c: 0, r: maxWordOffset, ok: true},
		{a: maxWordOffset, b: 2, c: 1, r: maxWordOffset + 3, ok: false},
		{a: maxWordOffset - 1, b: 0, c: 1, r: maxWordOffset, ok: true},
		{a: maxWordOffset, b: -maxWordOffset, c: 1, r: 1, ok: true},
		{a: maxInt32, b: 0, c: 0, r: maxInt32, ok: false},
		{a: maxInt32, b: 0, c: 1, r: minInt32, ok: false}, // Overflow.
		{a: maxInt32, b: maxInt32, c: 1, r: maxInt32*2 + 1, ok: false},
		{a: maxInt32, b: -maxInt32, c: 0, r: 0, ok: true},
		{a: minWordOffset, b: 0, c: 0, r: minWordOffset, ok: true},
		{a: minWordOffset, b: -1, c: 0, r: minWordOffset - 1, ok: false},
		{a: minWordOffset, b: -1, c: 1, r: minWordOffset, ok: true},
		{a: minWordOffset, b: -minWordOffset, c: 0, r: 0, ok: true},
		{a: minWordOffset, b: -minWordOffset, c: 1, r: 1, ok: true},
		{a: minInt32, b: 0, c: 0, r: minInt32, ok: false},
		{a: minInt32, b: maxInt32, c: 0, r: -1, ok: true},
		{a: minInt32, b: maxInt32, c: 1, r: 0, ok: true},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("%08x+%08x+%d", tc.a, tc.b, tc.c)
		t.Run(name, func(t *testing.T) {
			got, gotOk := addWordOffsetsWithCarry(tc.a, tc.b, tc.c)
			require.EqualValues(t, tc.r, got)
			require.EqualValues(t, tc.ok, gotOk)
		})
	}
}

// BenchmarkAddWordOffsetsWithCarry benchmarks the addWordOffsetsWithCarry
// function.
func BenchmarkAddWordOffsetsWithCarry(b *testing.B) {
	var r WordOffset
	var ok bool
	var seed = rand.Uint64()

	for i := range b.N {
		var a, b WordOffset
		var c uint64
		u := uint64(i)

		// Make the values unpredictable (to prevent compiler using
		// constants).
		a, seed = WordOffset(u^seed), bits.RotateLeft64(seed, i%64)
		b, seed = WordOffset(u^seed), bits.RotateLeft64(seed, i%64)
		c, seed = seed&1, bits.RotateLeft64(seed, i%64)

		r, ok = addWordOffsetsWithCarry(a, b, c)
	}

	if r == 666 {
		b.Logf("Result: r=%016x ok=%v", r, ok)
	}
}
