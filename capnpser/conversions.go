// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import "unsafe"

// bool2Word converts a bool to a Word. This relies on the compiler and runtime
// representing bools as bytes with value 1.
//
// This is generally unsafe.
//
// Credits: https://github.com/hasansino/gobasics/blob/master/gonuts/unsafepkg/unsafepkg.go
func boolToWord(b bool) Word {
	// (unsafe.Pointer(&b) - create new unsafe pointer with address of b
	// (*int8) - convert to pointer of type int8
	// * - take value of last pointer (int8)
	return Word(*(*byte)(unsafe.Pointer(&b)))
}

func BoolToWord(b bool) Word {
	return Word(*(*byte)(unsafe.Pointer(&b)))
}

func uintBytesToWordAligned(i uint) Word {
	return Word((i + (WordSize - 1)) / WordSize)
}

// Uint64BytesToWordCount converts a uint64 that represents a byte count into a
// word count.
//
// Returns true if the resulting word count is a valid word count (lower than
// the max allowed).
func Uint64BytesToWordCount(u uint64) (wc WordCount, valid bool) {
	wc = WordCount((u + (WordSize - 1)) / WordSize)
	valid = u < maxValidBytes
	return
}
