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

func BenchmarkStructGetInt64(b *testing.B) {
	buf := binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)

	b.Run("with RL", func(b *testing.B) {
		st := &SmallTestStruct{seg: &Segment{b: buf, rl: NewReadLimiter(maxReadOnReadLimiter)}}
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
		st := &SmallTestStruct{seg: &Segment{b: buf}}
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

func BenchmarkStructReadList(b *testing.B) {
	buf := appendWords(nil, 0x00000000fffffffd)

	for _, withRL := range []bool{false, true} {
		var rl *ReadLimiter
		if withRL {
			rl = NewReadLimiter(maxReadOnReadLimiter)
		}
		b.Run(fmt.Sprintf("rl=%v", withRL), func(b *testing.B) {
			seg := &Segment{b: buf, rl: rl}

			b.Run("single struct", func(b *testing.B) {
				st := &Struct{seg: seg, dataStartOffset: 0, dataSize: 0, pointerSize: 1}
				var ls List

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					if err := st.ReadList(0, &ls); err != nil {
						b.Fatal(err)
					}
				}

				require.Equal(b, WordOffset(0), ls.baseOffset)
			})

			// This test verifies if struct escapes to the heap when
			// reading a list from it.
			b.Run("struct per iter", func(b *testing.B) {
				var ls List

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					st := Struct{seg: seg, dataStartOffset: 0, dataSize: 0, pointerSize: 1}
					if err := st.ReadList(0, &ls); err != nil {
						b.Fatal(err)
					}
				}

				require.Equal(b, WordOffset(0), ls.baseOffset)
			})

			// This test verifies if list escapes to the heap when
			// reading it from a struct.
			b.Run("list and struct per iter", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					var ls List
					st := Struct{seg: seg, dataStartOffset: 0, dataSize: 0, pointerSize: 1}
					if err := st.ReadList(0, &ls); err != nil {
						b.Fatal(err)
					}
					if ls.baseOffset != 0 {
						b.Fatal("error")
					}
				}
			})

		})
	}
}
