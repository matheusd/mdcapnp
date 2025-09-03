// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"

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

type CapDescriptor struct {
	senderHosted ExportId
	// senderPromise ExportId
}

func (cp *CapDescriptor) IsSenderHosted() bool     { return cp.senderHosted > 0 }
func (cp *CapDescriptor) AsSenderHosted() ExportId { return cp.senderHosted }

type Payload struct {
	content  anyPointer
	capTable []CapDescriptor
}

func (p *Payload) Content() anyPointer       { return p.content }
func (p *Payload) CapTable() []CapDescriptor { return p.capTable }

type Exception struct {
	reason string
	typ    int
}

type Return struct {
	aid         AnswerId
	isResults   bool
	pay         Payload
	isException bool
	exc         Exception
}

func (r *Return) AnswerId() AnswerId { return r.aid }
func (r *Return) IsResults() bool    { return r.isResults }
func (r *Return) AsResults() Payload { return r.pay }

type PromisedAnswer struct {
	qid QuestionId
	// transform
}

type MessageTarget struct {
	isImportedCap    bool
	impcap           ImportId
	isPromisedAnswer bool
	pans             PromisedAnswer
}

type Call struct {
	qid    QuestionId
	target MessageTarget
	iid    uint64
	mid    uint16
	params Payload
}

type Bootstrap struct {
	qid QuestionId
}

func (bt *Bootstrap) QuestionId() QuestionId { return bt.qid }

type Message struct { // RPC message type
	isBootstrap bool
	boot        Bootstrap
	isReturn    bool
	ret         Return
	isCall      bool
	call        Call
}

func (m *Message) ReadFromRoot(msg *capnpser.Message) error { return nil }
func (m *Message) IsBootstrap() bool                        { return m.isBootstrap }
func (m *Message) AsBootstrap() Bootstrap                   { return m.boot }
func (m *Message) IsReturn() bool                           { return m.isReturn }
func (m *Message) AsReturn() Return                         { return m.ret }
func (m *Message) IsCall() bool                             { return m.isCall }
func (m *Message) AsCall() Call                             { return m.call }

type interfaceId uint64
type methodId uint16

type callReturnBuilder struct {
	payload Payload
}

func (crb *callReturnBuilder) setContent(content anyPointer) {
	crb.payload.content = content
}

type callExceptionError interface {
	ToException() Exception
}

type errUnimplemented struct {
	Iid interfaceId
	Mid methodId
}

func (err errUnimplemented) Error() string {
	return fmt.Sprintf("call %d.%d unimplemented", err.Iid, err.Mid)
}

func (err errUnimplemented) ToException() Exception {
	return Exception{typ: 3, reason: err.Error()}
}

type callHandlerArgs struct {
	iid    interfaceId
	mid    methodId
	params Payload
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
	msg Message
}

// To be generated from rpc.capnp

type QuestionId uint32
type AnswerId uint32
type ExportId uint32
type ImportId uint32

type question struct {
	pipe    *pipeline
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
