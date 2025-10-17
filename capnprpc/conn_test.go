// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"io"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

type testConnReceiver struct {
	f func() (capnpser.Message, error)
}

type testConn struct {
	th          *testHarness
	sent        chan *capnpser.MessageBuilder
	sentResult  chan error
	fillReceive chan testConnReceiver
}

// checkNextSentRpcMsg is called by test code to check the next message sent.
func (tc *testConn) checkNextSentRpcMsg(f func(types.Message) error) {
	var m *capnpser.MessageBuilder
	select {
	case m = <-tc.sent:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}

	serBytes, err := m.Serialize()
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

func (tc *testConn) fillNextReceiveWithSer(mb *capnpser.MessageBuilder) {
	tc.fillNextReceive(func() (capnpser.Message, error) {
		bytes, err := mb.Serialize()
		if err != nil {
			return capnpser.Message{}, err
		}
		arena, err := capnpser.DecodeArena(bytes)
		if err != nil {
			return capnpser.Message{}, err
		}
		arena.ReadLimiter().InitNoLimit()
		serMsg := capnpser.MakeMsg(arena)
		return serMsg, nil
	})
}

func (tc *testConn) fillNextReceive(f func() (capnpser.Message, error)) {
	select {
	case tc.fillReceive <- testConnReceiver{f: f}:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatal("Vat did not ask to receive message")
	}
}

// send is called by the vat end of this test conn. It waits until test code had
// a chance to decide what to do with the message.
func (tc *testConn) send(ctx context.Context, m OutMsg) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case tc.sent <- m.Msg:
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
func (tc *testConn) receive(ctx context.Context) (InMsg, error) {
	select {
	case <-ctx.Done():
		return InMsg{}, context.Cause(ctx)
	case tcr := <-tc.fillReceive:
		msg, err := tcr.f()
		return InMsg{Msg: msg}, err
	}
}

func (tc *testConn) remoteName() string {
	return "testconn"
}

type testPipeConn struct {
	remName  string
	remIndex int

	w     io.Writer
	r     io.Reader
	inBuf []byte

	recvArena  capnpser.Arena
	recvSerMsg capnpser.Message
}

func (tpc *testPipeConn) send(ctx context.Context, m OutMsg) error {
	// Adapt sending raw bytes.
	serBytes, err := m.Msg.Serialize()
	if err != nil {
		return err
	}

	_, err = tpc.w.Write(serBytes)
	return err
}

func (tpc *testPipeConn) receive(ctx context.Context) (InMsg, error) {
	n, err := tpc.r.Read(tpc.inBuf)
	if err != nil {
		return InMsg{}, err
	}

	tpc.recvArena.ReadLimiter().InitConcurrentUnsafe(64 * 1024 * 1024)
	err = tpc.recvArena.DecodeSingleSegment(tpc.inBuf[:n])
	if err != nil {
		return InMsg{}, err
	}
	tpc.recvSerMsg = capnpser.MakeMsg(&tpc.recvArena)

	return InMsg{Msg: tpc.recvSerMsg}, nil
}

func (tpc *testPipeConn) remoteName() string {
	return tpc.remName
}
