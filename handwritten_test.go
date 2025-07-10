// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

type SmallTestStruct Struct

func (st *SmallTestStruct) Siblings() int64 {
	return (*Struct)(st).Int64(0)
}

func (st *SmallTestStruct) ReadNameField(ls *List) error {
	return (*Struct)(st).ReadList(0, ls) // First pointer.
}
