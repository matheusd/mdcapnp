// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"slices"

	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"
	"matheusd.com/mdcapnp/capnpser"
)

type Vat struct {
	log *zerolog.Logger

	// testIDsOffset is only set during tests.
	testIDsOffset int

	newConn  chan *runningConn
	connDone chan *runningConn

	inMsg     chan inMsg
	pipelines chan runningPipeline
}

func NewVat(opts ...VatOption) *Vat {
	cfg := defaultVatConfig()
	cfg.applyOptions(opts...)

	v := &Vat{
		log:       cfg.vatLogger(),
		newConn:   make(chan *runningConn),
		connDone:  make(chan *runningConn),
		inMsg:     make(chan inMsg),
		pipelines: make(chan runningPipeline, 5),
	}
	return v
}

func (v *Vat) RunConn(c conn) *runningConn {
	rc := newRunningConn(c, v)

	// testIDsOffset is set during tests, to randomize the starting range of
	// every table id per vat. This ensures code isn't relying on specific
	// hardcoded low index IDs and can catch programming errors.
	if v.testIDsOffset > 0 {
		rc.questions.lastID = 10000 + QuestionId(v.testIDsOffset)
		rc.answers.lastID = 20000 + AnswerId(v.testIDsOffset)
		rc.exports.lastID = 30000 + ExportId(v.testIDsOffset)
		rc.imports.lastID = 40000 + ImportId(v.testIDsOffset)
	}

	// Set the bootstrap capability.
	// TODO: parametrize on vat creation.
	rc.bootExportId, _ = rc.exports.nextID() // First export, no need to check ok.
	rc.exports.set(rc.bootExportId, export{typ: exportTypeSenderHosted})

	v.newConn <- rc
	return rc
}

func (v *Vat) execPipeline(ctx context.Context, p *pipeline) error {
	rp, err := p.PrepareRunning()
	if err != nil {
		return err
	}

	// If this pipeline is a fork, wait until the its parent step is
	// running, which means it can proceed.
	if p.parent != nil {
		// Start parent if parent hasn't started yet.
		if p.parent.State() == pipelineStateBuilding {
			err := v.execPipeline(ctx, p.parent)

			// Ignore errPipelineNotBuildingState because it may
			// have changed since the State() call.
			if err != nil && !errors.Is(err, errPipelineNotBuildingState) {
				return err
			}
		}

		parentStep := p.parent.Step(p.parentStepIdx)
		_, err := parentStep.stepRunning.Wait(ctx)
		if err != nil {
			return err
		}
	}

	// Shortcicuit empty pipelines (caller is likely to wait on parent).
	if len(rp.steps) == 0 {
		return nil
	}

	// TODO: prepare messages for sending.
	for i := range rp.steps {
		step := &rp.steps[i]
		step.serMsg = capnpser.MakeMsg(nil) // Needed??

		step.rpcMsg = Message{}
		if step.step.interfaceId == 0 && step.step.methodId == 0 {
			// Bootstrap.
			step.rpcMsg.isBootstrap = true
		}
	}

	// Send the pipeline for processing by the Vat's goroutine. This cashes
	// out into Vat.startPipeline().
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

func (v *Vat) runConn(ctx context.Context, rc *runningConn) {
	rc.ctx, rc.cancel = context.WithCancelCause(ctx)

	connG := pool.New().WithContext(rc.ctx).WithCancelOnError().WithFirstError()
	connG.Go(func(ctx context.Context) error {
		for {
			var msg Message // TODO: obtain from pool in Vat
			err := rc.c.receive(ctx, &msg)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return context.Cause(ctx)
			case v.inMsg <- inMsg{rc: rc, msg: msg}:
			}
		}
	})

	connG.Go(func(ctx context.Context) error {
		for {
			select {
			case msg := <-rc.outQueue:
				err := rc.c.send(ctx, msg)
				// TODO: return msg to Vat pool
				if err != nil {
					return err
				}

			case <-ctx.Done():
				return context.Cause(ctx)
			}
		}
	})

	// Start the bootstrap pipeline (sends the bootstrap message).
	/*
		connG.Go(func(ctx context.Context) error {
			return v.execPipeline(ctx, rc.boot.pipe)
		})
	*/

	// Remove conn once it finishes processing.
	go func() {
		err := connG.Wait()
		if err != nil && !errors.Is(err, context.Canceled) {
			v.log.Debug().Err(err).Msg("Conn goroutines finished due to error")
		} else {
			v.log.Trace().Msg("Conn goroutines finished successfully")
		}
		v.connDone <- rc
	}()
}

// startPipeline starts processing a pipeline. This sends the entire pipeline to
// the respective remote Vats and modifies the local vat's state according to
// each step.
func (v *Vat) startPipeline(ctx context.Context, rp runningPipeline) error {
	v.log.Trace().Int("len", len(rp.steps)).Msg("Starting pipeline")

	type outBatch struct { // Needs to support multiple conns on a single pipeline?
		conn     *runningConn
		batch    msgBatch
		startIdx int
		endIdx   int
	}
	batches := make([]outBatch, 0, 1)

	// Determine how the local Vat will change in response to this
	// pipeline.
	var lastConn *runningConn
	for i := range rp.steps {
		step := &rp.steps[i]
		if err := v.prepareOutMessage(ctx, rp, i); err != nil {
			return err
		}
		conn := step.step.conn
		if conn != lastConn {
			batches = append(batches, outBatch{
				conn:     conn,
				batch:    msgBatch{msgs: make([]Message, 0, len(rp.steps)-i)},
				startIdx: i,
			})
		}
		outbToAppend := &batches[len(batches)-1]
		outbToAppend.batch.msgs = append(outbToAppend.batch.msgs, step.rpcMsg)
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

		// Commit changes of this batch the local Vat.
		for i := batch.startIdx; i <= batch.endIdx; i++ {
			// Commit the changes to the local Vat.
			if err := v.commitOutMessage(ctx, rp, i); err != nil {
				return err
			}

			// This step is now in flight. Allow forks from it to
			// start. The forks won't go out after the entirety of
			// this pipeline has processed, because this is running
			// within the Vat's main goroutine.
			step := &rp.steps[i]
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

func (s *vatRunState) delConn(target *runningConn) {
	s.conns = slices.DeleteFunc(s.conns, func(rc *runningConn) bool {
		if rc == target {
			rc.cancel(errors.New("conn deleted")) // Ok to call multiple times if that's the case.
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

func (v *Vat) runStep(rs *vatRunState) error {
	select {
	case rc := <-v.newConn:
		v.runConn(rs.ctx, rc)
		rs.conns = append(rs.conns, rc)

	case oc := <-v.connDone:
		rs.delConn(oc)

	case m := <-v.inMsg:
		// Process input msg.
		err := v.processInMessage(rs.ctx, m.rc, m.msg)
		if err != nil {
			m.rc.cancel(err)
		}

		// TODO: return m.msg to Vat's msg pool.

	case pipe := <-v.pipelines:
		err := v.startPipeline(rs.ctx, pipe)
		if err != nil {
			pipe.cancel(err)

			// Do some errors cause the Vat to error out?
		}

	case <-rs.ctx.Done():
		return rs.ctx.Err()
	}

	return nil
}

func (v *Vat) Run(ctx context.Context) (err error) {
	rs := &vatRunState{
		g: pool.New().WithContext(ctx).WithCancelOnError().WithFirstError(),
	}
	rs.g.Go(func(ctx context.Context) error {
		var err error
		rs.ctx = ctx
		v.log.Info().Msg("Vat is running")
		for err == nil {
			err = v.runStep(rs)
		}
		return err
	})

	return rs.g.Wait()
}
