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
	f func(*Message) error
}

type testConn struct {
	th          *testHarness
	sent        chan testConnBatch
	fillReceive chan testConnReceiver
}

// checkNextSent is called by test code to check the next message sent.
func (tc *testConn) checkNextSent(f func(msgBatch) error) {
	select {
	case tcb := <-tc.sent:
		tcb.res <- f(tcb.b)
	case <-tc.th.ctx.Done():
		tc.th.t.Fatalf("No message sent before context done")
	}
}

func (tc *testConn) fillNextReceiveWith(target Message) {
	tc.fillNextReceive(func(m *Message) error {
		*m = target
		return nil
	})
}

func (tc *testConn) fillNextReceive(f func(m *Message) error) {
	select {
	case tc.fillReceive <- testConnReceiver{f: f}:
	case <-tc.th.ctx.Done():
		tc.th.t.Fatal("Vat did not ask to receive message")
	}
}

// send is called by the vat end of this test conn. It waits until test code had
// a chance to decide what to do with the message.
func (tc *testConn) send(ctx context.Context, b msgBatch) error {
	tcb := testConnBatch{b: b, res: make(chan error, 1)}
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case tc.sent <- tcb:
	}

	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case err := <-tcb.res:
		return err
	}
}

// receive is called by the vat end of this test conn. It allows test code to
// set the next message to be received.
func (tc *testConn) receive(ctx context.Context, m *Message) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case tcr := <-tc.fillReceive:
		return tcr.f(m)
	}
}
