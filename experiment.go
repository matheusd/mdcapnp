// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"errors"
	"unsafe"
)

type Word uint64

const WordSize = 8

type WordCount Word

type SegmentID uint64

type ReaderArena interface {
	GetWord(seg SegmentID, offset Word) (Word, error)
	ReadWord(seg SegmentID, offset Word, out *Word) error
	// ReadWords(seg SegmentID, offset Word, b []Word) (WordCount, error)
	ReadBytes(seg SegmentID, offset Word, b []byte) error
	// SegmentCount() uint64
	// Slice(seg SegmentID, offset, size Word) ([]byte, error)
}

type SingleSegmentArena struct {
	b []byte
}

func (ssa *SingleSegmentArena) ReadBytes(seg SegmentID, offset Word, b []byte) error {
	if seg != 0 {
		return ErrInvalidSegmentOffset(seg)
	}

	byteOffset := int(offset * WordSize)
	if len(ssa.b) < byteOffset+len(b) {
		return ErrInvalidOffset{AvailableLen: len(ssa.b), EndOffset: byteOffset + len(b)}
	}

	copy(b, ssa.b[byteOffset:])
	return nil
}

func (ssa *SingleSegmentArena) Slice(seg SegmentID, offset, size Word) ([]byte, error) {
	if seg != 0 {
		return nil, errors.New("no segment")
	}

	byteOffset, lenBytes := int(offset*WordSize), int(size*WordSize)
	if len(ssa.b) < byteOffset+lenBytes {
		return nil, errors.New("invalid offset")
	}

	return ssa.b[byteOffset:lenBytes], nil
}

//go:noinline
func (ssa *SingleSegmentArena) ReadWord(seg SegmentID, offset Word, out *Word) error {
	if seg != 0 {
		return errors.New("no segment")
	}

	byteOffset := offset * WordSize
	if len(ssa.b) < int(byteOffset+WordSize) {
		return errors.New("invalid offset")
	}

	copy((*[8]byte)(unsafe.Pointer(&out))[:], ssa.b[byteOffset:byteOffset+WordSize])
	return nil
}

func (ssa *SingleSegmentArena) GetWord(seg SegmentID, offset Word) (res Word, err error) {
	if seg != 0 {
		return 0, errors.New("no segment")
	}

	byteOffset := offset * WordSize
	if len(ssa.b) < int(byteOffset+WordSize) {
		return 0, errors.New("invalid offset")
	}

	copy((*[8]byte)(unsafe.Pointer(&res))[:], ssa.b[byteOffset:byteOffset+WordSize])
	return
}

type Message struct {
	arena ReaderArena
}

type Struct struct {
	msg         *Message
	segID       SegmentID
	baseOffset  Word
	dataSize    Word
	pointerSize Word
}

/*
func (s *Struct) ReadInt64(dataOffset Word, res *int64) error {
	_, err := s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, (*[8]byte)(unsafe.Pointer(res))[:])
	return err
}
*/

// _go:noinline
func (s *Struct) Int64(dataOffset Word) (res int64) {
	// 87 instructions
	/*
		sl, err := s.msg.arena.Slice(s.segID, s.baseOffset+dataOffset, WordSize)
		if err != nil {
			return 0
		}

		return int64(binary.LittleEndian.Uint64(sl))
	*/

	/*
		// cost: 71
		// var aux Word
		//s.msg.arena.ReadWord(s.segID, s.baseOffset+dataOffset, &aux)
		// res = int64(aux)
		s.msg.arena.ReadWord(s.segID, s.baseOffset+dataOffset, (*Word)(unsafe.Pointer(&res)))
		return
	*/

	/*
		// Cost 73 (assumes a big endian version will be written).
		// s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, (*[8]byte)(unsafe.Pointer(&res))[:])
		s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, *(*[]byte)(unsafe.Pointer(&res)))
		return res
	*/

	/*
		// 77 instructions
		s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, (*[8]byte)(unsafe.Pointer(&res))[:])
		return wireLEToNativeEndianInt64(res)
	*/

	// return int64(binary.LittleEndian.Uint64(*(*[]byte)(unsafe.Pointer(&res))))

	/*
		// 86 instructions
		var aux [8]byte
		s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, aux[:])
		return wireLEBytesToNativeEndianInt64(&aux)
	*/

	/*
		// 80 instructions
		data, _ := s.msg.arena.GetWord(s.segID, s.baseOffset+dataOffset)
		return wireWordLEToNativeEndianInt64(data)
	*/

	// Assumes a BE version will be written.
	// 76 instructions
	data, _ := s.msg.arena.GetWord(s.segID, s.baseOffset+dataOffset)
	return int64(data)

}

func (s *Struct) Float64(dataOffset Word) (res float64) {
	s.msg.arena.ReadBytes(s.segID, s.baseOffset+dataOffset, (*[8]byte)(unsafe.Pointer(&res))[:])
	return
}

type SmallTestStruct Struct

func (st *SmallTestStruct) Siblings() int64 {
	return (*Struct)(st).Int64(0)
}
