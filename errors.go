// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "fmt"

type ErrInvalidSegmentOffset uint64

func (err ErrInvalidSegmentOffset) Error() string {
	return fmt.Sprintf("not a valid segment offset: %d", uint64(err))
}

type ErrInvalidOffset struct {
	EndOffset    int
	AvailableLen int
}

func (err ErrInvalidOffset) Error() string {
	return fmt.Sprintf("invalid offset: wanted to read up to offset %d when only %d bytes were available",
		err.EndOffset, err.AvailableLen)
}
