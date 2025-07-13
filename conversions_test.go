// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestMy(t *testing.T) {
	var original int64 = 0x1234567890abcdef
	t.Logf("original: %x", original)
	buf := binary.LittleEndian.AppendUint64(make([]byte, 0, 16), uint64(original))
	t.Logf("Buf: %x", buf)

	var res int64
	copy((*[8]byte)(unsafe.Pointer(&res))[:], buf[:8])
	t.Logf("%x res after copy", res)
	res = int64(binary.LittleEndian.Uint64((*[8]byte)(unsafe.Pointer(&res))[:]))
	t.Logf("%x final res", res)
}

func appendWords(b []byte, words ...Word) []byte {
	for _, w := range words {
		b = binary.LittleEndian.AppendUint64(b, uint64(w))
	}
	return b
}
