// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

// TestStructPointer tests conversion between structPointer and raw wire
// pointers.
func TestStructPointer(t *testing.T) {
	sp := func(dataOffset WordOffset, dataSize, ptrSize wordCount16) structPointer {
		return structPointer{dataOffset: dataOffset, dataSectionSize: dataSize, pointerSectionSize: ptrSize}
	}

	tests := []struct {
		name string
		ptr  pointer
		sp   structPointer
	}{
		{name: "zeroes", ptr: 0x0000000000000000, sp: sp(0, 0, 0)},
		{name: "values", ptr: 0x8923_4567_6af3788c, sp: sp(0x1abcde23, 0x4567, 0x8923)},
		{name: "zero off", ptr: 0x8923_4567_00000000, sp: sp(0, 0x4567, 0x8923)},
		{name: "one off", ptr: 0x8923_4567_00000004, sp: sp(1, 0x4567, 0x8923)},
		{name: "neg off", ptr: 0x8923_4567_fffffffc, sp: sp(-1, 0x4567, 0x8923)},
		{name: "min off", ptr: 0x8923_4567_80000000, sp: sp(minWordOffset, 0x4567, 0x8923)},
		{name: "max off", ptr: 0x8923_4567_7ffffffc, sp: sp(maxWordOffset, 0x4567, 0x8923)},
		{name: "ones", ptr: 0xfffffffffffffffc, sp: sp(-1, 0xffff, 0xffff)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSP := tc.ptr.toStructPointer()
			require.EqualValues(t, tc.sp, gotSP)
			gotPTR := gotSP.toPointer()
			require.EqualValues(t, tc.ptr, gotPTR)
		})
	}
}

func TestListPointer(t *testing.T) {
	lp := func(startOffset WordOffset, elSize listElementSize, lsSize listSize) listPointer {
		return listPointer{startOffset: startOffset, elSize: elSize, listSize: lsSize}
	}

	tests := []struct {
		name string
		ptr  pointer
		lp   listPointer
	}{
		{name: "zeroes", ptr: 0x0000000000000001, lp: lp(0, 0, 0)},
		{name: "values", ptr: 0xa2b3c4d4_6af3788d, lp: lp(0x1abcde23, 4, 0x1456789a)},
		{name: "zero off", ptr: 0xa2b3c4d4_00000001, lp: lp(0, 4, 0x1456789a)},
		{name: "zero ls", ptr: 0x00000004_6af3788d, lp: lp(0x1abcde23, 4, 0)},
		{name: "zero el", ptr: 0xa2b3c4d0_6af3788d, lp: lp(0x1abcde23, 0, 0x1456789a)},
		{name: "max off", ptr: 0xa2b3c4d4_7ffffffd, lp: lp(maxWordOffset, 4, 0x1456789a)},
		{name: "max ls", ptr: 0xfffffffc_6af3788d, lp: lp(0x1abcde23, 4, MaxListSize)},
		{name: "max el", ptr: 0xa2b3c4d7_6af3788d, lp: lp(0x1abcde23, 0b111, 0x1456789a)},
		{name: "neg off", ptr: 0xa2b3c4d4_fffffffd, lp: lp(-1, 4, 0x1456789a)},
		{name: "min off", ptr: 0xa2b3c4d4_80000001, lp: lp(minWordOffset, 4, 0x1456789a)},
		{name: "ones", ptr: 0xfffffffffffffffd, lp: lp(-1, 0b111, MaxListSize)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLP := tc.ptr.toListPointer()
			require.EqualValues(t, tc.lp, gotLP)
			gotPTR := gotLP.toPointer()
			require.EqualValues(t, tc.ptr, gotPTR)
		})
	}

}
