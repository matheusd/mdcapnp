// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// BenchmarkMsgReadRoot benchmarks reading the root struct of a message.
func BenchmarkMsgReadRoot(b *testing.B) {
	buf := appendWords(nil, 0x0000000100000000, 0x0000000000000000, 0x0000000000000000)
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
