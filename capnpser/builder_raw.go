// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"unsafe"
)

// TODO: too unsafe to export??
type unsafeRawBuilder struct {
	ptr unsafe.Pointer
}

func (rb *unsafeRawBuilder) SetWord(offset WordOffset, value Word) {
	*(*Word)(unsafe.Add(rb.ptr, offset*WordSize)) = Word(value)
}

func (rb *unsafeRawBuilder) maskAndMergeWord(offset WordOffset, mask, value Word) {
	ptr := (*Word)(unsafe.Add(rb.ptr, offset*WordSize))
	*ptr = *ptr&^mask | value
}

func (rb *unsafeRawBuilder) SetString(ptrOffset WordOffset, v string, startOffset WordOffset) (nextOffset WordOffset) {
	textLen := uint(len(v) + 1)
	copy(unsafe.Slice((*byte)(unsafe.Add(rb.ptr, startOffset*WordSize)), len(v)), v)
	nextOffset = startOffset + WordOffset(uintBytesToWordAligned(textLen))
	lsPtr := buildRawListPointer(startOffset-ptrOffset-1, listElSizeByte, listSize(textLen))
	*(*Word)(unsafe.Add(rb.ptr, ptrOffset)) = Word(lsPtr)
	return
}

func (rb *unsafeRawBuilder) SetStruct(ptrOff, structOff WordOffset, size StructSize) {
	*(*Word)(unsafe.Add(rb.ptr, ptrOff*WordSize)) = Word(buildRawStructPointer(structOff-ptrOff-1, size))
}

func (rb *unsafeRawBuilder) AliasChild(offset WordOffset, child *unsafeRawBuilder) {
	child.ptr = unsafe.Add(rb.ptr, offset*WordSize)
}

func (rb *unsafeRawBuilder) Child(offset WordOffset) unsafeRawBuilder {
	return unsafeRawBuilder{ptr: unsafe.Add(rb.ptr, offset*WordSize)}
}

type rawBuilder struct {
	b []Word
}

func (rb *rawBuilder) SetWord(offset WordOffset, value Word) {
	rb.b[offset] = value
}

func (rb *rawBuilder) SetString(ptrOffset WordOffset, v string, startOffset WordOffset) (nextOffset WordOffset) {
	textLen := uint(len(v) + 1)
	copy([]byte(unsafe.Slice((*byte)(unsafe.Pointer(&rb.b[startOffset])), len(v))), v)
	nextOffset = startOffset + WordOffset(uintBytesToWordAligned(textLen))
	lsPtr := buildRawListPointer(startOffset-ptrOffset-1, listElSizeByte, listSize(textLen))
	rb.b[ptrOffset] = Word(lsPtr)
	return
}

func (rb *rawBuilder) SetStringXXX(offset WordOffset, sizeBounds WordCount, v string) {
	textLen := Word(len(v) + 1)
	copy([]byte(unsafe.Slice((*byte)(unsafe.Pointer(&rb.b[offset])), len(v))), v)
	rb.b[offset+WordOffset(sizeBounds)-1] |= textLen << 56
}

func (rb *rawBuilder) SetStruct(ptrOff, structOff WordOffset, size StructSize) {
	rb.b[ptrOff] = Word(buildRawStructPointer(structOff-ptrOff-1, size))
}

func (rb *rawBuilder) AliasChild(offset WordOffset, child *rawBuilder) {
	child.b = rb.b[offset:]
}
