// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

func TestReadLimiterCorrectness(t *testing.T) {
	tests := []struct {
		name    string
		initial uint64
		value   WordCount
		want    error
	}{{
		name:    "zero on zero",
		initial: 0,
		value:   0,
		want:    nil,
	}, {
		name:    "zero on max",
		initial: maxReadOnReadLimiter,
		value:   0,
		want:    nil,
	}, {
		name:    "one on max",
		initial: maxReadOnReadLimiter,
		value:   1,
		want:    nil,
	}, {
		name:    "one on zero",
		initial: 0,
		value:   1,
		want:    ErrReadLimitExceeded{},
	}, {
		name:    "max on max",
		initial: maxReadOnReadLimiter,
		value:   MaxValidWordCount,
		want:    nil,
	}, {
		name:    "1001 on 1000",
		initial: 1000,
		value:   1001,
		want:    ErrReadLimitExceeded{},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewReadLimiter(tc.initial)
			got := rl.CanRead(tc.value)
			require.ErrorIs(t, got, tc.want)
		})
	}

	t.Run("panic on new over max", func(t *testing.T) {
		require.Panics(t, func() {
			NewReadLimiter(maxReadOnReadLimiter + 1)
		})
	})

	t.Run("valid after invalid", func(t *testing.T) {
		rl := NewReadLimiter(1000)
		got := rl.CanRead(1001)
		require.ErrorIs(t, got, ErrReadLimitExceeded{})
		got2 := rl.CanRead(1000)
		require.Nil(t, got2)

	})
}

// BenchmarkCanReadAlternatives benchmarks alternatives to the read limiter
// canRead() function.
func BenchmarkCanReadAlternatives(b *testing.B) {
	var readSz uint64 = 1000

	// This is setup so that the last check fails.

	// Original implementation in go-capnp.
	b.Run("CAS", func(b *testing.B) {
		var rlimit atomic.Uint64
		rlimit.Store(uint64(b.N-1) * readSz)

		for i := range b.N {
			ok := false
			for {
				curr := rlimit.Load()

				var new uint64
				if ok = curr >= readSz; ok {
					new = curr - readSz
				}

				if rlimit.CompareAndSwap(curr, new) {
					break
				}
			}
			if !(ok == (i < b.N-1)) {
				panic(fmt.Sprintf("invalid result at %d", i))
			}
		}
	})

	b.Run("MUTEX", func(b *testing.B) {
		var mtx sync.Mutex
		var rlimit uint64 = uint64(b.N-1) * readSz
		for i := range b.N {
			mtx.Lock()
			ok := rlimit >= readSz
			if ok {
				rlimit -= readSz
			}
			mtx.Unlock()

			if !(ok == (i < b.N-1)) {
				panic(fmt.Sprintf("invalid result at %d", i))
			}
		}
	})

	// Suggested implementation.
	b.Run("CHECK", func(b *testing.B) {
		var rlimit atomic.Uint64
		readLimit := uint64(b.N-1) * readSz
		for i := range b.N {
			curr := rlimit.Add(readSz)
			ok := curr <= readLimit
			if !(ok == (i < b.N-1)) {
				panic("invalid result")
			}
		}
	})
}

func BenchmarkCanReadLimiter(b *testing.B) {
	const readSz = 1000
	rl := NewReadLimiter(uint64((b.N - 1) * readSz))
	b.ResetTimer()

	for i := range b.N {
		got := rl.CanRead(readSz) == nil
		want := (i < b.N-1)
		if got != want {
			panic(fmt.Sprintf("invalid result at %d (got %v, want %v)", i, got, want))
		}
	}
}
