// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"bytes"
	"math"
	"testing"

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
}
