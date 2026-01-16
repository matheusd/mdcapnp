// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"

	"matheusd.com/mdcapnp/capnpser"
)

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
	return waitReturnResultsCapability(ctx, CallFuture(bc))
}

// Boostrappable is a helper interface that matches both [RemoteVat] and
// [Level0ClientVat] to call their respective Bootstrap() methods.
type Boostrappable interface {
	Bootstrap() BootstrapFuture
}

// WaitBootstrap is a helper to call the Bootstrap() method of a vat, wait for
// it complete, and return its callable capability client caset to a specific
// type.
func WaitBootstrap[T ~CallFutureType](ctx context.Context, b Boostrappable) (bootCap T, cc CapabilityClient, err error) {
	f := b.Bootstrap()
	bootCap = T(f)
	cc, err = f.Wait(ctx)
	return
}

type VoidFuture CallFuture

func (fv VoidFuture) Wait(ctx context.Context) error {
	return WaitVoidReturnResults(ctx, CallFuture(fv))
}

type StringFuture CallFuture

func (sf StringFuture) waitStringList(ctx context.Context) (capnpser.List, ReturnResults, error) {
	l, rr, err := waitReturnResultsList(ctx, CallFuture(sf))
	if err != nil {
		return l, rr, err
	}
	if !l.IsElementSizeByte() {
		rr.Release()
		return l, rr, errors.New("returned list is not a 1-byte-per-element list")
	}
	return l, rr, nil
}

// WaitDiscardResult waits for the call to complete but discards the result.
func (sf StringFuture) WaitDiscardResult(ctx context.Context) error {
	_, rr, err := sf.waitStringList(ctx)
	rr.Release() // Discard result.
	return err
}

// Wait waits for the call that produces this future to complete as a string.
func (sf StringFuture) Wait(ctx context.Context) (string, error) {
	l, rr, err := sf.waitStringList(ctx)
	if err != nil {
		return "", err
	}

	var res string
	if l.LenBytes() > 1 { // LenBytes() == 1 is only null char.
		res = l.String()
	}

	// Ok to release because string() uses the standard runtime method of
	// making a copy.
	rr.Release()
	return res, nil
}

// WaitUnsafe waits for the call that produces this future to complete as an
// unsafe string.
//
// This string is aliased to the contents of the incoming message, which means
// it is only valid for use while the message itself is valid (that is, until
// the [ReturnResults.Release] method is called).
//
// NOTE: Great care should be taken with these strings. Trying to hold on to
// them after the message is released, even indirectly (for example, in a
// closure) is likely to cause serious issues (up to and including memory
// corruption and security issues).
func (sf StringFuture) WaitUnsafe(ctx context.Context) (string, ReturnResults, error) {
	l, rr, err := sf.waitStringList(ctx)
	if err != nil {
		return "", rr, err
	}

	if err := l.CheckCanGetUnsafeString(); err != nil {
		return "", rr, err
	}

	var res string
	if l.LenBytes() > 1 { // LenBytes() == 1 is only null char.
		res = l.UnsafeString()
	}

	return res, rr, nil

}
