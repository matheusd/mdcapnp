// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
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
	tc.checkNextSentRpcMsg(func(m types.Message) error {
		boot, err := m.AsBootstrap()
		if err != nil {
			return err
		}
		bootQid = uint32(boot.QuestionId())
		return nil
	})

	// Remote replies with a Return.
	targetExportId := ExportId(666)
	msg, mb := th.newRpcMsg()
	ret, _ := msg.NewReturn()
	ret.SetAnswerId(AnswerId(bootQid))
	pay, _ := ret.NewResults()
	pay.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0))
	capTable, _ := pay.NewCapTable(1, 1)
	capDesc := capTable.At(0)
	capDesc.SetSenderHosted(targetExportId)
	tc.fillNextReceiveWithSer(mb)

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
	msg, mb := th.newRpcMsg()
	boot, _ := msg.NewBoostrap()
	boot.SetQuestionId(targetQid)
	tc.fillNextReceiveWithSer(mb)

	// Vat sends the Bootstrap cap.
	var bootQid QuestionId
	tc.checkNextSentRpcMsg(func(m types.Message) error {
		boot, err := m.AsBootstrap()
		if err != nil {
			return err
		}
		bootQid = boot.QuestionId()
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
	handler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		called.Store(true)
		return nil
	})

	th := newTestHarness(t)
	c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(handler))
	cc, _ := th.connectVats(c, s)

	// First call.
	api := testAPIAsBootstrap(cc.Bootstrap())
	voidCall1 := api.VoidCall()
	err := voidCall1.Wait(testctx.New(t))
	require.NoError(t, err)
	require.True(t, called.Load())

	// Second call (bootstrap should be an export already).
	voidCall2 := api.VoidCall()
	err = voidCall2.Wait(testctx.New(t))
	require.NoError(t, err)

	// TODO: verify questions and answers still exist.

	// After this point, voidCall1 and voidCall2 are not used anymore,
	// so they are free to be released.
	runtime.KeepAlive(voidCall2)
	runtime.KeepAlive(voidCall1)
	t.Logf("voidcall1 and voidcall2 no longer referenced")

	runtime.GC()
	time.Sleep(500 * time.Millisecond)

	// TODO: verify answers were deleted.
}

// TestAddCall tests an Add() call between two vats.
func TestAddCall(t *testing.T) {
	handler := CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		req, err := CallContextParamsStruct[addRequest](cc)
		if err != nil {
			return err
		}
		a, b := req.A(), req.B()
		res, err := RespondCallAsStruct[addResponseBuilder](cc, addResponseSize)
		if err != nil {
			return err
		}
		res.SetC(a + b)
		return nil
	})

	th := newTestHarness(t)
	c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(handler))
	cc, _ := th.connectVats(c, s)

	a, b := int64(11), int64(1100)
	addResFuture := testAPIAsBootstrap(cc.Bootstrap()).Add(a, b)
	res, err := addResFuture.wait(testctx.New(t))
	require.NoError(t, err)
	require.Equal(t, a+b, res)
}

// TestRemotePromiseWithCap performs a basic level 1 test (resolving a remote
// promise with a capability).
func TestRemotePromiseWithCap(t *testing.T) {
	var callHandled atomic.Bool
	callHandler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		if !callHandled.CompareAndSwap(false, true) {
			return errors.New("already called")
		}
		return nil
	})

	resolvePromiseChan := make(chan struct{}, 1)
	resolveErrChan := make(chan error, 1)
	bootHandler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		ap, err := rb.respondAsPromise()
		if err != nil {
			return err
		}

		go func() {
			<-resolvePromiseChan
			resolveErrChan <- ap.resolveToHandler(callHandler)
		}()
		return nil
	})

	// Setup harness.
	th := newTestHarness(t)
	c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(bootHandler))
	cc, _ := th.connectVats(c, s)

	// Wait for bootstrap to complete to ease log reviewing.
	api := testAPIAsBootstrap(cc.Bootstrap())
	require.NoError(t, api.WaitDiscardResult(testctx.New(t)))

	// Make a call that returns a capability.
	getCapErrChan := make(chan error, 1)
	getCapCall := api.GetAnotherAPICap()
	go func() {
		getCapErrChan <- getCapCall.WaitDiscardResult(testctx.New(t))
	}()

	// Call isn't done yet (waiting on remote promise).
	chantest.AssertNoRecv(t, getCapErrChan)
	// t.Logf("XXXXXXX %v", <-getCapErrChan)
	// t.FailNow()

	// Resolve.
	chantest.AssertSend(t, resolvePromiseChan, struct{}{})
	gotResolveErr := chantest.AssertRecv(t, resolveErrChan)
	require.Nil(t, gotResolveErr)

	// Call should complete.
	gotVoidErr := chantest.AssertRecv(t, getCapErrChan)
	require.Nil(t, gotVoidErr)

	// Call a method on the returned capability. Could've been pipelined,
	// but we're assessing resolution in this test, not pipelining.
	require.NoError(t, getCapCall.VoidCall().Wait(testctx.New(t)))
	require.True(t, callHandled.Load())
}

/*

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
		recvEcho := func(m message) error {
			if m.testEcho != i {
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
			recvEcho := func(m message) error {
				if m.testEcho != i {
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

	// Not a great test at the moment.
	/*
		b.Run("pre-filled", func(b *testing.B) {
			th := newTestHarness(b)
			v := th.newVat("server")
			tc := th.newTestConn()

			// Re-create and fill buffers to avoid having to run a second
			// goroutine.
			tc.fillReceive = make(chan testConnReceiver, b.N)
			tc.sent = make(chan message, b.N)
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
}
*/

// BenchmarkVoidCall benchmarks a basic void call under various circumstances.
func BenchmarkVoidCall(b *testing.B) {
	var callCount atomic.Uint64
	handler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		callCount.Add(1)
		return nil
	})

	b.Run("pipe", func(b *testing.B) {
		th := newTestHarness(b)
		c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(handler))
		cc, _ := th.connectVats(c, s)
		callCount.Store(0)
		ctx := testctx.New(b)

		// Wait for bootstrap.
		_, err := cc.Bootstrap().Wait(testctx.New(b))
		require.NoError(b, err)

		// Bootstrap resolved.
		api := testAPIAsBootstrap(cc.Bootstrap())

		b.ReportAllocs()
		for b.Loop() {
			err := api.VoidCall().Wait(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}

		// b.Logf("XXXXX sets %d max len %d", xxx_qtsets, xxx_maxqtsize)

		require.Equal(b, uint64(b.N), callCount.Load())
	})

	b.Run("conn", func(b *testing.B) {
		th := newTestHarness(b)
		c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(handler))
		cc, _ := th.connectVatsWithTCP(c, s)
		callCount.Store(0)
		ctx := testctx.New(b)

		// Wait for bootstrap.
		_, err := cc.Bootstrap().Wait(testctx.New(b))
		require.NoError(b, err)

		// Bootstrap resolved.
		api := testAPIAsBootstrap(cc.Bootstrap())

		b.ReportAllocs()
		for b.Loop() {
			err := api.VoidCall().Wait(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}

		// b.Logf("XXXXX sets %d max len %d", xxx_qtsets, xxx_maxqtsize)

		require.Equal(b, uint64(b.N), callCount.Load())
	})

	b.Run("level0", func(b *testing.B) {
		th := newTestHarness(b)
		s := th.newVat("server", WithBootstrapHandler(handler))

		io1, io2 := th.tcpTransportPair("client", "server")
		s.RunConn(io2)

		c := NewLevel0ClientVat(Level0ClientCfg{Conn: io1})

		callCount.Store(0)
		ctx := testctx.New(b)

		// Wait for bootstrap.
		require.NoError(b, c.WaitBootstrap(ctx))

		b.ReportAllocs()
		for b.Loop() {
			c.NextCallMsg(testAPI_InterfaceID, testAPI_Void_CallID)
			_, err := c.Call(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}

		// b.Logf("XXXXX sets %d max len %d", xxx_qtsets, xxx_maxqtsize)

		require.Equal(b, uint64(b.N), callCount.Load())
	})
}

// BenchmarkAddCall benchmarks a basic Add() call under various circumstances.
func BenchmarkAddCall(b *testing.B) {
	handler := CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		req, err := CallContextParamsStruct[addRequest](cc)
		if err != nil {
			return err
		}
		a, b := req.A(), req.B()
		res, err := RespondCallAsStruct[addResponseBuilder](cc, addResponseSize)
		if err != nil {
			return err
		}
		res.SetC(a + b)
		return nil
	})

	b.Run("conn", func(b *testing.B) {
		th := newTestHarness(b)
		c, s := th.newVat("client"), th.newVat("server", WithBootstrapHandler(handler))
		cc, _ := th.connectVatsWithTCP(c, s)
		ctx := testctx.New(b)

		// Wait for bootstrap.
		_, err := cc.Bootstrap().Wait(testctx.New(b))
		require.NoError(b, err)

		// Bootstrap resolved.
		api := testAPIAsBootstrap(cc.Bootstrap())

		var aa, bb int64 = 1, 3
		b.ReportAllocs()
		for b.Loop() {
			aa += 1
			bb = bb<<1 + aa
			res, err := api.Add(aa, bb).wait(ctx)
			if err != nil {
				b.Fatal(err)
			}
			if res != aa+bb {
				b.Fatal("wrong result")
			}
		}

		// b.Logf("XXXXX sets %d max len %d", xxx_qtsets, xxx_maxqtsize)
	})
}
