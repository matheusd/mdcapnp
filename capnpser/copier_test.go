// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// TestDeepCopyBasic tests a basic copy operation.
func TestDeepCopyBasic(t *testing.T) {
	strPtr := func(off WordOffset, dataSize wordCount16, ptrSize wordCount16) Word {
		return Word(buildRawStructPointer(off, StructSize{DataSectionSize: dataSize, PointerSectionSize: ptrSize}))
	}
	lsPtr := func(off WordOffset, elSize listElementSize, lsSize listSize) Word {
		return Word(buildRawListPointer(off, elSize, lsSize))
	}

	srcSeg := newSegmentData(
		strPtr(0, 3, 3),                                            // Root pointer (3 data, 3 ptr)
		0x1234567890abcdef, 0x1213141516171819, 0xa0b1c2d3e4f59687, // 3 data words
		strPtr(2, 1, 1),             // First root pointer is a sub-struct (sub #01)
		lsPtr(3, listElSizeByte, 8), // Second root pointer is a list (sub #02, 8 bytes)
		0x4321000000000003,          // Third root pointer is a cap pointer.
		0xfeed000000000066,          // sub #01 data 01
		0x0000000000000000,          // sub #01 ptr 01
		0x0102030405060708,          // sub #02 list contents

		0xfefefefefefefefe, // Trailer data (not part of root object, technically invalid).
	)

	t.Logf("Source seg: %x", srcSeg)

	arena := NewSingleSegmentArena(srcSeg)
	arena.ReadLimiter().InitNoLimit()
	msg := MakeMsg(arena)
	msg.RemoveDepthLimit()

	tests := []struct {
		name   string
		getSrc func(t *testing.T) AnyPointer
		want   []byte
	}{{
		name: "entire object",
		getSrc: func(t *testing.T) AnyPointer {
			root, err := msg.GetRoot()
			require.NoError(t, err)
			return root.AsAnyPointer()
		},
		want: srcSeg[:len(srcSeg)-8], // Full data minus invalid trailer.
	}, {
		name: "list within object",
		getSrc: func(t *testing.T) AnyPointer {
			root, err := msg.GetRoot()
			require.NoError(t, err)
			var ls List
			require.NoError(t, root.ReadList(1, &ls)) // Second pointer of root is a list
			return ls.AsAnyPointer()
		},
		want: newSegmentData(
			lsPtr(0, listElSizeByte, 8),
			0x0102030405060708,
		),
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dst, err := NewMessageBuilder(DefaultSimpleSingleAllocator)
			require.NoError(t, err)

			src := tc.getSrc(t)

			err = DeepCopyAndSetRoot(src, dst)
			if err != nil {
				t.Logf("Partial result: %x", dst.state.FirstSeg)
			}
			require.NoError(t, err)

			t.Logf("Final result: %x", dst.state.FirstSeg)
			got := dst.state.FirstSeg
			require.Equal(t, tc.want, got)
		})
	}

}
