// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"math"
	"sync/atomic"
)

const maxReadOnReadLimiter = math.MaxInt64

var errLimitOverMaxReadLimiter = fmt.Errorf("cannot read more than %d words", maxReadOnReadLimiter)

// ReadLimiter limits the amount of data read while traversing structures.
type ReadLimiter struct {
	limit atomic.Int64
}

// NewReadLimiter creates a new read limiter.
//
// NOTE: limit cannot be higher than [math.MaxInt64]. This is unlikely to be an
// actual limitation, during regular use.
func NewReadLimiter(limit WordCount) *ReadLimiter {
	if limit > maxReadOnReadLimiter {
		panic(errLimitOverMaxReadLimiter)
	}
	rl := &ReadLimiter{}
	rl.limit.Store(int64(limit))
	return rl
}

func (rl *ReadLimiter) CanRead(wc WordCount) (err error) {
	if wc > maxReadOnReadLimiter {
		return errLimitOverMaxReadLimiter
	}
	if rl == nil {
		return nil
	}

	// Loop to ensure concurrent calls are correct.
	for {
		limit := rl.limit.Load()
		newLimit := limit - int64(wc)
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
