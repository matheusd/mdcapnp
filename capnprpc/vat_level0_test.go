// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"fmt"
	"testing"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

// TestVoidCallLevel0 tests that a simple void call using a level 0 client
// works.
func TestVoidCallLevel0(t *testing.T) {
	for _, waitBoot := range []bool{true, false} {
		t.Run(fmt.Sprintf("waitBoot=%v", waitBoot), func(t *testing.T) {
			called := make(chan struct{}, 1)
			handler := CapabilityServerFunc(func(rb *CallContext) error {
				called <- struct{}{}
				return nil
			})

			th := newTestHarness(t)
			s := th.newVat("server", WithBootstrapServer(handler))

			io1, io2 := th.tcpTransportPair("client", "server")
			s.RunConn(io2)

			c := th.newLevel0Vat("client", &Level0ClientCfg{Conn: io1})

			ctx := testctx.New(t)

			// Wait for bootstrap.
			boot := c.Bootstrap()
			if waitBoot {
				_, err := boot.Wait(ctx)
				require.NoError(t, err)
			}

			// Execute void call.
			api := testAPIAsBootstrap(boot)
			err := api.VoidCall().Wait(ctx)
			require.NoError(t, err)

			// Ensure server was actually called.
			chantest.AssertRecv(t, called)

			// Execute second call.
			err = api.VoidCall().Wait(ctx)
			require.NoError(t, err)
			chantest.AssertRecv(t, called)
		})
	}
}
