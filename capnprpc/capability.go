// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "context"

type CapabilityType = struct {
	eid ExportId // Is this right???
}

type capability CapabilityType

type BootstrapFuture CallFuture

func (bc BootstrapFuture) Wait(ctx context.Context) (capability, error) {
	// return CastCallResultOrErr[capability](WaitReturn(ctx, CallFuture(bc)))
	return WaitReturnResultsCapability[capability](ctx, CallFuture(bc))
}

type VoidFuture CallFuture

func (fv VoidFuture) Wait(ctx context.Context) error {
	_, err := WaitReturn(ctx, CallFuture(fv))
	return err
}
