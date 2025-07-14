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

	if !AddWordOffsets(sp.dataOffset, 1, &sp.dataOffset) {
		return errWordOffsetSumOverflows{sp.dataOffset, 1}
	}

	fullSize := WordCount(sp.dataSectionSize) + WordCount(sp.pointerSectionSize)
	if !fullSize.Valid() {
		return errInvalidStructSectionSizes{sp.dataSectionSize, sp.pointerSectionSize, fullSize}
	}
	if err := seg.CheckBounds(sp.dataOffset, fullSize); err != nil {
		return err
	}
	if err := msg.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return err
	}

	s.seg = seg
	s.arena = msg.arena
	s.ptr = sp
	return nil
}
