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
	// TODO: split InterfaceId, MethodId and WantShallowReturnCopy into one
	// struct, then bundle with callOutMsg in another.

	callOutMsg rpcCallMsgBuilder // !!! Experimental

	InterfaceId InterfaceId
	MethodId    MethodId
}

type ReturnResults struct {
	inMsg *InMsg
	vat   vatImpl
}

func (rr *ReturnResults) Release() {
	if rr.inMsg != nil && rr.inMsg.pool != nil {
		rr.inMsg.pool.put(rr.inMsg)
		rr.inMsg = nil
	}
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
	exportNewCap(capServer CallHandler, parent CallFuture) (ExportId, error)
}

func (v *Vat) msgBuilderPool() *MessageBuilderPool             { return v.mbp }
func (v *Level0ClientVat) msgBuilderPool() *MessageBuilderPool { return v.mbp }
func (v *Level0ClientVat) exportNewCap(capServer CallHandler, parent CallFuture) (ExportId, error) {
	panic("boo")
}

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

type CallFutureType = struct {
	step *pipelineStep
	vat  vatImpl
}

type CallFuture CallFutureType

func setupCall(parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId) (CallSetup, types.PayloadBuilder) {
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

// SetupCallNoParams prepares to make a call that will not send any parameters.
func SetupCallNoParams(parent CallFuture, iid InterfaceId, mid MethodId) CallSetup {
	cs, _ := setupCall(parent, 0, iid, mid)
	return cs
}

// SetupCallWithParams prepares to make a call that will send parameters which
// are likely to take only up to an additional [payloadSizeHint] words as total
// size.
func SetupCallWithParams(parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId) (CallSetup, types.PayloadBuilder) {
	return setupCall(parent, payloadSizeHint, iid, mid)
}

// SetupCallWithStructParams prepares to make a call that will send a struct as
// parameters, with size [paramsSize] and (likely) only up to an additional
// [payloadSizeHint] words as extra size.
func SetupCallWithStructParams(parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId, paramsSize capnpser.StructSize) (CallSetup, capnpser.StructBuilder) {
	cs, pbuilder := SetupCallWithParams(parent, payloadSizeHint, iid, mid)
	sb, err := pbuilder.SetContentAsNewStruct(paramsSize)
	if err != nil {
		panic(err)
	}
	return cs, sb
}

// SetupCallWithStructParamsGeneric is similar to [SetupCallWithParams] but
// allows returning a generic StructBuilder to be used for building the request
// parameters.
func SetupCallWithStructParamsGeneric[T ~capnpser.StructBuilderType](parent CallFuture, payloadSizeHint capnpser.WordCount, iid InterfaceId, mid MethodId, paramsSize capnpser.StructSize) (CallSetup, T) {
	cs, pbuilder := SetupCallWithParams(parent, payloadSizeHint, iid, mid)
	sb, err := pbuilder.SetContentAsNewStruct(paramsSize)
	if err != nil {
		panic(err)
	}
	return cs, T(sb)
}

// SetupCallWithNewCap prepares to make a call that will export a new capability
// to the remote end.
//
// TODO: This might need to return a reference to the exported cap (so that *the
// same* cap can be sent to someone else, referenced when processing some other
// call made by the remote, etc).
func SetupCallWithNewCap(parent CallFuture, iid InterfaceId, mid MethodId, capServer CallHandler) CallSetup {
	// TODO: this should be tracked inside CallSetup to avoid possible
	// deadlock with rc.
	eid, err := parent.vat.exportNewCap(capServer, parent)
	if err != nil {
		panic(err)
	}
	cs, pbuilder := setupCall(parent, 0, iid, mid)
	ctbuilder, err := pbuilder.NewCapTable(1, 1)
	if err != nil {
		panic(err)
	}
	cdbuilder := ctbuilder.At(0)
	if err := cdbuilder.SetSenderHosted(eid); err != nil {
		panic(err)
	}
	pbuilder.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0))
	return cs
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
	// WaitReturn will directly access the vat to send the Call (based on
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

// waitReturn waits until the passed future is resolved. Returns either the
// resolved value or an error.
func waitReturn(ctx context.Context, cf CallFuture) (res any, err error) {
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

// WaitVoidReturnResults waits for the return of a call that is meant to return
// a void valule.
func WaitVoidReturnResults(ctx context.Context, cf CallFuture) (err error) {
	resAny, err := waitReturn(ctx, cf)
	if err != nil {
		return
	}

	switch resAny.(type) {
	case struct{}:
		return nil
	default:
		return fmt.Errorf("unexpected non-void return type %T", resAny)
	}
}

// waitReturnResultsContent waits for the return of a call that is meant to
// return a non-void Results content.
func waitReturnResultsContent(ctx context.Context, cf CallFuture) (content capnpser.AnyPointer, rr ReturnResults, err error) {
	var ok bool
	var resAny any

	resAny, err = waitReturn(ctx, cf)
	if err != nil {
		return
	}

	rr.vat = cf.vat
	if rr.inMsg, ok = resAny.(*InMsg); !ok {
		err = fmt.Errorf("unhandled result type (got %T instead of *InMsg)", resAny)
		return
	}

	// TODO: A few of these checks are now redundant...

	resRoot, err := rr.inMsg.Msg.GetRoot()
	if err != nil {
		err = fmt.Errorf("unable to get root: %v", err)
		return
	}
	resMsg := types.Message(resRoot)
	if resMsg.Which() != types.Message_Which_Return {
		err = errors.New("not the expected Return")
		return
	}
	resReturn, err := resMsg.AsReturn()
	if err != nil {
		return
	}
	resPayload, err := resReturn.AsResults()
	if err != nil {
		return
	}
	content, err = resPayload.Content()
	if err != nil {
		return
	}
	return
}

// waitReturnResultsStruct waits for the return of a call that is expected to
// produce a struct as result.
func waitReturnResultsStruct(ctx context.Context, cf CallFuture) (s capnpser.Struct, rr ReturnResults, err error) {
	var content capnpser.AnyPointer
	content, rr, err = waitReturnResultsContent(ctx, cf)
	if err != nil {
		return
	}

	if !content.IsStruct() {
		err = errors.New("Return.Results.Content is not the expected Struct")
		return
	}

	s = content.AsStruct()
	return
}

// WaitReturnResultsStruct waits for the return of a call that is expected to
// return a struct of type T.
func WaitReturnResultsStruct[T ~capnpser.StructType](ctx context.Context, cf CallFuture) (res T, rr ReturnResults, err error) {
	var s capnpser.Struct
	s, rr, err = waitReturnResultsStruct(ctx, cf)
	res = T(s)
	return
}

// waitReturnResultsList waits for the return of a call that is expected to
// produce a list as result.
func waitReturnResultsList(ctx context.Context, cf CallFuture) (l capnpser.List, rr ReturnResults, err error) {
	var content capnpser.AnyPointer
	content, rr, err = waitReturnResultsContent(ctx, cf)
	if err != nil {
		return
	}

	if !content.IsList() {
		err = errors.New("Return.Results.Content is not the expected List")
		return
	}

	l = content.AsList()
	return
}

// WaitReturnResultsList waits for the return of a call that is expected to
// produce a list of some particular type as result.
func WaitReturnResultsList[T ~capnpser.ListType](ctx context.Context, cf CallFuture) (l T, rr ReturnResults, err error) {
	var ll capnpser.List
	ll, rr, err = waitReturnResultsList(ctx, cf)
	l = T(ll)
	return
}

// waitReturnResultsCapability waits for the return of a remote call that should
// return a capability client.
func waitReturnResultsCapability(ctx context.Context, cf CallFuture) (cc CapabilityClient, err error) {
	var resAny any
	resAny, err = waitReturn(ctx, cf)
	if err != nil {
		return
	}

	var ok bool
	cc, ok = resAny.(CapabilityClient)
	if !ok {
		err = fmt.Errorf("result was not a capability (got %T)", resAny)
		return
	}
	return
}

// WaitReturnResultsCapability waits for the return of a remote call that should
// return a capability client.
func WaitReturnResultsCapability[T ~CapabilityClientType](ctx context.Context, cf CallFuture) (res T, err error) {
	cc, err := waitReturnResultsCapability(ctx, cf)
	res = T(cc)
	return
}
