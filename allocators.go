// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "slices"

type SimpleSingleAllocator struct {
	initialSize int
}

func (s SimpleSingleAllocator) Init(state *AllocState) (err error) {
	state.HeaderBuf = make([]byte, 8, s.initialSize)
	state.Segs = make([][]byte, 1)
	state.Segs[0] = state.HeaderBuf[8:16]
	return
}

func (s SimpleSingleAllocator) Allocate(state *AllocState, preferred SegmentID, size WordCount) (seg SegmentID, off WordOffset, err error) {
	segbuf := state.Segs[0]
	sizeBytes := int(size.ByteCount())
	freeCap := cap(segbuf) - len(segbuf)
	if freeCap < sizeBytes {
		// Resize needed.
		state.HeaderBuf = slices.Grow(state.HeaderBuf, len(segbuf)+sizeBytes)
		state.Segs[0] = state.HeaderBuf[8:len(segbuf)]
		segbuf = state.Segs[0]
	}

	// Increase len of segment 0.
	off = WordOffset(len(segbuf) / WordSize)
	state.Segs[0] = segbuf[:len(segbuf)+sizeBytes]
	return
}

func (s SimpleSingleAllocator) Reset(state *AllocState) (err error) {
	// Truncate segment 0 to root word.
	clear(state.HeaderBuf)
	clear(state.Segs[0])
	state.Segs[0] = state.Segs[0][:8]
	return
}

var DefaultSimpleSingleAllocator = SimpleSingleAllocator{initialSize: 1024}
