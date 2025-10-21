// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

const (
	sigModeChannel int = iota
	sigModePipe
)

type peer struct {
	r   io.Reader
	w   io.Writer
	rng io.Reader

	inBuf []byte

	echoPipeR *os.File
	echoPipeW *os.File
	inPipeSig []byte

	sigMode int

	pool      sync.Pool
	outQueue  chan *[]byte
	gotEcho   chan struct{}
	echoCount atomic.Uint64
}

func (p *peer) inLoop(ctx context.Context) {
	// runtime.LockOSThread()
	inBuf := p.inBuf
	pipeSig := make([]byte, 8)
	var i uint64
	for {
		_, err := io.ReadFull(p.r, inBuf)
		if err != nil {
			return
		}

		if inBuf[0] == 1 {
			// echo.
			outBufPtr := p.pool.Get().(*[]byte)
			copy((*outBufPtr), inBuf)
			(*outBufPtr)[0] = 0

			select {
			case p.outQueue <- outBufPtr:
			case <-ctx.Done():
			}
			p.echoCount.Add(1)
		} else {
			i++
			binary.LittleEndian.PutUint64(pipeSig, i)

			// Echo reply.
			switch p.sigMode {
			case sigModeChannel:
				p.gotEcho <- struct{}{}
			case sigModePipe:
				p.echoPipeW.Write(pipeSig)
			}
		}
	}
}

func (p *peer) outLoop(ctx context.Context) {
	var nextOut *[]byte
	for {
		select {
		case nextOut = <-p.outQueue:
		case <-ctx.Done():
			return
		}

		_, err := p.w.Write(*nextOut)
		if err != nil {
			return
		}

		p.pool.Put(nextOut)
	}
}

func (p *peer) echoLoop(ctx context.Context) {
	inBuf := p.inBuf
	for {
		_, err := io.ReadFull(p.r, inBuf)
		if err != nil {
			return
		}

		inBuf[0] = 0
		p.echoCount.Add(1)
		if _, err := p.w.Write(inBuf); err != nil {
			return
		}
	}

}

func (p *peer) sendEchoRequest(ctx context.Context) {
	outBufPtr := p.pool.Get().(*[]byte)
	outBuf := *outBufPtr
	p.rng.Read(outBuf)
	outBuf[0] = 1

	select {
	case p.outQueue <- outBufPtr:
	case <-ctx.Done():
	}
}

func (p *peer) waitEchoResponse(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-p.gotEcho:
		return
	}
}

func (p *peer) sendEchoRequestDirect(_ context.Context) {
	outBufPtr := p.pool.Get().(*[]byte)
	outBuf := *outBufPtr
	p.rng.Read(outBuf)
	outBuf[0] = 1

	p.w.Write(outBuf)
	p.pool.Put(outBufPtr)
}

func (p *peer) waitEchoResponseDirect(_ context.Context) {
	_, err := io.ReadFull(p.r, p.inBuf)
	if err != nil {
		return
	}
}

func (p *peer) waitEchoSigFromPipe(_ context.Context) {
	p.echoPipeR.Read(p.inPipeSig)
}

func newPeer(c net.Conn) *peer {
	var seed [32]byte
	crand.Read(seed[:])
	rng := rand.NewChaCha8(seed)
	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	return &peer{
		r:         c,
		w:         c,
		rng:       rng,
		inBuf:     make([]byte, 96),
		outQueue:  make(chan *[]byte, 100),
		gotEcho:   make(chan struct{}, 10),
		inPipeSig: make([]byte, 8),
		echoPipeR: pr,
		echoPipeW: pw,
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, 96)
				return &buf
			},
		},
	}
}

func BenchmarkTCPMultiplexAlternatives(b *testing.B) {
	var iterIndex uint64

	tests := []struct {
		name     string
		setup    func(b *testing.B, ctx context.Context, p1, p2 *peer)
		loopIter func(b *testing.B, ctx context.Context, p1, p2 *peer)
	}{{
		name: "std",
		setup: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			go p1.inLoop(ctx)
			go p1.outLoop(ctx)
			go p2.inLoop(ctx)
			go p2.outLoop(ctx)
		},
		loopIter: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			p1.sendEchoRequest(ctx)
			p1.waitEchoResponse(ctx)
		},
	}, {
		name: "svrecho",
		setup: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			go p1.inLoop(ctx)
			go p1.outLoop(ctx)
			go p2.echoLoop(ctx)
		},
		loopIter: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			p1.sendEchoRequest(ctx)
			p1.waitEchoResponse(ctx)
		},
	}, {
		name: "pipe",
		setup: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			go p2.echoLoop(ctx)
			p1.sigMode = sigModePipe
			go p1.inLoop(ctx)
		},
		loopIter: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			p1.sendEchoRequestDirect(ctx)
			p1.waitEchoSigFromPipe(ctx)
		},
	}, {
		name: "direct",
		setup: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			go p2.echoLoop(ctx)
		},
		loopIter: func(b *testing.B, ctx context.Context, p1, p2 *peer) {
			p1.sendEchoRequestDirect(ctx)
			p1.waitEchoResponseDirect(ctx)
		},
	}}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(b, err)
			b.Cleanup(func() { lis.Close() })

			connChan := make(chan net.Conn, 2)
			errChan := make(chan error, 2)
			go func() {
				c, err := lis.Accept()
				if err != nil {
					errChan <- err
				} else {
					connChan <- c
				}
			}()
			go func() {
				c, err := net.Dial("tcp", lis.Addr().String())
				if err != nil {
					errChan <- err
				} else {
					connChan <- c
				}
			}()

			conns := make([]net.Conn, 0, 2)
			for len(conns) < 2 {
				select {
				case c := <-connChan:
					conns = append(conns, c)
				case err := <-errChan:
					b.Fatal(err)
				}
			}

			ctx := testctx.New(b)

			p1 := newPeer(conns[0])
			p2 := newPeer(conns[1])
			tc.setup(b, ctx, p1, p2)

			b.ReportAllocs()
			var wantCount uint64
			iterIndex = 0
			for b.Loop() {
				tc.loopIter(b, ctx, p1, p2)
				iterIndex++

				wantCount++
				if ctx.Err() != nil {
					b.Fatal(ctx.Err())
				}
			}

			require.Equal(b, wantCount, p2.echoCount.Load())
		})
	}
}
