// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

type testConnBatch struct {
	b   msgBatch
	res chan error
}

type testConnReceiver struct {
	f func() (message, error)
}

type testConn struct {
	th          *testHarness
	sent        chan message
	sentResult  chan error
	fillReceive chan testConnReceiver
}

// checkNextSentRpcMsg is called by test code to check the next message sent.
func (tc *testConn) checkNextSentRpcMsg(f func(types.Message) error) {
	var m message
	select {
	case m = <-tc.sent:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}

	if m.rawSerMb == nil {
		tc.th.t.Fatalf("Did not receive expected rawSerMb")
	}
	serBytes, err := m.rawSerMb.Serialize()
	if err != nil {
		tc.th.t.Fatal(err)
	}
	arena, err := capnpser.DecodeArena(serBytes)
	if err != nil {
		tc.th.t.Fatal(err)
	}
	arena.ReadLimiter().InitNoLimit()
	msg := capnpser.MakeMsg(arena)
	var rpcMsg types.Message
	if err := rpcMsg.ReadFromRoot(&msg); err != nil {
		tc.th.t.Fatal(err)
	}

	select {
	case tc.sentResult <- f(rpcMsg):
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}
}

// checkNextSent is called by test code to check the next message sent.
func (tc *testConn) checkNextSent(f func(message) error) {
	var m message
	select {
	case m = <-tc.sent:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}
	select {
	case tc.sentResult <- f(m):
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}
}

func (tc *testConn) fillNextReceiveWith(target message) {
	tc.fillNextReceive(func() (message, error) {
		return target, nil
	})
}

func (tc *testConn) fillNextReceiveWithSer(mb *capnpser.MessageBuilder) {
	tc.fillNextReceive(func() (message, error) {
		bytes, err := mb.Serialize()
		if err != nil {
			return message{}, err
		}
		arena, err := capnpser.DecodeArena(bytes)
		if err != nil {
			return message{}, err
		}
		arena.ReadLimiter().InitNoLimit()
		serMsg := capnpser.MakeMsg(arena)
		return message{rawSerMsg: &serMsg}, nil
	})
}

func (tc *testConn) fillNextReceive(f func() (message, error)) {
	select {
	case tc.fillReceive <- testConnReceiver{f: f}:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatal("Vat did not ask to receive message")
	}
}

// send is called by the vat end of this test conn. It waits until test code had
// a chance to decide what to do with the message.
func (tc *testConn) send(ctx context.Context, m message, _ int) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case tc.sent <- m:
	}

	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case err := <-tc.sentResult:
		return err
	}
}

// receive is called by the vat end of this test conn. It allows test code to
// set the next message to be received.
func (tc *testConn) receive(ctx context.Context) (message, error) {
	select {
	case <-ctx.Done():
		return message{}, context.Cause(ctx)
	case tcr := <-tc.fillReceive:
		return tcr.f()
	}
}

func (tc *testConn) remoteName() string {
	return "testconn"
}

type testPipeConn struct {
	remName  string
	remIndex int

	nextOut   chan *message
	wroteOut  chan struct{}
	nextIn    chan *message
	inWritten chan struct{}

	recvArena  capnpser.Arena
	recvMsg    message
	recvSerMsg capnpser.Message
}

func (tpc *testPipeConn) send(ctx context.Context, m message, _ int) error {
	select {
	case nextMsg := <-tpc.nextOut:

		// Adapt sending raw bytes.
		if m.rawSerMb != nil {
			serBytes, err := m.rawSerMb.Serialize()
			if err != nil {
				return err
			}
			nextMsg.rawSerBytes = append(nextMsg.rawSerBytes[:0], serBytes...)
		} else {
			*nextMsg = m
		}
		select {
		case tpc.wroteOut <- struct{}{}:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	case <-ctx.Done():
		return context.Cause(ctx)
	}

	return nil
}

func (tpc *testPipeConn) receive(ctx context.Context) (message, error) {
	select {
	case tpc.nextIn <- &tpc.recvMsg:
	case <-ctx.Done():
		return message{}, context.Cause(ctx)
	}

	select {
	case <-tpc.inWritten:
		// Adapt receiving raw bytes.
		if tpc.recvMsg.rawSerBytes != nil {
			tpc.recvArena.ReadLimiter().InitConcurrentUnsafe(64 * 1024 * 1024)
			err := tpc.recvArena.DecodeSingleSegment(tpc.recvMsg.rawSerBytes)
			if err != nil {
				return message{}, err
			}
			tpc.recvMsg.rawSerBytes = tpc.recvMsg.rawSerBytes[:0]
			tpc.recvSerMsg = capnpser.MakeMsg(&tpc.recvArena)
			tpc.recvMsg.rawSerMsg = &tpc.recvSerMsg
		} else {
			tpc.recvMsg.rawSerMsg = nil
		}
	case <-ctx.Done():
		return message{}, context.Cause(ctx)
	}

	return tpc.recvMsg, nil
}

func (tpc *testPipeConn) remoteName() string {
	return tpc.remName
}
