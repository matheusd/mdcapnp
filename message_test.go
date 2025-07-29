// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"math"
	"testing"
	"time"

	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/internal/testdata"
)

// BenchmarkMsgReadRoot benchmarks reading the root struct of a message.
func BenchmarkMsgReadRoot(b *testing.B) {
	buf := appendWords(nil, 0x0000000100000000, 0x0000000000000000)
	arena := NewSingleSegmentArena(buf)
	arena.ReadLimiter().InitNoLimit()
	msg := MakeMsg(arena)

	var st Struct

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		err := msg.ReadRoot(&st)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Ensure st is not eliminated.
	require.Equal(b, WordOffset(1), st.ptr.dataOffset)
}

// BenchmarkMsgReadList benchmarks reading a list field from a message.
func BenchmarkMsgReadList(b *testing.B) {
	targetName := []byte("mynameisslimshady\u0000") // Text + null marker.
	buf := appendWords(nil,
		0x0001000000000000,
		// 0x0000000200000001,
		0x0000009200000001,
	)
	buf = append(buf, targetName...)
	buf = append(buf, []byte{5: 0}...) // Pad to word boundary

	benchmarkRLMatrix(b, func(b *testing.B, rlt readLimiterType) {
		arena := NewSingleSegmentArena(buf)
		rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)
		seg, _ := arena.Segment(0)
		st := &SmallTestStruct{
			seg:   seg,
			arena: arena,
			ptr:   structPointer{dataOffset: 1, pointerSectionSize: 1},
			dl:    noDepthLimit,
		}
		ls := new(List)
		nameBuf := make([]byte, 32)

		var n int

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			err := st.ReadNameField(ls)
			if err != nil {
				ls = nil
				b.Fatal(err)
			}

			n, err = ls.Read(nameBuf)
			if err != nil {
				b.Fatal(err)
			}
		}

		require.NotNil(b, ls)
		require.Equal(b, targetName, nameBuf[:n])
	})
}

// BenchmarkDecodeGoserbenchSmallStruct benchmarks decoding a goserbench
// SmallStruct under various configurations.
func BenchmarkDecodeGoserbenchSmallStruct(b *testing.B) {
	var oa goserbenchSmallStruct
	checkOA := func(b *testing.B) {
		require.Equal(b, "slimshady0123456", oa.Name)
		require.Equal(b, int64(0x1011121314151617), oa.BirthDay.Unix())
		require.Equal(b, "phone67890", oa.Phone) // FIXME phone67890
		require.Equal(b, int(0x66669999), oa.Siblings)
		require.Equal(b, true, oa.Spouse)
		require.Equal(b, uint64(0xabcd0000ef01), math.Float64bits(oa.Money))
	}

	serialBuf := testdata.GoserbenchSampleA
	segBuf := testdata.GoserbenchSampleA[8:] // Skip the header.

	tests := []struct {
		rlt    readLimiterType
		unsafe bool
	}{
		{rlt: rlTypeNoLimit, unsafe: true},
		{rlt: rlTypeNoLimit, unsafe: false},
		{rlt: rlTypeUnsafe, unsafe: true},
		{rlt: rlTypeUnsafe, unsafe: false},
		{rlt: rlTypeSafe, unsafe: true},
		{rlt: rlTypeSafe, unsafe: false},
	}

	for _, tc := range tests {
		b.Run(fmt.Sprintf("%v/unsafe=%v", tc.rlt, tc.unsafe), func(b *testing.B) {
			b.Run("reuse all", func(b *testing.B) {
				arena := NewSingleSegmentArena(segBuf)
				tc.rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)
				msg := MakeMsg(arena)
				var st GoserbenchSmallStruct

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					err := st.ReadFromRoot(&msg)
					if err != nil {
						b.Fatal(err)
					}

					if tc.unsafe {
						oa.Name = st.UnsafeName()
						oa.Phone = st.UnsafePhone()
					} else {
						oa.Name = st.Name()
						oa.Phone = st.Phone()
					}
					oa.BirthDay = time.Unix(st.BirthDay(), 0)
					oa.Siblings = int(st.Siblings())
					oa.Spouse = st.Spouse()
					oa.Money = st.Money()
				}

				checkOA(b)
			})

			b.Run("reuse arena", func(b *testing.B) {
				arena := NewSingleSegmentArena(segBuf)
				tc.rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					arena.Reset(segBuf)
					msg := MakeMsg(arena)
					var st GoserbenchSmallStruct

					err := st.ReadFromRoot(&msg)
					if err != nil {
						b.Fatal(err)
					}

					if tc.unsafe {
						oa.Name = st.UnsafeName()
						oa.Phone = st.UnsafePhone()
					} else {
						oa.Name = st.Name()
						oa.Phone = st.Phone()
					}

					oa.BirthDay = time.Unix(st.BirthDay(), 0)
					oa.Siblings = int(st.Siblings())
					oa.Spouse = st.Spouse()
					oa.Money = st.Money()
				}

				checkOA(b)
			})

			b.Run("reuse arena deserialize", func(b *testing.B) {
				arena := NewSingleSegmentArena(segBuf)
				tc.rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					if err := arena.DecodeSingleSegment(serialBuf); err != nil {
						b.Fatal(err)
					}
					msg := MakeMsg(arena)
					var st GoserbenchSmallStruct

					err := st.ReadFromRoot(&msg)
					if err != nil {
						b.Fatal(err)
					}

					if tc.unsafe {
						oa.Name = st.UnsafeName()
						oa.Phone = st.UnsafePhone()
					} else {
						oa.Name = st.Name()
						oa.Phone = st.Phone()
					}

					oa.BirthDay = time.Unix(st.BirthDay(), 0)
					oa.Siblings = int(st.Siblings())
					oa.Spouse = st.Spouse()
					oa.Money = st.Money()
				}

				checkOA(b)
			})

			b.Run("reuse none", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					arena := NewSingleSegmentArena(segBuf)
					tc.rlt.initRL(arena.ReadLimiter(), MaxReadLimiterLimit)
					msg := MakeMsg(arena)
					var st GoserbenchSmallStruct

					err := st.ReadFromRoot(&msg)
					if err != nil {
						b.Fatal(err)
					}

					if tc.unsafe {
						oa.Name = st.UnsafeName()
						oa.Phone = st.UnsafePhone()
					} else {
						oa.Name = st.Name()
						oa.Phone = st.Phone()
					}

					oa.BirthDay = time.Unix(st.BirthDay(), 0)
					oa.Siblings = int(st.Siblings())
					oa.Spouse = st.Spouse()
					oa.Money = st.Money()
				}

				checkOA(b)
			})
		})
	}
}
