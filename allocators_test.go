// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

// alwaysReallocAllocator is a test allocator. It has the following properties:
//
// - Always re-allocates buffers.
// - Buffers always have len == cap (no ability to append without realloc).
// - Allocate() calls determine the segment to allocate on:
//   - If usePreferredSeg is true, allocates in the preferred segment.
//   - If createNewSeg is true, creates a new segment and allocates in it.
//   - Otherwise, allocates in the last segment.
type alwaysReallocAllocator struct {
	segsCapacity    uint
	createNewSeg    bool
	usePreferredSeg bool
}

func (a *alwaysReallocAllocator) Init(state *AllocState) (err error) {
	state.Segs = make([][]byte, 1, max(1, a.segsCapacity))
	state.Segs[0] = make([]byte, WordSize, WordSize)
	return
}

func (a *alwaysReallocAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	// Determine which seg to use.
	switch {
	case a.usePreferredSeg:
		seg = preferred
	case a.createNewSeg:
		state.Segs = append(state.Segs, []byte{})
		seg = SegmentID(len(state.Segs) - 1)
	default:
		seg = SegmentID(len(state.Segs) - 1)
	}
	if a.createNewSeg {
	}

	oldLen := len(state.Segs[seg])
	newBuf := make([]byte, oldLen+int(size*WordSize), oldLen+int(size*WordSize))
	copy(newBuf, state.Segs[seg])
	state.Segs[seg] = newBuf

	off = WordOffset(oldLen / WordSize)
	return
}

func (a *alwaysReallocAllocator) Reset(state *AllocState) (err error) {
	// Works as if it had inited a brand new one (doesn't reuse anything).
	return a.Init(state)
}

var globalNopAllocatorSegs = [][]byte{make([]byte, WordSize, WordSize)}

type nopAllocator struct{}

func (n *nopAllocator) Init(state *AllocState) (err error) {
	state.Segs = globalNopAllocatorSegs
	return
}

func (n *nopAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	return
}

func (n *nopAllocator) Reset(state *AllocState) (err error) {
	state.Segs = globalNopAllocatorSegs
	return
}
