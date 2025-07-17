// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"math"
	"sync/atomic"
)

const MaxReadLimiterLimit = math.MaxInt64

var errLimitOverMaxReadLimiter = fmt.Errorf("cannot read more than %d words", MaxReadLimiterLimit)

// ReadLimiter limits the amount of data read while traversing structures.
//
// A nil ReadLimiter will always allow reads ([CanRead] will always return
// true). A zero-valued ReadLimiter will reject all read attempts ([CanRead]
// will always return false).
type ReadLimiter struct {
	limit            atomic.Int64
	unsafeLimit      int64
	originalLimit    int64
	concurrentUnsafe bool
}

// NewReadLimiter creates a new read limiter.
//
// NOTE: limit cannot be higher than [math.MaxInt64]. This is unlikely to be an
// actual limitation during regular use.
func NewReadLimiter(limit uint64) *ReadLimiter {
	if limit > MaxReadLimiterLimit {
		panic(errLimitOverMaxReadLimiter)
	}
	rl := &ReadLimiter{originalLimit: int64(limit)}
	rl.limit.Store(int64(limit))
	return rl
}

// NewConcurrentUnsafeReadLimiter creates a new read limiter that is NOT safe
// for concurrent access by multiple goroutines.
//
// This limiter may be used when the caller is certain that only a single
// goroutine accesses an arena/message (and any objects/structs/lists/unsafe
// strings derived from such).
func NewConcurrentUnsafeReadLimiter(limit uint64) *ReadLimiter {
	if limit > MaxReadLimiterLimit {
		panic(errLimitOverMaxReadLimiter)
	}
	return &ReadLimiter{
		originalLimit:    int64(limit),
		unsafeLimit:      int64(limit),
		concurrentUnsafe: true,
	}
}

// testName returns a description of this RL for tests.
func (rl *ReadLimiter) testName() string {
	if rl == nil {
		return "nil RL"
	}
	if rl.concurrentUnsafe {
		return "unsafe RL"
	}
	return "safe RL"
}

// Reset the read limiter to its original limit. This is valid even for nil
// read limiters.
func (rl *ReadLimiter) Reset() {
	if rl != nil {
		if rl.concurrentUnsafe {
			rl.limit.Store(rl.originalLimit)
		} else {
			rl.unsafeLimit = rl.originalLimit
		}
	}
}

// CanRead returns nil if [wc] words can be read or an error otherwise. If this
// ReadLimiter was created by a call to NewReadLimiter, then this is safe for
// concurrent access by multiple goroutines.
func (rl *ReadLimiter) CanRead(wc WordCount) (err error) {
	wcu := int64(wc)
	if wcu > MaxReadLimiterLimit {
		return errLimitOverMaxReadLimiter
	}

	if rl == nil {
		return
	}

	// Version used when concurrent safety is not necessary (i.e. only one
	// goroutine is assured to be using the arena). A simple test and dec.
	if rl.concurrentUnsafe {
		if rl.unsafeLimit < wcu {
			return ErrReadLimitExceeded{Target: wc}
		} else {
			rl.unsafeLimit -= wcu
		}
		return
	}

	// Version used when concurrent safety is required.
	//
	// Loop to ensure concurrent calls are correct.
	for {
		limit := rl.limit.Load()
		newLimit := limit - wcu
		if newLimit < 0 {
			return ErrReadLimitExceeded{Target: wc}
		}

		// This will be false if the limit changed between the Load()
		// call above and this point of the execution stack. In that
		// case, try again.
		if rl.limit.CompareAndSwap(limit, newLimit) {
			return
		}
	}
}

// depthLimit is the internal representation of the depth limit when reading and
// de-referencing pointers.
type depthLimit uint

const (
	// noDepthLimit is used as a flag to signal that no depth limit should
	// be applied.
	noDepthLimit depthLimit = math.MaxUint

	// maxDepthLimit is the maximum allowed depth limit.
	maxDepthLimit depthLimit = math.MaxUint - 1

	// defaultDepthLimit is the default depth limit applied when
	// initializing a message.
	defaultDepthLimit depthLimit = 64

	// MaxDepthLimit is the maximum allowed, valid depth limit.
	MaxDepthLimit = uint(maxDepthLimit)
)

func (dl depthLimit) dec() (newDL depthLimit, ok bool) {
	if dl == noDepthLimit {
		return noDepthLimit, true
	} else if dl == 0 {
		return 0, false
	}
	return dl - 1, true
}
