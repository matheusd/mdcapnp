// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

type testCap struct{}
type testCapFuture futureCap[testCap]

func (tcf testCapFuture) next() testCapFuture {
	return testCapFuture(remoteCall[testCap, testCap](futureCap[testCap](tcf), 1000, 11, nil))
}

//go:noinline
func (tcf testCapFuture) nextNoInline() testCapFuture {
	return testCapFuture(remoteCall[testCap, testCap](futureCap[testCap](tcf), 1000, 11, nil))
}

// BenchmarkAddPipeRemoteCall benchmarks adding a remote call to a pipeline
// under various circumstances.
func BenchmarkAddPipeRemoteCall(b *testing.B) {
	b.Run("no hint/inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap[testCap](0))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next()
		}
		require.Equal(b, b.N, f.stepIndex)
	})

	b.Run("no hint/no inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap[testCap](0))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.nextNoInline()
		}
		require.Equal(b, b.N, f.stepIndex)
	})

	b.Run("hint/inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap[testCap](b.N))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next()
		}
		require.Equal(b, b.N, f.stepIndex)
	})

	b.Run("hint/no inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap[testCap](b.N))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.nextNoInline()
		}
		require.Equal(b, b.N, f.stepIndex)
	})

	b.Run("fork/inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap[testCap](b.N))
		var final testCapFuture = f.next() // First one doesn't count.
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			final = f.next()
		}
		require.Equal(b, 0, final.stepIndex)
	})
}
