// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"unsafe"
)

type SegmentID uint32

type SegmentCount uint32

type Segment struct {
	b []byte
}

// uncheckedOpenSlice returns a slice starting at the provided offset up to the
// end of the buffer without checking bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedTailSlice(offset WordOffset) []byte {
	return ms.b[offset*WordSize:]
}

func (ms *Segment) intLen() int {
	return len(ms.b)
}

// wordLen returns the number of words in the segment. This assumes the segment
// has a valid length (i.e. <= MaxValidWordCount).
func (ms *Segment) wordLen() WordCount {
	return WordCount(len(ms.b) / WordSize)
}

func (ms *Segment) getWordAsPointer(offset WordOffset) (pointer, error) {
	w, err := ms.GetWord(offset)
	return pointer(w), err
}

// checkSliceBounds checks whether a subsequent call to [uncheckedSlice] with
// the same arguments will fail. If this function returns true, immediately
// calling [uncheckedSlice] will generate a valid slice.
func (ms *Segment) checkSliceBounds(offset WordOffset, size ByteCount) error {
	startOffset := int(offset * WordSize) // FIXME: check for overflows in 32bit archs
	endOffset := startOffset + int(size)
	if endOffset > ms.intLen() {
		return ErrInvalidMemOffset{AvailableLen: ms.intLen(), Offset: endOffset}
	}

	return nil
}

// uncheckedSlice returns a slice without checking for bounds. Bounds MUST be
// checked first by calling checkSliceBounds, otherwise this may panic.
//
// These functions are split to allow uncheckedSlice to be trivially inlineable.
func (ms *Segment) uncheckedSlice(offset WordOffset, size ByteCount) []byte {
	startOffset := int(offset * WordSize)
	return ms.b[startOffset : startOffset+int(size)]
}

// uncheckedUnsafeString returns a subslice of the segment as an unsafe string
// without performing a bounds check.
func (ms *Segment) uncheckedUnsafeString(offset WordOffset, size ByteCount) string {
	startOffset := int(offset * WordSize)
	buf := ms.b[startOffset : startOffset+int(size)]
	return *(*string)(unsafe.Pointer(&buf))
}

func (ms *Segment) hasRootPointer() bool {
	return len(ms.b) >= WordSize
}

func (ms *Segment) Read(offset WordOffset, b []byte) (int, error) {
	byteOffset := int(offset * WordSize)
	if byteOffset >= ms.intLen() {
		return 0, ErrInvalidMemOffset{AvailableLen: ms.intLen(), Offset: byteOffset}
	}

	n := copy(b, ms.b[byteOffset:])
	return n, nil
}

// checkBounds checks whether the given offset and size are within bounds of
// this segment.
func (ms *Segment) checkBounds(offset WordOffset, size WordCount) (err error) {
	if wordLen := ms.wordLen(); offset < 0 || WordCount(offset)+size > wordLen {
		err = ErrObjectOutOfBounds{Offset: offset, Size: size, WordLen: wordLen}
	}
	return
}

type Arena struct {
	// fb is the full, framed data for the arena (includes header and arena
	// size framing when != nil).
	fb []byte

	// s is the first segment. It is the only segment in single-segment
	// arenas.
	s Segment

	// segs are the additional segments in multi-segment arenas. The segment
	// at index 0 is the segment with id 1, and so on.
	segs *[]*Segment

	rl ReadLimiter

	// notResetable is true if this arena is not resettable from the public
	// API (i.e. it is a readerArena in a MessageBuilder).
	notResetable bool
}

func (arena *Arena) ReadLimiter() *ReadLimiter {
	return &arena.rl
}

// segment returns the given segment without bounds check.
func (arena *Arena) segment(id SegmentID) *Segment {
	if id == 0 {
		return &arena.s
	}

	index := int(id - 1)
	segs := *arena.segs
	return segs[index]
}

func (arena *Arena) Segment(id SegmentID) (*Segment, error) {
	if arena == nil {
		return nil, errArenaNotInitialized
	}

	if id == 0 {
		return &arena.s, nil
	}

	index := int(id - 1)
	segs := *arena.segs
	if index >= len(segs) {
		return nil, ErrUnknownSegment(id)
	}

	return segs[index], nil
}

// DecodeSingleSegment decodes the given buffer as a single segment arena.
func (arena *Arena) DecodeSingleSegment(fb []byte) error {
	b, err := decodeSingleSegmentStream(fb)
	if err != nil {
		return err
	}
	arena.Reset(b)
	arena.fb = fb
	return nil
}

// TotalSize returns the sum of data currently referenced by this Arena.
func (arena *Arena) TotalSize() WordCount {
	var bc ByteCount
	if arena.fb != nil {
		bc = ByteCount(len(arena.fb))
	} else if arena.segs == nil {
		bc = ByteCount(arena.s.intLen())
	} else {
		bc = ByteCount(arena.s.intLen())
		segs := *arena.segs
		for _, s := range segs {
			bc += ByteCount(s.intLen())
		}
	}

	res, _ := bc.StorageWordCount()
	return res
}

// RawDataCopy returns a copy of the underlying arena data. This is mostly
// useful for debugging issues.
func (arena *Arena) RawDataCopy() (res [][]byte) {
	var segs []*Segment
	if arena.segs != nil {
		segs = *arena.segs
	}

	res = make([][]byte, 1+len(segs))
	res[0] = append([]byte(nil), arena.s.b...)
	if arena.segs != nil {
		for i := range segs {
			res[i+1] = append([]byte(nil), segs[i].b...)
		}
	}
	return
}

// Reset the arena to the given single-segment buffer. This may panic if the
// arena is not resettable.
func (arena *Arena) Reset(b []byte) {
	if arena.notResetable {
		panic("arena is not resetable")
	}

	arena.s.b = b
	arena.fb = nil
	arena.rl.Reset()
}

func NewSingleSegmentArena(b []byte) *Arena {
	var arena Arena
	arena.Reset(b)
	return &arena
}

func DecodeArena(fb []byte) (*Arena, error) {
	// TODO: decode
	var arena Arena
	if err := arena.DecodeSingleSegment(fb); err != nil {
		return nil, err
	}
	return &arena, nil
}
