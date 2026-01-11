// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "context"

type CapabilityClientType = struct {
	eid ExportId // Is this right???
}

type CapabilityClient CapabilityClientType

// CallHandler is the lowest level handler for calls (including bootstrap).
type CallHandler interface {
	Call(ctx context.Context, cc *CallContext) error
}

type CallHandlerFunc func(ctx context.Context, rb *CallContext) error

func (f CallHandlerFunc) Call(ctx context.Context, rb *CallContext) error {
	return f(ctx, rb)
}

type BootstrapFuture CallFuture

func (bc BootstrapFuture) Wait(ctx context.Context) (CapabilityClient, error) {
	// return CastCallResultOrErr[capability](WaitReturn(ctx, CallFuture(bc)))
	return WaitReturnResultsCapability[CapabilityClient](ctx, CallFuture(bc))
}

type VoidFuture CallFuture

func (fv VoidFuture) Wait(ctx context.Context) error {
	_, err := WaitReturn(ctx, CallFuture(fv))
	return err
}
