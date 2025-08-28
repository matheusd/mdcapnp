// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

// This is a version of functions that forgo the binary.LittleEndian.Uint64
// function in favor of direct memory copy in architectures that are little
// endian (most of them).
//
// My (matheusd) original hope was that this would be faster than using the
// little endian function, and some investigations on the generated ASM and
// microbenchmarks would seem to indicate that this was true. But when running
// larger benchmarks, this version ends up slightly slower.
//
// I've opted to leave this code here for the moment, for future investigations.

//go:build nativeleperfimprov && !(mips || mips64 || ppc64)

package mdcapnp

import "unsafe"

// uncheckedGetWord returns the word at the given offset without checking for
// valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWord(offset WordOffset) Word {
	return Word(*(*uint64)(unsafe.Pointer(&ms.b[offset*WordSize])))
}

func (ms *Segment) GetWord(offset WordOffset) (res Word, err error) {
	if byteOffset := offset * WordSize; len(ms.b) < int(byteOffset+WordSize) {
		err = ErrInvalidMemOffset{AvailableLen: len(ms.b), Offset: int(byteOffset)}
	} else {
		res = Word(*(*uint64)(unsafe.Pointer(&ms.b[byteOffset])))
	}
	return
}

// uncheckedGetWord returns the word at the given offset as a pointer without
// checking for valid bounds.
//
// The assumption is that this method is only called in instances where the
// offset has already been determined to exist.
func (ms *Segment) uncheckedGetWordAsPointer(offset WordOffset) pointer {
	return pointer(*(*uint64)(unsafe.Pointer(&ms.b[offset*WordSize])))
}
