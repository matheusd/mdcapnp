// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"fmt"
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
