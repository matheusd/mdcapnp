// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "context"

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

	ctx    context.Context
	cancel func() // Closes runningConn.
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
