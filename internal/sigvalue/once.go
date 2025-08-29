// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package sigvalue

import (
	"context"
	"errors"
	"sync/atomic"
)

var errOnceAlreadySet = errors.New("value can only be set once")

// Once is a value that can only be set once and signals when it is set.
//
// After being set, it retains the same value forever.
type Once[T any] struct {
	isSet     atomic.Bool
	isSetChan chan struct{}
	v         T
}

// MaybeSet sets the value to v if it hasn't been set yet. Returns true if the
// value was modified.
//
// It is usually an error to use this function. Prefer using [Set].
func (o *Once[T]) MaybeSet(v T) bool {
	if !o.isSet.CompareAndSwap(false, true) {
		return false
	}

	o.v = v
	close(o.isSetChan)
	return true
}

// Set sets the value to v. It panics if the value has already been set.
func (o *Once[T]) Set(v T) {
	if !o.MaybeSet(v) {
		panic(errOnceAlreadySet)
	}
}

// IsSet returns whether the value has already been set.
//
// This function is prone to misuse. Prefer using [Wait].
func (o *Once[T]) IsSet() bool {
	return o.isSet.Load()
}

// Wait waits until the value is set or the context expires.
func (o *Once[T]) Wait(ctx context.Context) (T, error) {
	select {
	case <-o.isSetChan:
		return o.v, nil
	case <-ctx.Done():
		var empty T
		return empty, context.Cause(ctx)
	}
}

// OnceSetter is the setter side of a [Once] value. It can only set the value,
// not read it.
type OnceSetter[T any] struct {
	o *Once[T]
}

// MaybeSet sets the value to v if it hasn't been set yet. Returns true if the
// value was modified.
//
// It is usually an error to use this function. Prefer using [Set].
func (os OnceSetter[T]) MaybeSet(v T) bool {
	return os.o.MaybeSet(v)
}

// Set sets the value to v. It panics if the value has already been set.
func (os OnceSetter[T]) Set(v T) {
	os.o.Set(v)
}

// OnceGetter is the getter side of a [Once] value. It can only read the value,
// not set it.
type OnceGetter[T any] struct {
	o *Once[T]
}

// IsSet returns whether the value has already been set.
//
// This function is prone to misuse. Prefer using [Wait].
func (og OnceGetter[T]) IsSet() bool {
	return og.o.IsSet()
}

// Wait waits until the value is set or the context expires.
func (og OnceGetter[T]) Wait(ctx context.Context) (T, error) {
	return og.o.Wait(ctx)
}
