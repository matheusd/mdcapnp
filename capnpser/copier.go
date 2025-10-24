// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"errors"
	"fmt"
)

// copyStructListReferents goes through the items of a struct list and copies
// their referents (the values of any pointers).
//
// This is called after the tag word and the list items have been copied to the
// destination.
func copyStructListReferents(srcList StructList, dst *MessageBuilder, dstList AnyPointerBuilder) error {
	listLen := srcList.Len()

	dstItem := AnyPointerBuilder{
		mb:  dst,
		urb: dstList.urb,
		off: dstList.off + 1, // One past the tag word offset.
		sid: dstList.sid,
		ptr: buildRawStructPointer(dstList.off+1, srcList.itemSize),
	}

	itemSizeWords := srcList.itemSize.TotalSize()
	for i := range listLen {
		srcItem := srcList.At(i)

		if err := copyStructPointerReferents(srcItem, dst, dstItem); err != nil {
			return fmt.Errorf("error copying struct list referent %d: %v", i, err)
		}

		dstItem.off += WordOffset(itemSizeWords)
		dstItem.ptr = dstItem.ptr.withDataOffset(dstItem.off)
	}
	return nil
}

func copyList(src List, dst *MessageBuilder) (AnyPointerBuilder, error) {
	// TODO: read tag word if src.ptr.elSize == listElSizeComposite and
	// double check it is valid?
	/*
		if src.ptr.elSize == listElSizeComposite {
			// TODO: support this.
			return AnyPointerBuilder{}, errors.New("list of structs not supported in copyList()")
		}
	*/

	// Allocate space for list contents.
	totalWords := src.lenWords()
	seg, off, err := dst.allocate(0, totalWords)
	if err != nil {
		return AnyPointerBuilder{}, err
	}

	// Copy the entire list contents.
	seg.copyWordsFrom(&src.seg.b, src.ptr.startOffset, off, totalWords)

	dstList := AnyPointerBuilder{
		mb:  dst,
		urb: seg,
		off: off,
		ptr: buildRawListPointer(off, src.ptr.elSize, src.ptr.listSize),
		sid: seg.id,
	}

	// For struct lists, descend into the individual items to copy their
	// referents.
	if src.ptr.elSize == listElSizeComposite {
		strList, err := src.AsStructList()
		if err != nil {
			return AnyPointerBuilder{}, err
		}
		err = copyStructListReferents(strList, dst, dstList)
		if err != nil {
			return AnyPointerBuilder{}, err
		}
	}

	// All done.
	return dstList, nil
}

// copyStructPointerReferents goes through the list of pointers in the struct,
// copying the referents.
func copyStructPointerReferents(src Struct, dst *MessageBuilder, dstStruct AnyPointerBuilder) error {
	seg, off := dstStruct.urb, dstStruct.off

	// For each pointer, copy the referent and adjust offset.
	srcSubPtrOff := src.ptr.dataOffset + WordOffset(src.ptr.dataSectionSize)
	for i := wordCount16(0); i < src.ptr.pointerSectionSize; i, srcSubPtrOff = i+1, srcSubPtrOff+1 {
		// TODO: add src.ReadField(i) instead of reading pointer???
		ptr, err := src.seg.getWordAsPointer(srcSubPtrOff)
		if err != nil {
			return fmt.Errorf("unable to get value of ptr %d at %d: %v",
				i, srcSubPtrOff, err)
		}

		if ptr.isOtherPointer() || ptr.isZeroStruct() || ptr.isNullPointer() {
			// Empty pointers can be ignored (nothing to de-ref and
			// already copied by the caller, outside the loop).
			continue
		}

		if ptr.isFarPointer() {
			// TODO: support this.
			return errors.New("far pointers not supported in copyStruct()")
		}

		var subDst AnyPointerBuilder

		if ptr.isListPointer() {
			var sub List
			if err := src.ReadList(PointerFieldIndex(i), &sub); err != nil {
				return err
			}

			// Recurse into list.
			subDst, err = copyList(sub, dst)
		} else if ptr.isStructPointer() {
			var sub Struct
			if err := src.ReadStruct(PointerFieldIndex(i), &sub); err != nil {
				return err
			}

			// Recurse into it.
			subDst, err = copyStruct(sub, dst)
		} else {
			// Should not happen if we handled all cases.
			err = errors.New("unknown case in copyStructPointerReferents()")
		}

		// At this point, we recursed into the sub struct/list.
		if err != nil {
			return err
		}

		if subDst.sid != seg.id {
			// TODO: support this.
			return errors.New("point to far segments not supported in copyStructPointerReferents")
		}

		// Determine the new offset to this pointer field in dst.
		dstSubPtrOff := off + WordOffset(src.ptr.dataSectionSize) + WordOffset(i)

		// Modify the data offset of the current pointer to
		// point to the newly allocated child in dest.
		newPtr := ptr.withDataOffset(subDst.off - dstSubPtrOff - 1)

		// Finally, rewrite the pointer.
		seg.SetWord(dstSubPtrOff, Word(newPtr))
	}

	return nil
}

func copyStruct(src Struct, dst *MessageBuilder) (AnyPointerBuilder, error) {
	seg, off, err := dst.allocateValidSize(0, src.structSize())
	if err != nil {
		return AnyPointerBuilder{}, err
	}

	// Copy data and pointers.
	seg.copyWordsFrom(&src.seg.b, src.ptr.dataOffset, off, src.structSize())

	// Build the final object.
	dstStruct := AnyPointerBuilder{
		mb:  dst,
		off: off,
		ptr: buildRawStructPointer(off, src.ptr.structSize()),
		sid: seg.id,
		urb: seg,
	}

	// Copy the child referents.
	if err := copyStructPointerReferents(src, dst, dstStruct); err != nil {
		return AnyPointerBuilder{}, err
	}

	return dstStruct, nil
}

// DeepCopy a source object into a destination builder. This produces a
// partially-canonical object in dst.
//
// src may point to an object inside dst (i.e. src may be an object obtained
// from a reader in dst).
func DeepCopy(src AnyPointer, dst *MessageBuilder) (AnyPointerBuilder, error) {
	switch {
	case src.IsZeroStruct():
		// Nothing to do.
		return ZeroStructAsPointerBuilder(), nil

	case src.IsCapPointer():
		// Cap pointer is a single word pointer that doesn't really
		// point anywhere.
		seg, off, err := dst.allocateValidSize(0, 1)
		if err != nil {
			return AnyPointerBuilder{}, err
		}
		seg.SetWord(off, Word(src.ptr))
		return AnyPointerBuilder{
			mb:  dst,
			urb: seg,
			off: off,
			ptr: src.ptr,
			sid: seg.id,
		}, nil

	case src.IsStruct():
		// TODO: estimate total size?
		return copyStruct(src.AsStruct(), dst)

	case src.IsList():
		return copyList(src.AsList(), dst)

	default:
		return AnyPointerBuilder{}, errors.New("unsupported case in DeepCopy()")
	}
}

// DeepCopyAndSetRoot performs a [DeepCopy] of src to dst and sets dst root as
// the resulting object.
//
// Note: this produces a potentially non-standard message if src points to
// anything except a struct.
func DeepCopyAndSetRoot(src AnyPointer, dst *MessageBuilder) error {
	res, err := DeepCopy(src, dst)
	if err != nil {
		return err
	}

	return dst.NonStdSetRoot(&res)
}
