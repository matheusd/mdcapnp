// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

func isListPointer(p Word) bool {
	return (p & 0x03) == 1
}

func isStructPointer(p Word) bool {
	return (p & 0x03) == 0
}

type structPointer struct {
	dataOffset         listOrStructOffset
	dataSectionSize    wordCount16
	pointerSectionSize wordCount16
}

func (sp *structPointer) fromWord(w Word) {
	// 0x0000000080000000
	//         0x7ffffffc
	//         0x80000000
	sp.dataOffset = listOrStructOffset(w&0xfffffffc) >> 2
	sp.dataSectionSize = wordCount16(w & 0xffff00000000 >> 32)
	sp.pointerSectionSize = wordCount16(w >> 48)
}

type listPointer struct {
	startOffset listOrStructOffset
	elSize      listElementSize
	listSize    listSize
}

func (lp *listPointer) fromWord(w Word) {
	lp.startOffset = listOrStructOffset(w&0xfffffffc) >> 2
	lp.elSize = listElementSize(w & 0x300000000 >> 32)
	lp.listSize = listSize(w & 0xfffffff800000000 >> 35)
}
