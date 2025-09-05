// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
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
	remName string
	in      chan message
	out     chan message
}

func (tpc *testPipeConn) send(ctx context.Context, m message, _ int) error {
	select {
	case tpc.out <- m:
	case <-ctx.Done():
		return context.Cause(ctx)
	}

	return nil
}

func (tpc *testPipeConn) receive(ctx context.Context) (message, error) {
	select {
	case m := <-tpc.in:
		return m, nil
	case <-ctx.Done():
		return message{}, context.Cause(ctx)
	}
}

func (tpc *testPipeConn) remoteName() string {
	return tpc.remName
}
