// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
)

// appendWords appends little endian encoded words into a slice.
func appendWords(b []byte, words ...Word) []byte {
	for _, w := range words {
		b = binary.LittleEndian.AppendUint64(b, uint64(w))
	}
	return b
}
