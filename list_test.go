// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// TestListWordCount tests the various cases of the listWordCount function.
func TestListWordCount(t *testing.T) {
	maxsz := listSize(MaxListSize)

	tests := []struct {
		el listElementSize
		sz listSize
		w  WordCount
	}{
		{el: listElSizeVoid, sz: 0, w: 0},
		{el: listElSizeVoid, sz: 1, w: 0},
		{el: listElSizeVoid, sz: maxsz, w: 0},
		{el: listElSizeBit, sz: 0, w: 0},
		{el: listElSizeBit, sz: 1, w: 1},
		{el: listElSizeBit, sz: 64, w: 1},
		{el: listElSizeBit, sz: 65, w: 2},
		{el: listElSizeBit, sz: maxsz, w: WordCount(maxsz)/64 + 1},
		{el: listElSizeByte, sz: 0, w: 0},
		{el: listElSizeByte, sz: 1, w: 1},
		{el: listElSizeByte, sz: 8, w: 1},
		{el: listElSizeByte, sz: 9, w: 2},
		{el: listElSizeByte, sz: maxsz, w: WordCount(maxsz)/8 + 1},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("%d/%d", tc.el, tc.sz)
		t.Run(name, func(t *testing.T) {
			got := listWordCount(tc.el, tc.sz)
			require.EqualValues(t, tc.w, got)
		})
	}
}

func BenchmarkListGetUnsafeString(b *testing.B) {
	name := "mynameisslimsha"
	buf := appendWords(nil, 0x0000008200000001)
	buf = append(buf, []byte(name)...)
	buf = append(buf, 0) // Null mark for text

	benchmarkRLMatrix(b, func(b *testing.B, newRL newRLFunc) {
		arena := NewSingleSegmentArena(buf, false, newRL(MaxReadLimiterLimit))
		seg, _ := arena.Segment(0)

		// Tests only reading after already having checked for
		// validity.
		b.Run("from list no check", func(b *testing.B) {
			ls := &List{
				seg: seg,
				ptr: listPointer{elSize: listElSizeByte, listSize: listSize(len(name) + 1), startOffset: 1},
				dl:  noDepthLimit,
			}
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
			ls := &List{
				seg: seg,
				ptr: listPointer{elSize: listElSizeByte, listSize: listSize(len(name) + 1), startOffset: 1},
				dl:  noDepthLimit,
			}

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
			st := &Struct{
				seg:   seg,
				ptr:   structPointer{pointerSectionSize: 1},
				arena: arena,
				dl:    noDepthLimit,
			}

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
