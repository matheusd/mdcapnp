// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package sigvalue

import (
	"cmp"
	"context"
	"sync"
)

// StateType is the type of value that can be used as a state of a [Stateful]
// instance.
//
// It must be an ordered type to allow the [Stateful.AtLeast] and similar
// functions to work.
//
// No semantics are attached to the state itself: that is left to client code to
// define.
type StateType = cmp.Ordered

type stateTargetComparator[S StateType] func(state, target S) bool

func stateAtLeast[S StateType](state, target S) bool {
	return state >= target
}

func stateNotGreater[S StateType](state, target S) bool {
	return state <= target
}

func stateExactly[S StateType](state, target S) bool {
	return state == target
}

type stateChangeEvent[S StateType, V any] struct {
	state S
	value V
}

type stateChangeWaiterType int

const (
	waitExactly stateChangeWaiterType = iota
	waitAtLeast
	waitNotGreater
)

type stateChangeWaiter[S StateType] struct {
	c  chan struct{} // chan stateChangeEvent[S, V]
	wt stateChangeWaiterType

	targetState S
}

// Stateful is a signalled value that has a "current state" and value bound to
// that state. It allows client code to wait for and react to specific state
// changes, effected by other goroutines.
//
// Values must not be copied after initialization and are safe for use from
// multiple goroutines.
//
// An empty stateful value of a given type is ready for use, with "current
// state" and "current value" set to their corresponding empty values.
type Stateful[S StateType, V any] struct {
	mu sync.Mutex

	state   S
	value   V
	waiters []stateChangeWaiter[S]
}

// NewStateful creates a new Stateful value initialized to the passed initial
// state and value.
func NewStateful[S StateType, V any](initialState S, initialValue V) *Stateful[S, V] {
	return &Stateful[S, V]{state: initialState, value: initialValue}
}

func (s *Stateful[S, V]) matchesType(state, target S, wt stateChangeWaiterType) bool {
	switch wt {
	case waitExactly:
		return stateExactly(state, target)
	case waitAtLeast:
		return stateAtLeast(state, target)
	case waitNotGreater:
		return stateNotGreater(state, target)
	default:
		panic("missing case")
	}
}

// checkAndCallWaiters checks which waiters are fulfilled by the current state
// and calls them.
func (s *Stateful[S, V]) checkAndCallWaiters() {
	if len(s.waiters) == 0 {
		return
	}

	// Determine which waiters (if any) to signal and remove them.
	for i := 0; i < len(s.waiters); {
		w := &s.waiters[i]
		if !s.matchesType(s.state, w.targetState, w.wt) {
			i++
			continue
		}

		// Signal the channel. Note this is done under the s.mu.Lock(),
		// therefore w.c MUST be buffered to avoid deadlocking.
		// w.c <- event
		close(w.c)
		s.removeWaiterWithChan(w.c)
	}
}

// Set sets the state and value of the current state. Any waiters that match
// this new state will be alerted to it.
func (s *Stateful[S, V]) Set(newState S, newValue V) (oldState S, oldValue V) {
	s.mu.Lock()
	oldState, oldValue = s.state, s.value
	s.state, s.value = newState, newValue
	s.checkAndCallWaiters()
	s.mu.Unlock()
	return
}

// Modify modifies the current state according to the given function. Any errors
// returned by f are bubbled up to the caller.
//
// NOTE: this is called with the Stateful locked, so f should not block for
// long.
func (s *Stateful[S, V]) Modify(f func(oldState S, oldValue V) (newState S, newValue V, err error)) (err error) {
	s.mu.Lock()
	s.state, s.value, err = f(s.state, s.value)
	s.checkAndCallWaiters()
	s.mu.Unlock()
	return
}

// Get returns the current state and value.
func (s *Stateful[S, V]) Get() (state S, value V) {
	s.mu.Lock()
	state, value = s.state, s.value
	s.mu.Unlock()
	return
}

// GetValue returns the current value.
func (s *Stateful[S, V]) GetValue() (value V) {
	s.mu.Lock()
	value = s.value
	s.mu.Unlock()
	return
}

func (s *Stateful[S, V]) addWaiter(wt stateChangeWaiterType, ts S) chan struct{} {
	c := make(chan struct{})
	s.waiters = append(s.waiters, stateChangeWaiter[S]{c: c, wt: wt, targetState: ts})
	return c
}

func (s *Stateful[S, V]) removeWaiterWithChan(c chan struct{}) {
	for i := range s.waiters {
		if s.waiters[i].c != c {
			continue
		}

		lw := len(s.waiters)
		if i < lw-1 {
			s.waiters[i] = s.waiters[lw-1]
		}
		s.waiters[lw-1] = stateChangeWaiter[S]{}
		s.waiters = s.waiters[:lw-1]
		break
	}
}

func (s *Stateful[S, V]) waitStateFunc(ctx context.Context, targetState S, wt stateChangeWaiterType) (state S, value V, err error) {
	var c chan struct{}

	s.mu.Lock()
	if s.matchesType(s.state, targetState, wt) {
		state, value = s.state, s.value
	} else {
		c = s.addWaiter(wt, targetState)
	}
	s.mu.Unlock()

	if c != nil {
		select {
		case <-c:
			s.mu.Lock()
			state, value = s.state, s.value
			s.mu.Unlock()
		case <-ctx.Done():
			err = context.Cause(ctx)
			s.mu.Lock()
			s.removeWaiterWithChan(c)
			s.mu.Unlock()
		}
	}

	return
}

// WaitStateAtLeast waits until the state is changed to be at least (>=) the
// passed target state or the context is canceled. It returns immediately if the
// current state already matches this requirement.
func (s *Stateful[S, V]) WaitStateAtLeast(ctx context.Context, targetState S) (state S, value V, err error) {
	return s.waitStateFunc(ctx, targetState, waitAtLeast)
}

// WaitStateNotGreater waits until the state is changed to be no greater (<=)
// than the passed target state or the context is canceled. It returns
// immediately if the current state already matches this requirement.
func (s *Stateful[S, V]) WaitStateNotGreater(ctx context.Context, targetState S) (state S, value V, err error) {
	return s.waitStateFunc(ctx, targetState, waitNotGreater)
}

// WaitStateExactly waits until the state is changed to be exactly (==) equal to
// the passed target state or the context is canceled. It returns immediately if
// the current state already matches this requirement.
func (s *Stateful[S, V]) WaitStateExactly(ctx context.Context, targetState S) (state S, value V, err error) {
	return s.waitStateFunc(ctx, targetState, waitExactly)
}
