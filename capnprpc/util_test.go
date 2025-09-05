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

	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"
)

type testVat struct {
	name string
	*Vat
}

type testHarness struct {
	t        testing.TB
	vatCount int
	ctx      context.Context
	g        *pool.ContextPool
	logger   zerolog.Logger
}

func (th *testHarness) newVat(name string, opts ...VatOption) *testVat {
	var testVatOpts []VatOption
	testVatOpts = append(testVatOpts, WithName(name), WithLogger(&th.logger))
	testVatOpts = append(testVatOpts, opts...)

	v := NewVat(testVatOpts...)
	v.testIDsOffset = (th.vatCount + 1) * 1000
	th.vatCount++
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
		sent:        make(chan msgBatch, 5),
		sentResult:  make(chan error, 5),
		fillReceive: make(chan testConnReceiver),
	}
}

func (th *testHarness) connectVats(v1, v2 *testVat) (rc1, rc2 *runningConn) {
	c1 := testPipeConn{
		remName: v2.name,
		in:      make(chan message, 10),
		out:     make(chan message, 10),
	}
	c2 := testPipeConn{
		remName: v1.name,
		in:      c1.out,
		out:     c1.in,
	}

	rc1 = v1.RunConn(&c1)
	rc2 = v2.RunConn(&c2)
	return
}

func newTestHarness(t testing.TB) *testHarness {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	g := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	var logger zerolog.Logger
	if _, isBench := t.(*testing.B); isBench {
		logger = zerolog.Nop()
	} else {
		logger = zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)))
	}

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

	th := &testHarness{
		ctx:    ctx,
		t:      t,
		g:      g,
		logger: logger,
	}

	return th
}
