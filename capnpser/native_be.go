// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

// See native_le.go for an alternative version when the "nativeleperfimprov" tag
// is specified.

//go:build !nativeleperfimprov || mips || mips64 || ppc64

package mdcapnp

import (
	"encoding/binary"
)

// uncheckedGetWord returns the word at the given offset without checking for
// valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWord(offset WordOffset) Word {
	return Word(binary.LittleEndian.Uint64(ms.uncheckedTailSlice(offset)))
}

func (ms *Segment) GetWord(offset WordOffset) (res Word, err error) {
	if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(binary.LittleEndian.Uint64(ms.b[byteOffset:]))
	}
	return
}

// uncheckedGetWord returns the word at the given offset as a pointer without
// checking for valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWordAsPointer(offset WordOffset) pointer {
	return pointer(binary.LittleEndian.Uint64(ms.uncheckedTailSlice(offset)))
}
