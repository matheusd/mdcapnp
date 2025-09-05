// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

// TestBootstrapSendSide tests the client side of a vat performing bootsrap.
func TestBootstrapSendSide(t *testing.T) {
	th := newTestHarness(t)
	v := th.newVat("client")
	tc := th.newTestConn()
	rc := v.RunConn(tc)
	boot := rc.Bootstrap()
	errChan := make(chan error, 1)
	var finalBootCap capability
	go func() {
		var err error
		finalBootCap, err = boot.Wait(testctx.New(t))
		errChan <- err
	}()

	// Vat sends a Bootstrap message.
	var bootQid uint32
	tc.checkNextSent(func(mb msgBatch) error {
		boot := mb.msgs[0].AsBootstrap()
		bootQid = uint32(boot.QuestionId())
		return nil
	})

	// Remote replies with a Return.
	targetExportId := ExportId(666)
	resMsg := message{
		isReturn: true,
		ret: rpcReturn{
			aid:       AnswerId(bootQid),
			isResults: true,
			pay: payload{
				content: anyPointer{
					isCapPointer: true,
					cp:           capPointer{index: 0},
				},
				capTable: []capDescriptor{
					{senderHosted: targetExportId},
				},
			},
		},
	}
	tc.fillNextReceiveWith(resMsg)

	// Bootstrap() fulfilled.
	require.Nil(t, chantest.Before(time.Second).AssertRecv(t, errChan))
	require.Equal(t, targetExportId, finalBootCap.eid)
}

// TestBootstrapReceiveSide tests the bootstrap process from the receiver side.
func TestBootstrapReceiveSide(t *testing.T) {
	th := newTestHarness(t)
	v := th.newVat("server")
	tc := th.newTestConn()
	_ = v.RunConn(tc)

	// Vat receives a Bootstrap message.
	targetQid := QuestionId(666)
	tc.fillNextReceiveWith(message{isBootstrap: true, boot: bootstrap{qid: targetQid}})

	// Vat sends the Bootstrap cap.
	var bootQid QuestionId
	tc.checkNextSent(func(mb msgBatch) error {
		ret := mb.single.ret
		bootQid = QuestionId(ret.aid)
		return nil
	})

	require.Equal(t, targetQid, bootQid)
}

// TestBootstrapBothSides tests the bootstrap process from both sides.
func TestBootstrapBothSides(t *testing.T) {
	th := newTestHarness(t)
	c, s := th.newVat("client"), th.newVat("server")
	cc, cs := th.connectVats(c, s)

	// Request the bootstrap cap on client.
	boot, err := cc.Bootstrap().Wait(testctx.New(t))
	require.NoError(t, err)

	require.Equal(t, boot.eid, cs.bootExportId)
}

// TestVoidCallBothSides tests a void call between two vats.
func TestVoidCallBothSides(t *testing.T) {
	var called atomic.Bool
	handler := callHandlerFunc(func(ctx context.Context, args callHandlerArgs, rb *callReturnBuilder) error {
		called.Store(true)
		return nil
	})

	th := newTestHarness(t)
	c, s := th.newVat("client"), th.newVat("server", withBootstrapHandler(handler))
	cc, _ := th.connectVats(c, s)

	// First call.
	api := testAPIAsBootstrap(cc.Bootstrap())
	err := api.VoidCall().Wait(testctx.New(t))
	require.NoError(t, err)
	require.True(t, called.Load())

	// Second call (bootstrap should be an export already).
	err = api.VoidCall().Wait(testctx.New(t))
	require.NoError(t, err)
}

// BenchmarkVatRunOverhead benchmarks the overhead of a step in Run().
func BenchmarkVatRunOverhead(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		th := newTestHarness(b)
		v := th.newVat("server")
		tc := th.newTestConn()
		_ = v.RunConn(tc)

		var i uint64 = 1
		sendEcho := func() (message, error) {
			return message{testEcho: i}, nil
		}
		recvEcho := func(mb msgBatch) error {
			if mb.single.testEcho != i {
				return errors.New("wrong testEcho number")
			}
			return nil
		}

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			tc.fillNextReceive(sendEcho)
			tc.checkNextSent(recvEcho)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		th := newTestHarness(b)
		v := th.newVat("server")

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			tc := th.newTestConn()
			_ = v.RunConn(tc)

			var i uint64 = 1
			sendEcho := func() (message, error) {
				return message{testEcho: i}, nil
			}
			recvEcho := func(mb msgBatch) error {
				if mb.single.testEcho != i {
					return errors.New("wrong testEcho number")
				}
				return nil
			}

			for pb.Next() {
				tc.fillNextReceive(sendEcho)
				tc.checkNextSent(recvEcho)
			}
		})
	})

	/*
		// Not a great test at the moment.
		b.Run("pre-filled", func(b *testing.B) {
			th := newTestHarness(b)
			v := th.newVat("server")
			tc := th.newTestConn()

			// Re-create and fill buffers to avoid having to run a second
			// goroutine.
			tc.fillReceive = make(chan testConnReceiver, b.N)
			tc.sent = make(chan msgBatch, b.N)
			close(tc.sentResult) // Always returns nil

			var i uint64 = 0
			sendEcho := func() (message, error) {
				i += 1
				return message{testEcho: i}, nil
			}

			for range b.N {
				tc.fillReceive <- testConnReceiver{f: sendEcho}
			}

			b.ReportAllocs()
			b.ResetTimer()

			// Run conn, which processes all messages.
			_ = v.RunConn(tc)

			for i := 0; i < 10000; i++ {
				if len(tc.sent) == b.N {
					return
				}
				time.Sleep(time.Millisecond)
			}
			b.Fatalf("Final sent len: %d", len(tc.sent))
		})
	*/
}

// BenchmarkVoidCall benchmarks a basic void call under various circumstances.
func BenchmarkVoidCall(b *testing.B) {
	var callCount atomic.Uint64
	handler := callHandlerFunc(func(ctx context.Context, args callHandlerArgs, rb *callReturnBuilder) error {
		callCount.Add(1)
		return nil
	})

	b.Run("both", func(b *testing.B) {
		th := newTestHarness(b)
		c, s := th.newVat("client"), th.newVat("server", withBootstrapHandler(handler))
		cc, _ := th.connectVats(c, s)
		callCount.Store(0)
		ctx := testctx.New(b)

		// Wait for bootstrap.
		_, err := cc.Bootstrap().Wait(testctx.New(b))
		require.NoError(b, err)

		// Bootstrap resolved.
		api := testAPIAsBootstrap(cc.Bootstrap())

		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			err := api.VoidCall().Wait(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}

		require.Equal(b, uint64(b.N), callCount.Load())
	})
}
