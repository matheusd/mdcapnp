// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func BenchmarkStructGetInt64(b *testing.B) {
	arena := &SingleSegmentArena{b: binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)}
	msg := &Message{arena: arena}
	st := &SmallTestStruct{msg: msg, seg: &MemSegment{b: arena.b}}
	b.ResetTimer()
	b.ReportAllocs()
	var v int64
	for range b.N {
		v = st.Siblings()
	}

	if v == 666 {
		panic("boo")
	}
}

func BenchmarkGenStructGetInt64(b *testing.B) {
	arena := &SingleSegmentArena{b: binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)}
	msg := &GenMessage[*SingleSegmentArena]{arena: arena}
	st := &GenSmallTestStruct[*SingleSegmentArena]{msg: msg}
	b.ResetTimer()
	b.ReportAllocs()
	var v int64
	for range b.N {
		v = st.Siblings()
	}

	if v == 666 {
		panic("boo")
	}
}

func BenchmarkConcreteGetInt64(b *testing.B) {
	var arena ReaderArena = &SingleSegmentArena{b: binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)}

	b.ResetTimer()
	b.ReportAllocs()

	// Ensure seg and offset are not compilte-time constants.
	var seg SegmentID
	var offset Word
	if rand.Int31n(100) == 200 {
		seg = 1
		offset = 2
	}

	var v int64
	for range b.N {
		data, _ := arena.GetWord(seg, offset)
		v = int64(data)
	}

	if v == 666 {
		panic("boo")
	}
}

func BenchmarkMemGetInt64(b *testing.B) {
	arena := &SingleSegmentArena{b: binary.LittleEndian.AppendUint64(nil, 0x1234567890abcdef)}
	msg := &MemMessage{arena: arena}
	st := &MemSmallTestStruct{msg: msg, arena: arena, seg: &MemSegment{b: arena.b}}
	b.ResetTimer()
	b.ReportAllocs()
	var v int64
	for range b.N {
		v = st.Siblings()
	}

	if v == 666 {
		panic("boo")
	}

}
