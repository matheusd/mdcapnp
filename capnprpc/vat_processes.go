// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
)

func (v *Vat) processBootstrap(ctx context.Context, rc *runningConn, msg Message) error {
	// TODO: get the bootstrap capability from Vat.
	// TODO: send as reply to remote Vat.
	panic("fixme")
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

	// Go through cap table, modify exports table based on what was
	// exported by this call.
	//
	// TODO: only do this if the cap is referenced in the content?
	payload := ret.AsResults()
	capTable := payload.CapTable()
	for _, entry := range capTable {
		if !entry.IsSenderHosted() {
			return fmt.Errorf("only senderHosted capabilities supported")
		}
		eid := entry.AsSenderHosted()
		rc.exports.set(eid, export{typ: exportTypeSenderHosted})
		rc.log.Debug().Int("eid", int(eid)).Msg("Set export as senderHosted")
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
	} else {
		// TODO: copy if its a struct? Or release serialized message if
		// content is just a cap (because it's not needed anymore)?
		stepResult = content.AsStruct()
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

// processInMessage processes an incoming message from a remote Vat.
func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg Message) error {
	switch {
	case msg.IsBootstrap():
		return v.processBootstrap(ctx, rc, msg)
	case msg.IsReturn():
		return v.processReturn(ctx, rc, msg.AsReturn())
	default:
		return errors.New("unknown Message type")
	}
}

// prepareOutMessage prepares an outgoing Message message that is part of a
// pipeline to be sent to the remote Vat.
//
// Note: this does _not_ commit the changes to the conn's tables yet.
func (v *Vat) prepareOutMessage(_ context.Context, pipe runningPipeline, stepIdx int) error {
	// TODO: what about resolves, returns, etc?

	step := &pipe.steps[stepIdx]
	if step.rpcMsg.isBootstrap {
		var ok bool
		step.qid, ok = step.step.conn.questions.nextID()
		if !ok {
			return errors.New("too many open questions")
		}

		step.rpcMsg.boot.qid = step.qid
		step.step.conn.log.Debug().Int("qid", int(step.qid)).Msg("Prepared Bootstrap() message")
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
		step.step.conn.log.Debug().Int("qid", int(step.qid)).Msg("Comitted Bootstrap() message")
	}

	return nil
}
