// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
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
	resMsg := Message{
		isReturn: true,
		ret: Return{
			aid:       AnswerId(bootQid),
			isResults: true,
			pay: Payload{
				content: anyPointer{
					isCapPointer: true,
					cp:           capPointer{index: 0},
				},
				capTable: []CapDescriptor{
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
	tc.fillNextReceiveWith(Message{isBootstrap: true, boot: Bootstrap{qid: targetQid}})

	// Vat sends the Bootstrap cap.
	var bootQid QuestionId
	tc.checkNextSent(func(mb msgBatch) error {
		ret := mb.msgs[0].ret
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
