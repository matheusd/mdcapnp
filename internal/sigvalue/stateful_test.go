// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package sigvalue

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"matheusd.com/testctx"
)

// TestStatefulWait tests the various Wait() functions of the Stateful signalled
// value.
func TestStatefulWait(t *testing.T) {
	t.Parallel()

	tests := []struct {
		f       string
		is      int32 // Initial state
		ts      int32 // Target state
		ss      int32 // Set() state
		want    bool
		skipSet bool
	}{
		{f: "atLeast", is: 0, ts: 0, ss: 0, want: true},
		{f: "atLeast", is: 0, ts: 10, ss: 0, want: false},
		{f: "atLeast", is: 0, ts: 10, ss: -1, want: false},
		{f: "atLeast", is: 0, ts: 10, ss: +1, want: false},
		{f: "atLeast", is: 0, ts: 10, ss: 9, want: false},
		{f: "atLeast", is: 0, ts: 10, ss: 10, want: true},
		{f: "atLeast", is: 0, ts: 10, ss: 20, want: true},
		{f: "atLeast", is: 100, ts: 200, ss: 201, want: true},
		{f: "atLeast", is: 100, ts: 200, ss: 199, want: false},
		{f: "atLeast", is: 100, ts: 200, skipSet: true, want: false},
		{f: "atLeast", is: 100, ts: 100, skipSet: true, want: true},
		{f: "atLeast", is: 100, ts: 90, skipSet: true, want: true},

		{f: "notGreater", is: 0, ts: 0, ss: 0, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: 0, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: -1, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: +1, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: 9, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: 10, want: true},
		{f: "notGreater", is: 0, ts: 10, ss: 20, want: false},
		{f: "notGreater", is: 100, ts: 200, ss: 201, want: false},
		{f: "notGreater", is: 100, ts: 200, ss: 199, want: true},

		{f: "exactly", is: 0, ts: 0, ss: 0, want: true},
		{f: "exactly", is: 0, ts: 10, ss: 10, want: true},
		{f: "exactly", is: 0, ts: 0, ss: 9, want: false},
		{f: "exactly", is: 0, ts: 0, ss: 11, want: false},
		{f: "exactly", is: 100, ts: 200, ss: 199, want: false},
		{f: "exactly", is: 100, ts: 200, ss: 200, want: true},
		{f: "exactly", is: 100, ts: 200, ss: 201, want: false},
		{f: "exactly", is: 0, ts: 0, skipSet: true, want: true},
		{f: "exactly", is: 100, ts: 200, skipSet: true, want: false},
		{f: "exactly", is: 100, ts: 99, skipSet: true, want: false},
		{f: "exactly", is: 100, ts: 100, skipSet: true, want: true},
		{f: "exactly", is: 100, ts: 101, skipSet: true, want: false},
	}

	const initVal = "initial"
	const setVal = "set"

	for _, tc := range tests {
		tc := tc
		setStr := fmt.Sprintf("set %d", tc.ss)
		if tc.skipSet {
			setStr = "skip set"
		}
		name := fmt.Sprintf("from %d %s wait %s(%d)", tc.is, setStr, tc.f, tc.ts)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			st := NewStateful(tc.is, initVal)
			var f func(ctx context.Context) (string, error)
			switch tc.f {
			case "atLeast":
				f = func(ctx context.Context) (string, error) {
					_, v, err := st.WaitStateAtLeast(ctx, tc.ts)
					return v, err
				}
			case "notGreater":
				f = func(ctx context.Context) (string, error) {
					_, v, err := st.WaitStateNotGreater(ctx, tc.ts)
					return v, err
				}
			case "exactly":
				f = func(ctx context.Context) (string, error) {
					_, v, err := st.WaitStateExactly(ctx, tc.ts)
					return v, err
				}
			default:
				panic("implement f")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			vChan := make(chan string, 1)
			errChan := make(chan error, 1)
			go func() {
				v, err := f(ctx)
				if err != nil {
					errChan <- err
				} else {
					vChan <- v
				}
			}()

			want := initVal
			if !tc.skipSet {
				st.Set(tc.ss, setVal)
				want = setVal
			}

			select {
			case err := <-errChan:
				if tc.want {
					t.Fatalf("Received error %v when expected result", err)
				}
			case got := <-vChan:
				if !tc.want {
					t.Fatalf("Got value %s when expected error", got)
				}
				if got != want {
					t.Fatalf("Unexpected value: got %s, want %s", got, want)
				}
			}
		})
	}
}

// BenchmarkStatefulSet benchmarks setting values in a Stateful under various
// circumstances.
func BenchmarkStatefulSet(b *testing.B) {
	b.Run("no waiters", func(b *testing.B) {
		st := NewStateful(0, 0)

		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			st.Set(i, i)
		}
	})

	b.Run("one waiter always match", func(b *testing.B) {
		// Make a channel big enough to receive every event without
		// needing a new goroutine to process.
		c := make(chan stateChangeEvent[int, int], b.N)
		waiters := []stateChangeWaiter[int, int]{{}}

		st := NewStateful(0, 0)

		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			// Manually set the waiters on every iteration.
			waiters[0] = stateChangeWaiter[int, int]{
				wt:          waitAtLeast,
				targetState: -1,
				c:           c,
			}
			st.waiters = waiters
			st.Set(i, i)
		}
	})

	b.Run("one waiter never match", func(b *testing.B) {
		// Make a channel big enough to receive every event without
		// needing a new goroutine to process.
		c := make(chan stateChangeEvent[int, int], b.N)
		waiters := []stateChangeWaiter[int, int]{{}}

		st := NewStateful(0, 0)

		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			// Manually set the waiters on every iteration.
			waiters[0] = stateChangeWaiter[int, int]{
				wt:          waitNotGreater,
				targetState: -1,
				c:           c,
			}
			st.waiters = waiters
			st.Set(i, i)
		}
	})
}

// BenchmarkStatefulWait benchmarks waiting for a Stateful change in various
// ways.
func BenchmarkStatefulWait(b *testing.B) {
	b.Run("already set", func(b *testing.B) {
		st := NewStateful(0, 0)
		ctx := testctx.New(b)

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			_, _, err := st.WaitStateAtLeast(ctx, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Uses a context that is already done to ensure we go through the
	// entire function without having to create a goroutine.
	b.Run("context done", func(b *testing.B) {
		st := NewStateful(0, 0)
		targetErr := errors.New("guard error")
		ctx, cancel := context.WithCancelCause(context.Background())
		cancel(targetErr)

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			_, _, err := st.WaitStateAtLeast(ctx, 1)
			if !errors.Is(err, targetErr) {
				b.Fatal(err)
			}
		}

	})

	b.Run("need wait", func(b *testing.B) {
		st := NewStateful(0, 0)
		ctx := testctx.New(b)

		chanAdvance := make(chan int, 1)
		go func() {
			for i := range chanAdvance {
				time.Sleep(time.Microsecond) // Yield.
				st.Set(i, i)
			}
		}()

		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			chanAdvance <- i
			gotState, gotValue, err := st.WaitStateAtLeast(ctx, i)
			if err != nil {
				b.Fatal(err)
			}
			if gotState != i || gotValue != i {
				b.Fatalf("unexpected got: %d %d %d", i, gotState, gotValue)
			}
		}

		close(chanAdvance)
	})

}
