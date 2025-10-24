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
		strPtr(0, 3, 3),                  // Root pointer (3 data, 3 ptr)
		0x1234567890abcdef,               // root data word 0
		0x1213141516171819,               // root data word 1
		0xa0b1c2d3e4f59687,               // root data word 2
		strPtr(2, 1, 1),                  // First root pointer is a sub-struct (sub #01)
		lsPtr(3, listElSizeByte, 8),      // Second root pointer is a list (sub #02, 8 bytes)
		0x4321000000000003,               // Third root pointer is a cap pointer.
		0xfeed000000000066,               // sub #01 data 01
		strPtr(1, 0, 2),                  // sub #01 ptr 01 (defines sub#01#01)
		0x0102030405060708,               // sub #02 list contents
		strPtr(1, 0, 3),                  // sub #0101 ptr 1 (defines sub#01#01#01)
		0x789a000000000003,               // sub #0101 ptr 2 (cap pointer)
		lsPtr(8, listElSizeComposite, 6), // sub #01#01#01 ptr 1 (defines sub#01#01#01#02 struct list) // 0xcdef000000000003
		strPtr(5, 1, 1),                  // sub #01#01#01 ptr 2 (defines sub#01#01#01#01)
		0x2345000000000003,               // sub #01#01#01 ptr 3

		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data

		0xccdd0000aabb1122, // sub#01#01#01#01 data 1
		0x0000123400000003, // sub#01#01#01#01 ptr 1

		strPtr(3, 1, 1),    // sub#01#01#01#02 struct list tag word
		0x328d89b9281c7b80, // struct list item 0 data 0
		strPtr(9, 2, 2),    // struct list item 0 ptr 0 (defines #slc01)
		0xac47823f84aded95, // struct list item 1 data 0
		0x0000000000000000, // struct list item 1 ptr 0
		0xa19dd7155b398c38, // struct list item 2 data 0
		0x0000000000000000, // struct list item 2 ptr 0

		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data
		0x0000000000000000, // Orphaned data

		0x122fbac990ffd613, // #slc01 data 0
		0xd811108a437567b3, // #slc01 data 1
		0xe5d625e900000003, // #slc01 ptr 0 (cap pointer)
		0x9b55815700000003, // #slc01 ptr 1 (cap pointer)

		0xfefefefefefefefe, // Trailer data (not part of root object, technically invalid).
	)

	// t.Logf("Source seg: %x", srcSeg)

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
		want: newSegmentData(
			strPtr(0, 3, 3),                  // Root pointer (3 data, 3 ptr)
			0x1234567890abcdef,               // Root data 0
			0x1213141516171819,               // Root data 1
			0xa0b1c2d3e4f59687,               // Root data 2
			strPtr(2, 1, 1),                  // First root pointer is a sub-struct (sub #01)
			lsPtr(21, listElSizeByte, 8),     // Second root pointer is a list (sub #02, 8 bytes)
			0x4321000000000003,               // Third root pointer is a cap pointer.
			0xfeed000000000066,               // sub #01 data 01
			strPtr(0, 0, 2),                  // sub #01 ptr 01 (defines sub#01#01)
			strPtr(1, 0, 3),                  // sub #0101 ptr 1 (defines sub#01#01#01)
			0x789a000000000003,               // sub #0101 ptr 2 (cap pointer)
			lsPtr(2, listElSizeComposite, 6), // sub #01#01#01 ptr 1 (defines sub#01#01#01#02 struct list)
			strPtr(12, 1, 1),                 // sub #01#01#01 ptr 2 (defines sub#01#01#01#01)
			0x2345000000000003,               // sub #01#01#01 ptr 3

			strPtr(3, 1, 1),    // sub#01#01#01#02 struct list tag word
			0x328d89b9281c7b80, // struct list item 0 data 0
			strPtr(4, 2, 2),    // struct list item 0 ptr 0 (defines #slc01)
			0xac47823f84aded95, // struct list item 1 data 0
			0x0000000000000000, // struct list item 1 ptr 0
			0xa19dd7155b398c38, // struct list item 2 data 0
			0x0000000000000000, // struct list item 2 ptr 0

			0x122fbac990ffd613, // #slc01 data 0
			0xd811108a437567b3, // #slc01 data 1
			0xe5d625e900000003, // #slc01 ptr 0 (cap pointer)
			0x9b55815700000003, // #slc01 ptr 1 (cap pointer)

			0xccdd0000aabb1122, // sub#01#01#01#01 data 1
			0x0000123400000003, // sub#01#01#01#01 ptr 1

			0x0102030405060708, // sub #02 list contents

			// No trailer data or orphans.
		),
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
	}, {
		name: "struct within object",
		getSrc: func(t *testing.T) AnyPointer {
			root, err := msg.GetRoot()
			require.NoError(t, err)
			var st Struct
			require.NoError(t, root.ReadStruct(0, &st)) // sub#01
			require.NoError(t, st.ReadStruct(0, &st))   // sub#0101
			require.NoError(t, st.ReadStruct(0, &st))   // sub#010101
			require.NoError(t, st.ReadStruct(1, &st))   // sub#01010102
			return st.AsAnyPointer()
		},
		want: newSegmentData(
			strPtr(0, 1, 1),
			0xccdd0000aabb1122, // sub#01#01#01#01 data 1
			0x0000123400000003, // sub#01#01#01#01 ptr 1
		),
	}, {
		name: "struct list",
		getSrc: func(t *testing.T) AnyPointer {
			root, err := msg.GetRoot()
			require.NoError(t, err)
			var st Struct
			var stl StructList
			require.NoError(t, root.ReadStruct(0, &st))    // sub#01
			require.NoError(t, st.ReadStruct(0, &st))      // sub#0101
			require.NoError(t, st.ReadStruct(0, &st))      // sub#010101
			require.NoError(t, st.ReadStructList(0, &stl)) // sub#01010102
			return stl.AsAnyPointer()
		},
		want: newSegmentData(
			lsPtr(0, listElSizeComposite, 6),
			strPtr(3, 1, 1),    // sub#01#01#01#02 struct list tag word
			0x328d89b9281c7b80, // struct list item 0 data 0
			strPtr(4, 2, 2),    // struct list item 0 ptr 0 (defines #slc01)
			0xac47823f84aded95, // struct list item 1 data 0
			0x0000000000000000, // struct list item 1 ptr 0
			0xa19dd7155b398c38, // struct list item 2 data 0
			0x0000000000000000, // struct list item 2 ptr 0

			0x122fbac990ffd613, // #slc01 data 0
			0xd811108a437567b3, // #slc01 data 1
			0xe5d625e900000003, // #slc01 ptr 0 (cap pointer)
			0x9b55815700000003, // #slc01 ptr 1 (cap pointer)
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

			// t.Logf("Final result: %x", dst.state.FirstSeg)
			// t.Logf("want        : %x", tc.want)
			got := dst.state.FirstSeg
			require.Equal(t, tc.want, got)
		})
	}
}

// TestDeepCopyInvalidData tests that trying to copy various invalid data fails
// as expected.
func TestDeepCopyInvalidData(t *testing.T) {
	strPtr := func(off WordOffset, dataSize wordCount16, ptrSize wordCount16) Word {
		return Word(buildRawStructPointer(off, StructSize{DataSectionSize: dataSize, PointerSectionSize: ptrSize}))
	}

	tests := []struct {
		name       string
		srcSeg     []byte
		depthLimit uint
		wantErr    error
	}{{
		name: "out of bounds struct",
		srcSeg: newSegmentData(
			strPtr(0, 0, 1), // Root pointer
			strPtr(2, 0, 1), // Out of bounds.
		),
		wantErr: ErrObjectOutOfBounds{},
	}, {
		name:       "overflows depth limit",
		depthLimit: 3,
		srcSeg: newSegmentData(
			strPtr(0, 0, 1), // Root pointer
			strPtr(0, 0, 1), // Depth 1
			strPtr(0, 0, 1), // Depth 2
			strPtr(0, 0, 1), // Depth 3
			strPtr(0, 0, 1), // Depth 4
			strPtr(0, 0, 0), // Depth 5
		),
		wantErr: errDepthLimitExceeded,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arena := NewSingleSegmentArena(tc.srcSeg)
			arena.ReadLimiter().InitNoLimit()
			msg := MakeMsg(arena)
			if tc.depthLimit == 0 {
				msg.RemoveDepthLimit()
			} else {
				msg.SetDepthLimit(tc.depthLimit)
			}

			srcRoot, err := msg.GetRoot()
			require.NoError(t, err)

			dst, err := NewMessageBuilder(DefaultSimpleSingleAllocator)
			require.NoError(t, err)

			err = DeepCopyAndSetRoot(srcRoot.AsAnyPointer(), dst)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
