// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"slices"

	"github.com/sourcegraph/conc/pool"
	"matheusd.com/mdcapnp/capnpser"
)

type vat struct {
	newConn  chan *runningConn
	connDone chan conn

	inMsg     chan inMsg
	pipelines chan runningPipeline
}

func (v *vat) RunConn(c conn) *runningConn {
	rc := newRunningConn(c, v)
	v.newConn <- rc
	return rc
}

func (v *vat) execPipeline(ctx context.Context, p *pipeline) error {
	rp, err := p.PrepareRunning()
	if err != nil {
		return err
	}

	// If this pipeline is a fork, wait until the its parent step is
	// running, which means it can proceed.
	if p.parent != nil {
		parentStep := p.parent.Step(p.parentStepIdx)
		_, err := parentStep.stepRunning.Wait(ctx)
		if err != nil {
			return err
		}
	}

	// TODO: prepare messages for sending.
	for i := range rp.steps {
		rp.steps[i].serMsg = capnpser.MakeMsg(nil)
		rp.steps[i].rpcMsg = Message{}
	}

	// Send the pipeline for processing by the vat's goroutine. This cashes
	// out into vat.startPipeline().
	ctx, rp.cancel = context.WithCancelCause(ctx)
	select {
	case v.pipelines <- rp:
	case <-ctx.Done():
		return ctx.Err()
	}

	// Wait until pipeline has started processing completely. This is
	// signalled by the pipeline's context getting done.
	<-ctx.Done()
	err = context.Cause(ctx)
	if errors.Is(err, errPipelineStarted) {
		err = nil
	}

	return err
}

func (v *vat) runConn(g *pool.ContextPool, ctx context.Context, rc *runningConn) {
	rc.ctx, rc.cancel = context.WithCancel(ctx)

	connG := pool.New().WithContext(rc.ctx).WithCancelOnError().WithFirstError()
	connG.Go(func(ctx context.Context) error {
		for {
			var msg capnpser.Message // TODO: obtain from pool in vat
			err := rc.c.receive(ctx, msg)
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

	connG.Go(func(ctx context.Context) error {
		for {
			select {
			case msg := <-rc.outQueue:
				err := rc.c.send(ctx, msg)
				// TODO: return msg to vat pool
				if err != nil {
					return err
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	// Start the bootstrap pipeline (sends the bootstrap message).
	connG.Go(func(ctx context.Context) error {
		return v.execPipeline(ctx, rc.boot.pipe)
	})

	// Add this running conn to the vat's running pool.
	g.Go(func(ctx context.Context) error {
		err := connG.Wait()
		// Consider canceled a graceful conn closure.
		if errors.Is(err, context.Canceled) {
			err = nil
		}

		// TODO: returning any errors here will make the vat error out.
		// Add options to log and not fail.
		return err
	})
}

// startPipeline starts processing a pipeline. This sends the entire pipeline to
// the respective remote vats and modifies the local vat's state according to
// each step.
func (v *vat) startPipeline(ctx context.Context, rp runningPipeline) error {
	type outBatch struct { // Needs to support multiple conns on a single pipeline?
		conn     *runningConn
		batch    msgBatch
		startIdx int
		endIdx   int
	}
	batches := make([]outBatch, 0, 1)

	// Determine how the local vat will change in response to this
	// pipeline.
	var lastConn *runningConn
	for i, step := range rp.steps {
		if err := v.prepareOutMessage(ctx, rp.pipe, &step); err != nil {
			return err
		}
		conn := step.step.conn
		if conn != lastConn {
			batches = append(batches, outBatch{
				conn:     conn,
				batch:    msgBatch{msgs: make([]capnpser.Message, 0, len(rp.steps)-i)},
				startIdx: i,
			})
		}
		outbToAppend := &batches[len(batches)-1]
		outbToAppend.batch.msgs = append(outbToAppend.batch.msgs, step.serMsg)
		outbToAppend.endIdx = i
		lastConn = conn
	}

	// Send resulting batch to remote sides.
	for _, batch := range batches {
		// Generally, this fails only if ctx is done or if the outbound
		// queue for this conn is full.
		err := batch.conn.queue(ctx, batch.batch)
		if err != nil {
			return err
		}

		// Commit changes of this batch the local vat.
		for i := batch.startIdx; i <= batch.endIdx; i++ {
			// Commit the changes to the local vat.
			step := &rp.steps[i]
			if err := v.commitOutMessage(ctx, rp.pipe, i, step); err != nil {
				return err
			}

			// This step is now in flight. Allow forks from it to
			// start. The forks won't go out after the entirety of
			// this pipeline has processed, because this is running
			// within the vat's main goroutine.
			step.step.stepRunning.Set(step.qid)
		}
	}

	// Entire pipeline is in flight.
	return errPipelineStarted
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
		if rc.c == c {
			rc.cancel()
			return true
		}
		return false
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
		err := v.processInMessage(rs.ctx, m.rc, m.msg)
		if err != nil {
			// TODO: should the error cancel the vat or just the
			// conn?
			return err
		}

		// TODO: return m.msg to vat's msg pool.

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
