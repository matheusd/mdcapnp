// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"context"
	crand "crypto/rand"
	"io"
	"math/rand/v2"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

type peer struct {
	r   io.Reader
	w   io.Writer
	rng io.Reader

	inBuf []byte

	pool      sync.Pool
	outQueue  chan *[]byte
	gotEcho   chan struct{}
	echoCount atomic.Uint64
}

func (p *peer) inLoop(ctx context.Context) {
	inBuf := p.inBuf
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
			// Echo reply.
			select {
			case p.gotEcho <- struct{}{}:
			case <-ctx.Done():
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
	case <-p.gotEcho:
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

func (p *peer) waitEchoResponseDirect(ctx context.Context) {
	_, err := io.ReadFull(p.r, p.inBuf)
	if err != nil {
		return
	}
}

func newPeer(c net.Conn) *peer {
	var seed [32]byte
	crand.Read(seed[:])
	rng := rand.NewChaCha8(seed)
	return &peer{
		r:        c,
		w:        c,
		rng:      rng,
		inBuf:    make([]byte, 96),
		outQueue: make(chan *[]byte, 100),
		gotEcho:  make(chan struct{}, 100),
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, 96)
				return &buf
			},
		},
	}
}

func BenchmarkTCPMultiplexIO(b *testing.B) {
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
	// go p1.inLoop(ctx)
	// go p1.outLoop(ctx)
	// go p2.inLoop(ctx)
	// go p2.outLoop(ctx)
	go p2.echoLoop(ctx)

	b.ReportAllocs()
	var wantCount uint64
	for b.Loop() {
		// p1.sendEchoRequest(ctx)
		//p1.waitEchoResponse(ctx)
		p1.sendEchoRequestDirect(ctx)
		p1.waitEchoResponseDirect(ctx)
		wantCount++
	}

	require.Equal(b, wantCount, p2.echoCount.Load())
}
