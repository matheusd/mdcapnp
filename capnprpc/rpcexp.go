// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "matheusd.com/mdcapnp/capnpser"

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
	isCapPointer bool
	cp           capPointer
	st           serStruct
}

func (ap *anyPointer) IsStruct() bool           { return ap.isStruct }
func (ap *anyPointer) IsCapPointer() bool       { return ap.isCapPointer }
func (ap *anyPointer) AsCapPointer() capPointer { return ap.cp }
func (ap *anyPointer) AsStruct() serStruct      { return ap.st }

type CapDescriptor struct {
	senderHosted ExportId
}

func (cp *CapDescriptor) IsSenderHosted() bool     { return cp.senderHosted > 0 }
func (cp *CapDescriptor) AsSenderHosted() ExportId { return cp.senderHosted }

type Payload struct {
	content  anyPointer
	capTable []CapDescriptor
}

func (p *Payload) Content() anyPointer       { return p.content }
func (p *Payload) CapTable() []CapDescriptor { return p.capTable }

type Return struct {
	aid       AnswerId
	isResults bool
	pay       Payload
}

func (r *Return) AnswerId() AnswerId { return r.aid }
func (r *Return) IsResults() bool    { return r.isResults }
func (r *Return) AsResults() Payload { return r.pay }

type Bootstrap struct {
	qid QuestionId
}

func (bt *Bootstrap) QuestionId() QuestionId { return bt.qid }

type Message struct { // RPC message type
	isBootstrap bool
	boot        Bootstrap
	isReturn    bool
	ret         Return
}

func (m *Message) ReadFromRoot(msg *capnpser.Message) error { return nil }
func (m *Message) IsBootstrap() bool                        { return m.isBootstrap }
func (m *Message) AsBootstrap() Bootstrap                   { return m.boot }
func (m *Message) IsReturn() bool                           { return m.isReturn }
func (m *Message) AsReturn() Return                         { return m.ret }

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
type answer struct{}
type imprt struct{}

type exportType uint

const (
	exportTypeSenderHosted exportType = iota
	exportTypeSenderPromise
)

type export struct {
	typ exportType

	// Track calls that must be fulfilled once this is fulfilled.

	// TODO: refcount to send Finish()?
}
