// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

type messageBuilderPool struct {
	alloc capnpser.Allocator // Alloc strategy for outbound rpc messages.
	p     *sync.Pool
}

func (mp *messageBuilderPool) getRawMessageBuilder() *capnpser.MessageBuilder {
	return mp.p.Get().(*capnpser.MessageBuilder)
}

func (mp *messageBuilderPool) getForPayloadSize(extraPayloadSize int) (outMsg, error) {
	// TODO: calculate the size hint.
	serMb := mp.getRawMessageBuilder()
	mb, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		return outMsg{}, err
	}
	return outMsg{serMsg: serMb, mb: mb}, nil
}

func (mp *messageBuilderPool) get() (outMsg, error) {
	return mp.getForPayloadSize(0)
}

func (mp *messageBuilderPool) put(serMb *capnpser.MessageBuilder) {
	err := serMb.Reset()
	if err != nil {
		panic(err) // Simple allocator never errors on Reset().
	}
	mp.p.Put(serMb)
}

func newMessageBuilderPool() *messageBuilderPool {
	alloc := capnpser.NewSimpleSingleAllocator(16, false)
	return &messageBuilderPool{
		alloc: alloc,
		p: &sync.Pool{
			New: func() any {
				serMb, err := capnpser.NewMessageBuilder(alloc)
				if err != nil {
					// SimpleSingleAlloc never errors here.
					panic(err)
				}
				return serMb
			},
		},
	}
}

type rpcCallMsgBuilder struct {
	outMsg
	builder     capnpser.StructBuilder
	isBootstrap bool
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

func (ap answerPromise) resolveToThirdPartyCap(tpRc *runningConn, cap capability) error {
	return ap.resolveToThirdPartyImport(tpRc, ImportId(cap.eid))
}

type CallParamsBuilder func(types.MessageBuilder) error

type CallContext struct {
	rc    *runningConn
	pb    types.PayloadBuilder
	serMb *capnpser.MessageBuilder // Root reply message builder

	iid InterfaceId
	mid MethodId
}

func (cc *CallContext) InterfaceId() InterfaceId {
	return cc.iid
}

func (cc *CallContext) MethodId() MethodId {
	return cc.mid
}

// readReturnCapTable returns the capTable from the payload.
//
// NOTE: this assumes the message being built is a Return with Results and cap
// table list.
func (crb *CallContext) readReturnCapTable() capnpser.GenericStructList[types.CapDescriptor] {
	var rpcMsg types.Message

	// Errors don't need checking, because this is assumed to be called
	// during return building, where the structures were allocated.
	serMsg := crb.serMb.MessageReader()
	rpcMsg.ReadFromRoot(&serMsg)
	ret, _ := rpcMsg.AsReturn()
	res, _ := ret.AsResults()
	capTable, _ := res.CapTable()
	return capTable
}

// TODO: Add crb.respondAsStruct()
// TODO: Add crb.respondAsSingleCap()

// MUST be called with crb.rc.mu locked.
// Does NOT create the export, only determines the ID it should have.
func (crb *CallContext) respondAsPromise() (answerPromise, error) {
	eid, ok := crb.rc.exports.nextID()
	if !ok {
		return answerPromise{}, errors.New("no more exports allowed")
	}

	// TODO: Track caps somewhere else? See corresponding in processCall.
	capTable, err := crb.pb.NewCapTable(1, 1)
	if err != nil {
		return answerPromise{}, err
	}
	capDesc := capTable.At(0)
	capDesc.SetSenderPromise(eid)

	crb.pb.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0))

	return answerPromise{rc: crb.rc, eid: eid}, nil
}

type exception struct { // TODO: improve.
	typ    int
	reason string
}

type callExceptionError interface {
	ToException() exception
}

type errUnimplemented struct {
	Iid InterfaceId
	Mid MethodId
}

func (err errUnimplemented) Error() string {
	return fmt.Sprintf("call %d.%d unimplemented", err.Iid, err.Mid)
}

func (err errUnimplemented) ToException() exception {
	return exception{typ: 3, reason: err.Error()}
}

// CallHandler is the lowest level handler for calls (including bootstrap).
type CallHandler interface {
	Call(ctx context.Context, cc *CallContext) error
}

type CallHandlerFunc func(ctx context.Context, rb *CallContext) error

func (f CallHandlerFunc) Call(ctx context.Context, rb *CallContext) error {
	return f(ctx, rb)
}

type allUnimplementedCallHandler struct{}

func (h allUnimplementedCallHandler) Call(ctx context.Context, res *CallContext) error {
	return errUnimplemented{Iid: res.iid, Mid: res.mid}
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

	handler CallHandler // Set when this is senderHosted.

	// Set when this was a promise that got resolved.
	resolvedToExport ExportId
	rc               *runningConn

	// Set when this is resolved to a third party imported cap.
	thirdPartyRC           *runningConn
	thirdPartyIid          ImportId
	thirdPartyCapDescIdAlt capnpser.AnyPointer
	thirdPartyVineId       ExportId
	thirdPartyProvideQid   QuestionId

	// Track calls that must be fulfilled once this is fulfilled.

	// TODO: refcount to send Finish()?
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
