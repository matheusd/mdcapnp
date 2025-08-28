// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

func newNoReadLimiter(uint64) *ReadLimiter {
	var rl ReadLimiter
	rl.InitNoLimit()
	return &rl
}

func newUnsafeReadLimiter(limit uint64) *ReadLimiter {
	var rl ReadLimiter
	rl.InitConcurrentUnsafe(limit)
	return &rl
}

func newSafeReadLimiter(limit uint64) *ReadLimiter {
	var rl ReadLimiter
	rl.Init(limit)
	return &rl
}

type newRLFunc func(uint64) *ReadLimiter

func rlTestName(newRL newRLFunc) string {
	rl := newRL(0)
	return rl.testName()
}

func (rlt readLimiterType) newRL(limit uint64) *ReadLimiter {
	return rlt.initRL(new(ReadLimiter), limit)
}

func (rlt readLimiterType) initRL(rl *ReadLimiter, limit uint64) *ReadLimiter {
	switch rlt {
	case rlTypeNoLimit:
		rl.InitNoLimit()
	case rlTypeUnsafe:
		rl.InitConcurrentUnsafe(limit)
	case rlTypeSafe:
		rl.Init(limit)
	}
	return rl
}

// rlTestCases is the test matrix for tests and benchmarks that use different
// ReadLimiters.
var rlTestCases = []readLimiterType{rlTypeNoLimit, rlTypeUnsafe, rlTypeSafe}

// benchmarkRLMatrix executes a benchmark using the various possible read
// limiters.
func benchmarkRLMatrix(b *testing.B, f func(b *testing.B, rlType readLimiterType)) {
	for _, rltc := range rlTestCases {
		b.Run(rltc.String(), func(b *testing.B) {
			f(b, rltc)
		})
	}
}

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
		initial: MaxReadLimiterLimit,
		value:   0,
		want:    nil,
	}, {
		name:    "one on max",
		initial: MaxReadLimiterLimit,
		value:   1,
		want:    nil,
	}, {
		name:    "one on zero",
		initial: 0,
		value:   1,
		want:    ErrReadLimitExceeded{},
	}, {
		name:    "max on max",
		initial: MaxReadLimiterLimit,
		value:   MaxValidWordCount,
		want:    nil,
	}, {
		name:    "1001 on 1000",
		initial: 1000,
		value:   1001,
		want:    ErrReadLimitExceeded{},
	}}

	rlTypes := []struct {
		name  string
		newRL func(uint64) *ReadLimiter
	}{
		{name: "safe", newRL: newSafeReadLimiter},
		{name: "unsafe", newRL: newUnsafeReadLimiter},
	}

	for _, rltc := range rlTypes {
		for _, tc := range tests {
			t.Run(rltc.name+"/"+tc.name, func(t *testing.T) {
				rl := rltc.newRL(tc.initial)
				got := rl.CanRead(tc.value)
				require.ErrorIs(t, got, tc.want)
			})
		}

		t.Run(rltc.name+"/panic on new over max", func(t *testing.T) {
			require.Panics(t, func() {
				rltc.newRL(MaxReadLimiterLimit + 1)
			})
		})

		t.Run(rltc.name+"/valid after invalid", func(t *testing.T) {
			rl := rltc.newRL(1000)
			got := rl.CanRead(1001)
			require.ErrorIs(t, got, ErrReadLimitExceeded{})
			got2 := rl.CanRead(1000)
			require.Nil(t, got2)

		})
	}
}

func BenchmarkCanReadLimiter(b *testing.B) {
	const readSz = 1000
	var rl *ReadLimiter // Ensure it escapes to the heap for the test to be fair.

	// This MUST be the first test to ensure rl is nil.
	b.Run("nil limiter", func(b *testing.B) {
		rl := newNoReadLimiter(0)
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			got := rl.CanRead(readSz) == nil
			if got != true {
				panic("unexpected result")
			}
		}
	})

	b.Run("concurrent unsafe", func(b *testing.B) {
		rl = newUnsafeReadLimiter(uint64((b.N - 1) * readSz))
		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			got := rl.CanRead(readSz) == nil
			want := (i < b.N-1)
			if got != want {
				panic(fmt.Sprintf("invalid result at %d (got %v, want %v)", i, got, want))
			}
		}
	})

	b.Run("concurrent safe", func(b *testing.B) {
		rl = newSafeReadLimiter(uint64((b.N - 1) * readSz))
		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			got := rl.CanRead(readSz) == nil
			want := (i < b.N-1)
			if got != want {
				panic(fmt.Sprintf("invalid result at %d (got %v, want %v)", i, got, want))
			}
		}
	})
}

// TestDepthLimit tests the correctness of the depth limiter.
func TestDepthLimit(t *testing.T) {
	tests := []struct {
		dl     depthLimit
		want   depthLimit
		wantOk bool
	}{
		{dl: 0, want: 0, wantOk: false},
		{dl: 1, want: 0, wantOk: true},
		{dl: 2, want: 1, wantOk: true},
		{dl: maxDepthLimit, want: maxDepthLimit - 1, wantOk: true},
		{dl: noDepthLimit, want: noDepthLimit, wantOk: true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%016x", uint(tc.dl)), func(t *testing.T) {
			var dl depthLimit = tc.dl
			got, gotOk := dl.dec()
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantOk, gotOk)
		})
	}
}

// BenchmarkDepthLimitDec benchmarks the dec() function of the depth limiter.
func BenchmarkDepthLimitDec(b *testing.B) {
	const initial = maxDepthLimit
	var dl depthLimit = initial
	var ok bool

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		dl, ok = dl.dec()
	}

	require.Equal(b, initial-depthLimit(b.N), dl)
	require.True(b, ok)
}
