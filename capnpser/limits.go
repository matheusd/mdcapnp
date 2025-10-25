// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import (
	"fmt"
	"math"
	"sync/atomic"
)

// MaxReadLimiterLimit is the maximum allowed limit that can be used when
// initializing a [ReadLimiter].
const MaxReadLimiterLimit = math.MaxInt64

var errLimitOverMaxReadLimiter = fmt.Errorf("cannot read more than %d words", MaxReadLimiterLimit)

// readLimiterType determines the type of the ReadLimiter.
type readLimiterType int

const (
	rlTypeSafe readLimiterType = iota
	rlTypeUnsafe
	rlTypeNoLimit
)

func (rlt readLimiterType) String() string {
	switch rlt {
	case rlTypeSafe:
		return "safe RL"
	case rlTypeUnsafe:
		return "unsafe RL"
	case rlTypeNoLimit:
		return "no RL"
	default:
		panic(fmt.Sprintf("unknown read limiter type %d", rlt))
	}
}

// ReadLimiter limits the amount of data read while traversing structures.
//
// A zero-valued ReadLimiter will reject all read attempts ([CanRead] will
// always return false).
//
// A read limit may be imposed by calling one of the Init* functions. These
// functions are not safe for concurrent access and should in general be called
// only once per read limiter.
type ReadLimiter struct {
	limit         atomic.Int64
	unsafeLimit   int64
	originalLimit int64
	rlType        readLimiterType
}

// Init sets up the read limiter so that it is concurrent safe for reads and up
// to limit words can be read.
func (rl *ReadLimiter) Init(limit uint64) {
	if limit > MaxReadLimiterLimit {
		panic(errLimitOverMaxReadLimiter)
	}

	rl.rlType = rlTypeSafe
	rl.originalLimit = int64(limit)
	rl.limit.Store(int64(limit))
}

// InitConcurrentUnsafe sets up the read limiter so that it is NOT safe for
// concurrent access and up to limit words can be read.
func (rl *ReadLimiter) InitConcurrentUnsafe(limit uint64) {
	if limit > MaxReadLimiterLimit {
		panic(errLimitOverMaxReadLimiter)
	}

	rl.rlType = rlTypeUnsafe
	rl.originalLimit = int64(limit)
	rl.unsafeLimit = int64(limit)
}

// InitNoLimit sets up the read limiter so that no limit is imposed.
func (rl *ReadLimiter) InitNoLimit() {
	rl.rlType = rlTypeNoLimit
}

// InitFrom copies the settings from another ReadLimiter. If the other
// ReadLimiter has consumed bytes from the limit, this limiter will also reflect
// that.
func (rl *ReadLimiter) InitFrom(other *ReadLimiter) {
	rl.limit.Store(other.limit.Load())
	rl.unsafeLimit = other.unsafeLimit
	rl.originalLimit = other.originalLimit
	rl.rlType = other.rlType
}

// testName returns a description of this RL for tests.
func (rl *ReadLimiter) testName() string {
	if rl == nil {
		return "nil RL"
	}
	if rl.rlType == rlTypeNoLimit {
		return "no RL"
	}
	if rl.rlType == rlTypeUnsafe {
		return "unsafe RL"
	}
	return "safe RL"
}

// Reset the read limiter to its original limit. This is valid even for nil
// read limiters.
func (rl *ReadLimiter) Reset() {
	if rl.rlType != rlTypeNoLimit {
		if rl.rlType == rlTypeSafe {
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

	// No limit imposed.
	if rl.rlType == rlTypeNoLimit {
		return
	}

	// Version used when concurrent safety is not necessary (i.e. only one
	// goroutine is assured to be using the arena). A simple test and dec.
	if rl.rlType == rlTypeUnsafe {
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
