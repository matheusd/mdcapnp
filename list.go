// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "unsafe"

type listElementSize byte

const (
	listElSizeVoid      listElementSize = 0
	listElSizeBit       listElementSize = 1
	listElSizeByte      listElementSize = 2
	listElSizeComposite listElementSize = 7
)

type listSize uint32

func listWordCount(elSize listElementSize, lsSize listSize) WordCount {
	switch elSize {
	case listElSizeVoid:
		return 0
	case listElSizeBit:
		return WordCount(lsSize) // FIXME calc and align to word
	case listElSizeByte:
		return WordCount(lsSize / WordSize) // FIXME: align to word
	case listElSizeComposite:
		return WordCount(lsSize) + 1 // +1 because of tag word
	default:
		panic("unknown el size")
	}

}

type List struct {
	seg        *Segment
	baseOffset WordOffset
	elSize     listElementSize
	listSize   listSize
}

/*
func (ls *List) fromPointerWord(pointerOffset WordOffset, w Word) {
	ls.baseOffset = pointerOffset + (WordOffset(w) & 0xfffffffc >> 2) + 1
	ls.elSize = ListElementSize(w & 0x300000000 >> 32)
	ls.listSize = ListSize(w & 0xfffffff800000000 >> 35)
}
*/

func (ls *List) LenBytes() ByteCount {
	return ByteCount(listWordCount(ls.elSize, ls.listSize)) * WordSize
}

// Read this list into a slice. Only valid for one-byte-per-element lists.
func (ls *List) Read(b []byte) (n int, err error) {
	if ls.elSize != listElSizeByte {
		return 0, errNotOneByteElList
	}
	n = min(len(b), int(ls.listSize)) // FIXME: check if conversion valid in 32bit archs
	return ls.seg.Read(ls.baseOffset, b[:n])
}

func (ls *List) String() string {
	buf := ls.seg.uncheckedSlice(ls.baseOffset, ls.LenBytes())
	return string(buf)
}

// CheckCanGetUnsafeString returns nil if a subsequent call to [UnsafeString]
// will work correctly. This is only valid for as long as the underlying arena
// is not modified or invalidated.
//
// TODO: is this really needed? Struct.ReadList() already checks for list
// validity.
func (ls *List) CheckCanGetUnsafeString() error {
	if ls.elSize != listElSizeByte {
		return errNotOneByteElList
	}
	if err := ls.seg.checkSliceBounds(ls.baseOffset, ls.LenBytes()); err != nil {
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
	buf := ls.seg.uncheckedSlice(ls.baseOffset, ls.LenBytes())
	return *(*string)(unsafe.Pointer(&buf))
}

func (ls *List) Foo() {
	v := ls.seg
	if v == nil {
		panic("boo")
	}
	v.Read(0, nil)
}
