// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"

	"matheusd.com/mdcapnp/capnpser"
)

func (v *vat) processBootstrap(ctx context.Context, rc *runningConn, msg Message) error {
	// TODO: get the bootstrap capability from vat.
	// TODO: send as reply to remote vat.
	panic("fixme")
}

func (v *vat) processReturn(ctx context.Context, rc *runningConn, ret Return) error {
	qid := QuestionId(ret.AnswerId())
	q, ok := rc.questions.get(qid)
	if !ok {
		return fmt.Errorf("question %d not found", qid)
	}

	// TODO: go to pipeline item and fulfill it.
	_ = q

	panic("fixme")
}

// processInMessage processes an incoming message from a remote vat.
func (v *vat) processInMessage(ctx context.Context, rc *runningConn, serMsg capnpser.Message) error {
	var msg Message
	if err := msg.ReadFromRoot(&serMsg); err != nil {
		return err
	}

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
// pipeline to be sent to the remote vat.
//
// Note: this does _not_ commit the changes to the conn's tables yet.
func (v *vat) prepareOutMessage(ctx context.Context, pipe *pipeline, step *runningPipelineStep) error {
	var ok bool
	step.qid, ok = step.step.conn.questions.nextID()
	if !ok {
		return errors.New("too many open questions")
	}

	return nil
}

// commitOutMessage commits the changes of the pipeline step to the local vat's
// state, under the assumption that the given pipeline step was successfully
// sent to the remote vat.
func (v *vat) commitOutMessage(ctx context.Context, pipe *pipeline, step *runningPipelineStep) error {
	q := question{pipe: pipe /*, stepIdx: stepIdx*/}
	step.step.conn.questions.set(step.qid, q)
	panic("boo")
}
