// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "slices"

type SimpleSingleAllocator struct {
	initialSize int
}

func (s SimpleSingleAllocator) Init() (initState AllocState, err error) {
	initState.HeaderBuf = make([]byte, 8, s.initialSize)
	initState.Segs = make([][]byte, 1)
	initState.Segs[0] = initState.HeaderBuf[8:16]
	return
}

func (s SimpleSingleAllocator) Allocate(prevState AllocState, preferred SegmentID, size WordCount) (nextState AllocState, seg SegmentID, off WordOffset, err error) {
	segbuf := prevState.Segs[0]
	sizeBytes := int(size.ByteCount())
	freeCap := cap(segbuf) - len(segbuf)
	nextState = prevState
	if freeCap < sizeBytes {
		// Resize needed.
		nextState.HeaderBuf = slices.Grow(nextState.HeaderBuf, len(segbuf)+sizeBytes)
		nextState.Segs[0] = nextState.HeaderBuf[8:len(segbuf)]
		segbuf = nextState.Segs[0]
	}

	// Increase len of segment 0.
	off = WordOffset(len(segbuf) / WordSize)
	nextState.Segs[0] = segbuf[:len(segbuf)+sizeBytes]
	return
}

func (s SimpleSingleAllocator) Reset(lastState AllocState) (blankState AllocState, err error) {
	// Truncate segment 0.
	blankState = lastState
	clear(blankState.HeaderBuf)
	clear(blankState.Segs[0])
	blankState.Segs[0] = blankState.Segs[0][:8]
	return
}

var DefaultSimpleSingleAllocator = SimpleSingleAllocator{initialSize: 1024}
