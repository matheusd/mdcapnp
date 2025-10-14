// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"unsafe"
)

type UnsafeRawBuilder struct {
	ptr unsafe.Pointer
}

func (rb *UnsafeRawBuilder) SetWord(offset WordOffset, value Word) {
	*(*Word)(unsafe.Add(rb.ptr, offset*WordSize)) = Word(value)
}

func (rb *UnsafeRawBuilder) maskAndMergeWord(offset WordOffset, mask, value Word) {
	ptr := (*Word)(unsafe.Add(rb.ptr, offset*WordSize))
	*ptr = *ptr&^mask | value
}

func (rb *UnsafeRawBuilder) SetString(ptrOffset WordOffset, v string, startOffset WordOffset) (nextOffset WordOffset) {
	textLen := uint(len(v) + 1)
	copy(unsafe.Slice((*byte)(unsafe.Add(rb.ptr, startOffset*WordSize)), len(v)), v)
	nextOffset = startOffset + WordOffset(uintBytesToWordAligned(textLen))
	lsPtr := buildRawListPointer(startOffset-ptrOffset-1, listElSizeByte, listSize(textLen))
	*(*Word)(unsafe.Add(rb.ptr, ptrOffset)) = Word(lsPtr)
	return
}

func (rb *UnsafeRawBuilder) SetStruct(ptrOff, structOff WordOffset, size StructSize) {
	*(*Word)(unsafe.Add(rb.ptr, ptrOff*WordSize)) = Word(buildRawStructPointer(structOff-ptrOff-1, size))
}

func (rb *UnsafeRawBuilder) AliasChild(offset WordOffset, child *UnsafeRawBuilder) {
	child.ptr = unsafe.Add(rb.ptr, offset*WordSize)
}

type RawBuilder struct {
	b []Word
}

func (rb *RawBuilder) SetWord(offset WordOffset, value Word) {
	rb.b[offset] = value
}

func (rb *RawBuilder) SetString(ptrOffset WordOffset, v string, startOffset WordOffset) (nextOffset WordOffset) {
	textLen := uint(len(v) + 1)
	copy([]byte(unsafe.Slice((*byte)(unsafe.Pointer(&rb.b[startOffset])), len(rb.b)*WordSize)), v)
	nextOffset = startOffset + WordOffset(uintBytesToWordAligned(textLen))
	lsPtr := buildRawListPointer(startOffset-ptrOffset-1, listElSizeByte, listSize(textLen))
	rb.b[ptrOffset] = Word(lsPtr)
	return
}

func (rb *RawBuilder) SetStringXXX(offset WordOffset, sizeBounds WordCount, v string) {
	textLen := Word(len(v) + 1)
	copy([]byte(unsafe.Slice((*byte)(unsafe.Pointer(&rb.b[offset])), len(rb.b)*WordSize)), v)
	rb.b[offset+WordOffset(sizeBounds)-1] |= textLen << 56
}

func (rb *RawBuilder) SetStruct(ptrOff, structOff WordOffset, size StructSize) {
	rb.b[ptrOff] = Word(buildRawStructPointer(structOff-ptrOff-1, size))
}

func (rb *RawBuilder) AliasChild(offset WordOffset, child *RawBuilder) {
	child.b = rb.b[offset:]
}
