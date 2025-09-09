// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"runtime/trace"
	"slices"
	"time"

	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"
)

type Vat struct {
	cfg vatConfig
	log *zerolog.Logger

	// testIDsOffset is only set during tests.
	testIDsOffset int

	newConn  chan *runningConn
	connDone chan *runningConn

	inMsg     chan inMsg
	pipelines chan *pipeline

	crb callReturnBuilder
}

func NewVat(opts ...VatOption) *Vat {
	cfg := defaultVatConfig()
	cfg.applyOptions(opts...)

	v := &Vat{
		cfg:       cfg,
		log:       cfg.vatLogger(),
		newConn:   make(chan *runningConn),
		connDone:  make(chan *runningConn),
		inMsg:     make(chan inMsg),
		pipelines: make(chan *pipeline, 5),
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
	rc.bootExportId, _ = rc.exports.nextID() // First export, no need to check ok.
	if v.cfg.bootstrapHandler != nil {
		rc.exports.set(rc.bootExportId, export{typ: exportTypeLocallyHosted, handler: v.cfg.bootstrapHandler})
	}

	v.newConn <- rc
	return rc
}

// execPipeline starts the execution of a pipeline. All parent pipelines are
// started (if not yet started) and preconditional steps are waited for.
//
// This blocks until the pipeline has fully started: the local vat has changed
// its state in response to the pipeline being in-flight.
func (v *Vat) execPipeline(ctx context.Context, p *pipeline) error {
	p.mu.Lock()
	if p.state != pipelineStateBuilding && p.state != pipelineStateBuilt {
		p.mu.Unlock()
		return errPipelineNotBuildingState
	}
	p.state = pipelineStateRunning
	p.mu.Unlock()

	// If this pipeline is a fork, wait until the its parent step is
	// running, which means it can proceed.
	if p.parent != nil {
		// Start parent if parent hasn't started yet.
		parentState := p.parent.State()
		if parentState == pipelineStateBuilding || parentState == pipelineStateBuilt {
			err := v.execPipeline(ctx, p.parent)

			// Ignore errPipelineNotBuildingState because it may
			// have changed since the State() call.
			if err != nil && !errors.Is(err, errPipelineNotBuildingState) {
				return err
			}
		}

		parentStep := p.parent.Step(p.parentStepIdx)
		parentStepState, parentStepVal, err := parentStep.value.WaitStateAtLeast(ctx, pipeStepStateRunning)
		if err != nil {
			return err
		}
		if parentStepState == pipelineStepFailed {
			return parentStepVal.err
		}
	}

	// Shortcicuit empty pipelines (caller is likely to wait on parent).
	if p.isEmpty() {
		return nil
	}

	// Prepare (as much as possible) messages for sending.
	var lastStep *pipelineStep
	for i := range p.numSteps() {
		step := p.step(i)
		step.rpcMsg = message{}
		if step.interfaceId == 0 && step.methodId == 0 {
			// Bootstrap.
			step.rpcMsg.isBootstrap = true
		} else {
			step.rpcMsg.isCall = true
			step.rpcMsg.call = call{
				// target must be set in vat's run().
				iid: step.interfaceId,
				mid: step.methodId,
				// TODO: params?
			}
		}
		lastStep = step
	}

	// Send the pipeline for processing by the Vat's goroutine. This cashes
	// out into Vat.startPipeline().
	select {
	case v.pipelines <- p:
	case <-ctx.Done():
		return context.Cause(ctx)
	}

	// Wait until pipeline has started processing completely. This is
	// signalled by the last step having been set as running.
	stepState, stepVal, err := lastStep.value.WaitStateAtLeast(ctx, pipeStepStateRunning)
	if err != nil {
		return err
	}
	if stepState == pipelineStepFailed {
		return stepVal.err
	}
	return nil
}

func (v *Vat) runConn(ctx context.Context, rc *runningConn) {
	rc.ctx, rc.cancel = context.WithCancelCause(ctx)

	connG := pool.New().WithContext(rc.ctx).WithCancelOnError().WithFirstError()
	connG.Go(func(ctx context.Context) error {
		for {
			msg, err := rc.c.receive(ctx)
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
			case mb := <-rc.outQueue:
				err := rc.c.send(ctx, mb.msg, mb.remainingInBatch)
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
		select {
		case v.connDone <- rc:
		case <-time.After(time.Second):
		}
	}()
}

func (v *Vat) stopConn(rc *runningConn) {
	rc.cancel(errors.New("conn stopped")) // Ok to call multiple times if that's the case.

	// Every non-answered question is answered with an error.
	for qid, q := range rc.questions.entries {
		if q.pipe != nil {
			rc.log.Trace().Int("qid", int(qid)).Msg("Cancelling pipeline step due to conn done")
			q.pipe.mu.Lock()
			q.pipe.state = pipelineStateConnDone
			for i := range q.pipe.numSteps() {
				q.pipe.step(i).value.Set(pipelineStepFailed, pipelineStepStateValue{err: errConnDone})
			}
			q.pipe.mu.Unlock()
		}
	}
}

// startPipeline starts processing a pipeline. This sends the entire pipeline to
// the respective remote Vats and modifies the local vat's state according to
// each step.
func (v *Vat) startPipeline(ctx context.Context, p *pipeline) error {
	v.log.Trace().Int("len", p.numSteps()).Msg("Starting pipeline")

	// Determine how the local Vat will change in response to this
	// pipeline.
	var prevQid QuestionId
	for i := range p.numSteps() {
		var err error
		if prevQid, err = v.prepareOutMessage(ctx, p, i, prevQid); err != nil {
			return err
		}
	}

	// Generally, this fails only if ctx is done or if the outbound
	// queue for this conn is full.
	for i := range p.numSteps() {
		step := p.step(i)
		err := p.conn.queue(ctx, outMsg{msg: step.rpcMsg, remainingInBatch: p.numSteps() - i})
		if err != nil {
			return err
		}
	}

	// Commit the changes to the local Vat.
	for i := range p.numSteps() {
		if err := v.commitOutMessage(ctx, p, i); err != nil {
			return err
		}
	}

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

func (s *vatRunState) delConn(target *runningConn) {
	s.conns = slices.DeleteFunc(s.conns, func(rc *runningConn) bool {
		if rc == target {
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
	traceReg := trace.StartRegion(rs.ctx, "runStep")
	defer traceReg.End()

	select {
	case rc := <-v.newConn:
		v.runConn(rs.ctx, rc)
		rs.conns = append(rs.conns, rc)

	case oc := <-v.connDone:
		v.stopConn(oc)
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
			pipe.failAllSteps(err)

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

		// Remove every remaining conn.
		for _, rc := range rs.conns {
			if rc == nil {
				panic("XXXXX rc is nil here")
			}
			v.stopConn(rc)
		}

		// Wait until the remaining conns have signalled their
		// termination.
		for range rs.conns {
			select {
			case <-v.connDone:
			case <-time.After(time.Second): // TODO: improve this.
			}
		}

		return err
	})

	err = rs.g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		v.log.Err(err).Msg("Vat finished running with unexpected error")
	} else {
		v.log.Info().Msg("Vat finished running successfully")
	}

	return
}
