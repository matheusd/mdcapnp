// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"
)

func BenchmarkMsgGetRoot(b *testing.B) {
	buf := appendWords(nil, 0x00000000fffffffc)
	arena := MakeSingleSegmentMemArena(buf, false)
	msg := MakeMsg(&arena)

	var st Struct

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		err := msg.ReadRoot(&st)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Ensure st is not eliminated.
	if st.dataSize != 0 {
		panic("error")
	}
}
