// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"errors"
	"fmt"
)

type Message struct {
	arena *Arena
	dl    depthLimit
}

func MakeMsg(arena *Arena) Message {
	return Message{arena: arena, dl: defaultDepthLimit}
}

// Arena returns the underlying Arena associated with this message.
func (msg *Message) Arena() *Arena {
	return msg.arena
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

func (msg *Message) readRootPtr() (seg *Segment, ptr pointer, ptrType pointerType, newDL depthLimit, segWordLen WordCount, err error) {
	var ok bool
	newDL, ok = msg.dl.dec()
	if !ok {
		err = errDepthLimitExceeded
		return
	}

	seg, err = msg.arena.Segment(0)
	if err != nil {
		return
	}

	segWordLen = seg.wordLen()
	if segWordLen < 1 {
		err = errNoRootPointer
		return
	}
	ptr = seg.uncheckedGetWordAsPointer(0)

	// De-ref far pointers into the concrete list segment and near pointer.
	ptrType = ptr.pointerType()
	if ptrType == pointerTypeFarPointer {
		seg, ptr, newDL, err = derefFarPointer(msg.arena, newDL, ptr)
		if err != nil {
			return
		}
		ptrType = ptr.pointerType()
		segWordLen = seg.wordLen()
	}
	return
}

func (msg *Message) ReadRoot(s *Struct) error {
	var ok bool

	// TODO: abstract this function to extract struct fields from a struct.
	//
	// This implementation has been extensively reviewed to yield the
	// best performance for extracting the root struct, so when
	// generalizing, be mindful not to downgrade it.
	seg, ptr, ptrType, newDL, segWordLen, err := msg.readRootPtr()
	if err != nil {
		return err
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
		dl:    newDL,
	}
	return nil
}

// GetRoot returns the root struct of the message.
func (msg *Message) GetRoot() (res Struct, err error) {
	err = msg.ReadRoot(&res)
	return
}

// NonStdRootAsAnyPointer returns the root of the message as an [AnyPointer],
// independently of whether it is a struct.
//
// This is a non-standard usage of a message, which is supposed to have a
// struct as root.
//
// Note: This de-refs any indirections in the actual root pointer, to get the
// first concrete object (list, struct or capPointer).
func (msg *Message) NonStdRootAsAnyPointer() (res AnyPointer, err error) {
	var ok bool

	seg, ptr, _, newDL, segWordLen, err := msg.readRootPtr()
	if err != nil {
		return AnyPointer{}, err
	}

	var totalWordSize WordCount
	var dataOffset WordOffset
	switch {
	case ptr.isOtherPointer() || ptr.isNullPointer() || ptr.isZeroStruct():
		// totalWordCount is zero (only pointer).

	case ptr.isFarPointer():
		// Shouldn't happen because readRootPtr de-refs.
		return AnyPointer{}, errors.New("unhandled far pointer case in NonStdRootAsAnyPointer()")

	case ptr.isStructPointer():
		sp := ptr.toStructPointer()
		dataOffset = sp.dataOffset
		totalWordSize = sp.structSize().TotalSize()

	case ptr.isListPointer():
		lp := ptr.toListPointer()
		dataOffset = lp.startOffset
		totalWordSize = listWordCount(lp.elSize, lp.listSize)

	default:
		return AnyPointer{}, errors.New("unhandled case in NonStdRootAsAnyPointer()")
	}

	// The only negative value allowed as an offset is -1, to denote a
	// zero-sized struct.
	if dataOffset < -1 {
		return AnyPointer{}, errInvalidNegativeStructOffset
	}

	// Bounds check that the entire object is readable.
	//
	// Concrete offset is always one more than the encoded start offset.
	if dataOffset, ok = addWordOffsets(dataOffset, 1); !ok {
		return AnyPointer{}, errWordOffsetSumOverflows{dataOffset, 1}
	}

	// The end offset must be valid word offset and must come before the end
	// of the segment.
	if endOffset, ok := addWordOffsets(dataOffset, WordOffset(totalWordSize)); !ok || endOffset > WordOffset(segWordLen) {
		return AnyPointer{}, ErrObjectOutOfBounds{Offset: dataOffset, Size: totalWordSize, WordLen: segWordLen}
	}

	// Safety check that the caller allows us to read this many words.
	if err := msg.arena.ReadLimiter().CanRead(totalWordSize); err != nil {
		return AnyPointer{}, err
	}

	return AnyPointer{
		seg:   seg,
		arena: msg.arena,
		dl:    newDL,
		ptr:   ptr,
	}, nil
}
