// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import "matheusd.com/mdcapnp/capnpser"

type msgBuilder struct{} // Alias to a serializer MessageBuilder

type Return capnpser.Struct

func (r *Return) AnswerId() AnswerId { panic("fixme") }

type Message capnpser.Struct // RPC message type

func (m *Message) ReadFromRoot(msg *capnpser.Message) error { panic("fixme") }
func (m *Message) IsBootstrap() bool                        { panic("fixme") }
func (m *Message) IsReturn() bool                           { panic("fixme") }
func (m *Message) QuestionId() QuestionId                   { panic("fixme") }
func (m *Message) AsReturn() Return                         { panic("fixme") }

type callable struct {
	// promise || local-callable || remote-capability
	// pipelinable
}

type callParamsBuilder func(*msgBuilder) error

type inMsg struct {
	rc  *runningConn
	msg capnpser.Message
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
type export struct{}
