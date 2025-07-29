// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/internal/testdata"
)

// BenchmarkStructGetInt64 benchmarks indirectly calling the GetInt64 of a
// Struct.
func BenchmarkStructGetInt64(b *testing.B) {
	buf := appendWords(nil, 0x0000000000000100, 0x1234567890abcdef)

	// No need to test multiple ReadLimiter values because GetInt64 does not
	// depend on them.
	arena := NewSingleSegmentArena(buf)
	arena.ReadLimiter().InitNoLimit()
	seg, _ := arena.Segment(0)
	st := &SmallTestStruct{
		seg:   seg,
		arena: arena,
		ptr:   structPointer{dataOffset: 1, dataSectionSize: 1},
		dl:    noDepthLimit,
	}

	b.ResetTimer()
	b.ReportAllocs()

	var v int64
	for range b.N {
		v = st.Siblings()
	}
	require.Equal(b, int64(0x1234567890abcdef), v)
}

// BenchmarkStructReadList benchmarks calling the ReadList call of a struct.
func BenchmarkStructReadList(b *testing.B) {
	buf := appendWords(nil, 0x00000000fffffffd)

	benchmarkRLMatrix(b, func(b *testing.B, rlt readLimiterType) {
		arena := NewSingleSegmentArena(buf)
		rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)
		seg, _ := arena.Segment(0)

		b.Run("single struct", func(b *testing.B) {
			st := &Struct{
				seg:   seg,
				arena: arena,
				ptr:   structPointer{dataOffset: 0, dataSectionSize: 0, pointerSectionSize: 1},
				dl:    noDepthLimit,
			}
			var ls List

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if err := st.ReadList(0, &ls); err != nil {
					b.Fatal(err)
				}
			}

			require.Equal(b, WordOffset(0), ls.ptr.startOffset)
		})

		// This test verifies if struct escapes to the heap when
		// reading a list from it.
		b.Run("struct per iter", func(b *testing.B) {
			var ls List

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				st := Struct{
					arena: arena,
					seg:   seg,
					dl:    noDepthLimit,
					ptr:   structPointer{dataOffset: 0, dataSectionSize: 0, pointerSectionSize: 1},
				}
				if err := st.ReadList(0, &ls); err != nil {
					b.Fatal(err)
				}
			}

			require.Equal(b, WordOffset(0), ls.ptr.startOffset)
		})

		// This test verifies if list escapes to the heap when
		// reading it from a struct.
		b.Run("list and struct per iter", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				var ls List
				st := Struct{
					seg:   seg,
					arena: arena,
					ptr:   structPointer{dataOffset: 0, dataSectionSize: 0, pointerSectionSize: 1},
					dl:    noDepthLimit,
				}
				if err := st.ReadList(0, &ls); err != nil {
					b.Fatal(err)
				}
				if ls.ptr.startOffset != 0 {
					b.Fatal("error")
				}
			}
		})
	})
}

func BenchmarkStructUnsafeString(b *testing.B) {
	benchmarkRLMatrix(b, func(b *testing.B, rlt readLimiterType) {
		segBuf := testdata.GoserbenchSampleA[8:]
		arena := NewSingleSegmentArena(segBuf)
		rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)
		msg := MakeMsg(arena)
		var st GoserbenchSmallStruct

		var s string

		err := st.ReadFromRoot(&msg)
		require.Nil(b, err)

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			s = st.UnsafeName()
		}

		require.Equal(b, "slimshady0123456", s)
	})
}

//go:noinline
func testStructFuncHeap(s *Struct) {
	if s.seg != nil {
		panic("failed")
	}
}

//go:noinline
func testStructFuncStack(s Struct) {
	if s.seg != nil {
		panic("failed")
	}
}

var globalStructHeapVsStackTest *Struct

// BenchmarkStructHeapVsStack benchmarks passing a struct around as either a
// heap or stack allocated reference.
//
// This is used to determine how to build the API.
//
// Note: this only measures the overhead in calling Struct functions using value
// (ensured to be on the stack) vs pointer (which may be on the stack or on the
// heap), _not_ of the heap allocation itself.
func BenchmarkStructHeapVsStack(b *testing.B) {
	st := new(Struct)
	b.Run("heap", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			testStructFuncHeap(st)
		}
	})
	b.Run("stack", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			testStructFuncStack(*st)
		}
	})

	// Ensure st is in the heap.
	globalStructHeapVsStackTest = st
}
