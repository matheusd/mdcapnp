// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "fmt"

type Message struct {
	arena Arena
	dl    depthLimit
}

func MakeMsg(arena Arena) Message {
	return Message{arena: arena, dl: defaultDepthLimit}
}

// SetDepthLimit sets the depth limit for reading structs and de-referencing
// pointers from this message. To have an effect, this must be called before
// [ReadRoot].
//
// Note that using a large depth limit may introduce vulnerabilities when
// decoding messages from untrusted inputs.
//
// The maximum depth limit allowable is MaxDepthLimit. Calling this function
// with values greater than that causes panics.
//
// This is NOT safe for concurrent access.
func (msg *Message) SetDepthLimit(dl uint) {
	if dl > MaxDepthLimit {
		panic(fmt.Sprintf("value %d not allowable as a depth limit", dl))
	}

	msg.dl = depthLimit(dl)
}

// RemoveDepthLimit disables the depth limit check for this message. To have an
// effect, this must be called before [ReadRoot].
//
// Note that removing the depth limit check may introduce vulnerabilities when
// decoding messages from untrusted inputs.
//
// This is NOT safe for concurrent access.
func (msg *Message) RemoveDepthLimit() {
	msg.dl = noDepthLimit
}

func (msg *Message) ReadRoot(s *Struct) error {
	structDL, ok := msg.dl.dec()
	if !ok {
		return errDepthLimitExceeded
	}

	seg, err := msg.arena.Segment(0)
	if err != nil {
		return err
	}

	ptr, err := seg.getWordAsPointer(0)
	if err != nil {
		return err
	}

	// De-ref far pointers into the concrete list segment and near pointer.
	if ptr.isFarPointer() {
		seg, ptr, structDL, err = derefFarPointer(s.arena, structDL, ptr)
		if err != nil {
			return err
		}
	}

	// The resulting pointer (after de-ref) MUST be a struct pointer.
	if !ptr.isStructPointer() {
		return ErrNotStructPointer
	}

	sp := ptr.toStructPointer()
	fullSize := WordCount(sp.dataSectionSize) + WordCount(sp.pointerSectionSize)

	// The only negative value allowed as an offset is -1, to denote a
	// zero-sized struct.
	//
	// TODO: double check if this is true.
	if sp.dataOffset < -1 {
		return errInvalidNegativeStructOffset
	}

	// Perform a bounds check when either a data offset or fields were
	// specified (i.e. non-null struct).
	//
	// Note that the gap for the case where (dataOffset == -1) && (fullSize
	// == 0) does not need to be checked because for that case (empty
	// struct), the offset points to the pointer itself (which was already
	// obtained and thus necessarily in bounds). Therefore we elide the
	// redundant check.
	if sp.dataOffset > 0 || fullSize > 0 {
		// TODO: abstract and add base struct offset when obtaining a
		// struct field from a struct.

		// Concrete offset is always one more than the encoded start
		// offset.
		if !addWordOffsets(sp.dataOffset, 1, &sp.dataOffset) {
			return errWordOffsetSumOverflows{sp.dataOffset, 1}
		}
		if !fullSize.Valid() {
			return errInvalidStructSectionSizes{sp.dataSectionSize, sp.pointerSectionSize, fullSize}
		}
		if err := seg.CheckBounds(sp.dataOffset, fullSize); err != nil {
			return err
		}
		if err := msg.arena.ReadLimiter().CanRead(fullSize); err != nil {
			return err
		}
	}

	s.seg = seg
	s.arena = msg.arena
	s.ptr = sp
	s.dl = structDL
	return nil
}
