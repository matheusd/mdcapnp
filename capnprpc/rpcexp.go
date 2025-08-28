// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

type msgBuilder struct{} // Alias to a serializer MessageBuilder

type message struct {
	// serialized msg?
}

func (m *message) HasBootstrap() bool {
	panic("fixme")
}

func (m *message) HasReturn() bool {
	panic("fixme")
}

func (m *message) QuestionId() QuestionId {
	panic("fixme")
}

type callable struct {
	// promise || local-callable || remote-capability
	// pipelinable
}

type callParamsBuilder func(*msgBuilder) error

type inMsg struct {
	rc  *runningConn
	msg *message
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
