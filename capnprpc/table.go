// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "math"

type table[T ~uint32, U any] struct {
	lastID  T
	entries map[T]U
}

func makeTable[T ~uint32, U any]() table[T, U] {
	return table[T, U]{
		entries: make(map[T]U),
	}
}

// nextID returns the next free id. Does NOT track it as used.
func (t *table[T, U]) nextID() (id T, ok bool) {
	// TODO: track and reuse low numbered IDs.

	if t.lastID == math.MaxUint32 {
		// TODO: handle overflows.
		ok = false
		return
	}

	t.lastID++
	id = t.lastID
	ok = true
	return
}

func (t *table[T, U]) set(id T, v U) {
	t.entries[id] = v
}

func (t *table[T, U]) get(id T) (res U, ok bool) {
	res, ok = t.entries[id]
	return
}

func (t *table[T, U]) has(id T) bool {
	_, ok := t.entries[id]
	return ok
}
