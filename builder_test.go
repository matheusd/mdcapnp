// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"
	"unsafe"

	"matheusd.com/depvendoredtestify/require"
)

// TestAllocStateHeaderBufSeg0Prefix verifies that the headerBufPrefixesSeg0Buf
// which verifies whether the HeaderBuf slice prefixes seg0's slice is correct.
func TestAllocStateHeaderBufSeg0Prefix(t *testing.T) {
	// Allocate two slices which are contiguous in memory in order to ensure
	// the tests below assert the false cases between the two slices
	// correctly.
	var buf0, buf1 []byte
	const testRange = 1000000
	for i := range testRange {
		buf0 = make([]byte, 16, 32)
		buf1 = make([]byte, 16)
		contiguous := unsafe.Pointer(unsafe.SliceData(buf1)) ==
			unsafe.Add(unsafe.Pointer(unsafe.SliceData(buf0)), cap(buf0))
		if contiguous {
			break
		}
		if i == testRange-1 {
			panic("could not make contiguous slices for test")
		}
	}

	tests := []struct {
		name string
		hb   []byte // Header buf
		sb   []byte // seg0 buf
		want bool
	}{{
		name: "non-contiguous bufffers",
		hb:   buf1,
		sb:   buf0,
		want: false,
	}, {
		name: "different contiguous bufffers wrong cap", // Fails because the cap is wrong
		hb:   buf0[31:32],
		sb:   buf1,
		want: false,
	}, {
		name: "different contiguous bufffers wrong ptr", // Fails because the ptr is wrong
		hb:   buf0[16:],
		sb:   buf1,
		want: false,
	}, {
		name: "aliased buffers",
		hb:   buf0[:8],
		sb:   buf0[8:],
		want: true,
	}, {
		name: "aliased without cap on hb",
		hb:   buf0[:8:8],
		sb:   buf0[8:],
		want: false,
	}, {
		name: "aliased without cap on sb",
		hb:   buf0[:8],
		sb:   buf0[8:10:10],
		want: true,
	}, {
		name: "aliased with gap",
		hb:   buf0[:4],
		sb:   buf0[8:10],
		want: false,
	}, {
		name: "aliased in the middle",
		hb:   buf0[2:8],
		sb:   buf0[8:10],
		want: true,
	}, {
		name: "aliases overlap", // Technically invalid to have as a state.
		hb:   buf0[:8],
		sb:   buf0[4:12],
		want: false,
	}, {
		name: "aliases overlap zero len hb", // Technically invalid to have as a state.
		hb:   buf0[2:2],
		sb:   buf0[2:12],
		want: true,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			as := AllocState{
				HeaderBuf: tc.hb,
				FirstSeg:  tc.sb,
			}
			got := as.headerBufPrefixesSeg0Buf()
			require.Equal(t, tc.want, got)
		})
	}
}

// TestSegmentBuilderPreservesBufAfterRealloc tests that re-allocating the
// underlying buffer maintains the correct references within SegmentBuilder
// objects.
func TestSegmentBuilderPreservesBufAfterRealloc(t *testing.T) {
	alloc := &alwaysReallocAllocator{segsCapacity: 2}
	mb, err := NewMessageBuilder(alloc)
	require.NoError(t, err)

	// Allocate first word.
	v1 := Word(0x1122333455667788)
	seg1, off1, err := mb.allocate(0, 1)
	require.NoError(t, err)
	require.EqualValues(t, 0, seg1.id) // Sanity check.
	require.EqualValues(t, 1, off1)    // Sanity check.
	seg1.uncheckedSetWord(off1, v1)

	// Allocate a large amount to cause a re-allocation in the internal
	// buffer.
	const extraWords = 256
	_, _, err = mb.allocate(0, extraWords)
	require.NoError(t, err)

	// Allocate second word.
	v2 := Word(0x99aabbccddeeff00)
	seg2, off2, err := mb.allocate(0, 1)
	require.NoError(t, err)
	seg2.uncheckedSetWord(off2, v2)

	// Ensure the second offset was made far from the first offset.
	require.EqualValues(t, seg1.id, seg2.id)
	require.Equal(t, off1+extraWords+1, off2)

	// Ensure both segment builders can read each other's values.
	require.Equal(t, v1, seg2.uncheckedGetWord(off1))
	require.Equal(t, v2, seg1.uncheckedGetWord(off2))

	// Create a new segment and allocate the third word on it.
	alloc.createNewSeg = true
	v3 := Word(0x5566778899001122)
	seg3, off3, err := mb.allocate(0, 4)
	require.NoError(t, err)
	require.EqualValues(t, 1, seg3.id) // Sanity check.
	seg3.uncheckedSetWord(off3, v3)
	require.NotEqualValues(t, seg3.id, seg1.id)

	// Allocate a new word back in the first segment.
	alloc.usePreferredSeg = true
	v4 := Word(0x7788990011223344)
	seg4, off4, err := mb.allocate(0, 1)
	require.NoError(t, err)
	require.EqualValues(t, 0, seg4.id) // Sanity check.
	seg4.uncheckedSetWord(off4, v4)

	// Ensure segment builders on segment 0 can read each other's values.
	//
	// This is the critical test: ensures that even after allocating new
	// segments and new segment buffers, old segment builders (created
	// before the reallocs) can still read new values (i.e. they are all
	// pointing to the _same_ buffer).
	require.Equal(t, v1, seg1.uncheckedGetWord(off1))
	require.Equal(t, v1, seg2.uncheckedGetWord(off1))
	require.Equal(t, v1, seg4.uncheckedGetWord(off1))
	require.Equal(t, v2, seg1.uncheckedGetWord(off2))
	require.Equal(t, v2, seg2.uncheckedGetWord(off2))
	require.Equal(t, v2, seg4.uncheckedGetWord(off2))
	require.Equal(t, v4, seg1.uncheckedGetWord(off4))
	require.Equal(t, v4, seg2.uncheckedGetWord(off4))
	require.Equal(t, v4, seg4.uncheckedGetWord(off4))
}

// BenchmarkBuilderSetInt64 benchmarks the SetInt64 function.
func BenchmarkBuilderSetInt64(b *testing.B) {
	alloc := NewSimpleSingleAllocator(10, false)

	b.Run("reuse all", func(b *testing.B) {
		mb, err := NewMessageBuilder(alloc)
		require.NoError(b, err)

		st, err := NewGoserbenchSmallStruct(mb)
		require.NoError(b, err)

		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			st.SetBirthDay(int64(i))
		}

		ser, err := mb.Serialize()
		require.NoError(b, err)
		// b.Logf("%x", ser)
		_ = ser
	})

	b.Run("reuse mb", func(b *testing.B) {
		mb, err := NewMessageBuilder(alloc)
		require.NoError(b, err)

		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			st, err := NewGoserbenchSmallStruct(mb)
			if err != nil {
				b.Fatal(err)
			}
			st.SetBirthDay(int64(i))
			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}

		ser, err := mb.Serialize()
		require.NoError(b, err)
		// b.Logf("%x", ser)
		_ = ser
	})

	b.Run("reuse none", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			mb, err := NewMessageBuilder(alloc)
			if err != nil {
				b.Fatal(err)
			}

			st, err := NewGoserbenchSmallStruct(mb)
			if err != nil {
				b.Fatal(err)
			}
			st.SetBirthDay(int64(i))
			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMsgBuilderAllocate benchmarks the overhead in
// MessageBuilder.allocate using a simulated test allocator that doesn't do
// anything.
func BenchmarkMsgBuilderAllocate(b *testing.B) {
	mb, err := NewMessageBuilder(&nopAllocator{})
	require.NoError(b, err)

	var sb SegmentBuilder
	var off WordOffset

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		sb, off, err = mb.allocate(0, 0)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Ensure off and sb are not eliminated by the compiler.
	if off == 666 {
		b.Logf("%v", sb)
	}
}
