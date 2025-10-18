// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"testing"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

func TestVoidCallLevel0(t *testing.T) {
	called := make(chan struct{})
	handler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		close(called)
		return nil
	})

	th := newTestHarness(t)
	s := th.newVat("server", WithBootstrapHandler(handler))

	io1, io2 := th.tcpTransportPair("client", "server")
	s.RunConn(io2)

	c := NewLevel0ClientVat(Level0ClientCfg{Conn: io1})

	ctx := testctx.New(t)

	// Wait for bootstrap.
	require.NoError(t, c.WaitBootstrap(ctx))

	c.NextCallMsg(testAPI_InterfaceID, testAPI_Void_CallID)
	_, err := c.Call(ctx)
	require.NoError(t, err)
	chantest.AssertRecv(t, called)
}
