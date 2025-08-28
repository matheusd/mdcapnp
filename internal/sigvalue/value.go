// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package sigvalue

import (
	"context"
	"sync"
)

// Value is a value that emits a signal every time it changes. It supports
// concurrent setters and multiple waiters.
type Value[T any] struct {
	mu        sync.Mutex
	v         T
	waitChans []chan T
}

// NewValue returns a new [Value] initialized to v.
func NewValue[T any](v T) *Value[T] {
	return &Value[T]{v: v}
}

// Set the value.
func (v *Value[T]) Set(newV T) {
	v.mu.Lock()
	v.v = newV
	chans := v.waitChans
	v.waitChans = nil
	v.mu.Unlock()

	for _, c := range chans {
		if c != nil {
			c <- newV
		}
	}
}

// Get returns the current value.
func (v *Value[T]) Get() T {
	v.mu.Lock()
	res := v.v
	v.mu.Unlock()
	return res
}

// Wait for a change in the value or for the context to be done.
func (v *Value[T]) Wait(ctx context.Context) (res T, err error) {
	if ctx.Err() != nil { // Quick sanity check.
		err = ctx.Err()
		return
	}

	// Add a new channel to the list of waiting channels.
	c := make(chan T, 1)
	v.mu.Lock()
	v.waitChans = append(v.waitChans, c)
	i := len(v.waitChans) - 1
	v.mu.Unlock()

	select {
	case res = <-c:
	case <-ctx.Done():
		err = context.Cause(ctx)

		// Remove c from list of waiting channels. Note that c may have
		// been written between <-ctx.Done and this (and thus removed
		// from the list of waiting channels).
		v.mu.Lock()
		if i < len(v.waitChans) && v.waitChans[i] == c {
			v.waitChans[i] = nil
		}
		v.mu.Unlock()
	}

	return
}
