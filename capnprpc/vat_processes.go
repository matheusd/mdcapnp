// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"weak"

	"github.com/rs/zerolog"
)

func (v *Vat) processBootstrap(ctx context.Context, rc *runningConn, msg message) error {
	bootMsg := msg.AsBootstrap()
	reply := message{
		isReturn: true,
		ret: rpcReturn{
			aid:       AnswerId(bootMsg.qid),
			isResults: true,
			pay: payload{
				content: anyPointer{
					isCapPointer: true,
					cp:           capPointer{index: 0},
				},
				capTable: []capDescriptor{
					{senderHosted: rc.bootExportId},
				},
			},
		},
	}

	// Modify answer table to track the bootrap export. bootExportId is set
	// during conn setup automatically.
	rc.answers.set(AnswerId(bootMsg.qid), answer{eid: rc.bootExportId})

	rc.log.Debug().
		Int("qid", int(bootMsg.qid)).
		Int("eid", int(rc.bootExportId)).
		Msg("Exported Bootstrap")

	return rc.queue(ctx, singleMsgBatch(reply))
}

func (v *Vat) processReturn(ctx context.Context, rc *runningConn, ret rpcReturn) error {
	qid := QuestionId(ret.AnswerId())
	q, ok := rc.questions.get(qid)
	if !ok {
		return fmt.Errorf("question %d not found", qid)
	}

	// TODO: support exception, cancel, etc
	if !ret.IsResults() {
		return fmt.Errorf("only results supported")
	}

	pipe := q.pipe.Value()
	if pipe == nil {
		// This pipeline isn't used anymore (was released and a Finish
		// should've been sent, or will be shortly), so nothing to do.
		rc.log.Debug().Int("qid", int(qid)).Msg("Received Return message for released pipeline")
		return nil
	}

	// Go through cap table, modify imports table based on what was
	// exported by this call.
	//
	// TODO: only do this if the cap is referenced in the content?
	payload := ret.AsResults()
	capTable := payload.CapTable()
	for _, entry := range capTable {
		if !entry.IsSenderHosted() {
			return fmt.Errorf("only senderHosted capabilities supported")
		}
		iid := ImportId(entry.AsSenderHosted())
		rc.imports.set(iid, imprt{typ: importTypeSenderHosted})
		rc.log.Debug().
			Int("qid", int(qid)).
			Int("iid", int(iid)).
			Msg("Imported cap in Return as senderHosted")
	}

	// Get contents of result.
	var stepResult any
	content := payload.Content()
	if content.IsCapPointer() {
		// NOT GOOD. Must have a new type to pass along instead of
		// parsing like this (maybe). Think about embedded caps.
		cp := content.AsCapPointer()
		capIndex := cp.Index()
		if int(capIndex) >= len(capTable) {
			return fmt.Errorf("capability referenced index outside cap table")
		}
		stepResult = capability{eid: ExportId(capTable[capIndex].AsSenderHosted())}
	} else if content.IsStruct() {
		// TODO: copy if its a struct? Or release serialized message if
		// content is just a cap (because it's not needed anymore)?
		stepResult = content.AsStruct()
	} else if content.IsVoid() {
		stepResult = struct{}{}
	} else {
		return errors.New("unknown/unimplemented content type")
	}

	rc.log.Debug().Int("qid", int(qid)).Msg("Processed Return message")

	// Fulfill pieline waiting for this result.
	step := pipe.Step(q.stepIdx)
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}
		ov.value = stepResult
		return pipeStepStateDone, ov, nil
	})
}

var errCallWithoutId = errors.New("call with zero interfaceId and methodId")

func (v *Vat) processCall(ctx context.Context, rc *runningConn, c call) error {
	if c.iid == 0 && c.mid == 0 {
		// Only bootstrap is allowed to have iid+mid == 0.
		return errCallWithoutId
	}

	// Determine the target of this call (either an exported cap or a
	// promised answer).
	var eid ExportId
	if c.target.isPromisedAnswer {
		// Promised answers are in the answer table.
		//
		// TODO: Recursively track it down if the answer is another
		// promise.
		q, ok := rc.answers.get(AnswerId(c.target.pans.qid))
		if !ok {
			return fmt.Errorf("call referenced unknown promised answer %d", c.target.pans.qid)
		}

		eid = q.eid // What about promises?
	} else if c.target.isImportedCap {
		eid = ExportId(c.target.impcap)
	} else {
		return errors.New("unsupported call target")
	}

	exp, ok := rc.exports.get(eid)
	if !ok {
		return fmt.Errorf("call message target determined to be unset export %d", eid)
	}

	if exp.typ != exportTypeLocallyHosted {
		return fmt.Errorf("unsupported export type %d", exp.typ)
	}

	// Make the call!
	callArgs := callHandlerArgs{
		iid:    interfaceId(c.iid),
		mid:    methodId(c.mid),
		params: c.params,
		rc:     rc,
	}

	// Start preparing reply.
	reply := message{
		isReturn: true,
		ret:      rpcReturn{aid: AnswerId(c.qid)},
	}
	crb := &v.crb // Reuse on vat (this is running on the vat's main goroutine).
	crb.payload = payload{content: anyPointer{
		isVoid: true, // Void result by default on non-error.
	}}

	rc.log.Trace().
		Int("qid", int(c.qid)).
		Int("eid", int(eid)).
		Msg("Locally handling call")

	err := exp.handler.Call(rc.ctx, callArgs, crb)
	if ex, ok := err.(callExceptionError); ok {
		// Turn the error into a returned exception.
		reply.ret.isException = true
		reply.ret.exc = ex.ToException()

		rc.log.Debug().
			Int("qid", int(c.qid)).
			Int("eid", int(eid)).
			Dict("ex", zerolog.Dict().
				Int("type", reply.ret.exc.typ).
				Str("reason", reply.ret.exc.reason)).
			Msg("Processed call into exception")

	} else if err != nil {
		// Fatal connection error.
		return err
	} else {
		reply.ret.isResults = true
		reply.ret.pay = crb.payload

		// TODO: Save this in the answers table.

		rc.log.Debug().
			Int("qid", int(c.qid)).
			Int("eid", int(eid)).
			Msg("Processed call into payload result")
	}

	// TODO: Go through capDescriptors and setup exports.

	return rc.queue(ctx, singleMsgBatch(reply))
}

// processInMessage processes an incoming message from a remote Vat.
func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg message) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	var err error
	switch {
	case msg.IsBootstrap():
		err = v.processBootstrap(ctx, rc, msg)
	case msg.IsReturn():
		err = v.processReturn(ctx, rc, msg.AsReturn())
	case msg.IsCall():
		err = v.processCall(ctx, rc, msg.AsCall())
	case msg.testEcho != 0:
		rc.queue(ctx, singleMsgBatch(msg))
	default:
		err = errors.New("unknown Message type")
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		logEvent := rc.log.Err(err)
		if rc.log.GetLevel() < zerolog.InfoLevel {
			logEvent.Any("msg", msg)
		}
		logEvent.Msg("Error while processing inbound message")
	}

	return err
}

var errTooManyOpenQuestions = errors.New("too many open questions")

// prepareOutMessage prepares an outgoing Message message that is part of a
// pipeline to be sent to the remote Vat.
//
// Note: this does _not_ commit the changes to the conn's tables yet.
func (v *Vat) prepareOutMessage(ctx context.Context, pipe *pipeline,
	stepIdx int, prevQid QuestionId) (thisQid QuestionId, err error) {

	var ok bool

	step := pipe.step(stepIdx)
	conn := pipe.conn
	if step.rpcMsg.IsBootstrap() {
		if thisQid, ok = conn.questions.nextID(); !ok {
			return 0, errTooManyOpenQuestions
		}

		step.rpcMsg.boot.qid = thisQid
		conn.log.Debug().
			Int("qid", int(thisQid)).
			Msg("Prepared Bootstrap message")
	} else if step.rpcMsg.IsCall() {
		if thisQid, ok = conn.questions.nextID(); !ok {
			return 0, errTooManyOpenQuestions
		}
		step.rpcMsg.call.qid = thisQid

		// Find the parent step. If it's in this pipeline, access it
		// directly. Otherwise, check whether the parent completed
		// already or not (to determine if this is call pipelined to a
		// an incomplete local call or to a remote promise).
		if stepIdx > 0 {
			step.rpcMsg.call.target = messageTarget{
				isPromisedAnswer: true,
				pans:             promisedAnswer{qid: prevQid},
			}

			conn.log.Debug().
				Int("qid", int(thisQid)).
				Int("pans", int(prevQid)).
				Msg("Prepared call for in-pipeline promised answer")
		} else if pipe.parent == nil {
			// Should never happen, but avoid a panic below.
			return 0, errors.New("non-bootstrap call without a parent pipeline")
		} else {
			// execPipeline() already ensured the required parent
			// step was sent before allowing this pipeline to be
			// started, so it is safe to read from stepRunning
			// without risk of blocking.
			//
			// Nevertheless, we err on side of caution and check
			// that isSet is true.
			//
			// TODO: check if step already completed (stepDone) and
			// turned into an exported cap.
			parentStep := pipe.parent.Step(pipe.parentStepIdx)
			parentStepState, parentStepVal, err := parentStep.value.WaitStateAtLeast(ctx, pipeStepStateRunning)
			if err != nil {
				return 0, err
			}
			if parentStepState == pipelineStepFailed {
				return 0, parentStepVal.err
			}
			parentQid := parentStepVal.qid
			step.rpcMsg.call.target = messageTarget{
				isPromisedAnswer: true,
				pans:             promisedAnswer{qid: parentQid},
			}

			conn.log.Debug().
				Int("qid", int(thisQid)).
				Int("pans", int(parentQid)).
				Msg("Prepared call for parent pipeline promised answer")
		}
	}

	return thisQid, nil
}

// commitOutMessage commits the changes of the pipeline step to the local Vat's
// state, under the assumption that the given pipeline step was successfully
// sent to the remote Vat.
func (v *Vat) commitOutMessage(_ context.Context, pipe *pipeline, stepIdx int) error {
	step := pipe.step(stepIdx)
	conn := pipe.conn
	var qid QuestionId
	var q question
	if step.rpcMsg.isBootstrap {
		qid = step.rpcMsg.boot.qid
		conn.log.Debug().Int("qid", int(qid)).Msg("Comitted Bootstrap message")
	} else if step.rpcMsg.isCall {
		qid = step.rpcMsg.call.qid
		conn.log.Debug().Int("qid", int(qid)).Msg("Comitted Call message")
	} else {
		// Guard against errors while developing.
		return errors.New("unimplemented commitment of message")
	}

	// runtime.AddCleanup(step, conn.cleanupQuestionIdDueToUnref, qid) // TODO: Save cleanup in question in case of early finish?
	runtime.SetFinalizer(step, finalizePipelineStep)
	q = question{pipe: weak.Make(pipe), stepIdx: stepIdx}
	conn.questions.set(qid, q)

	// This step is now in flight. Allow forks from it to start. The forks
	// won't go out after the entirety of this pipeline has processed,
	// because this is running within the Vat's main goroutine.
	if qid > 0 {
		return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
			if os != pipeStepStateBuilding {
				return os, ov, fmt.Errorf("invalid precondition state: %v", os)
			}
			ov.qid = qid
			ov.conn = pipe.conn // TODO: set this earlier in the call stack.
			return pipeStepStateRunning, ov, nil
		})
	}

	return nil
}
