// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"slices"

	"github.com/sourcegraph/conc/pool"
)

type vat struct {
	conns []*conn

	newConn  chan *runningConn
	connDone chan conn

	inMsg     chan inMsg
	outMsg    chan outMsg
	pipelines chan *pipeline
}

func (v *vat) RunConn(c conn) *runningConn {
	rc := &runningConn{
		vat: v,
	}

	v.newConn <- rc
	return rc
}

func (v *vat) ExecPipeline(ctx context.Context, p *pipeline) error {
	select {
	case <-p.ctx.Done():
		// This is a programmer error. Pipelines aren't concurrent safe
		// and MUST NOT be executed twice. This is only a soft
		// protection, because the pipeline could be running but not
		// done yet.
		panic("pipeline already executed")
	default:
	}

	// If this pipeline is a fork, wait until the its parent step is
	// running, which means it can proceed.
	if p.parent != nil {
		parentStep := &p.parent.steps[p.parentStepIdx]
		_, err := parentStep.stepRunning.Wait(ctx)
		if err != nil {
			return err
		}
	}

	// Send the pipeline for processing by the vat's goroutine. This cashes
	// out into vat.startPipeline().
	p.ctx, p.cancel = context.WithCancelCause(ctx)
	select {
	case v.pipelines <- p:
	case <-p.ctx.Done():
		return p.ctx.Err()
	}

	// Wait until pipeline has finished processing. This is signalled by the
	// pipeline's context getting done.
	<-p.ctx.Done()
	err := context.Cause(p.ctx)
	if errors.Is(err, errPipelineSuccessful) {
		err = nil
	}

	return err
}

func (v *vat) runConn(g *pool.ContextPool, ctx context.Context, rc *runningConn) {
	rc.ctx, rc.cancel = context.WithCancel(ctx)

	connG := pool.New().WithContext(rc.ctx).WithCancelOnError().WithFirstError()
	connG.Go(func(ctx context.Context) error {
		for {
			var msg message // TODO: pool in vat
			err := rc.c.receive(ctx, &msg)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case v.inMsg <- inMsg{rc: rc, msg: msg}:
			}
		}
	})

	// Add this running conn to the vat's running pool.
	g.Go(func(ctx context.Context) error {
		err := connG.Wait()
		// Consider canceled a graceful conn closure.
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		return err
	})
}

func (v *vat) startPipeline(_ context.Context, pipe *pipeline) error {
	for _, step := range pipe.steps {
		// TODO: run step.
		_ = step

		// This step is now in flight. Allow forks from it to start. The
		// forks won't go out after the entirety of this pipeline has
		// processed, because this is running within the vat's main
		// goroutine.
		if step.stepRunning != nil {
			step.stepRunning.Set(struct{}{})
		}
	}

	// Pipeline is in flight.
	return nil
}

// vatRunState is the running state of the vat.
//
// This MUST NOT be accessed directly or stored outside of Run() and runStep(),
// under penalty of race conditions.
//
// Individual fields _may_ be accessed on other goroutines, as long as they have
// been properly captured by a closure.
type vatRunState struct {
	ctx   context.Context
	g     *pool.ContextPool
	conns []*runningConn
}

func (s *vatRunState) delConn(c conn) {
	s.conns = slices.DeleteFunc(s.conns, func(rc *runningConn) bool {
		return rc.c == c
	})
}

func (s *vatRunState) findConn(c conn) *runningConn {
	i := slices.IndexFunc(s.conns, func(rc *runningConn) bool { return rc.c == c })
	if i < 0 {
		return nil
	}
	return s.conns[i]
}

func (v *vat) runStep(rs *vatRunState) error {
	select {
	case rc := <-v.newConn:
		v.runConn(rs.g, rs.ctx, rc)
		rs.conns = append(rs.conns, rc)

	case oc := <-v.connDone:
		rs.delConn(oc)

	case m := <-v.inMsg:
		// Process input msg.
		_ = m

	case m := <-v.outMsg:
		// Queue outgoing msg
		_ = m // ????

	case pipe := <-v.pipelines:
		err := v.startPipeline(rs.ctx, pipe)
		if err != nil {
			pipe.cancel(err)

			// Do some errors cause the vat to error out?
		}

	case <-rs.ctx.Done():
		return rs.ctx.Err()
	}

	return nil
}

func (v *vat) Run(ctx context.Context) (err error) {
	rs := &vatRunState{
		g: pool.New().WithContext(ctx).WithCancelOnError().WithFirstError(),
	}
	rs.g.Go(func(ctx context.Context) error {
		var err error
		rs.ctx = ctx
		for err == nil {
			err = v.runStep(rs)
		}
		return err
	})

	return rs.g.Wait()
}
