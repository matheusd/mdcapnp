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

var errNotOneByteElList = errors.New("list is not a one-byte-per-element list")
