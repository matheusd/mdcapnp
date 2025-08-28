// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
)

type conn interface {
	send(context.Context, *message) error
	receive(context.Context, *message) error

	// TODO: Allow conn-owned buffer (io_uring)?
	// usesReceiverBuffer() bool
	// receiveMsg(context.Context) (*message, error)
}

type runningConn struct {
	c   conn
	vat *vat

	outQueue chan *message

	ctx    context.Context
	cancel func() // Closes runningConn.
}

func (rc *runningConn) queue(ctx context.Context, msg *message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case rc.outQueue <- msg:
		return nil

	default:
		return errors.New("outbound queue is full")
	}
}

type _bootstrapCap struct{}
type bootstrapCap futureCap[_bootstrapCap]

func castBootstrap[T any](bc bootstrapCap) futureCap[T] {
	return futureCap[T]{pipe: bc.pipe, stepIndex: bc.stepIndex}
}

func (rc *runningConn) Bootstrap() bootstrapCap {
	res := bootstrapCap(newRootFutureCap[_bootstrapCap](defaultPipelineSizeHint))
	res.pipe.steps[0].conn = rc

	// TODO: what if bootstrap already resolved?

	return res
}
