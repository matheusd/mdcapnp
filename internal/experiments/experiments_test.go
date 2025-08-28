// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/bits"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

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

//go:noinline
func setWordPutLEUint64(b []byte, off int, u uint64) {
	binary.LittleEndian.PutUint64(b[off*8:], u)
}

//go:noinline
func setWordCopyUint64(b []byte, off int, u uint64) {
	*(*uint64)(unsafe.Pointer(&b[off*8])) = u
}

//go:noinline
func setWordCopyUint64PtrSlice(b *[]byte, off int, u uint64) {
	*(*uint64)(unsafe.Pointer(&(*b)[off*8])) = u
}

//go:noinline
func setWordPutUint64Slice(b []uint64, off int, u uint64) {
	b[off] = u
}

//go:noinline
func setWordPutUint64PtrSlice(b *[]uint64, off int, u uint64) {
	(*b)[off] = u
}

//go:noinline
func setWordUnsafePointer(b unsafe.Pointer, off int, u uint64) {
	*(*uint64)(unsafe.Add(b, off*8)) = u
}

//go:noinline
func setWordUnsafePtrUint64(b *uint64, off int, u uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(b), off*8)) = u
}

type customWordSlice struct {
	b *uint64
	l int
	c int
}

//go:noinline
func (cws customWordSlice) setWord(off int, u uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(cws.b), off*8)) = u
}

type customWordSliceNoLC struct {
	b *uint64
}

//go:noinline
func (cws customWordSliceNoLC) setWord(off int, u uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(cws.b), off*8)) = u
}

//go:noinline
func (cws *customWordSliceNoLC) setWordPtr(off int, u uint64) {
	*(*uint64)(unsafe.Add(unsafe.Pointer(cws.b), off*8)) = u
}

// BenchmarkSetWordAlternatives benchmarks alternatives for setting a uint64
// inside a memory block specified in various ways.
//
// The functions are purposefully not inlined in order to be able to easily see
// their generated ASM. Check it out with:
//
// go test -run Bench -bench BenchmarkSetWordAlternatives -cpuprofile /tmp/cpu.pprof
// go tool objdump -S -s setWord experiments.test
func BenchmarkSetWordAlternatives(b *testing.B) {
	var seed = rand.Uint64()
	const nbWords = 16

	doBench := func(b *testing.B, f func(off int, u uint64)) {
		var off uint
		var u uint64

		b.ResetTimer()
		for i := range b.N {
			off, seed = uint(i^int(seed))%nbWords, bits.RotateLeft64(seed, i%64)
			u, seed = uint64(i^int(seed)), bits.RotateLeft64(seed, i%64)
			f(int(off), u)
		}
	}

	byteBuf := make([]byte, nbWords*8)
	byteBufPtr := &byteBuf
	wordBuf := make([]uint64, nbWords)
	wordBufPtr := &wordBuf
	unsafeUint64Ptr := unsafe.SliceData(wordBuf)
	unsafePtr := unsafe.Pointer(&wordBuf[0])
	cws := customWordSlice{b: unsafeUint64Ptr, l: nbWords, c: nbWords}
	cwsNoLC := customWordSliceNoLC{b: unsafeUint64Ptr}
	ptrCwsNoLC := &cwsNoLC

	b.Run("le.put []byte", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordPutLEUint64(byteBuf, off, u)
		}
		doBench(b, f)
	})

	b.Run("copy []byte", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordCopyUint64(byteBuf, off, u)
		}
		doBench(b, f)
	})

	b.Run("copy *[]byte", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordCopyUint64PtrSlice(byteBufPtr, off, u)
		}
		doBench(b, f)
	})

	b.Run("put []word", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordPutUint64Slice(wordBuf, off, u)
		}
		doBench(b, f)
	})

	b.Run("put *[]word", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordPutUint64PtrSlice(wordBufPtr, off, u)
		}
		doBench(b, f)
	})

	b.Run("unsafe ptr", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordUnsafePointer(unsafePtr, off, u)
		}
		doBench(b, f)
	})

	b.Run("unsafe *word", func(b *testing.B) {
		f := func(off int, u uint64) {
			setWordUnsafePtrUint64(unsafeUint64Ptr, off, u)
		}
		doBench(b, f)
	})

	b.Run("custom ws", func(b *testing.B) {
		f := func(off int, u uint64) {
			cws.setWord(off, u)
		}
		doBench(b, f)
	})

	b.Run("custom wsNoLC", func(b *testing.B) {
		f := func(off int, u uint64) {
			cwsNoLC.setWord(off, u)
		}
		doBench(b, f)
	})

	b.Run("custom *wsNoLC", func(b *testing.B) {
		f := func(off int, u uint64) {
			ptrCwsNoLC.setWordPtr(off, u)
		}
		doBench(b, f)
	})
}

const CHANINXBUFSIZE = 256

// BenchmarkChanInStack verifies passing values through channels don't cause
// heap allocations.
func BenchmarkChanInStack(b *testing.B) {
	const BUFSIZE = CHANINXBUFSIZE
	type testStruct struct {
		b [BUFSIZE]byte
	}
	c := make(chan testStruct)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		rng := rand.NewChaCha8([32]byte{0: 0x01})
		var b [BUFSIZE]byte
		for {
			rng.Read(b[:])
			select {
			case <-ctx.Done():
				return
			case c <- testStruct{b: b}:
			}
		}
	}()

	var out, zero [BUFSIZE]byte
	time.Sleep(time.Millisecond)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		v := <-c
		out = v.b
	}

	if bytes.Equal(out[:], zero[:]) {
		b.Fatal("boo")
	}
}

// BenchmarkChanHeapWithSyncPool verifies passing heap allocated values in the
// channel.
func BenchmarkChanHeapWithSyncPool(b *testing.B) {
	const BUFSIZE = CHANINXBUFSIZE
	type testStruct struct {
		b *[BUFSIZE]byte
	}
	c := make(chan testStruct)

	pool := sync.Pool{
		New: func() any {
			b := new([BUFSIZE]byte)
			return b
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		rng := rand.NewChaCha8([32]byte{0: 0x01})
		for {
			b := pool.Get().(*[BUFSIZE]byte)
			rng.Read((*b)[:])
			select {
			case <-ctx.Done():
				return
			case c <- testStruct{b: b}:
			}
		}
	}()

	var out, zero [BUFSIZE]byte
	time.Sleep(time.Millisecond)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		v := <-c
		out = *v.b
		pool.Put(v.b)
	}

	if bytes.Equal(out[:], zero[:]) {
		b.Fatal("boo")
	}

}
