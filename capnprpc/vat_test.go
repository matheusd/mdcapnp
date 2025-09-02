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
	v := th.newVat("vat")
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
