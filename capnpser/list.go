// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"errors"
	"fmt"
	"unsafe"
)

type listElementSize int

const (
	listElSizeVoid      listElementSize = 0
	listElSizeBit       listElementSize = 1
	listElSizeByte      listElementSize = 2
	listElSizeComposite listElementSize = 7
)

type listSize uint32

// listWordCount is the total number of words for storage of a  list of the
// given type and element size (as read from a pointer).
func listWordCount(elSize listElementSize, lsSize listSize) WordCount {
	switch elSize {
	case listElSizeVoid:
		return 0
	case listElSizeBit:
		return WordCount((lsSize + 63) / 64)
	case listElSizeByte:
		return WordCount((lsSize + 7) / 8)
	case listElSizeComposite:
		return WordCount(lsSize + 1) // + tag word
	// FIXME: add missing types

	default:
		panic("unknown el size")
	}
}

const (
	// MaxListSize is the max number of elements in a list.
	MaxListSize = (1 << 29) - 1
)

func listByteCount(elSize listElementSize, lsSize listSize) ByteCount {
	return ByteCount(listWordCount(elSize, lsSize)) * WordSize
}

type List struct {
	seg   *Segment
	arena *Arena
	ptr   listPointer
	dl    depthLimit
}

func (ls *List) lenWords() WordCount {
	return listWordCount(ls.ptr.elSize, ls.ptr.listSize)
}

func (ls *List) LenBytes() ByteCount {
	return listByteCount(ls.ptr.elSize, ls.ptr.listSize)
}

func (ls *List) AsAnyPointer() AnyPointer {
	return AnyPointer{
		seg:   ls.seg,
		arena: ls.arena,
		dl:    ls.dl,
		ptr:   ls.ptr.toPointer(),
	}
}

func (ls *List) AsStructList() (StructList, error) {
	if ls.ptr.elSize != listElSizeComposite {
		return StructList{}, errors.New("list is not a struct list")
	}

	if ls.ptr.toPointer().isNullPointer() {
		// Empty struct list.
		return StructList{
			l: *ls,
		}, nil
	}

	// Get the tag word and convert it into a struct pointer.
	tagWord, err := ls.seg.getWordAsPointer(ls.ptr.startOffset)
	if err != nil {
		return StructList{}, fmt.Errorf("unable to read tag word: %v", err)
	}
	if !tagWord.isStructPointer() {
		return StructList{}, fmt.Errorf("tag word is not a struct pointer")
	}
	itemsStruct := tagWord.toStructPointer()

	// Check if the total items size as encoded in the tag word is the
	// expected one given the list word count.
	totalItemsSize, ok := mulWordCounts(WordCount(itemsStruct.dataOffset), itemsStruct.structSize().TotalSize())
	if !ok {
		return StructList{}, errors.New("item count * item size overflows valid max word count")
	}
	if totalItemsSize != WordCount(ls.ptr.listSize) {
		return StructList{}, fmt.Errorf("incongruent word counts when converting list to StructList (%d vs %d)",
			ls.ptr.listSize, totalItemsSize)
	}

	// All good.
	return StructList{
		l:        *ls,
		itemSize: itemsStruct.structSize(),
		listLen:  listSize(itemsStruct.dataOffset),
	}, nil
}

// Read this list into a slice. Only valid for one-byte-per-element lists.
func (ls *List) Read(b []byte) (n int, err error) {
	if ls.ptr.elSize != listElSizeByte {
		return 0, errNotOneByteElList
	}
	n = min(len(b), int(ls.ptr.listSize)) // FIXME: check if conversion valid in 32bit archs
	return ls.seg.Read(ls.ptr.startOffset, b[:n])
}

func (ls *List) String() string {
	if ls.ptr.listSize == 0 {
		return ""
	}
	buf := ls.seg.uncheckedSlice(ls.ptr.startOffset, ByteCount(ls.ptr.listSize-1)) // -1 to skip final null
	return string(buf)
}

// Bytes returns the list as bytes. Prefer using [Read] to avoid allocating the
// byte slice.
func (ls *List) Bytes() ([]byte, error) {
	if ls.ptr.elSize != listElSizeByte {
		return nil, errNotOneByteElList
	}
	res := make([]byte, ls.LenBytes())
	n, err := ls.seg.Read(ls.ptr.startOffset, res)
	if err != nil {
		return nil, err
	}
	return res[:n], nil
}

// CheckCanGetUnsafeString returns nil if a subsequent call to [UnsafeString]
// will work correctly. This is only valid for as long as the underlying arena
// is not modified or invalidated.
//
// TODO: is this really needed? Struct.ReadList() already checks for list
// validity.
func (ls *List) CheckCanGetUnsafeString() error {
	if ls.ptr.elSize != listElSizeByte {
		return errNotOneByteElList
	}
	if err := ls.seg.checkSliceBounds(ls.ptr.startOffset, ByteCount(ls.ptr.listSize)); err != nil {
		return err
	}
	return nil
}

// UnsafeString returns this list as an unsafe string. Before calling this
// function, [CheckCanGetUnsafeString] should be called to ensure this list can
// be converted to a string.
//
// The returned string is only safe for use while the underlying arena is valid.
// If the arena is modified, attempting to use the string may result in
// undefined behavior.
func (ls *List) UnsafeString() string {
	if ls.ptr.listSize == 0 {
		return ""
	}
	buf := ls.seg.uncheckedSlice(ls.ptr.startOffset, ByteCount(ls.ptr.listSize-1)) // -1 to skip final null
	return *(*string)(unsafe.Pointer(&buf))
}

type StructList struct {
	l        List
	itemSize StructSize
	listLen  listSize
}

func (sl *StructList) AsAnyPointer() AnyPointer {
	return sl.l.AsAnyPointer()
}

// Len returns the number of elements in this list.
func (sl *StructList) Len() int {
	return int(sl.listLen)
}

// At returns the i'th element of the list. Panics if the item is out of bounds.
func (sl *StructList) At(i int) Struct {
	if i < 0 {
		panic("i is out of bounds (< 0)")
	}
	if i > int(sl.listLen) {
		panic("i is out of bounds (> len)")
	}

	ptr := structPointer{
		dataOffset:         sl.l.ptr.startOffset + WordOffset(i)*WordOffset(sl.itemSize.TotalSize()) + 1, // +1 tag word
		dataSectionSize:    sl.itemSize.DataSectionSize,
		pointerSectionSize: sl.itemSize.PointerSectionSize,
	}

	return Struct{
		seg:   sl.l.seg,
		arena: sl.l.arena,
		dl:    sl.l.dl,
		ptr:   ptr,
	}
}

func ReadGenericStructList[T ~StructType](s *Struct, ptrIndex PointerFieldIndex) (GenericStructList[T], error) {
	var sl StructList
	err := s.ReadStructList(ptrIndex, &sl)
	return GenericStructList[T]{sl}, err
}

type GenericStructList[T ~StructType] struct {
	sl StructList
}

func (gsl *GenericStructList[T]) Len() int   { return gsl.sl.Len() }
func (gsl *GenericStructList[T]) At(i int) T { return T(gsl.sl.At(i)) }
