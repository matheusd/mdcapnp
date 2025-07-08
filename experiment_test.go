// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"bytes"
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

func BenchmarkReadList(b *testing.B) {
	targetName := []byte("mynameisslimshady   ")
	buf := appendWords(nil,
		0x0001000000000000,
		// 0x0000000200000001,
		0x000000ba00000001,
	)
	buf = append(buf, targetName...)

	arena := &SingleSegmentArena{b: buf}
	msg := &Message{arena: arena}
	st := &SmallTestStruct{msg: msg, seg: &MemSegment{b: arena.b}, dataStartOffset: 1, pointerSize: 1}

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

	if ls == nil {
		panic("boo")
	}
	if !bytes.Equal(targetName, nameBuf[:n]) {
		panic("wrong targetName")
	}
}
