// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"testing"
)

func BenchmarkStructGetInt64(b *testing.B) {
	buf := binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)

	b.Run("with RL", func(b *testing.B) {
		st := &SmallTestStruct{seg: &MemSegment{b: buf, rl: NewReadLimiter(maxReadOnReadLimiter)}}
		b.ResetTimer()
		b.ReportAllocs()
		var v int64
		for range b.N {
			v = st.Siblings()
		}

		if v == 666 {
			panic("boo")
		}
	})

	b.Run("no RL", func(b *testing.B) {
		st := &SmallTestStruct{seg: &MemSegment{b: buf}}
		b.ResetTimer()
		b.ReportAllocs()
		var v int64
		for range b.N {
			v = st.Siblings()
		}

		if v == 666 {
			panic("boo")
		}
	})
}
