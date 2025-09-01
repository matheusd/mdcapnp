// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"

	"matheusd.com/mdcapnp/capnpser"
)

type msgBatch struct {
	// TODO: add a `first Message` and use it when only a single message?
	msgs []capnpser.Message
}

type conn interface {
	send(context.Context, msgBatch) error
	receive(context.Context, capnpser.Message) error

	// TODO: Allow conn-owned buffer (io_uring)?
	// usesReceiverBuffer() bool
	// receiveMsg(context.Context) (*message, error)
}

// runningConn is a connection that is running to another vat.
type runningConn struct {
	// Design note: most of the fields are only meant to be accessed from
	// within a vat's runStep() call. They are not safe for concurrent
	// access from within client code.
	//
	// TODO: maybe convert the public runningConn into a handle instead of
	// pointer?

	c   conn
	vat *vat

	boot bootstrapCap

	outQueue chan msgBatch

	questions table[QuestionId, question]
	answers   table[AnswerId, answer]
	imports   table[ImportId, imprt]
	exports   table[ExportId, export]

	ctx    context.Context
	cancel func() // Closes runningConn.
}

func (rc *runningConn) queue(ctx context.Context, batch msgBatch) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case rc.outQueue <- batch:
		return nil

	default:
		return errors.New("outbound queue is full")
	}
}

func newRunningConn(c conn, v *vat) *runningConn {
	// TODO: prepare boot message.
	boot := bootstrapCap(newRootFutureCap[_bootstrapCap](1))

	rc := &runningConn{
		c:   c,
		vat: v,

		boot: boot,

		outQueue:  make(chan msgBatch, 1000), // TODO: Parametrize buffer size.
		questions: makeTable[QuestionId, question](),
		answers:   makeTable[AnswerId, answer](),
		imports:   makeTable[ImportId, imprt](),
		exports:   makeTable[ExportId, export](),
	}

	return rc
}

type _bootstrapCap struct{}
type bootstrapCap futureCap[_bootstrapCap]

func castBootstrap[T any](bc bootstrapCap) futureCap[T] {
	return futureCap[T]{pipe: bc.pipe, stepIndex: bc.stepIndex}
}

func (rc *runningConn) Bootstrap() bootstrapCap {
	// Fork the root bootstrap future into a new pipeline.
	return bootstrapCap(forkFuture[_bootstrapCap](futureCap[_bootstrapCap](rc.boot), defaultPipelineSizeHint))
}
