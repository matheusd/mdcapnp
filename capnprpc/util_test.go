// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sourcegraph/conc/pool"
)

type testVat struct {
	name string
	*Vat
}

type testHarness struct {
	t   testing.TB
	ctx context.Context
	g   *pool.ContextPool
}

func (th *testHarness) newVat(name string) *testVat {
	v := NewVat()
	th.g.Go(func(ctx context.Context) error {
		err := v.Run(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("%s Run() errored: %w", name, err)
		}
		return err
	})
	return &testVat{Vat: v, name: name}
}

func (th *testHarness) newTestConn() *testConn {
	return &testConn{
		th:          th,
		sent:        make(chan testConnBatch),
		fillReceive: make(chan testConnReceiver),
	}
}

func newTestHarness(t testing.TB) *testHarness {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	g := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Add g.Wait first to the cleanup because cancel() should be called
	// first (FILO).
	t.Cleanup(func() {
		err := g.Wait()
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Logf("Harness run group errored: %v", err)
			if !t.Failed() {
				t.FailNow()
			}
		}
	})

	t.Cleanup(cancel)

	go func() {
	}()

	th := &testHarness{
		ctx: ctx,
		t:   t,
		g:   g,
	}

	return th
}
