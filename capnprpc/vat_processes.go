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
		var iid ImportId
		var imp imprt
		if entry.IsSenderHosted() {
			iid = ImportId(entry.AsSenderHosted())
			imp = imprt{typ: importTypeSenderHosted}
		} else if entry.IsSenderPromise() {
			iid = ImportId(entry.AsSenderPromise())
			imp = imprt{typ: importTypeRemotePromise, pipe: q.pipe, stepIdx: q.stepIdx}
		} else {
			return fmt.Errorf("unsupported capability type")
		}

		rc.imports.set(iid, imp)
		rc.log.Debug().
			Int("qid", int(qid)).
			Int("iid", int(iid)).
			Str("ityp", imp.typ.String()).
			Msg("Imported cap contained in Return")
	}

	// Get contents of result.
	var stepResult any
	var stepResultPromise ExportId
	var stepResultType string
	var stepImportId ImportId
	content := payload.Content()
	if content.IsCapPointer() {
		// NOT GOOD. Must have a new type to pass along instead of
		// parsing like this (maybe). Think about embedded caps.
		cp := content.AsCapPointer()
		capIndex := cp.Index()
		if int(capIndex) >= len(capTable) {
			return fmt.Errorf("capability referenced index outside cap table")
		}
		capEntry := capTable[capIndex]
		if capEntry.IsSenderHosted() {
			stepImportId = ImportId(capEntry.AsSenderHosted())
			stepResult = capability{eid: ExportId(stepImportId)}
			stepResultType = "senderHostedCap"
		} else if capEntry.IsSenderPromise() {
			stepImportId = ImportId(capEntry.AsSenderPromise())
			stepResultPromise = ExportId(stepImportId)
			stepResultType = "senderPromise"
		} else {
			return errors.New("unknown cap entry type")
		}

	} else if content.IsStruct() {
		// TODO: copy if its a struct? Or release serialized message if
		// content is just a cap (because it's not needed anymore)?
		stepResult = content.AsStruct()
		stepResultType = "struct"
	} else if content.IsVoid() {
		stepResult = struct{}{}
		stepResultType = "void"
	} else {
		return errors.New("unknown/unimplemented content type")
	}

	// Fulfill pieline waiting for this result.
	step := pipe.Step(q.stepIdx)
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}

		rc.log.Debug().
			Int("qid", int(qid)).
			Str("rtyp", stepResultType).
			Msg("Processed Return message")

		ov.iid = stepImportId
		if stepResultPromise > 0 {
			// This step isn't resolved into a concrete exported
			// cap or struct yet. Keep it running.
			return os, ov, nil
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
	crb := &rc.crb // Ok to reuse (rc is locked).
	crb.payload = payload{content: anyPointer{
		isVoid: true, // Void result by default on non-error.
	}}

	rc.log.Trace().
		Int("qid", int(c.qid)).
		Int("eid", int(eid)).
		Msg("Locally handling call")

	// Make the call!
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

		// Save this in the answers table.
		rc.answers.set(AnswerId(c.qid), answer{})

		// Track all exported caps.
		for _, cp := range crb.payload.capTable {
			if !cp.hasExportId() {
				continue
			}
			capEid := cp.exportId()
			if capExp, ok := rc.exports.get(capEid); ok {
				// TODO: take a pointer instead?
				capExp.refCount++
				rc.exports.set(capEid, capExp)
			} else if cp.IsSenderPromise() { // here sender == local vat.
				capExp = export{typ: exportTypePromise, refCount: 1}
				rc.exports.set(capEid, capExp)

				rc.log.Trace().
					Int("eid", int(capEid)).
					Str("typ", "senderPromise").
					Msg("Exporting capability")
			} else {
				return errors.New("other types of exports not implemented")
			}
		}

		rc.log.Debug().
			Int("qid", int(c.qid)).
			Int("eid", int(eid)).
			Msg("Processed call into payload result")
	}

	return rc.queue(ctx, singleMsgBatch(reply))
}

func (v *Vat) processFinish(ctx context.Context, rc *runningConn, fin finish) error {
	var err error
	aid := AnswerId(fin.qid)

	if !rc.answers.has(aid) {
		err = fmt.Errorf("answer %d not in answers table", aid)
	} else {
		rc.answers.del(AnswerId(fin.qid))
	}

	// TODO: release exported caps?

	if err == nil {
		rc.log.Debug().
			Int("aid", int(aid)).
			Msg("Removed answer due to Finish message")
	}

	return err
}

func (v *Vat) processResolve(ctx context.Context, rc *runningConn, res resolve) error {
	iid := ImportId(res.pid)
	imp, ok := rc.imports.get(iid)
	if !ok {
		return fmt.Errorf("import id %d not found", iid)
	}

	if imp.typ != importTypeRemotePromise {
		return fmt.Errorf("import %d is not a remote promise to be resolved", iid)
	}

	pipe := imp.pipe.Value()
	if pipe == nil {
		// Pipeline already canceled?
		return fmt.Errorf("pipeline of import %d already released", iid)
	}

	// Similar to code in processReturn. Unify?
	capEntry := res.cap
	var resolveCap capability
	var resolvePromise ExportId
	var resolveId ImportId
	var resImport imprt
	if capEntry.IsSenderHosted() {
		// Resolved into a remote capability.
		resolveCap = capability{eid: ExportId(capEntry.AsSenderHosted())}
		resolveId = ImportId(capEntry.AsSenderHosted())
		resImport = imprt{typ: importTypeSenderHosted}
	} else if capEntry.IsSenderPromise() {
		// Resolved into another remote promise.
		resolvePromise = capEntry.AsSenderPromise()
		resolveId = ImportId(resolvePromise)
		resImport = imprt{typ: importTypeRemotePromise, pipe: imp.pipe, stepIdx: imp.stepIdx}
	} else {
		return errors.New("unknown cap entry type")
	}

	step := pipe.Step(imp.stepIdx)
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}

		rc.log.Debug().
			Int("iid", int(iid)).
			Int("resIid", int(resolveId)).
			Str("resTyp", resImport.typ.String()).
			Msg("Resolved remote promise")

		ov.iid = resolveId
		if resolvePromise > 0 {
			// This step isn't resolved into a concrete exported
			// cap or struct yet. Keep it running.
			return os, ov, nil
		}

		ov.value = resolveCap
		return pipeStepStateDone, ov, nil
	})
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
	case msg.IsFinish():
		err = v.processFinish(ctx, rc, msg.AsFinish())
	case msg.IsResolve():
		err = v.processResolve(ctx, rc, msg.AsResolve())
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

			if parentStepVal.iid > 0 {
				step.rpcMsg.call.target = messageTarget{
					isImportedCap: true,
					impcap:        parentStepVal.iid,
				}

				conn.log.Debug().
					Int("qid", int(thisQid)).
					Int("iid", int(parentStepVal.iid)).
					Msg("Prepared call for exported cap")
			} else {
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

func (v *Vat) sendFinish(ctx context.Context, rc *runningConn, qid QuestionId) error {
	rc.mu.Lock()
	rc.questions.del(qid)
	rc.mu.Unlock()

	msg := message{
		isFinish: true,
		finish:   finish{qid: qid},
	}

	return rc.queue(ctx, singleMsgBatch(msg))
}

func (v *Vat) sendResolve(ctx context.Context, rc *runningConn, eid ExportId, exp export, resolution export) error {
	msg := message{
		isResolve: true,
		resolve:   resolve{pid: eid},
	}

	if resolution.typ == exportTypePromise {
		msg.resolve.cap.senderPromise = exp.resolvedToExport
	} else if resolution.typ == exportTypeLocallyHosted {
		msg.resolve.cap.senderHosted = exp.resolvedToExport
	}

	return rc.queue(ctx, singleMsgBatch(msg))
}
