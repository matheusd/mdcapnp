// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/internal/testdata"
)

// TestGoserbenchMarshal tests if serialization of the handwritten goserbench
// small struct is correct.
func TestGoserbenchMarshal(t *testing.T) {
	alloc := NewSimpleSingleAllocator(12, false)

	mb, err := NewMessageBuilder(alloc)
	require.NoError(t, err)

	// Generate the exact same structure that was serialized on testdata
	// sample.
	st, err := NewGoserbenchSmallStruct(mb)
	require.NoError(t, err)

	st.SetBirthDay(0x1011121314151617)
	st.SetSiblings(0x66669999)
	st.SetSpouse(true)
	st.SetMoney(math.Float64frombits(0xabcd0000ef01))
	st.SetName("slimshady0123456")
	st.SetPhone("phone67890")
	st.SetAsRoot(mb)

	ser, err := mb.Serialize()
	require.NoError(t, err)

	if !bytes.Equal(ser, testdata.GoserbenchSampleA) {
		t.Logf("    ser %x", ser)
		t.Logf(" sample %x", testdata.GoserbenchSampleA)
		t.Fatal("Generated and sample are not equal")
	}
}

// TestGoserbenchWrite tests writing a goserbench struct.
func TestGoserbenchWrite(t *testing.T) {
	alloc := NewSimpleSingleAllocator(12, false)

	mb, err := NewMessageBuilder(alloc)
	require.NoError(t, err)

	var sst GoserbenchSmallStructType
	sst.BirthDay = 0x1011121314151617
	sst.Siblings = 0x66669999
	sst.Spouse = true
	sst.Money = math.Float64frombits(0xabcd0000ef01)
	sst.Name = "slimshady0123456"
	sst.Phone = "phone67890"

	WriteRootGoserbenchSmallStructType(&sst, mb)

	ser, err := mb.Serialize()
	require.NoError(t, err)

	if !bytes.Equal(ser, testdata.GoserbenchSampleA) {
		t.Logf("    ser %x", ser)
		t.Logf(" sample %x", testdata.GoserbenchSampleA)
		t.Fatal("Generated and sample are not equal")
	}
}

// BenchmarkGoserbenchUnmarhsmal benchmarks decoding a goserbench SmallStruct
// under various configurations.
func BenchmarkGoserbenchUnmarshal(b *testing.B) {
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

// BenchmarkGoserbenchMarshal simulates the marshal goserbench benchmark.
func BenchmarkGoserbenchMarshal(b *testing.B) {
	alloc := NewSimpleSingleAllocator(16, false)

	b.Run("reuse buffer", func(b *testing.B) {
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
			st.SetSiblings(int32(i))
			st.SetSpouse(i%2 == 0)
			st.SetMoney(float64(i) * 10.5)
			st.SetName("slimshady0123456")
			st.SetPhone("phone678")
			if err := st.SetAsRoot(mb); err != nil {
				b.Fatal(err)
			}

			_, err = mb.Serialize()
			if err != nil {
				b.Fatal(err)
			}

			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reuse mb", func(b *testing.B) {
		alloc := NewSimpleSingleAllocator(16, true)
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
			st.SetSiblings(int32(i))
			st.SetSpouse(i%2 == 0)
			st.SetMoney(float64(i) * 10.5)
			st.SetName("slimshady0123456")
			st.SetPhone("phone678")
			if err := st.SetAsRoot(mb); err != nil {
				b.Fatal(err)
			}

			_, err = mb.Serialize()
			if err != nil {
				b.Fatal(err)
			}

			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reuse none", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			mb, err := NewMessageBuilder(alloc)
			require.NoError(b, err)

			st, err := NewGoserbenchSmallStruct(mb)
			if err != nil {
				b.Fatal(err)
			}

			st.SetBirthDay(int64(i))
			st.SetSiblings(int32(i))
			st.SetSpouse(i%2 == 0)
			st.SetMoney(float64(i) * 10.5)
			st.SetName("slimshady0123456")
			st.SetPhone("phone678")
			if err := st.SetAsRoot(mb); err != nil {
				b.Fatal(err)
			}

			_, err = mb.Serialize()
			if err != nil {
				b.Fatal(err)
			}

			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("write reuse all", func(b *testing.B) {
		mb, err := NewMessageBuilder(alloc)
		require.NoError(b, err)

		var sst GoserbenchSmallStructType
		sst.BirthDay = 0x1011121314151617
		sst.Siblings = 0x66669999
		sst.Spouse = true
		sst.Money = math.Float64frombits(0xabcd0000ef01)
		sst.Name = "slimshady0123456"
		sst.Phone = "phone67890"

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			if err := WriteRootGoserbenchSmallStructType(&sst, mb); err != nil {
				b.Fatal(err)
			}
			_, err = mb.Serialize()
			if err != nil {
				b.Fatal(err)
			}

			err = mb.Reset()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
