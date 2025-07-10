// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

type Message struct {
	arena Arena
}

func MakeMsg(arena Arena) Message {
	return Message{arena: arena}
}

func (msg *Message) ReadRoot(s *Struct) error {
	seg, err := msg.arena.Segment(0)
	if err != nil {
		return err
	}

	ptr, err := seg.GetWord(0)
	if err != nil {
		return err
	}
	if !isStructPointer(ptr) {
		return ErrNotStructPointer
	}

	// TODO: check null pointer? zero sized struct?

	var sp structPointer
	sp.fromWord(ptr)

	dataStartOffset := WordOffset(0 + sp.dataOffset + 1) // 0 == root pointer offset
	fullSize := WordCount(sp.dataSectionSize) + WordCount(sp.pointerSectionSize)
	if err := seg.CheckBounds(dataStartOffset, fullSize); err != nil {
		return err
	}

	s.seg = seg
	s.dataStartOffset = dataStartOffset
	s.dataSize = WordCount(sp.dataSectionSize)
	s.pointerSize = WordCount(sp.pointerSectionSize)

	return nil
}
