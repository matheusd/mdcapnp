// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "unsafe"

// bool2byte converts a bool to a Word. This relies on the compiler and runtime
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

func uintBytesToWordAligned(i uint) Word {
	return Word((i + (WordSize - 1)) / WordSize)
}
