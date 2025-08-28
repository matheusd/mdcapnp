// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// TestNativeSegmentWordCalls tests that the functions that use native versions
// of calls work as expected.
func TestNativeSegmentWordCalls(t *testing.T) {
	seg := &Segment{
		b: []byte{
			// Add one word to fetch from a non-zero offset.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

			// Little endian encoded version of 0x1234567890abcdef
			0xef, 0xcd, 0xab, 0x90, 0x78, 0x56, 0x34, 0x12,
		},
	}

	const offset = 1
	target := uint64(0x1234567890abcdef)

	require.Equal(t, Word(target), seg.uncheckedGetWord(offset))
	gotW, err := seg.GetWord(offset)
	require.NoError(t, err)
	require.Equal(t, Word(target), gotW)
	require.Equal(t, pointer(target), seg.uncheckedGetWordAsPointer(offset))
}
