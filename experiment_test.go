// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

func BenchmarkReadList(b *testing.B) {
	targetName := []byte("mynameisslimshady   ")
	buf := appendWords(nil,
		0x0001000000000000,
		// 0x0000000200000001,
		0x000000ba00000001,
	)
	buf = append(buf, targetName...)

	st := &SmallTestStruct{seg: &Segment{b: buf}, dataStartOffset: 1, pointerSize: 1}

	ls := new(List)

	nameBuf := make([]byte, 32)

	b.ResetTimer()
	b.ReportAllocs()
	var n int
	for range b.N {
		err := st.ReadNameField(ls)
		if err != nil {
			ls = nil
			b.Fatal(err)
		}

		n, err = ls.Read(nameBuf)
		if err != nil {
			b.Fatal(err)
		}
	}

	require.NotNil(b, ls)
	require.Equal(b, targetName, nameBuf[:n])
}
