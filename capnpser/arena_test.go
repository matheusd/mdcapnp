// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"math/bits"
	"math/rand/v2"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// BenchmarkSegmentCheckBounds benchmarks the checkBounds() function of Segment.
func BenchmarkSegmentCheckBounds(b *testing.B) {
	const segLen = 16
	var seed = rand.Uint64()

	seg := &Segment{
		b: make([]byte, segLen*WordSize),
	}

	var err error

	for i := range b.N {
		var off WordOffset
		var size WordCount

		// Force a random offset, but clamp size for the check to always
		// pass.
		u := uint64(i)
		off, seed = WordOffset(u^seed%segLen), bits.RotateLeft64(seed, i%64)
		size = segLen - WordCount(off)

		err = seg.checkBounds(off, size)
	}

	require.NoError(b, err)
}
