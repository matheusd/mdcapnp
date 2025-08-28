// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import "fmt"

type Message struct {
	arena *Arena
	dl    depthLimit
}

func MakeMsg(arena *Arena) Message {
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
	// TODO: abstract this function to extract struct fields from a struct.
	//
	// This implementation has been extensively reviewed to yield the
	// best performance for extracting the root struct, so when
	// generalizing, be mindful not to downgrade it.

	structDL, ok := msg.dl.dec()
	if !ok {
		return errDepthLimitExceeded
	}

	seg, err := msg.arena.Segment(0)
	if err != nil {
		return err
	}

	segWordLen := seg.wordLen()
	if segWordLen < 1 {
		return errNoRootPointer
	}
	ptr := seg.uncheckedGetWordAsPointer(0)

	// De-ref far pointers into the concrete list segment and near pointer.
	ptrType := ptr.pointerType()
	if ptrType == pointerTypeFarPointer {
		seg, ptr, structDL, err = derefFarPointer(s.arena, structDL, ptr)
		if err != nil {
			return err
		}
		ptrType = ptr.pointerType()
		segWordLen = seg.wordLen()
	}

	// The resulting pointer (after de-ref) MUST be a struct pointer.
	if ptrType != pointerTypeStruct {
		return ErrNotStructPointer
	}

	// All fields in sp are necessarily valid, because we're decoding it
	// from a generic pointer word.
	sp := ptr.toStructPointer()

	// The only negative value allowed as an offset is -1, to denote a
	// zero-sized struct.
	//
	// For the root pointer, this is also the only valid negative offset,
	// because any other would put the struct contents outside the bounds of
	// the first segment.
	//
	// TODO: double check if this is true for other structs.
	if sp.dataOffset < -1 {
		return errInvalidNegativeStructOffset
	}

	// Calculate full size of struct in words. Doesn't need an overflow
	// check because the section sizes are both uint16, so their sum can't
	// overflow uint32.
	//
	// The resulting size is also necessarily a valid word count, because it
	// can only be up to 2^17-2 which is < 2^29-1.
	fullSize := WordCount(sp.dataSectionSize) + WordCount(sp.pointerSectionSize)

	// Bounds check that the entire struct (all data fields + all pointers)
	// is readable.
	//
	// Concrete offset is always one more than the encoded start offset.
	if sp.dataOffset, ok = addWordOffsets(sp.dataOffset, 1); !ok {
		return errWordOffsetSumOverflows{sp.dataOffset, 1}
	}

	// The end offset must be valid word offset and must come before the end
	// of the segment.
	if endOffset, ok := addWordOffsets(sp.dataOffset, WordOffset(fullSize)); !ok || endOffset > WordOffset(segWordLen) {
		return ErrObjectOutOfBounds{Offset: sp.dataOffset, Size: fullSize, WordLen: segWordLen}
	}

	// Safety check that the caller allows us to read this many words.
	if err := msg.arena.ReadLimiter().CanRead(fullSize); err != nil {
		return err
	}

	*s = Struct{
		seg:   seg,
		arena: msg.arena,
		ptr:   sp,
		dl:    structDL,
	}
	return nil
}
