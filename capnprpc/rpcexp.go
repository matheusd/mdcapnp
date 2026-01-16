// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weak"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

type InterfaceId uint64
type MethodId uint16

type messageTarget struct {
	impCap  ImportId
	pansQid QuestionId

	isImportedCap    bool
	isPromisedAnswer bool
}

type introductionInfo struct {
	sendToRecipientAlt capnpser.AnyPointer
	sendToTargetAlt    capnpser.AnyPointer
}

type rpcCallMsgBuilder struct {
	outMsg
	builder       capnpser.StructBuilder
	paramsBuilder capnpser.StructBuilder
	isBootstrap   bool
}

func (rb *rpcCallMsgBuilder) initAsBootstrap() error {
	b, err := rb.mb.NewBoostrap()
	if err != nil {
		return err
	}
	rb.builder = capnpser.StructBuilder(b)
	return nil
}

func (rb *rpcCallMsgBuilder) bootstrapBuilder() types.BootstrapBuilder {
	return types.BootstrapBuilder(rb.builder)
}

func (rb *rpcCallMsgBuilder) callBuilder() types.CallBuilder {
	return types.CallBuilder(rb.builder)
}

type interfaceId uint64
type methodId uint16

type answerPromise struct {
	rc  *runningConn
	eid ExportId
}

func (ap answerPromise) resolveToHandler(handler CallHandler) error {
	var err error

	// resolution is the new export ap is resolving to.
	resolution := export{
		typ:      exportTypeLocallyHosted,
		handler:  handler,
		refCount: 1,
	}

	ap.rc.mu.Lock()
	exp, ok := ap.rc.exports.get(ap.eid)
	if !ok {
		fmt.Println(ap.rc.vat.cfg.name, ap.rc.exports)
		err = fmt.Errorf("export %d not found", ap.eid)
	} else if exp.resolved() {
		err = fmt.Errorf("promised export %d already resolved to %s",
			ap.eid, exp.resolveString())
	} else if resolveEid, ok := ap.rc.exports.nextID(); !ok {
		err = fmt.Errorf("could not obtain new export id for resolution of %d", ap.eid)
	} else {
		exp.resolvedToExport = resolveEid
		ap.rc.exports.set(ap.eid, exp)
		ap.rc.exports.set(resolveEid, resolution)

		ap.rc.log.Trace().
			Int("eid", int(ap.eid)).
			Int("resEid", int(resolveEid)).
			Str("rtyp", resolution.typ.String()).
			Msg("Resolving previously exported promise")

		err = ap.rc.vat.queueResolve(ap.rc.ctx, ap.rc, ap.eid, exp, resolution)
	}
	ap.rc.mu.Unlock()

	return err
}

func (ap answerPromise) resolveToThirdPartyImport(tpRc *runningConn, tpIid ImportId) error {
	var err error

	ap.rc.log.Debug().
		Int("eid", int(ap.eid)).
		Str("src", ap.rc.String()).
		Str("target", tpRc.String()).
		Int("targetIid", int(tpIid)).
		Msg("Resolving promised answer to third party")

	vat := ap.rc.vat
	if tpRc.vat != vat {
		return errors.New("cannot solve to third party in different local vat")
	}

	// TODO: Determine this by calling rc.c.introduceTo(ap.rc.c) to get an
	// IntroductionInfo.
	iinfo, err := vat.getNetworkIntroduction(ap.rc, tpRc)
	if err != nil {
		return err
	}

	// Send Provide to remote.
	rpcProv, prov, err := vat.newProvide(iinfo.sendToTargetAlt)
	if err != nil {
		return err
	}
	target, err := prov.NewTarget()
	if err != nil {
		return err
	}

	if err := target.SetImportedCap(tpIid); err != nil {
		return err
	}

	var ok bool
	var provideQid QuestionId
	tpRc.mu.Lock()
	if !tpRc.imports.has(tpIid) {
		err = fmt.Errorf("third party %s does not have import %d", tpRc, tpIid)
	} else if provideQid, ok = tpRc.questions.nextID(); !ok {
		err = fmt.Errorf("could not generate new question id for %s", tpRc)
	} else {
		// provide.qid = qid
		if err := prov.SetQuestionId(provideQid); err != nil {
			return err
		}
		tpRc.questions.set(provideQid, question{typ: questionTypeProvide})
	}
	tpRc.mu.Unlock()
	if err != nil {
		return err
	}

	// CHECK: ok to do here, outside tpRc.mu.Lock()?
	if err := tpRc.vat.sendProvide(tpRc.ctx, tpRc, rpcProv, prov); err != nil {
		return err
	}

	// vine is the export that proxies calls to rc[iid].
	//
	// TODO: skip if both rc and ap.rc are level 3?
	vine := export{
		typ:           exportTypeThirdPartyExport,
		thirdPartyRC:  tpRc,
		thirdPartyIid: tpIid,
		refCount:      1,
	}

	// Warning: There's a potential race condition in the protocol here.
	// We're expected to send a Resolve to the original caller ("Alice"),
	// assuming that the third party ("Carol") has fully processed the
	// Provide sent above and is ready to receive the Accept that the caller
	// will send in response to the Resolve.
	//
	// On regular network and processing conditions and under the assumption
	// that the network layer implementation uses TCP, which acks the
	// received bytes, this will be true.
	//
	// However, in tests where processing is very fast or if Carol is
	// congested, this may not be true. In that case, it may be necessary to
	// introduce some delay in this point.
	if vat.cfg.delayResolveIn3PH > 0 {
		startSleep := time.Now()
		select {
		case endSleep := <-time.After(vat.cfg.delayResolveIn3PH):
			ap.rc.log.Trace().
				Dur("dur", endSleep.Sub(startSleep)).
				Msg("Slept to delay Resolve after Provide in 3PH resolution")
		case <-ap.rc.ctx.Done():
			return ap.rc.ctx.Err()
		}
	}

	// Send resolve back to caller, pointing to the third party.
	ap.rc.mu.Lock()
	exp, ok := ap.rc.exports.get(ap.eid)
	if !ok {
		err = fmt.Errorf("export %d not found", ap.eid)
	} else if exp.resolved() {
		err = fmt.Errorf("promised export %d already resolved to %s",
			ap.eid, exp.resolveString())
	} else if vineEid, ok := ap.rc.exports.nextID(); !ok {
		err = fmt.Errorf("could not obtain new export id for vine of %d", ap.eid)
	} else {
		exp.thirdPartyRC = tpRc
		exp.thirdPartyIid = tpIid
		exp.thirdPartyVineId = vineEid
		exp.thirdPartyCapDescIdAlt = iinfo.sendToRecipientAlt
		exp.thirdPartyProvideQid = provideQid
		exp.typ = exportTypeThirdPartyExport
		ap.rc.exports.set(ap.eid, exp)
		ap.rc.exports.set(vineEid, vine)

		ap.rc.log.Trace().
			Int("eid", int(ap.eid)).
			Str("thirdParty", tpRc.c.remoteName()).
			Int("thirdPartyIid", int(tpIid)).
			Msg("Resolving previously exported promise")

		err = vat.queueResolve(ap.rc.ctx, ap.rc, ap.eid, exp, exp)
	}

	// After unlocking, any calls received from the remote to this answer
	// will be forwarded to the remote party. Disembargo will be the last
	// message, which should be followed by a finish.
	ap.rc.mu.Unlock()

	return err
}

func (ap answerPromise) resolveToThirdPartyCap(tpRc *runningConn, cap CapabilityClient) error {
	return ap.resolveToThirdPartyImport(tpRc, ImportId(cap.eid))
}

type CallParamsBuilder func(types.PayloadBuilder) error

type CallResultsParser func(types.Payload) (any, error)

type CallContext struct {
	vat *Vat
	rc  *runningConn

	serMb *capnpser.MessageBuilder // Root reply message builder
	mb    types.MessageBuilder
	rb    types.ReturnBuilder
	pb    types.PayloadBuilder

	aid    AnswerId
	iid    InterfaceId
	mid    MethodId
	params types.Payload // Request Call.Params field
}

func (cc *CallContext) InterfaceId() InterfaceId {
	return cc.iid
}

func (cc *CallContext) MethodId() MethodId {
	return cc.mid
}

// IsIidMid is a simple helper to verify if the call context is for a given
// interface and method id. Returns an error if they don't match.
func (cc *CallContext) IsIidMid(iid InterfaceId, mid MethodId) error {
	if cc.iid != iid {
		return ErrUnexpectedInterfaceId{Want: iid, Got: cc.iid}
	}
	if cc.mid != mid {
		return ErrUnexpectedMethodId{Want: mid, Got: cc.mid}
	}
	return nil
}

func (cc *CallContext) Params() types.Payload {
	return cc.params
}

// CallContextParamsStruct extracts the params of the given call context as a
// struct of type T.
func CallContextParamsStruct[T ~capnpser.StructType](cc *CallContext) (res T, err error) {
	pay := cc.Params()
	payContent, err := pay.Content()
	if err != nil {
		err = fmt.Errorf("could not get payload content: %v", err)
		return
	}

	if !payContent.IsStruct() && !payContent.IsZeroStruct() {
		err = fmt.Errorf("payload Content is not a struct")
		return
	}

	res = T(payContent.AsStruct())
	return
}

func CallContextParamsCapability[T ~CallFutureType](cc *CallContext) (res T, err error) {
	cf, err := cc.vat.importRemoteCapFromPayload(cc.rc, cc.Params())
	if err == nil {
		res = T(cf)
	}
	return
}

func (cc *CallContext) initResponsePayload(payloadSizeHint capnpser.WordCount) (err error) {
	totalSizeHint := returnMessageSizeOverhead + payloadSizeHint
	cc.serMb = cc.vat.mbp.getRawMessageBuilder(totalSizeHint)
	cc.mb, err = types.NewRootMessageBuilder(cc.serMb)
	if err != nil {
		return
	}
	cc.rb, err = cc.mb.NewReturn()
	if err != nil {
		return
	}
	cc.rb.SetAnswerId(cc.aid)
	cc.pb, err = cc.rb.NewResults()
	if err != nil {
		return
	}
	return
}

// RespondAsStruct responds this call with a struct. It sets up the content
// response struct and (optinally) allocates extra space for an estimated
// payload size.
//
// The caller should fill the response struct as required by the call.
func (cc *CallContext) RespondAsStruct(structSize capnpser.StructSize, payloadSizeHint capnpser.WordCount) (res capnpser.StructBuilder, err error) {
	if cc.serMb == nil {
		payloadSizeHint += structSize.TotalSize()
		err = cc.initResponsePayload(payloadSizeHint)
		if err != nil {
			return
		}
	}
	res, err = cc.pb.SetContentAsNewStruct(structSize)
	if err != nil {
		return
	}
	return
}

func (cc *CallContext) RespondAsString(s string) (err error) {
	if cc.serMb == nil {
		payloadSizeHint, ok := capnpser.ByteCount(len(s) + 1).StorageWordCount()
		if !ok {
			return errors.New("trying to respond with a string that is too large for CapNProto")
		}
		err = cc.initResponsePayload(payloadSizeHint)
		if err != nil {
			return
		}
	}

	l, err := cc.serMb.CopyToNewText(s)
	if err != nil {
		return err
	}
	return cc.pb.SetContent(l)
}

// RespondAsSenderHostedCap responds this call with a new sender-hosted
// capability.
func (cc *CallContext) RespondAsSenderHostedCap(capServer CallHandler) error {
	if cc.serMb == nil {
		err := cc.initResponsePayload(3) // CapTable + entry + pointer FIXME: double check.
		if err != nil {
			return nil
		}
	}

	// This is OK to do because in the context of a CallContext,
	// RespondAsSenderHostedCap is called within a processCall() stack, thus
	// rc is locked.
	eid, ok := cc.rc.exports.nextID()
	if !ok {
		return errors.New("no more exports allowed")
	}
	cc.rc.exports.set(eid, export{
		typ:     exportTypeLocallyHosted,
		handler: capServer,
	})

	// Similar to respondAsPromise (unify somehow?)
	capTable, err := cc.pb.NewCapTable(1, 1)
	if err != nil {
		return err
	}
	capDesc := capTable.At(0)
	capDesc.SetSenderHosted(eid)
	cc.pb.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0))
	return nil
}

func RespondCallAsStruct[T ~capnpser.StructBuilderType](cc *CallContext, size capnpser.StructSize, payloadSizeHint capnpser.WordCount) (T, error) {
	res, err := cc.RespondAsStruct(size, payloadSizeHint)
	return T(res), err
}

// readReturnCapTable returns the capTable from the payload.
//
// NOTE: this assumes the message being built is a Return with Results and cap
// table list.
func (cc *CallContext) readReturnCapTable() capnpser.GenericStructList[types.CapDescriptor] {
	if cc.serMb == nil {
		// Empty list.
		return capnpser.GenericStructList[types.CapDescriptor]{}
	}

	res := cc.pb.AsReader()
	capTable, err := res.CapTable()
	if err != nil {
		panic(err) // Should never happen.
	}
	return capTable
}

// TODO: Add crb.respondAsSingleCap()

// MUST be called with crb.rc.mu locked.
// Does NOT create the export, only determines the ID it should have.
func (cc *CallContext) respondAsPromise() (answerPromise, error) {
	eid, ok := cc.rc.exports.nextID()
	if !ok {
		return answerPromise{}, errors.New("no more exports allowed")
	}

	if cc.serMb == nil {
		err := cc.initResponsePayload(returnMessageSizeOverhead)
		if err != nil {
			return answerPromise{}, err
		}
	}

	// TODO: Track caps somewhere else? See corresponding in processCall.
	capTable, err := cc.pb.NewCapTable(1, 1)
	if err != nil {
		return answerPromise{}, err
	}
	capDesc := capTable.At(0)
	capDesc.SetSenderPromise(eid)
	cc.pb.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0))

	return answerPromise{rc: cc.rc, eid: eid}, nil
}

type StringReplyCallContext struct {
	cc *CallContext
}

func (scc *StringReplyCallContext) SetResponse(s string) {
	scc.cc.RespondAsString(s)
}

func MakeStringReplyCallContext(cc *CallContext) StringReplyCallContext {
	return StringReplyCallContext{cc: cc}
}

// StructReplyCallContext is a [CallContext] that is used in server method
// definitions that reply with a struct.
//
// In general, instances of this type are only valid for use during the handling
// of individual Calls.
type StructReplyCallContext[T ~capnpser.StructBuilderType] struct {
	cc              *CallContext
	structSize      capnpser.StructSize
	payloadSizeHint capnpser.WordCount
}

// MakeStructReplyCallContext creates a [StructReplyCallContext] for a given
// struct builder type and call context.
func MakeStructReplyCallContext[T ~capnpser.StructBuilderType](cc *CallContext, structSize capnpser.StructSize) StructReplyCallContext[T] {
	return StructReplyCallContext[T]{
		cc:         cc,
		structSize: structSize,
	}
}

// InitResults initializes the response message with a Return.Results (Payload)
// field set to a struct with the previously configured struct size. If any
// payload size hints were added, then the backing message buffer is
// additionally allocated (but not yet consumed) those many words to fit
// internal reply objects before having to grow the buffer.
func (gcc *StructReplyCallContext[T]) InitResults() T {
	s, err := gcc.cc.RespondAsStruct(gcc.structSize, gcc.payloadSizeHint)
	if err != nil {
		panic(err)
	}
	return T(s)
}

// AddWordsToSizeHint adds w many words to the size hint used to initialize a
// response. It only makes sense to call this method before [InitResults] is
// called.
func (gcc *StructReplyCallContext[T]) AddWordsToSizeHint(w capnpser.WordCount) {
	gcc.payloadSizeHint += w
}

// AddBytesToSizeHint adds as many words as needed to store b bytes to the size
// hint used to initialize a response. It only makes sense to call this method
// before [InitResults] is called.
func (gcc *StructReplyCallContext[T]) AddBytesToSizeHint(b int) {
	w, ok := capnpser.ByteCount(b).StorageWordCount()
	if !ok {
		panic(fmt.Sprintf("%b bytes cannot be stored as words", b))
	}
	gcc.payloadSizeHint += capnpser.WordCount(w)
}

// CC returns the underlying CallContext.
func (gcc *StructReplyCallContext[T]) CC() *CallContext {
	return gcc.cc
}

type exception struct { // TODO: improve.
	typ    int
	reason string
}

type callExceptionError interface {
	ToException() exception
}

type allUnimplementedCallHandler struct{}

func (h allUnimplementedCallHandler) Call(ctx context.Context, res *CallContext) error {
	return ErrUnknownMethodId{IId: res.iid, Mid: res.mid}
}

type inMsg struct {
	rc  *runningConn
	msg *capnpser.Message
}

// To be generated from rpc.capnp

type QuestionId = types.QuestionId
type AnswerId = types.AnswerId
type ExportId = types.ExportId
type ImportId = types.ImportId

type questionType int

const (
	questionTypeCall questionType = iota
	questionTypeProvide
)

type question struct {
	typ        questionType
	weakStep   weak.Pointer[pipelineStep]
	strongStep *pipelineStep
}

func (q *question) step() *pipelineStep {
	if q.strongStep != nil {
		return q.strongStep
	}
	return q.weakStep.Value()
}

type answerType int

const (
	answerTypeCall answerType = iota
	answerTypeBootstrap
	answerTypeProvide
	answerTypeAccept
)

type answer struct {
	typ answerType // What this is an anwer to (Call, Provide or Accept).
	eid ExportId   // The export that was the answer to the question.
}

type importType int

const (
	importTypeSenderHosted importType = iota
	importTypeRemotePromise
)

func (it importType) String() string {
	switch it {
	case importTypeSenderHosted:
		return "senderHosted"
	case importTypeRemotePromise:
		return "remotePromise"
	default:
		return fmt.Sprintf("[unknown %d]", it)
	}
}

type imprt struct {
	typ importType

	// Set when this import is a remote promise (this is the target call
	// waiting for the resolution).
	step weak.Pointer[pipelineStep] // TODO: may need same as question (weak/strong pipe).
}

type exportType uint

const (
	exportTypeLocallyHosted exportType = iota
	exportTypePromise
	exportTypeThirdPartyExport
)

func (typ exportType) String() string {
	switch typ {
	case exportTypeLocallyHosted:
		return "locallyHosted"
	case exportTypePromise:
		return "promise"
	case exportTypeThirdPartyExport:
		return "thirdPartyExp"
	default:
		return fmt.Sprintf("[unknown %d]", typ)
	}
}

type export struct {
	typ      exportType
	refCount int

	// Set when this is senderHosted.
	handler CallHandler

	// Set when this was a promise that got resolved.
	resolvedToExport ExportId
	rc               *runningConn

	// Set when this is resolved to a third party imported cap.
	thirdPartyRC           *runningConn
	thirdPartyIid          ImportId
	thirdPartyCapDescIdAlt capnpser.AnyPointer
	thirdPartyVineId       ExportId
	thirdPartyProvideQid   QuestionId
}

func (e *export) resolved() bool {
	return e.resolvedToExport > 0 || e.thirdPartyRC != nil
}

func (e *export) resolveString() string {
	if e.resolvedToExport > 0 {
		return fmt.Sprintf("export %d", e.resolvedToExport)
	} else if e.thirdPartyRC != nil && e.thirdPartyIid > 0 {
		return fmt.Sprintf("third party %s import %d", e.thirdPartyRC.c.remoteName(),
			e.thirdPartyIid)
	} else {
		return "unresolved"
	}
}
