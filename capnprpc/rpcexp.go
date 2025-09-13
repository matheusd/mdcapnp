// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"weak"

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

type capDescriptor struct {
	senderHosted  ExportId
	senderPromise ExportId
}

func (cp *capDescriptor) IsSenderHosted() bool      { return cp.senderHosted > 0 }
func (cp *capDescriptor) AsSenderHosted() ExportId  { return cp.senderHosted }
func (cp *capDescriptor) IsSenderPromise() bool     { return cp.senderPromise > 0 }
func (cp *capDescriptor) AsSenderPromise() ExportId { return cp.senderPromise }
func (cp *capDescriptor) hasExportId() bool         { return cp.IsSenderHosted() || cp.IsSenderPromise() }
func (cp *capDescriptor) exportId() ExportId        { return cp.senderHosted + cp.senderPromise } // Only one will be set.

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
	aid         AnswerId
	isResults   bool
	pay         payload
	isException bool
	exc         exception
}

func (r *rpcReturn) AnswerId() AnswerId { return r.aid }
func (r *rpcReturn) IsResults() bool    { return r.isResults }
func (r *rpcReturn) AsResults() payload { return r.pay }

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
	qid    QuestionId
	target messageTarget
	iid    uint64
	mid    uint16
	params payload
}

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

type message struct { // RPC message type
	isBootstrap bool
	isReturn    bool
	isCall      bool
	isFinish    bool
	isResolve   bool

	boot    bootstrap
	ret     rpcReturn
	call    call
	finish  finish
	resolve resolve

	testEcho uint64 // Special test message.
}

func (m *message) ReadFromRoot(msg *capnpser.Message) error { return nil }
func (m *message) IsBootstrap() bool                        { return m.isBootstrap }
func (m *message) AsBootstrap() bootstrap                   { return m.boot }
func (m *message) IsReturn() bool                           { return m.isReturn }
func (m *message) AsReturn() rpcReturn                      { return m.ret }
func (m *message) IsCall() bool                             { return m.isCall }
func (m *message) AsCall() call                             { return m.call }
func (m *message) IsFinish() bool                           { return m.isFinish }
func (m *message) AsFinish() finish                         { return m.finish }
func (m *message) IsResolve() bool                          { return m.isResolve }
func (m *message) AsResolve() resolve                       { return m.resolve }

type interfaceId uint64
type methodId uint16

type answerPromise struct {
	rc  *runningConn
	eid ExportId
}

func (ap answerPromise) resolveToHandler(handler callHandler) error {
	var err error

	resolution := export{
		typ:      exportTypeLocallyHosted,
		handler:  handler,
		refCount: 1,
	}

	ap.rc.mu.Lock()
	exp, ok := ap.rc.exports.get(ap.eid)
	if !ok {
		err = fmt.Errorf("export %d not found", ap.eid)
	} else if exp.resolvedToExport > 0 {
		err = fmt.Errorf("promised export %d already resolved to export %d",
			ap.eid, exp.resolvedToExport)
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

		err = ap.rc.vat.sendResolve(ap.rc.ctx, ap.rc, ap.eid, exp, resolution)
	}
	ap.rc.mu.Unlock()

	return err
}

type callReturnBuilder struct {
	rc      *runningConn
	payload payload
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

	// TODO: find a zero alloc way of representing this kind of return.
	crb.payload.capTable = []capDescriptor{{senderPromise: eid}}
	crb.payload.content = anyPointer{isCapPointer: true, cp: capPointer{index: 0}}

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

type QuestionId uint32
type AnswerId uint32
type ExportId uint32
type ImportId uint32

type question struct {
	pipe    weak.Pointer[pipeline]
	stepIdx int
}

type answer struct {
	eid ExportId // answered with an export.
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
	pipe    weak.Pointer[pipeline]
	stepIdx int
}

type exportType uint

const (
	exportTypeLocallyHosted exportType = iota
	exportTypePromise
)

func (typ exportType) String() string {
	switch typ {
	case exportTypeLocallyHosted:
		return "locallyHosted"
	case exportTypePromise:
		return "promise"
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

	// Track calls that must be fulfilled once this is fulfilled.

	// TODO: refcount to send Finish()?
}
