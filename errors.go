// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"errors"
	"fmt"
)

type ErrInvalidSegmentOffset uint64

func (err ErrInvalidSegmentOffset) Error() string {
	return fmt.Sprintf("not a valid segment offset: %d", uint64(err))
}

type ErrInvalidMemOffset struct {
	Offset       int
	AvailableLen int
}

func (err ErrInvalidMemOffset) Error() string {
	return fmt.Sprintf("invalid offset: wanted to read offset %d when only %d bytes were available",
		err.Offset, err.AvailableLen)
}

type ErrReadLimitExceeded struct {
	Target WordCount
}

func (err ErrReadLimitExceeded) Error() string {
	return fmt.Sprintf("read limit exceeded when attempting to read %d words", err.Target)
}

func (err ErrReadLimitExceeded) Is(target error) bool {
	_, ok := target.(ErrReadLimitExceeded)
	return ok
}

var ErrNotStructPointer = errors.New("pointer is not a struct pointer")

type ErrObjectOutOfBounds struct {
	Offset WordOffset
	Size   WordCount
	Len    int
}

func (err ErrObjectOutOfBounds) Error() string {
	return fmt.Sprintf("object at offset 0x%016x with size %d is out of bounds (segment length is %d)",
		err.Offset, err.Size, err.Len)
}

type ErrUnknownSegment SegmentID

func (err ErrUnknownSegment) Error() string {
	return fmt.Sprintf("segment with ID %d does not exist in arena", uint64(err))
}

var errSegmentNotInitialized = errors.New("segment was not initialized")

var errArenaNotInitialized = errors.New("arena was not initialized")

var errNotOneByteElList = errors.New("list is not a one-byte-per-element list")

var errDepthLimitExceeded = errors.New("depth limit exceeded")

var errNotListPointer = errors.New("not a list pointer")

var errInvalidNegativeStructOffset = errors.New("invalid negative struct offset")

type errWordOffsetSumOverflows struct {
	a, b WordOffset
}

func (err errWordOffsetSumOverflows) Error() string {
	return fmt.Sprintf("word offset sum between %d and %d overflows",
		err.a, err.b)
}

type errInvalidStructSectionSizes struct {
	dataSectionSize    wordCount16
	pointerSectionSize wordCount16
	sum                WordCount
}

func (err errInvalidStructSectionSizes) Error() string {
	return fmt.Sprintf("struct data (%d) and pointer (%d) sections sum to invalid value %d",
		err.dataSectionSize, err.pointerSectionSize, err.sum)
}

var errShortSingleSegmentStream = errors.New("single segment stream smaller than expected")

var errStreamNotSingleSegment = errors.New("stream is not a single segment stream")

type errShortStreamSegSize struct {
	segSize   ByteCount
	streamLen int
}

func (err errShortStreamSegSize) Error() string {
	return fmt.Sprintf("remaining stream length %d smaller than expected segment size %d",
		err.streamLen, err.segSize)
}

var errAllocNoFirstSeg = errors.New("allocator did not initialize the first segment")

var errAllocNoRootWord = errors.New("allocator did not initialize the root struct pointer")

var errAllocStateNoRootWord = errors.New("builder state does not have space for root pointer")

var errMsgBuilderNoSegments = errors.New("message builder does not have any segments")

var errMsgBuilderNoSegData = errors.New("message builder does not have any segment data")

var errAllocOverMaxWordCount = errors.New("cannot allocate more than the max word count")

var errStructSizeTooLarge = errors.New("struct size is too large")

var errDifferentMsgBuilders = errors.New("different message builders")

var errCannotChangeSegsCap = errors.New("cannot change capacity of segments during Allocate")

var errStringTooLarge = errors.New("string is too large")

type errStructBuilderDoesNotContainDataField uint16

func (err errStructBuilderDoesNotContainDataField) Error() string {
	return fmt.Sprintf("struct builder did not allocate space for field %d", int(err))
}

type errStructBuilderDoesNotContainPointerField uint16

func (err errStructBuilderDoesNotContainPointerField) Error() string {
	return fmt.Sprintf("struct builder did not allocate space for pointer %d", int(err))
}
