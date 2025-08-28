// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"math"
)

type UnsafeStructBuilder struct {
	off WordOffset // Concrete offset into segment where struct data begins.
	sz  StructSize
	seg SegmentBuilder
}

func (sb *UnsafeStructBuilder) StructSize() WordCount {
	return sb.sz.TotalSize()
}

func (sb *UnsafeStructBuilder) SetSelfPointer(segb SegmentBuilder, ptrOff WordOffset) error {
	/*
		if sb.seg.mb != segb.mb {
			return fmt.Errorf("sb.mb vs mb: %w", errDifferentMsgBuilders)
		}

		// TODO: handle inter-segment pointer.
		if sb.seg.id != segb.id {
			panic("needs handling")
		}

		// Write the first word of seg0.
		sp := structPointer{
			// The wire offset is relative to the end of first word.
			dataOffset:         sb.off - ptrOff - 1,
			dataSectionSize:    sb.sz.DataSectionSize,
			pointerSectionSize: sb.sz.PointerSectionSize,
		}

		segb.uncheckedSetWord(ptrOff, Word(sp.toPointer()))
		return nil
	*/

	segb.uncheckedSetWord(ptrOff, Word(buildRawStructPointer(sb.off-ptrOff-1, sb.sz)))
	return nil
}

func (sb *UnsafeStructBuilder) SetInt64(dataIndex DataFieldIndex, v int64) {
	// TODO: check if dataIndex < sb.sz.DataSectionSize and ignore if not?
	// This would allow building a canonical structure.
	sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(v))
}

func (sb *UnsafeStructBuilder) SetInt32(dataIndex DataFieldIndex, mask Int32DataFieldSetMask, v int32) {
	sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), Word(mask), Word(v))
}

func (sb *UnsafeStructBuilder) SetBool(dataIndex DataFieldIndex, bit int, v bool) {
	/*
		if v {
			sb.seg.uncheckedSetBit(dataIndex.uncheckedWordOffset(sb.off), bit)
		} else {
			sb.seg.uncheckedClearBit(dataIndex.uncheckedWordOffset(sb.off), bit)
		}
	*/
	sb.seg.uncheckedMaskAndMergeWord(dataIndex.uncheckedWordOffset(sb.off), ^(1 << bit), boolToWord(v)<<bit)
}

func (sb *UnsafeStructBuilder) SetFloat64(dataIndex DataFieldIndex, v float64) {
	// Structure already fully allocated, no need to check for
	// bounds.
	sb.seg.uncheckedSetWord(dataIndex.uncheckedWordOffset(sb.off), Word(math.Float64bits(v)))
}

func (sb *UnsafeStructBuilder) SetString(ptrIndex PointerFieldIndex, v string, startOffset WordOffset) (nextOffset WordOffset) {
	textLen := uint(len(v) + 1)
	sb.seg.copyStringTo(startOffset, v)
	nextOffset = startOffset + WordOffset(uintBytesToWordAligned(textLen))
	ptrOff := ptrIndex.uncheckedWordOffset(sb.off + WordOffset(sb.sz.DataSectionSize))
	sb.seg.uncheckedSetWord(ptrOff, Word(buildRawListPointer(startOffset-ptrOff-1, listElSizeByte, listSize(textLen))))
	return
}
