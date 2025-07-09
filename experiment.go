// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"errors"
	"math"
)

type Word uint64

const WordSize = 8

type WordCount Word

type SegmentID uint64

type ReaderArena interface{}

type SingleSegmentArena struct{}

type Segment interface {
	GetWord(offset Word) (res Word, err error)
	Read(offset Word, b []byte) (int, error)
}

type MemSegment struct {
	b  []byte
	rl *ReadLimiter
}

func (ms *MemSegment) GetWord(offset Word) (res Word, err error) {
	if err = ms.rl.CanRead(1); err != nil {
	} else if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.BigEndian.Uint64(ms.b[byteOffset:]))

		// copy((*[8]byte)(unsafe.Pointer(&res))[:], ms.b[byteOffset:])

		// Assumes a big endian version is written. Note: this is
		// counterintuitive, double check.
		// res = Word(binary.BigEndian.Uint64((*[8]byte)(unsafe.Pointer(&res))[:]))
	}
	return
}

func (ms *MemSegment) Read(offset Word, b []byte) (int, error) {
	if err := ms.rl.CanRead(1); err != nil {
		return 0, err
	}

	byteOffset := int(offset * WordSize)
	if byteOffset >= len(ms.b) {
		return 0, ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: byteOffset}
	}

	n := copy(b, ms.b[byteOffset:])
	return n, nil
}

type Message struct {
	arena ReaderArena
}

type Struct struct {
	seg             Segment
	dataStartOffset Word
	dataSize        Word
	pointerSize     Word
}

func (s *Struct) Int64(dataOffset Word) (res int64) {
	data, _ := s.seg.GetWord(s.dataStartOffset + dataOffset)
	return int64(data)
}

func (s *Struct) Float64(dataOffset Word) (res float64) {
	data, _ := s.seg.GetWord(s.dataStartOffset + dataOffset)
	return math.Float64frombits(uint64(data))
}

func (s *Struct) ReadList(pointerOffset Word, ls *List) error {
	if s.dataSize+pointerOffset >= s.pointerSize {
		return errors.New("pointer at offset not set in struct")
	}

	finalPointerOffset := s.dataStartOffset + s.dataSize + pointerOffset
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

type ListElementSize byte

const (
	ListElSizeVoid      ListElementSize = 0
	ListElSizeBit       ListElementSize = 1
	ListElSizeByte      ListElementSize = 2
	ListElSizeComposite ListElementSize = 7
)

type ListSize uint64

type List struct {
	seg        Segment
	baseOffset Word
	elSize     ListElementSize
	listSize   ListSize
}

func isPointerList(p Word) bool {
	return (p & 0x03) == 1
}

func (ls *List) fromPointerWord(pointerOffset, w Word) {
	ls.baseOffset = pointerOffset + (w & 0xfffffffc >> 2) + 1
	ls.elSize = ListElementSize(w & 0x300000000 >> 32)
	ls.listSize = ListSize(w & 0xfffffff800000000 >> 35)
}

func (ls *List) LenBytes() int {
	switch ls.elSize {
	case ListElSizeVoid:
		return 0
	case ListElSizeBit:
		return int(ls.listSize)
	case ListElSizeByte:
		return int(ls.listSize)
	case ListElSizeComposite:
		return int(ls.listSize * WordSize)
	default:
		panic("unknown el size")
	}
}

func (ls *List) Read(b []byte) (n int, err error) {
	n = min(len(b), ls.LenBytes())
	return ls.seg.Read(ls.baseOffset, b[:n])
}

type SmallTestStruct Struct

func (st *SmallTestStruct) Siblings() int64 {
	return (*Struct)(st).Int64(0)
}

func (st *SmallTestStruct) ReadNameField(ls *List) error {
	return (*Struct)(st).ReadList(0, ls) // First pointer.
}
