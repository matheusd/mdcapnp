// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
)

type msgBatch struct {
	isSingle bool
	single   message
	msgs     []message
}

func singleMsgBatch(msg message) msgBatch {
	return msgBatch{isSingle: true, single: msg}
}

type conn interface {
	send(context.Context, msgBatch) error
	receive(context.Context, *message) error
	remoteName() string

	// TODO: Allow conn-owned buffer (io_uring)?
	// usesReceiverBuffer() bool
	// receiveMsg(context.Context) (*message, error)
}

var errConnDone = errors.New("conn is done")

// runningConn is a connection that is running to another vat.
type runningConn struct {
	// Design note: most of the fields are only meant to be accessed from
	// within a vat's runStep() call. They are not safe for concurrent
	// access from within client code.
	//
	// TODO: maybe convert the public runningConn into a handle instead of
	// pointer?

	c   conn
	vat *Vat
	log zerolog.Logger

	boot bootstrapCap

	// bootExportId is the export id of the bootstrap cap offered by the vat
	// on this conn.
	bootExportId ExportId

	outQueue chan msgBatch

	// TODO: question and export IDs are set by local vat, answer and import
	// ids are set by the remote vat. Split table type into two
	// (incoming/outgoing table) to protect from remote misuse and restrict
	// API.
	questions table[QuestionId, question]
	answers   table[AnswerId, answer]
	imports   table[ImportId, imprt]
	exports   table[ExportId, export]

	ctx    context.Context
	cancel func(error) // Closes runningConn.
}

func (rc *runningConn) queue(ctx context.Context, batch msgBatch) error {
	rc.log.Trace().Int("len", len(batch.msgs)).Msg("Queueing batch of outgoing messages")
	select {
	case <-ctx.Done():
		return context.Cause(ctx)

	case rc.outQueue <- batch:
		rc.log.Trace().Int("len", len(batch.msgs)).Msg("Queued batch of outgoing messages")
		return nil

	default:
		return errors.New("outbound queue is full")
	}
}

func newRunningConn(c conn, v *Vat) *runningConn {
	log := v.log.With().Str("remote", c.remoteName()).Logger()

	rc := &runningConn{
		c:   c,
		vat: v,
		log: log,

		boot: bootstrapCap(newRootFutureCap[capability](1)),

		outQueue:  make(chan msgBatch, 1000), // TODO: Parametrize buffer size.
		questions: makeTable[QuestionId, question](),
		answers:   makeTable[AnswerId, answer](),
		imports:   makeTable[ImportId, imprt](),
		exports:   makeTable[ExportId, export](),
	}

	// TODO: prepare boot message.
	rc.boot.pipe.vat = v
	rc.boot.pipe.steps[0].conn = rc
	rc.boot.pipe.state = pipelineStateBuilt

	return rc
}

type bootstrapCap futureCap[capability]

func (bc bootstrapCap) Wait(ctx context.Context) (capability, error) {
	return waitResult(ctx, futureCap[capability](bc))
}

func castBootstrap[T any](bc bootstrapCap) futureCap[T] {
	return futureCap[T]{pipe: bc.pipe, stepIndex: bc.stepIndex}
}

func (rc *runningConn) Bootstrap() bootstrapCap {
	return rc.boot // Any calls fork the pipeline.
}
