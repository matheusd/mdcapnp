// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
	"matheusd.com/mdcapnp/internal/sigvalue"
)

// CallSetup are the requirements to build and send a Call message.
type CallSetup struct {
	callOutMsg rpcCallMsgBuilder // !!! Experimental

	InterfaceId InterfaceId
	MethodId    MethodId

	// ParamsBuilder is called with the outbound message builder when the
	// Call message is being built. This can be used to fill the Params
	// field of the Call.
	ParamsBuilder CallParamsBuilder

	// ResultsParser is called when a Return.Results is received in response
	// to a Call. It can parse the encoded results into some Go value.
	ResultsParser CallResultsParser

	copyReturnResults bool

	WantReturnResults bool
}

type ReturnResults struct {
	content capnpser.AnyPointer
	copyMb  *capnpser.MessageBuilder
	vat     vatImpl
}

func (rr *ReturnResults) Release() {
	if rr.copyMb != nil {
		rr.vat.msgBuilderPool().put(rr.copyMb)
	}
}

func (rr *ReturnResults) Content() capnpser.AnyPointer {
	return rr.content
}

func ReturnResultsStruct[T ~capnpser.StructType](rr ReturnResults) T {
	if !rr.content.IsStruct() {
		panic("ReturnResults content is not a struct")
	}
	return T(rr.content.AsStruct())
}

func NewCallParamsStruct[T ~capnpser.StructBuilderType](payload types.PayloadBuilder, size capnpser.StructSize) (T, error) {
	res, err := payload.SetContentAsNewStruct(size)
	return T(res), err
}

func ResultsStruct[T ~capnpser.StructType](payload types.Payload) (res T, err error) {
	var content capnpser.AnyPointer
	content, err = payload.Content()
	if err != nil {
		return
	}
	if !content.IsStruct() {
		err = errors.New("payload contents is not a struct in results")
		return
	}
	res = T(content.AsStruct())
	return
}

type pipelineStepState int

const (
	pipeStepStateBuilding pipelineStepState = iota
	pipeStepStateSending
	pipeStepStateRunning
	pipeStepStateDone
	pipeStepStateFinished
	pipeStepFailed // Must be last to be > all other states.
)

type pipelineStepStateValue struct {
	qid   QuestionId   // Set if step state is >= running.
	iid   ImportId     // Set if this has been resolved into an import (either hosted or promise)
	value any          // Set if step state is >= done.
	err   error        // Set if step state is >= failed.
	conn  *runningConn // May be changed as step is resolved (3PH).
}

type pipelineStep struct {
	value  sigvalue.Stateful[pipelineStepState, pipelineStepStateValue]
	parent *pipelineStep
	csetup CallSetup
}

func finalizePipelineStep(step *pipelineStep) {
	var qid QuestionId
	var conn *runningConn
	_ = step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepFailed {
			qid = ov.qid
			ov.qid = 0
			os = pipeStepStateFinished
			conn = ov.conn
		}

		return os, ov, nil
	})

	if qid > 0 {
		conn.cleanupQuestionIdDueToUnref(qid) // TODO: How does this conflict with other uses?
	}
}

func newChildStep(parent *pipelineStep) *pipelineStep {
	return &pipelineStep{
		parent: parent,
	}
}

type vatImpl interface {
	msgBuilderPool() *MessageBuilderPool
}

func (v *Vat) msgBuilderPool() *MessageBuilderPool             { return v.mbp }
func (v *Level0ClientVat) msgBuilderPool() *MessageBuilderPool { return v.mbp }

func (v *Vat) getCallMessageBuilder(payloadSizeHint capnpser.WordCount) *capnpser.MessageBuilder {
	return v.mbp.getRawMessageBuilder(payloadSizeHint)
}

func (v *Vat) putCallMessageBuilder(mb *capnpser.MessageBuilder) {
	v.mbp.put(mb)
}

func (v *Level0ClientVat) getCallMessageBuilder(payloadSizeHint capnpser.WordCount) *capnpser.MessageBuilder {
	panic("todo")
}

func (v *Level0ClientVat) putCallMessageBuilder(mb *capnpser.MessageBuilder) {
	panic("todo")
}

type CallFuture struct {
	step *pipelineStep
	vat  vatImpl
}

func setupCall(parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId, wantParams bool) (CallSetup, types.PayloadBuilder) {
	serMb := parent.vat.msgBuilderPool().getRawMessageBuilder(callMessageSizeOverhead + payloadSizeHint)
	mb, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		panic(err)
	}

	call, err := mb.NewCall()
	if err != nil {
		panic(err)
	}
	call.SetInterfaceId(uint64(iid))
	call.SetMethodId(uint16(mid))
	pbuilder, err := call.NewParams()
	if err != nil {
		panic(err)
	}

	return CallSetup{
		InterfaceId: iid,
		MethodId:    mid,
		callOutMsg: rpcCallMsgBuilder{
			outMsg:      outMsg{serMsg: serMb, mb: mb},
			isBootstrap: false,
			builder:     capnpser.StructBuilder(call),
		},
	}, pbuilder
}

func SetupCallNoParams(parent CallFuture, iid InterfaceId, mid MethodId) CallSetup {
	cs, _ := setupCall(parent, 0, iid, mid, false)
	return cs
}

func SetupCallWithParams(parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId) (CallSetup, types.PayloadBuilder) {
	return setupCall(parent, payloadSizeHint, iid, mid, true)
}

func SetupCallWithStructParams[T ~capnpser.StructBuilderType](parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId, paramsSize capnpser.StructSize) (CallSetup, T) {
	cs, pbuilder := SetupCallWithParams(parent, payloadSizeHint, iid, mid)
	sb, err := pbuilder.SetContentAsNewStruct(paramsSize)
	if err != nil {
		panic(err)
	}
	return cs, T(sb)
}

// FIXME: rename or remove.
func newRootFutureCap(v vatImpl) CallFuture {
	return CallFuture{
		step: &pipelineStep{},
		vat:  v,
	}
}

func RemoteCall(parent CallFuture, csetup CallSetup) (res CallFuture) {
	if parent.step == nil {
		panic("root futures cannot be built by RemoteCall (nil parent step)")
	}
	if parent.vat == nil {
		panic("trying to issue RemoteCall without vat")
	}
	parentStep := parent.step

	// Level 0 vats only perform sync calls one at a time and don't support
	// pipelining, therefore tracking of pipeline steps isn't needed.
	//
	// WaitResult will directly access the vat to send the Call (based on
	// the contents of csetup) and will wait for the corresponding Return.
	if l0vat, isLevel0 := parent.vat.(*Level0ClientVat); isLevel0 {
		l0vat.setupNextCall(csetup)
		res.vat = l0vat
		return
	}

	res.step = newChildStep(parentStep)
	res.step.csetup = csetup
	res.vat = parent.vat

	return res
}

func CastCallResult[T any](callResult any) (res T, err error) {
	// Check if result is already the expected return type.
	var ok bool
	if res, ok = callResult.(T); ok {
		return
	}

	// Not an error, and not type T, so must be an AnyPointer with the
	// Content field of a Payload result.
	var resAnyPointer capnpser.AnyPointer
	if resAnyPointer, ok = callResult.(capnpser.AnyPointer); !ok {
		err = fmt.Errorf("future expected AnyPointer but got %T", callResult)
		return
	}

	// Content may be either a Struct or a CapPointer, and T will be an
	// alias to one of these (depending on what's expected based on the
	// schema).
	//
	// TODO: better way to convert to T?
	var contentAny any
	if resAnyPointer.IsStruct() {
		contentAny = resAnyPointer.AsStruct()
	} else if resAnyPointer.IsCapPointer() {
		// TODO: How to convert capPointer to importId???
		contentAny = resAnyPointer.AsCapPointer()
	} else {
		err = fmt.Errorf("content is not struct or capPointer")
		return
	}

	res, ok = contentAny.(T)
	if !ok {
		err = fmt.Errorf("future expected %T but got back %T", res, contentAny)
	}

	return
}

func CastCallResultOrErr[T any](callResult any, inErr error) (res T, err error) {
	if inErr != nil {
		err = inErr
	} else {
		res, err = CastCallResult[T](callResult)
	}
	return
}

func WaitReturn(ctx context.Context, cf CallFuture) (res any, err error) {
	if cf.vat == nil {
		err = errors.New("could not find vat for pipeline")
		return
	}

	// Level 0 vats only perform calls in sequence.
	if l0vat, isL0Vat := cf.vat.(*Level0ClientVat); isL0Vat {
		return l0vat.execNextCall(ctx)
	}

	// Vat has level > 0 (standard *Vat).
	vat := cf.vat.(*Vat)

	// Run pipeline.
	if err = vat.execStep(ctx, cf.step); err != nil {
		return
	}

	// Wait until the required step of the pipeline completes or fails.
	stepState, stepValue, err := cf.step.value.WaitStateAtLeast(ctx, pipeStepStateDone)
	if err != nil {
		return
	}
	if stepState == pipeStepFailed {
		err = stepValue.err
		if err == nil {
			err = fmt.Errorf("unknown error in failed pipeline step state: %v", stepState)
		} else {
			err = fmt.Errorf("step resolved to failure: %w", err)
		}
		return
	}

	// Determine if the result was an error.
	pipeRes := stepValue.value
	var ok bool
	if err, ok = pipeRes.(error); ok && err != nil {
		return
	}

	res = pipeRes
	return
}

func WaitReturnResults(ctx context.Context, cf CallFuture) (res ReturnResults, err error) {
	resAny, err := WaitReturn(ctx, cf)
	if err != nil {
		return ReturnResults{}, fmt.Errorf("WaitResult() errored: %v", err)
	}

	var ok bool
	res.copyMb, ok = resAny.(*capnpser.MessageBuilder)
	if !ok {
		err = fmt.Errorf("result was not a *capnpser.MessageBuilder (got %T)", resAny)
		return
	}

	resReader := res.copyMb.MessageReader()
	resRoot, err := resReader.GetRoot()
	if err != nil {
		return
	}
	res.vat = cf.vat
	res.content = resRoot.AsAnyPointer()
	return
}

func WaitReturnResultsStruct[T ~capnpser.StructType](ctx context.Context, cf CallFuture) (res T, rr ReturnResults, err error) {
	rr, err = WaitReturnResults(ctx, cf)
	if err != nil {
		return
	}
	if !rr.content.IsStruct() {
		err = errors.New("Return.Results.Content is not an expected Struct")
		return
	}
	res = T(rr.content.AsStruct())
	return
}

func WaitReturnResultsCapability[T ~CapabilityType](ctx context.Context, cf CallFuture) (res T, err error) {
	resAny, err := WaitReturn(ctx, cf)
	if err != nil {
		return
	}

	resCap, ok := resAny.(capability)
	if !ok {
		err = fmt.Errorf("result was not a capability (got %T)", resAny)
		return
	}

	res = T(resCap)
	return
}

func releaseStepResultMsgBuilder(cf CallFuture) {
	step := cf.step
	step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		// TODO: make assertions on type, state, etc
		anyp := ov.value.(capnpser.AnyPointerBuilder)
		vat := ov.conn.vat
		vat.mbp.put(anyp.MsgBuilder())
		ov.value = nil
		return os, ov, nil
	})
}

func releaseFuture(ctx context.Context, cap CallFuture) {
	finalizePipelineStep(cap.step)
}
