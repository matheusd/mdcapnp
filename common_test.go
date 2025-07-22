// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// TestValidWordOffsets tests the range of valid word offsets.
func TestValidWordOffsets(t *testing.T) {
	var minw uint32 = 0xe0000000
	var maxw uint32 = 0x1fffffff

	tests := []struct {
		w    WordOffset
		want bool
	}{
		{w: 0, want: true},
		{w: -1, want: true},
		{w: WordOffset(minw), want: true},
		{w: WordOffset(maxw), want: true},
		{w: WordOffset(maxw + 1), want: false},
		{w: WordOffset(minw - 1), want: false},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("%08x", uint32(tc.w))
		t.Run(name, func(t *testing.T) {
			got := tc.w.Valid()
			require.Equal(t, tc.want, got)
		})
	}
}

// TestAddWordOffsets verifies whether adding two word offsets detects the
// correct edge cases.
func TestAddWordOffsets(t *testing.T) {
	tests := []struct {
		a      WordOffset
		b      WordOffset
		want   WordOffset
		wantOk bool
	}{
		{a: 0, b: 0, want: 0, wantOk: true},
		{a: 0, b: 1, want: 1, wantOk: true},
		{a: 0, b: -1, want: -1, wantOk: true},
		{a: minWordOffset, b: 0, want: minWordOffset, wantOk: true},                  // min + 0
		{a: maxWordOffset, b: 0, want: maxWordOffset, wantOk: true},                  // max + 0
		{a: minWordOffset, b: +1, want: minWordOffset + 1, wantOk: true},             // min + 1
		{a: maxWordOffset, b: -1, want: maxWordOffset - 1, wantOk: true},             // max - 1
		{a: minWordOffset, b: -1, want: minWordOffset - 1, wantOk: false},            // min - 1
		{a: maxWordOffset, b: +1, want: maxWordOffset + 1, wantOk: false},            // max + 1
		{a: minWordOffset, b: minWordOffset, want: 2 * minWordOffset, wantOk: false}, // min + min = 2*min
		{a: maxWordOffset, b: maxWordOffset, want: 2 * maxWordOffset, wantOk: false}, // max + max = 2*max
		{a: maxWordOffset, b: minWordOffset, want: -1, wantOk: true},                 // max - min = 1
	}

	for _, tc := range tests {
		name := fmt.Sprintf("%08x+%08x", uint32(tc.a), uint32(tc.b))
		t.Run(name, func(t *testing.T) {
			got, gotOk := addWordOffsets(tc.a, tc.b)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantOk, gotOk)
		})
	}
}

// BenchmarkAddWordOffsets benchmarks the addWordOffsets function.
func BenchmarkAddWordOffsets(b *testing.B) {
	var wa, wb WordOffset
	rng := rand.New(rand.NewPCG(0, 0x1701d))
	wa = WordOffset(rng.Int32())
	wb = WordOffset(rng.Int32())

	var wc WordOffset
	var ok bool

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		wc, ok = addWordOffsets(wa, wb)
		wa = wc
	}

	require.NotEqual(b, WordOffset(666), wc)
	b.Logf("Ok: %v", ok) // Ensure ok value is tracked.
}
