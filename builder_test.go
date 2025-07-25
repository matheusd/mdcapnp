// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"
	"unsafe"

	"matheusd.com/depvendoredtestify/require"
)

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
				Segs:      [][]byte{tc.sb},
			}
			got := as.headerBufPrefixesSeg0Buf()
			require.Equal(t, tc.want, got)
		})
	}
}

func BenchmarkBuilderSetInt64(b *testing.B) {
	b.Run("reuse all", func(b *testing.B) {
		mb, err := NewMessageBuilder(DefaultSimpleSingleAllocator)
		require.NoError(b, err)

		st, err := NewGoserbenchSmallStruct(mb)
		require.NoError(b, err)

		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			st.SetSiblings(int64(i))
		}

		ser, err := mb.Serialize()
		require.NoError(b, err)
		// b.Logf("%x", ser)
		_ = ser
	})

	b.Run("reuse mb", func(b *testing.B) {
		mb, err := NewMessageBuilder(DefaultSimpleSingleAllocator)
		require.NoError(b, err)

		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			st, err := NewGoserbenchSmallStruct(mb)
			if err != nil {
				b.Fatal(err)
			}
			st.SetSiblings(int64(i))
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
			mb, err := NewMessageBuilder(DefaultSimpleSingleAllocator)
			if err != nil {
				b.Fatal(err)
			}

			st, err := NewGoserbenchSmallStruct(mb)
			if err != nil {
				b.Fatal(err)
			}
			st.SetSiblings(int64(i))
			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
