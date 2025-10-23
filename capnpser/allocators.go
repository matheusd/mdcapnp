// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"slices"
)

// SimpleSingleAllocator is a simple, single-segment allocator. It uses the
// standard go procedure to allocate and resize the segment data, without
// additional logic for trying to determine how to best size the segment.
//
// It can be configured to either truncate or re-allocate the segment buffer on
// resets.
type SimpleSingleAllocator struct {
	initialSize    WordCount
	reallocOnReset bool
}

func NewSimpleSingleAllocator(initialSize WordCount, reallocOnReset bool) *SimpleSingleAllocator {
	// One word for header + one word for root pointer.
	if initialSize < 2 {
		panic("minimum initial size is 2 words")
	}
	if initialSize > MaxValidWordCount {
		panic("initial size is larger than max valid word count")
	}
	return &SimpleSingleAllocator{initialSize: initialSize, reallocOnReset: reallocOnReset}
}

func (s *SimpleSingleAllocator) Init(state *AllocState) (err error) {
	state.HeaderBuf = make([]byte, WordSize, s.initialSize*WordSize)
	state.FirstSeg = state.HeaderBuf[8:16]
	return
}

func (s *SimpleSingleAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	sizeBytes := int(size.ByteCount())
	lenSeg0 := len(state.FirstSeg)
	newLenSeg0 := lenSeg0 + sizeBytes
	freeCap := cap(state.FirstSeg) - lenSeg0
	off = WordOffset(lenSeg0 / WordSize)
	if freeCap < sizeBytes {
		// Resize needed.
		state.HeaderBuf = slices.Grow(state.HeaderBuf, newLenSeg0)
		state.FirstSeg = state.HeaderBuf[8 : 8+newLenSeg0]
	} else {
		state.FirstSeg = state.FirstSeg[:newLenSeg0]
	}

	return
}

func (s *SimpleSingleAllocator) Reset(state *AllocState) (err error) {
	if s.reallocOnReset {
		s.Init(state)
	} else {
		clear(state.HeaderBuf[:len(state.FirstSeg)+8])

		// Truncate segment 0 to root word.
		state.FirstSeg = state.FirstSeg[:8]
	}
	return
}

var DefaultSimpleSingleAllocator = &SimpleSingleAllocator{initialSize: 1024 / WordSize}

type PoolableAllocatorPoolIntf interface {
	Get(size WordCount) []byte
	Put(b []byte)
}

type SingleSegmentPoolableAllocator struct {
	initialSize WordCount
	pool        PoolableAllocatorPoolIntf
}

func (p *SingleSegmentPoolableAllocator) Init(state *AllocState) (err error) {
	buf := p.pool.Get(p.initialSize)
	state.SetHeaderAndSeg0(buf[:16], 1)
	return
}

func (p *SingleSegmentPoolableAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	sizeBytes := int(size.ByteCount())
	lenSeg0 := len(state.FirstSeg)
	newLenSeg0 := lenSeg0 + sizeBytes
	freeCap := cap(state.FirstSeg) - lenSeg0
	off = WordOffset(lenSeg0 / WordSize)
	if freeCap < sizeBytes {
		// Resize needed.
		// state.HeaderBuf = slices.Grow(state.HeaderBuf, newLenSeg0)
		//state.FirstSeg = state.HeaderBuf[8 : 8+newLenSeg0]
		oldBuf := state.HeaderBuf[:8+lenSeg0]
		newBuf := p.pool.Get(WordCount(newLenSeg0/WordSize) + 1)

		copy(newBuf, oldBuf)
		state.SetHeaderAndSeg0(newBuf[:newLenSeg0+8], 1)
		clear(oldBuf)
		p.pool.Put(oldBuf)
	} else {
		state.FirstSeg = state.FirstSeg[:newLenSeg0]
	}

	return

}

func (p *SingleSegmentPoolableAllocator) Reset(state *AllocState) (err error) {
	oldBuf := state.HeaderBuf[:8+len(state.FirstSeg)]
	clear(oldBuf)
	if cap(oldBuf) == int(p.initialSize) {
		// Can still use the same slice.
		state.SetHeaderAndSeg0(oldBuf[:16], 1)
	} else {
		// Cannot use the same slice. Fetch one with the initial size.
		p.pool.Put(oldBuf)
		p.Init(state)
	}
	return nil
}

func NewSingleSegmentPoolableAllocator(initialSize WordCount, pool PoolableAllocatorPoolIntf) *SingleSegmentPoolableAllocator {
	return &SingleSegmentPoolableAllocator{
		initialSize: initialSize,
		pool:        pool,
	}
}
