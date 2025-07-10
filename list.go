// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

type ListElementSize byte

const (
	ListElSizeVoid      ListElementSize = 0
	ListElSizeBit       ListElementSize = 1
	ListElSizeByte      ListElementSize = 2
	ListElSizeComposite ListElementSize = 7
)

type ListSize uint64

type List struct {
	seg        Segment
	baseOffset WordOffset
	elSize     ListElementSize
	listSize   ListSize
}

func (ls *List) fromPointerWord(pointerOffset WordOffset, w Word) {
	ls.baseOffset = pointerOffset + (WordOffset(w) & 0xfffffffc >> 2) + 1
	ls.elSize = ListElementSize(w & 0x300000000 >> 32)
	ls.listSize = ListSize(w & 0xfffffff800000000 >> 35)
}

func (ls *List) LenBytes() int {
	switch ls.elSize {
	case ListElSizeVoid:
		return 0
	case ListElSizeBit:
		return int(ls.listSize)
	case ListElSizeByte:
		return int(ls.listSize)
	case ListElSizeComposite:
		return int(ls.listSize * WordSize)
	default:
		panic("unknown el size")
	}
}

func (ls *List) Read(b []byte) (n int, err error) {
	n = min(len(b), ls.LenBytes())
	return ls.seg.Read(ls.baseOffset, b[:n])
}

func (ls *List) Foo() {
	v := ls.seg
	if v == nil {
		panic("boo")
	}
	v.Read(0, nil)
}
