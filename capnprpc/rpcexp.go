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

type msgBuilder struct{} // Alias to a serializer MessageBuilder

type capability struct { // Right type?
	eid ExportId
}

type capPointer struct { // Alias to capnpser.CapPointer
	index uint32
}

func (cp *capPointer) Index() uint32 { return cp.index }

type serStruct struct { // Alias to capnpser.Struct
	rawData []byte
}

type anyPointer struct { // Alias to capnpser.AnyPointer
	isStruct     bool
	st           serStruct
	isCapPointer bool
	cp           capPointer
	isVoid       bool
}

func (ap *anyPointer) IsStruct() bool           { return ap.isStruct }
func (ap *anyPointer) AsStruct() serStruct      { return ap.st }
func (ap *anyPointer) IsCapPointer() bool       { return ap.isCapPointer }
func (ap *anyPointer) AsCapPointer() capPointer { return ap.cp }
func (ap *anyPointer) IsVoid() bool             { return ap.isVoid }

type thirdPartyCapId anyPointer
type provisionId anyPointer
type recipientId anyPointer

type introductionInfo struct {
	sendToRecipient thirdPartyCapId
	sendToTarget    recipientId
}

type thirdPartyCapDescriptor struct {
	id     thirdPartyCapId
	vineId ExportId
}

type capDescriptor struct {
	senderHosted     ExportId
	senderPromise    ExportId
	thirdPartyHosted thirdPartyCapDescriptor
}

func (cp *capDescriptor) IsSenderHosted() bool                        { return cp.senderHosted > 0 }
func (cp *capDescriptor) AsSenderHosted() ExportId                    { return cp.senderHosted }
func (cp *capDescriptor) IsSenderPromise() bool                       { return cp.senderPromise > 0 }
func (cp *capDescriptor) AsSenderPromise() ExportId                   { return cp.senderPromise }
func (cp *capDescriptor) hasExportId() bool                           { return cp.IsSenderHosted() || cp.IsSenderPromise() }
func (cp *capDescriptor) exportId() ExportId                          { return cp.senderHosted + cp.senderPromise } // Only one will be set.
func (cp *capDescriptor) IsThirdPartyHosted() bool                    { return cp.thirdPartyHosted.vineId > 0 }     // Asumes it is always set in this case.
func (cp *capDescriptor) AsThirdPartyHosted() thirdPartyCapDescriptor { return cp.thirdPartyHosted }

type payload struct {
	content  anyPointer
	capTable []capDescriptor
}

func (p *payload) Content() anyPointer       { return p.content }
func (p *payload) CapTable() []capDescriptor { return p.capTable }

type exception struct {
	reason string
	typ    int
}

type rpcReturn struct {
	aid            AnswerId
	isResults      bool
	pay            payload
	isException    bool
	exc            exception
	noFinishNeeded bool
}

func (r *rpcReturn) AnswerId() AnswerId   { return r.aid }
func (r *rpcReturn) IsResults() bool      { return r.isResults }
func (r *rpcReturn) AsResults() payload   { return r.pay }
func (r *rpcReturn) NoFinishNeeded() bool { return r.noFinishNeeded }

type promisedAnswer struct {
	qid QuestionId
	// transform
}

type messageTarget struct {
	isImportedCap    bool
	impcap           ImportId
	isPromisedAnswer bool
	pans             promisedAnswer
}

type call struct {
	qid                 QuestionId
	target              messageTarget
	iid                 uint64
	mid                 uint16
	params              payload
	noPromisePipelining bool
}

func (c *call) NoPromisePipelining() bool { return c.noPromisePipelining }

type bootstrap struct {
	qid QuestionId
}

func (bt *bootstrap) QuestionId() QuestionId { return bt.qid }

type finish struct {
	qid QuestionId
}

func (f *finish) QuestionId() QuestionId { return f.qid }

type resolve struct {
	pid ExportId
	cap capDescriptor
}

type disembargo struct {
	target    messageTarget
	isAccept  bool
	isProvide bool
	provide   QuestionId
}

type accept struct {
	qid       QuestionId
	provision provisionId
	embargo   bool
}

type provide struct {
	qid       QuestionId
	target    messageTarget
	recipient recipientId
}

type message_which int

const (
	message_which_bootstrap message_which = iota
	message_which_return
	message_which_call
	message_which_finish
	message_which_resolve
	message_which_disembargo
	message_which_accept
	message_which_provide

	message_which_rawRmb   message_which = 999998
	message_which_testecho message_which = 999999
)

func (mw message_which) String() string {
	switch mw {
	case message_which_bootstrap:
		return "bootstrap"
	case message_which_return:
		return "return"
	case message_which_call:
		return "call"
	case message_which_finish:
		return "finish"
	case message_which_resolve:
		return "resolve"
	case message_which_disembargo:
		return "disembargo"
	case message_which_accept:
		return "accept"
	case message_which_provide:
		return "provide"
	default:
		return "unknown"
	}
}

type message struct { // RPC message type
	//isBootstrap  bool
	// isReturn bool
	//isCall       bool
	//isFinish     bool
	isResolve    bool
	isDisembargo bool
	isAccept     bool
	isProvide    bool

	//boot       bootstrap
	// ret rpcReturn
	//call       call
	// finish     finish
	resolve    resolve
	disembargo disembargo
	accept     accept
	provide    provide

	testEcho uint64 // Special test message.

	rawSerBytes []byte
	rawSerMb    *capnpser.MessageBuilder
	rawSerMsg   *capnpser.Message
}

func (m *message) Which() message_which {
	switch {
	//case m.isBootstrap:
	//	return message_which_bootstrap
	// case m.isReturn:
	//	return message_which_return
	//case m.isCall:
	//	return message_which_call
	//case m.isFinish:
	//	return message_which_finish
	case m.isResolve:
		return message_which_resolve
	case m.isDisembargo:
		return message_which_disembargo
	case m.isAccept:
		return message_which_accept
	case m.isProvide:
		return message_which_provide
	case m.testEcho > 0:
		return message_which_testecho
	case m.rawSerMb != nil:
		return message_which_rawRmb
	default:
		panic("unknown message which")
	}
}
func (m *message) ReadFromRoot(msg *capnpser.Message) error { return nil }

// func (m *message) IsBootstrap() bool                        { return m.isBootstrap }
// func (m *message) AsBootstrap() bootstrap                   { return m.boot }
// func (m *message) IsReturn() bool      { return m.isReturn }
// func (m *message) AsReturn() rpcReturn { return m.ret }

// func (m *message) IsCall() bool                             { return m.isCall }
// func (m *message) AsCall() call             { return m.call }
// func (m *message) IsFinish() bool           { return m.isFinish }
// func (m *message) AsFinish() finish         { return m.finish }
func (m *message) IsResolve() bool          { return m.isResolve }
func (m *message) AsResolve() resolve       { return m.resolve }
func (m *message) IsDisembargo() bool       { return m.isDisembargo }
func (m *message) AsDisembargo() disembargo { return m.disembargo }
func (m *message) IsAccept() bool           { return m.isAccept }
func (m *message) AsAccept() accept         { return m.accept }
func (m *message) IsProvide() bool          { return m.isProvide }
func (m *message) AsProvide() provide       { return m.provide }

type messagePool struct {
	p *sync.Pool
}

func (mp *messagePool) getForPayloadSize(extraPayloadSize int) *message {
	// TODO: size the message based on the size of call args.
	return mp.p.Get().(*message)
}

func (mp *messagePool) get() *message {
	return mp.p.Get().(*message)
}

func (mp *messagePool) put(m *message) {
	*m = message{}
	mp.p.Put(m)
}

func newMessagePool() *messagePool {
	return &messagePool{
		p: &sync.Pool{
			New: func() any {
				return &message{}
			},
		},
	}
}

type rpcMsgBuilder struct {
	serMb *capnpser.MessageBuilder
	mb    types.MessageBuilder
}

type messageBuilderPool struct {
	alloc capnpser.Allocator // Alloc strategy for outbound rpc messages.
	p     *sync.Pool
}

func (mp *messageBuilderPool) getForPayloadSize(extraPayloadSize int) (rpcMsgBuilder, error) {
	// TODO: reuse allocator.
	// TODO: calculate the size hint.
	serMb := mp.p.Get().(*capnpser.MessageBuilder)
	mb, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		return rpcMsgBuilder{}, err
	}
	return rpcMsgBuilder{serMb: serMb, mb: mb}, nil
}

func (mp *messageBuilderPool) get() (rpcMsgBuilder, error) {
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
	rpcMsgBuilder
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

func (ap answerPromise) resolveToHandler(handler callHandler) error {
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
	provide := provide{
		target:    messageTarget{isImportedCap: true, impcap: tpIid},
		recipient: iinfo.sendToTarget,
	}
	tpRc.mu.Lock()
	if !tpRc.imports.has(tpIid) {
		err = fmt.Errorf("third party %s does not have import %d", tpRc, tpIid)
	} else if qid, ok := tpRc.questions.nextID(); !ok {
		err = fmt.Errorf("could not generate new question id for %s", tpRc)
	} else {
		provide.qid = qid
		tpRc.questions.set(qid, question{typ: questionTypeProvide})
	}
	tpRc.mu.Unlock()
	if err != nil {
		return err
	}

	// CHECK: ok to do here, outside tpRc.mu.Lock()?
	if err := tpRc.vat.sendProvide(tpRc.ctx, tpRc, provide); err != nil {
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
		exp.thirdPartyCapDescId = iinfo.sendToRecipient
		exp.thirdPartyProvideQid = provide.qid
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

type callReturnBuilder struct {
	rc      *runningConn
	payload payload
	pb      types.PayloadBuilder
	serMb   *capnpser.MessageBuilder // Root reply message builder
}

// readReturnCapTable returns the capTable from the payload.
//
// NOTE: this assumes the message being built is a Return with Results and cap
// table list.
func (crb *callReturnBuilder) readReturnCapTable() capnpser.GenericStructList[types.CapDescriptor] {
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

func (crb *callReturnBuilder) setContent(content anyPointer) {
	crb.payload.content = content
}

// TODO: Add crb.respondAsStruct()
// TODO: Add crb.respondAsSingleCap()

// MUST be called with crb.rc.mu locked.
// Does NOT create the export, only determines the ID it should have.
func (crb *callReturnBuilder) respondAsPromise() (answerPromise, error) {
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

type callExceptionError interface {
	ToException() exception
}

type errUnimplemented struct {
	Iid interfaceId
	Mid methodId
}

func (err errUnimplemented) Error() string {
	return fmt.Sprintf("call %d.%d unimplemented", err.Iid, err.Mid)
}

func (err errUnimplemented) ToException() exception {
	return exception{typ: 3, reason: err.Error()}
}

type callHandlerArgs struct {
	iid    interfaceId
	mid    methodId
	params payload
	rc     *runningConn
}

// callHandler is the lowest handler for calls (including bootstrap).
type callHandler interface {
	Call(ctx context.Context, args callHandlerArgs, rb *callReturnBuilder) error
}

type callHandlerFunc func(ctx context.Context, args callHandlerArgs, rb *callReturnBuilder) error

func (f callHandlerFunc) Call(ctx context.Context, args callHandlerArgs, rb *callReturnBuilder) error {
	return f(ctx, args, rb)
}

type allUnimplementedCallHandler struct{}

func (h allUnimplementedCallHandler) Call(ctx context.Context, args callHandlerArgs, res *callReturnBuilder) error {
	return errUnimplemented{Iid: args.iid, Mid: args.mid}
}

type callable struct {
	// promise || local-callable || remote-capability
	// pipelinable
}

type callParamsBuilder func(*msgBuilder) error

type inMsg struct {
	rc  *runningConn
	msg message
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

	handler callHandler // Set when this is senderHosted.

	// Set when this was a promise that got resolved.
	resolvedToExport ExportId
	rc               *runningConn

	// Set when this is resolved to a third party imported cap.
	thirdPartyRC         *runningConn
	thirdPartyIid        ImportId
	thirdPartyCapDescId  thirdPartyCapId
	thirdPartyVineId     ExportId
	thirdPartyProvideQid QuestionId

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
