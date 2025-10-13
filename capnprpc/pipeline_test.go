// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"testing"
)

type testCapFuture callFuture

func (tcf testCapFuture) next() testCapFuture {
	return testCapFuture(remoteCall(
		callFuture(tcf),
		callSetup{
			interfaceId: 1000,
			methodId:    11,
		},
	))
}

//go:noinline
func (tcf testCapFuture) nextNoInline() testCapFuture {
	return testCapFuture(remoteCall(
		callFuture(tcf),
		callSetup{
			interfaceId: 1000,
			methodId:    11,
		},
	))
}

// BenchmarkAddPipeRemoteCall benchmarks adding a remote call to a pipeline
// under various circumstances.
func BenchmarkAddPipeRemoteCall(b *testing.B) {
	v := NewVat()

	b.Run("inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap(v))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next()
		}
	})

	b.Run("no inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap(v))
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.nextNoInline()
		}
	})

	b.Run("fork/inline", func(b *testing.B) {
		f := testCapFuture(newRootFutureCap(v))
		var final testCapFuture = f.next()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			final = f.next()
		}
		_ = final
	})
}
