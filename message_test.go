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
	arena := NewSingleSegmentArena(buf, false, nil)
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

// BenchmarkDecodeGoserbenchSmallStruct benchmarks decoding a goserbench
// SmallStruct under various configurations.
func BenchmarkDecodeGoserbenchSmallStruct(b *testing.B) {
	var oa goserbenchSmallStruct
	checkOA := func(b *testing.B) {
		require.Equal(b, "slimshady0123456", oa.Name)
		require.Equal(b, int64(0x1011121314151617), oa.BirthDay.Unix())
		require.Equal(b, "phone678", oa.Phone) // FIXME phone67890
		require.Equal(b, int(0x66669999), oa.Siblings)
		require.Equal(b, true, oa.Spouse)
		require.Equal(b, uint64(0xabcd0000ef01), math.Float64bits(oa.Money))
	}

	// Skip the header.
	segBuf := testdata.GoserbenchSampleA[8:]

	tests := []struct {
		rl     newRLFunc
		unsafe bool
	}{
		{rl: nilReadLimiter, unsafe: true},
		{rl: nilReadLimiter, unsafe: false},
		{rl: NewConcurrentUnsafeReadLimiter, unsafe: true},
		{rl: NewConcurrentUnsafeReadLimiter, unsafe: false},
		{rl: NewReadLimiter, unsafe: true},
		{rl: NewReadLimiter, unsafe: false},
	}

	for _, tc := range tests {
		b.Run(fmt.Sprintf("%v/unsafe=%v", rlTestName(tc.rl), tc.unsafe), func(b *testing.B) {
			b.Run("reuse all", func(b *testing.B) {
				rl := tc.rl(MaxReadLimiterLimit)
				arena := NewSingleSegmentArena(segBuf, false, rl)
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
				rl := tc.rl(MaxReadLimiterLimit)
				arena := NewSingleSegmentArena(segBuf, false, rl)

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					arena.Reset(segBuf, false)
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
					rl := tc.rl(MaxReadLimiterLimit)
					arena := NewSingleSegmentArena(segBuf, false, rl)
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
