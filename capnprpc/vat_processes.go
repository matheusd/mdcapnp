// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

func (v *Vat) processBootstrap(ctx context.Context, rc *runningConn, msg Message) error {
	bootMsg := msg.AsBootstrap()
	reply := Message{
		isReturn: true,
		ret: Return{
			aid:       AnswerId(bootMsg.qid),
			isResults: true,
			pay: Payload{
				content: anyPointer{
					isCapPointer: true,
					cp:           capPointer{index: 0},
				},
				capTable: []CapDescriptor{
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

func (v *Vat) processReturn(ctx context.Context, rc *runningConn, ret Return) error {
	qid := QuestionId(ret.AnswerId())
	q, ok := rc.questions.get(qid)
	if !ok {
		return fmt.Errorf("question %d not found", qid)
	}

	// TODO: support exception, cancel, etc
	if !ret.IsResults() {
		return fmt.Errorf("only results supported")
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
	step := q.pipe.Step(q.stepIdx)
	if !step.stepDone.Set(stepResult) {
		// Can it ever be set twice on a return? I don't think so.
		return errors.New("question resolved twice")
	}

	return nil
}

var errCallWithoutId = errors.New("call with zero interfaceId and methodId")

func (v *Vat) processCall(ctx context.Context, rc *runningConn, c Call) error {
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
	reply := Message{
		isReturn: true,
		ret:      Return{aid: AnswerId(c.qid)},
	}
	crb := callReturnBuilder{ // Reuse on vat (this is running on the vat's main goroutine).
		payload: Payload{content: anyPointer{
			isVoid: true, // Void result by default on non-error.
		}},
	}

	rc.log.Trace().
		Int("qid", int(c.qid)).
		Int("eid", int(eid)).
		Msg("Locally handling call")

	err := exp.handler.Call(rc.ctx, callArgs, &crb)
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
func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg Message) error {
	var err error
	switch {
	case msg.IsBootstrap():
		err = v.processBootstrap(ctx, rc, msg)
	case msg.IsReturn():
		err = v.processReturn(ctx, rc, msg.AsReturn())
	case msg.isCall:
		err = v.processCall(ctx, rc, msg.AsCall())
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
func (v *Vat) prepareOutMessage(_ context.Context, pipe runningPipeline, stepIdx int) error {
	var ok bool

	step := &pipe.steps[stepIdx]
	if step.rpcMsg.IsBootstrap() {
		if step.qid, ok = step.step.conn.questions.nextID(); !ok {
			return errTooManyOpenQuestions
		}

		step.rpcMsg.boot.qid = step.qid
		step.step.conn.log.Debug().
			Int("qid", int(step.qid)).
			Msg("Prepared Bootstrap message")
	} else if step.rpcMsg.IsCall() {
		if step.qid, ok = step.step.conn.questions.nextID(); !ok {
			return errTooManyOpenQuestions
		}
		step.rpcMsg.call.qid = step.qid

		// Find the parent step. If it's in this pipeline, access it
		// directly. Otherwise, check whether the parent completed
		// already or not (to determine if this is call pipelined to a
		// an incomplete local call or to a remote promise).
		if stepIdx > 0 {
			parentStep := &pipe.steps[stepIdx-1]
			step.rpcMsg.call.target = MessageTarget{
				isPromisedAnswer: true,
				pans:             PromisedAnswer{qid: parentStep.qid},
			}

			step.step.conn.log.Debug().
				Int("qid", int(step.qid)).
				Int("pans", int(parentStep.qid)).
				Msg("Prepared call for in-pipeline promised answer")
		} else if pipe.pipe.parent == nil {
			// Should never happen, but avoid a panic below.
			return errors.New("non-bootstrap call without a parent pipeline")
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
			parentStep := pipe.pipe.parent.Step(pipe.pipe.parentStepIdx)
			if !parentStep.stepRunning.IsSet() {
				return errors.New("failed precondition: stepRunning should've been set")
			}
			parentQid := parentStep.stepRunning.Value()
			step.rpcMsg.call.target = MessageTarget{
				isPromisedAnswer: true,
				pans:             PromisedAnswer{qid: parentQid},
			}

			step.step.conn.log.Debug().
				Int("qid", int(step.qid)).
				Int("pans", int(parentQid)).
				Msg("Prepared call for parent pipeline promised answer")
		}
	}

	return nil
}

// commitOutMessage commits the changes of the pipeline step to the local Vat's
// state, under the assumption that the given pipeline step was successfully
// sent to the remote Vat.
func (v *Vat) commitOutMessage(_ context.Context, pipe runningPipeline, stepIdx int) error {
	step := &pipe.steps[stepIdx]
	if step.rpcMsg.isBootstrap {
		q := question{pipe: pipe.pipe, stepIdx: stepIdx}
		qid := pipe.steps[stepIdx].qid
		conn := pipe.steps[stepIdx].step.conn
		conn.questions.set(qid, q)
		step.step.conn.log.Debug().Int("qid", int(step.qid)).Msg("Comitted Bootstrap message")
	} else if step.rpcMsg.isCall {
		q := question{pipe: pipe.pipe, stepIdx: stepIdx}
		qid := pipe.steps[stepIdx].qid
		conn := pipe.steps[stepIdx].step.conn
		conn.questions.set(qid, q)
		step.step.conn.log.Debug().Int("qid", int(step.qid)).Msg("Comitted Call message")
	} else {
		// Guard against errors while developing.
		return errors.New("unimplemented commitment of message")
	}

	return nil
}
