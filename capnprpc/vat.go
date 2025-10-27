// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"runtime/trace"
	"slices"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/sourcegraph/conc/pool"
	"matheusd.com/mdcapnp/capnpser"
)

type connDone struct {
	rc  *runningConn
	err error
}

type Vat struct {
	cfg     vatConfig
	log     *zerolog.Logger
	rcCount atomic.Uint64

	// testIDsOffset is only set during tests.
	testIDsOffset int

	mbp *MessageBuilderPool

	newConn    chan *runningConn
	connDone   chan connDone
	expAccepts chan expectedAccept
	getAccepts chan getExpectedAccept
}

func NewVat(opts ...VatOption) *Vat {
	cfg := defaultVatConfig()
	cfg.applyOptions(opts...)

	v := &Vat{
		cfg:        cfg,
		log:        cfg.vatLogger(),
		mbp:        NewMessageBuilderPool(),
		newConn:    make(chan *runningConn),
		connDone:   make(chan connDone),
		expAccepts: make(chan expectedAccept, 5),    // Buffered to reduce locking contention on caller
		getAccepts: make(chan getExpectedAccept, 5), // Buffered to reduce locking contention on caller
	}
	return v
}

// isLevel0Sync returns true if this is a level 0 sync client vat.
func (v *Vat) isLevel0Sync() bool {
	return v.cfg.isLevel0Sync
}

func (v *Vat) GetCallMessageBuilder(payloadSizeHint capnpser.WordCount) rpcCallMsgBuilder {
	outMsg, _ := v.mbp.getForPayloadSize(payloadSizeHint)
	return rpcCallMsgBuilder{outMsg: outMsg}
}

// TODO: merge with RunConn??
func (v *Vat) UseRemoteVat(c conn) RemoteVat {
	return RemoteVat{rc: v.RunConn(c)}
}

func (v *Vat) RunConn(c conn) *runningConn {
	rc := newRunningConn(c, v)
	rc.rid = v.rcCount.Add(1)

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

	v.newConn <- rc // TODO: need context here too.
	return rc
}

func (v *Vat) execStep(ctx context.Context, step *pipelineStep) error {

	// Sanity check pipe hasn't run yet.
	var stepConn *runningConn
	err := step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os == pipeStepStateFinished {
			return os, ov, errPipeStepAlreadyFinished
		}
		if os == pipeStepFailed {
			return os, ov, ov.err
		}
		if os != pipeStepStateBuilding {
			return os, ov, errPipelineNotBuildingState
		}
		os = pipeStepStateSending
		stepConn = ov.conn
		return os, ov, nil
	})
	if err != nil {
		return err
	}

	// Wait for parent to be running (if this is a child step)
	var parentStepConn *runningConn
	if step.parent != nil {
		// Start parent if parent hasn't started yet.
		parentState := step.parent.value.GetState()
		if parentState == pipeStepStateBuilding {
			err := v.execStep(ctx, step.parent)

			// Ignore errPipelineNotBuildingState because it may
			// have changed since the State() call (happens if
			// multiple children are starting a single parent
			// concurrently).
			if err != nil && !errors.Is(err, errPipelineNotBuildingState) {
				return err
			}
		}

		parentState, parentVal, err := step.parent.value.WaitStateAtLeast(ctx, pipeStepStateRunning)
		if err != nil {
			return err
		}
		if parentState == pipeStepFailed {
			return parentVal.err
		}
		if parentState == pipeStepStateFinished {
			return errPipeParentStepAlreadyFinished
		}

		// Ok to continue (parent is running or completed).
		parentStepConn = parentVal.conn
	}

	// Prepare the RPC message for this single step
	var cmb rpcCallMsgBuilder
	if step.csetup.callOutMsg.serMsg != nil {
		cmb = step.csetup.callOutMsg
	} else if outMsg, err := v.mbp.getForPayloadSize(0); err != nil { // TODO: size hint?
		return err
	} else if step.parent == nil {
		cmb = rpcCallMsgBuilder{outMsg: outMsg, isBootstrap: step.parent == nil}
		bb, err := cmb.mb.NewBoostrap()
		if err != nil {
			return err
		}
		cmb.builder = capnpser.StructBuilder(bb)
	} else {
		cmb = rpcCallMsgBuilder{outMsg: outMsg, isBootstrap: step.parent == nil}
		call, err := cmb.mb.NewCall()
		if err != nil {
			return err
		}
		_ = call.SetInterfaceId(uint64(step.csetup.InterfaceId))
		_ = call.SetMethodId(uint16(step.csetup.MethodId))
		cmb.builder = capnpser.StructBuilder(call)
		/*
			pbuilder := step.csetup.ParamsBuilder
			if pbuilder != nil {
				pb, err := call.NewParams()
				if err != nil {
					return err
				}
				err = pbuilder(pb)
				if err != nil {
					return fmt.Errorf("param building errored: %w", err)
				}
			}
		*/
	}

	// Determine connection from parent or step itself
	if step.parent != nil {
		stepConn = parentStepConn
	}

	return v.sendStep(ctx, step, stepConn, cmb)
}

func (v *Vat) runConn(ctx context.Context, rc *runningConn) {
	rc.ctx, rc.cancel = context.WithCancelCause(ctx)

	connG := pool.New().WithContext(rc.ctx).WithCancelOnError().WithFirstError()
	connG.Go(rc.inLoop)
	connG.Go(rc.outLoop)
	connG.Go(rc.waitToClose)

	// Start the bootstrap pipeline (sends the bootstrap message).
	// Parametrize doing this automatically?
	/*
		connG.Go(func(ctx context.Context) error {
			return v.execPipeline(ctx, rc.boot.pipe)
		})
	*/

	// Remove conn once it finishes processing.
	go func() {
		err := connG.Wait()
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, errConnStopped) {
			v.log.Debug().Err(err).Msg("Conn goroutines finished due to error")
		} else {
			v.log.Trace().Msg("Conn goroutines finished successfully")
			err = nil
		}
		select {
		case v.connDone <- connDone{rc: rc, err: err}:
		case <-time.After(time.Second): // TODO: parametrize?
			v.log.Warn().Msg("Conn Wait() goroutine timed out waiting to send on vat.connDone")
		}
	}()
}

func (v *Vat) stopConn(rc *runningConn) {
	rc.cancel(errConnStopped) // Ok to call multiple times if that's the case.

	// Every non-answered question is answered with an error.
	for qid, q := range rc.questions.entries {
		step := q.step()
		if step != nil {
			rc.log.Trace().Int("qid", int(qid)).Msg("Cancelling step due to conn done")
			step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
				if ov.err == nil {
					ov.err = errConnDone
				}
				return pipeStepFailed, ov, nil
			})
		}
	}
}

// sendStep sends a pipeline step to the remote conn (i.e. sends a Call
// message).
func (v *Vat) sendStep(ctx context.Context, step *pipelineStep, conn *runningConn, cmb rpcCallMsgBuilder) error {
	// Lock the conn to start modifying its tables. Note: this is harder
	// than it looks because the conn (which was derived from a parent
	// pipeline step) may change from under us during 3PH. See the comment
	// tagged with 3PHCONNISSUE.
	//
	// What can happen is that, during 3PH, the pipeline step which is the
	// parent of this pipeline may have its conn modified and a Disembargo
	// may have been sent for it already. Thus, this conn (and associated
	// message target) is no longer usable as the parent of this pipeline.
	// But we can't know that until we lock the conn and check if the parent
	// didn't change from under us. Thus the need to loop and obtain the
	// conn again _inside_ the lock, to double check.
	//
	// The critical scenario that is fixed by this convoluted lock algo is
	// the following:
	//
	// - sendStep() is called with conn1 (conn1 is unlocked)
	// - resolveThirdPartyCapForPipeStep changes the parent pipeline step to
	// conn2 due to 3PH
	// - resolveThirdPartyCapForPipeStep sends a Disembargo on conn1
	// - conn1 is locked (inside sendStep) to send the Call message
	// - Call is sent after Disembargo: protocol violation!
	for {
		conn.mu.Lock()
		var otherConn *runningConn
		if step.parent != nil {
			otherConn = step.parent.value.GetValue().conn
		} else {
			otherConn = step.value.GetValue().conn
		}
		if otherConn != conn {
			conn.mu.Unlock()
			conn = otherConn
			continue
		}
		break
	}
	defer conn.mu.Unlock()

	stepQid, err := v.prepareOutMessageForStep(ctx, step, conn, cmb)
	if err != nil {
		return err
	}

	// err = conn.queue(ctx, cmb.outMsg)
	err = conn.send(ctx, cmb.outMsg)
	if err != nil {
		return err
	}

	return v.commitOutMessageForStep(ctx, step, conn, stepQid)
}

// vatRunState is the running state of the vat.
//
// This MUST NOT be accessed directly or stored outside of Run() and runStep(),
// under penalty of race conditions.
//
// Individual fields _may_ be accessed on other goroutines, as long as they have
// been properly captured by a closure.
type vatRunState struct {
	ctx        context.Context
	g          *pool.ContextPool
	conns      []*runningConn
	expAccepts map[VatNetworkUniqueID]expectedAccept
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

	case cd := <-v.connDone:
		v.stopConn(cd.rc)
		rs.delConn(cd.rc)
		if v.cfg.failOnConnErr && cd.err != nil {
			return cd.err
		}

	case expAc := <-v.expAccepts:
		rs.expAccepts[expAc.id] = expAc // TODO: How to timeout?
		expAc.srcConn.log.Trace().Hex("uniqueKey", expAc.id[:]).Msg("Registered expected accept in vat")

	case getAc := <-v.getAccepts:
		if expAc, ok := rs.expAccepts[getAc.id]; !ok {
			getAc.replyChan <- err3PHExpectedAcceptNotFound
		} else {
			getAc.replyChan <- expAc        // replyChan is buffered.
			delete(rs.expAccepts, getAc.id) // Picked up cap.
		}

	case <-rs.ctx.Done():
		return rs.ctx.Err()
	}

	return nil
}

func (v *Vat) Run(ctx context.Context) (err error) {
	rs := &vatRunState{
		g:          pool.New().WithContext(ctx).WithCancelOnError().WithFirstError(),
		expAccepts: make(map[VatNetworkUniqueID]expectedAccept),
	}
	rs.g.Go(func(ctx context.Context) error {
		var err error
		rs.ctx = ctx
		v.log.Info().Timestamp().Msg("Vat is running")
		for err == nil {
			err = v.runStep(rs)
		}

		v.log.Trace().Msg("Vat begining stop procedure")

		// Remove every remaining conn.
		for _, rc := range rs.conns {
			v.stopConn(rc)
		}

		v.log.Trace().Msg("All conns commanded to stop")

		// Wait until the remaining conns have signalled their
		// termination.
		for _, rc := range rs.conns {
			select {
			case <-v.connDone:
			case <-time.After(time.Second): // TODO: improve this.
				rc.log.Warn().Msg("vat.Run() timed out waiting for conn to be done")
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
