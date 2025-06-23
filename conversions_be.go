// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

//go:build native_be

package mdcapnp

import (
	"encoding/binary"
	"unsafe"
)

func wireLEToNativeEndian(u uint64) uint64 {
	return binary.LittleEndian.Uint64((*[8]byte)(unsafe.Pointer(&res))[:])
}

func wireLEToNativeEndianInt64(v int64) int64 {
	return int64(binary.LittleEndian.Uint64((*[8]byte)(unsafe.Pointer(&v))[:]))
}

func wireLEBytesToNativeEndianInt64(v *[8]byte) int64 {
	return int64(binary.LittleEndian.Uint64(v))
}
