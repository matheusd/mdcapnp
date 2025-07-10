// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"errors"
	"math"
)

type Struct struct {
	seg             Segment
	dataStartOffset WordOffset
	dataSize        WordCount
	pointerSize     WordCount
}

func (s *Struct) Int64(dataOffset WordOffset) (res int64) {
	data, _ := s.seg.GetWord(s.dataStartOffset + dataOffset)
	return int64(data)
}

func (s *Struct) Float64(dataOffset WordOffset) (res float64) {
	data, _ := s.seg.GetWord(s.dataStartOffset + dataOffset)
	return math.Float64frombits(uint64(data))
}

func (s *Struct) ReadList(pointerOffset WordOffset, ls *List) error {
	if s.dataSize+WordCount(pointerOffset) >= s.pointerSize {
		return errors.New("pointer at offset not set in struct")
	}

	finalPointerOffset := s.dataStartOffset + WordOffset(s.dataSize) + pointerOffset
	pointer, err := s.seg.GetWord(finalPointerOffset)
	if err != nil {
		return err
	}

	if !isPointerList(pointer) {
		return errors.New("not a list pointer")
	}

	ls.seg = s.seg
	ls.fromPointerWord(finalPointerOffset, pointer)
	return nil
}
