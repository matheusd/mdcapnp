// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"encoding/binary"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type cleanupStruct struct {
	b  []byte
	pb *[]byte
}

// BenchmarkCleanup benchmarks some cleanup scenarios.
func BenchmarkCleanup(b *testing.B) {
	var global *cleanupStruct // Ensure alloc escapes to heap.

	const sliceSize = 1<<14 - 1

	debug.SetMemoryLimit(1 << 31) // 2GB

	checkCleanupSize := func(b *testing.B, cleanupSize *atomic.Int64) {
		wantSize := int64((b.N - 1) * sliceSize)
		var gotSize int64
		for i := 0; i < 1000; i++ {
			runtime.GC()
			gotSize = cleanupSize.Load()
			if gotSize == wantSize {
				return // Test done.
			}
			time.Sleep(time.Microsecond)
		}

		b.Fatalf("Wrong size: got %d, want %d", gotSize, wantSize)
	}

	// Baseline for comparison.
	b.Run("baseline", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			global = &cleanupStruct{
				b: make([]byte, sliceSize),
			}
		}
	})

	// Add a cleanup func.
	b.Run("cleanup", func(b *testing.B) {
		var cleanupSize atomic.Int64
		cleanupFunc := func(b []byte) {
			cleanupSize.Add(int64(len(b)))
		}

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			global = &cleanupStruct{
				b: make([]byte, sliceSize),
			}
			runtime.AddCleanup(global, cleanupFunc, global.b)
		}

		b.StopTimer()

		checkCleanupSize(b, &cleanupSize)
	})

	// Set a finalizer func.
	b.Run("finalizer", func(b *testing.B) {
		var cleanupSize atomic.Int64
		finalFunc := func(p *cleanupStruct) {
			cleanupSize.Add(int64(len(p.b)))
		}

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			global = &cleanupStruct{
				b: make([]byte, sliceSize),
			}
			runtime.SetFinalizer(global, finalFunc)
		}
		b.StopTimer()

		checkCleanupSize(b, &cleanupSize)
	})
}

// BenchmarkReuseBufAfterGC benchmarks scenarios about reusing a buffer after GC
// releases it.
func BenchmarkReuseBufAfterGC(b *testing.B) {
	var global *cleanupStruct // Ensure alloc escapes to heap.

	const sliceSize = 1<<14 - 1

	debug.SetMemoryLimit(1 << 31) // 2GB

	checkCleanupSize := func(b *testing.B, cleanupSize *atomic.Int64) {
		wantSize := int64((b.N - 1) * sliceSize)
		var gotSize int64
		for i := 0; i < 1000; i++ {
			runtime.GC()
			gotSize = cleanupSize.Load()
			if gotSize == wantSize {
				return // Test done.
			}
			time.Sleep(time.Microsecond)
		}

		b.Fatalf("Wrong size: got %d, want %d", gotSize, wantSize)
	}

	// Baseline for comparison.
	b.Run("baseline no reuse", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			global = &cleanupStruct{
				b: make([]byte, sliceSize),
			}
		}
	})

	// Baseline where buf is returned to pool immediately (because we know
	// it won't be used and immediately discarded).
	b.Run("baseline reuse immediately", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, sliceSize)
				return &b
			},
		}

		var cleanupSize atomic.Int64

		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			if i > 0 {
				buf := global.pb
				pool.Put(buf)
				cleanupSize.Add(int64(len(*buf)))
			}

			buf := pool.Get().(*[]byte)
			global = &cleanupStruct{
				pb: buf,
			}
		}
		b.StopTimer()

		checkCleanupSize(b, &cleanupSize)
	})

	// Use a runtime.AddCleanup function to return to the pool.
	b.Run("cleanup", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, sliceSize)
				return &b
			},
		}

		var cleanupSize atomic.Int64
		cleanupFunc := func(b *[]byte) {
			pool.Put(b)
			cleanupSize.Add(int64(len(*b)))
		}

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			// runtime.GC()
			buf := pool.Get().(*[]byte)
			global = &cleanupStruct{
				pb: buf,
			}

			runtime.AddCleanup(global, cleanupFunc, buf)
		}
		b.StopTimer()

		checkCleanupSize(b, &cleanupSize)
	})

	b.Run("finalizer", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				b := make([]byte, sliceSize)
				return &b
			},
		}

		var cleanupSize atomic.Int64
		finalFunc := func(p *cleanupStruct) {
			cleanupSize.Add(int64(len(*p.pb)))
			pool.Put(p.pb)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			// runtime.GC()
			buf := pool.Get().(*[]byte)
			global = &cleanupStruct{
				pb: buf,
			}

			runtime.SetFinalizer(global, finalFunc)
		}
		b.StopTimer()

		checkCleanupSize(b, &cleanupSize)
	})
}

var globalTestFinCorrect *cleanupStruct // Ensure allocs go to the heap.

// TestFinalizerCorrectness tests that basic cleanup of byte slices with
// SetFinalizer works as expected.
func TestFinalizerCorrectness(t *testing.T) {
	const memLimit = 1 << 31 // 2GB
	const sliceSize = 1 << 14
	const trials = memLimit / sliceSize * 10
	debug.SetMemoryLimit(memLimit)

	// Track every index that has been finalized as a bitmap.
	var seenMtx sync.Mutex
	seenBuf := make([]uint64, trials/64+1)

	// Uncomment to double check if test would catch double sets.
	// seenBuf[666] = 1 << 7

	var bugged atomic.Bool
	var finalCount atomic.Int64

	pool := sync.Pool{
		New: func() any {
			b := make([]byte, sliceSize)
			return &b
		},
	}
	finalFunc := func(p *cleanupStruct) {
		buf := *p.pb
		i := binary.LittleEndian.Uint64(buf)
		seenIdx, seenBit := i/64, i%64
		seenMask := uint64(1 << seenBit)

		// Uncomment to double check if the test would catch unseen.
		/*
			if seenIdx == 666 && seenBit == 7 {
				return
			}
		*/

		// Fetch the byte and mark the corresponding bit as seeen.
		seenMtx.Lock()
		seenByte := seenBuf[seenIdx]
		seenBuf[seenIdx] |= seenMask
		seenMtx.Unlock()

		alreadySeen := seenByte&seenMask != 0
		if alreadySeen {
			t.Logf("Seen twice at byte %d bit %d", seenIdx, seenBit)
			bugged.Store(true)
		}

		clear((*p.pb)[:8])
		pool.Put(p.pb)
		finalCount.Add(1)
	}

	var i uint64
	for i = range trials {
		buf := pool.Get().(*[]byte)
		globalTestFinCorrect = &cleanupStruct{
			pb: buf,
		}
		binary.LittleEndian.PutUint64(*buf, i)
		runtime.SetFinalizer(globalTestFinCorrect, finalFunc)
	}

	if bugged.Load() {
		t.Fatalf("Someone set twice")
	}

	globalTestFinCorrect = nil
	for i := 0; i <= 1000; i++ {
		if finalCount.Load() == trials {
			break
		} else if i == 1000 {
			t.Fatalf("Not everyone cleaned: cleaned %d, wanted %d", finalCount.Load(), trials)
		}
		time.Sleep(time.Microsecond)
		runtime.GC()
	}

	for i := 0; i < trials; i += 64 {
		seenIdx := i / 64
		seenByte := seenBuf[seenIdx]
		if seenByte != 0xffffffffffffffff {
			t.Fatalf("Not everyone seen at index %d: %016x", seenIdx, seenByte)
		}
	}

	t.Logf("Seen %d finalized", trials)
}
