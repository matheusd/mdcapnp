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
	dl    depthLimit
	ptr   structPointer
}

// HasData returns true if the specified word in the data section is set in this
// struct.
func (s *Struct) HasData(dataIndex DataFieldIndex) bool {
	return s.ptr.dataOffset > 0 && dataIndex < DataFieldIndex(s.ptr.dataSectionSize)
}

func (s *Struct) Int64(dataIndex DataFieldIndex) (res int64) {
	if s.HasData(dataIndex) {
		data, _ := s.seg.GetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = int64(data)
	}
	return
}

func (s *Struct) Float64(dataIndex DataFieldIndex) (res float64) {
	if s.HasData(dataIndex) {
		data, _ := s.seg.GetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = math.Float64frombits(uint64(data))
	}
	return
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
func (s *Struct) Int32(dataIndex DataFieldIndex, shift Int32DataFieldShift) (res int32) {
	if s.HasData(dataIndex) {
		data, _ := s.seg.GetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = int32(data >> shift)
	}
	return
}

// Bool returns a data field as a bool. fieldIndex points to the data word
// within the struct, while bit determines which bit (within the word)
// corresponds to the target field.
func (s *Struct) Bool(dataIndex DataFieldIndex, bit byte) (res bool) {
	if s.HasData(dataIndex) {
		data, _ := s.seg.GetWord(dataIndex.uncheckedWordOffset(s.ptr.dataOffset))
		res = data&(1<<bit) != 0
	}
	return res
}

func (s *Struct) ReadList(ptrIndex PointerFieldIndex, ls *List) error {
	// Check if we can descend further into the struct (to fetch the first
	// list pointer).
	listDL, ok := s.dl.dec()
	if !ok {
		return errDepthLimitExceeded
	}

	seg := s.seg

	// Check if this pointer is set within the pointer section.
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
	ptr, err := seg.getWordAsPointer(pointerOffset)
	if err != nil {
		return err
	}

	// De-ref far pointers into the concrete list segment and near pointer.
	var err error
	ptrType := ptr.pointerType()
	if ptrType == pointerTypeFarPointer {
		seg, ptr, listDL, err = derefFarPointer(s.arena, listDL, ptr)
		if err != nil {
			return err
		}
		ptrType = ptr.pointerType()
	}

	// Check if it is a list pointer.
	if ptrType != pointerTypeList {
		return errNotListPointer
	}
	lp := ptr.toListPointer()

	// Determine concrete offset into segment of where the list actually
	// starts.
	if pointerOffset, ok = addWordOffsets(pointerOffset, 1); !ok {
		return errWordOffsetSumOverflows{pointerOffset, 1}
	}
	if lp.startOffset, ok = addWordOffsets(lp.startOffset, pointerOffset); !ok {
		return errWordOffsetSumOverflows{lp.startOffset, pointerOffset}
	}

	// Check if entire list is readable.
	fullSize := listWordCount(lp.elSize, lp.listSize)
	if err := s.seg.CheckBounds(lp.startOffset, fullSize); err != nil {
		return err
	}

	// If list elements have zero size, count them as one byte per element,
	// to avoid vulnerabilities where large simulated lists are iterated.
	if lp.elSize == listElSizeVoid {
		fullSize = WordCount(lp.listSize / WordSize)
	}
	if err := s.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return err
	}

	// All good.
	ls.seg = seg
	ls.arena = s.arena
	ls.ptr = lp
	ls.dl = listDL
	return nil
}
