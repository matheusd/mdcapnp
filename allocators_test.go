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
	state.FirstSeg = make([]byte, WordSize, WordSize)
	if a.segsCapacity > 0 {
		state.Segs = make([][]byte, 0, a.segsCapacity)
	}
	return
}

func (a *alwaysReallocAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	// Determine which seg to use.
	switch {
	case a.usePreferredSeg:
		seg = preferred
	case a.createNewSeg:
		state.Segs = append(state.Segs, []byte{})
		seg = SegmentID(len(state.Segs))
	default:
		seg = SegmentID(len(state.Segs))
	}

	var oldBuf *[]byte
	if seg == 0 {
		oldBuf = &state.FirstSeg
	} else {
		oldBuf = &state.Segs[seg-1]
	}

	oldLen := len(*oldBuf)
	newBuf := make([]byte, oldLen+int(size*WordSize), oldLen+int(size*WordSize))
	copy(newBuf, *oldBuf)
	*oldBuf = newBuf

	off = WordOffset(oldLen / WordSize)
	return
}

func (a *alwaysReallocAllocator) Reset(state *AllocState) (err error) {
	// Works as if it had inited a brand new one (doesn't reuse anything).
	return a.Init(state)
}

var globalNopFirstSeg = make([]byte, WordSize, WordSize)

type nopAllocator struct{}

func (n *nopAllocator) Init(state *AllocState) (err error) {
	state.FirstSeg = globalNopFirstSeg
	return
}

func (n *nopAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	return
}

func (n *nopAllocator) Reset(state *AllocState) (err error) {
	state.FirstSeg = globalNopFirstSeg
	return
}
