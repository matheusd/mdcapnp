// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"errors"
	"math"
)

type Struct struct {
	seg   *Segment
	arena Arena
	ptr   structPointer
}

func (s *Struct) Int64(dataOffset WordOffset) (res int64) {
	data, _ := s.seg.GetWord(WordOffset(s.ptr.dataOffset) + dataOffset)
	return int64(data)
}

func (s *Struct) Float64(dataOffset WordOffset) (res float64) {
	data, _ := s.seg.GetWord(WordOffset(s.ptr.dataOffset) + dataOffset)
	return math.Float64frombits(uint64(data))
}

type Int32DataFieldShift int

const (
	Int32FieldLo Int32DataFieldShift = 0
	Int32FieldHi Int32DataFieldShift = 32
)

// Int32 returns a data field as an int32. Given that an int32 field occupies
// either the low or high end of data word, the second parameter disambiguates
// between the two.
//
// TODO: review if this is the way to go.
func (s *Struct) Int32(fieldIndex DataFieldIndex, shift Int32DataFieldShift) int32 {
	data, _ := s.seg.GetWord(WordOffset(s.ptr.dataOffset) + WordOffset(fieldIndex))
	return int32(data >> shift)
}

// Bool returns a data field as a bool. fieldIndex points to the data word
// within the struct, while bit determines which bit (within the word)
// corresponds to the target field.
func (s *Struct) Bool(fieldIndex DataFieldIndex, bit byte) bool {
	data, _ := s.seg.GetWord(s.ptr.dataOffset + WordOffset(fieldIndex))
	return data&(1<<bit) != 0
}

func (s *Struct) ReadList(ptrIndex PointerFieldIndex, ls *List) error {
	if ptrIndex >= PointerFieldIndex(s.ptr.pointerSectionSize) {
		// TODO: return default if it exists? Or handle this at a higher
		// level?
		return errors.New("pointer at offset not set in struct")
	}

	// Determine the offset of the pointer word (given its index) and fetch
	// it.
	//
	// TODO: check if sum won't overflow?
	pointerOffset := s.ptr.dataOffset + WordOffset(s.ptr.dataSectionSize) + WordOffset(ptrIndex)
	pointer, err := s.seg.GetWord(pointerOffset)
	if err != nil {
		return err
	}

	if !isListPointer(pointer) {
		return errors.New("not a list pointer")
	}

	var lp listPointer
	lp.fromWord(pointer)

	// Determine concrete offset into segment of where the list actually
	// starts.
	if !AddWordOffsets(pointerOffset, 1, &pointerOffset) {
		return errWordOffsetSumOverflows{pointerOffset, 1}
	}
	if !AddWordOffsets(lp.startOffset, pointerOffset, &lp.startOffset) {
		return errWordOffsetSumOverflows{lp.startOffset, pointerOffset}
	}

	// Check if entire list is readable.
	fullSize := listWordCount(lp.elSize, lp.listSize)
	if err := s.seg.CheckBounds(lp.startOffset, fullSize); err != nil {
		return err
	}
	if err := s.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return err
	}

	// All good.
	ls.seg = s.seg
	ls.arena = s.arena
	ls.ptr = lp
	return nil
}
