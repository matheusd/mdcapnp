// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
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
	senderHosted ExportId
	// senderPromise ExportId
}

func (cp *capDescriptor) IsSenderHosted() bool     { return cp.senderHosted > 0 }
func (cp *capDescriptor) AsSenderHosted() ExportId { return cp.senderHosted }

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

type message struct { // RPC message type
	isBootstrap bool
	boot        bootstrap
	isReturn    bool
	ret         rpcReturn
	isCall      bool
	call        call
	isFinish    bool
	finish      finish

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

type interfaceId uint64
type methodId uint16

type callReturnBuilder struct {
	payload payload
}

func (crb *callReturnBuilder) setContent(content anyPointer) {
	crb.payload.content = content
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
)

type imprt struct {
	typ importType
}

type exportType uint

const (
	exportTypeLocallyHosted exportType = iota
)

type export struct {
	typ exportType

	handler callHandler // Set when this is senderHosted.

	// Track calls that must be fulfilled once this is fulfilled.

	// TODO: refcount to send Finish()?
}
