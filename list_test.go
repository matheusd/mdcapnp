// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"testing"
)

func BenchmarkListGetUnsafeString(b *testing.B) {
	name := "mynameisslimshad"
	buf := appendWords(nil, 0x0000008200000001)
	buf = append(buf, []byte(name)...)

	for _, withRL := range []bool{false, true} {
		var rl *ReadLimiter
		if withRL {
			rl = NewReadLimiter(maxReadOnReadLimiter)
		}
		seg := &Segment{b: buf, rl: rl}

		b.Run(fmt.Sprintf("rl=%v", withRL), func(b *testing.B) {
			// Tests only reading after already having checked for
			// validity.
			b.Run("from list no check", func(b *testing.B) {
				ls := &List{seg: seg, ptr: listPointer{elSize: listElSizeByte, listSize: listSize(len(name)), startOffset: 1}}
				if err := ls.CheckCanGetUnsafeString(); err != nil {
					b.Fatal(err)
				}

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					got := ls.UnsafeString()
					if got != name {
						b.Fatalf("Unexpected name: got %q, want %q", got, name)
					}
				}
			})

			// Tests both checking for validity and reading.
			b.Run("from list", func(b *testing.B) {
				ls := &List{seg: seg, ptr: listPointer{elSize: listElSizeByte, listSize: listSize(len(name)), startOffset: 1}}

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					if err := ls.CheckCanGetUnsafeString(); err != nil {
						b.Fatal(err)
					}
					got := ls.UnsafeString()
					if got != name {
						b.Fatalf("Unexpected name: got %q, want %q", got, name)
					}
				}
			})

			// Tests reading from struct.
			b.Run("from struct", func(b *testing.B) {
				st := &Struct{seg: seg, ptr: structPointer{pointerSectionSize: 1}}

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					var ls List
					if err := st.ReadList(0, &ls); err != nil {
						b.Fatal(err)
					}
					if err := ls.CheckCanGetUnsafeString(); err != nil {
						b.Fatal(err)
					}
					got := ls.UnsafeString()
					if got != name {
						b.Fatalf("Unexpected name: got %q, want %q", got, name)
					}
				}
			})

		})
	}

}
