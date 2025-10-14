// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpctypes

import (
	"fmt"

	"matheusd.com/mdcapnp/capnpser"
)

const (
	bootstrap_questionId_dataField      = 0
	bootstrap_questionId_dataFieldShift = capnpser.Uint32FieldLo
)

type Bootstrap capnpser.Struct

func (s *Bootstrap) QuestionId() QuestionId {
	return QuestionId((*capnpser.Struct)(s).Uint32(bootstrap_questionId_dataField, bootstrap_questionId_dataFieldShift))
}

type BootstrapBuilder capnpser.StructBuilder

func (b *BootstrapBuilder) SetQuestionId(v QuestionId) error {
	return (*capnpser.StructBuilder)(b).SetUint32(bootstrap_questionId_dataField, bootstrap_questionId_dataFieldShift, uint32(v))
}

type PromisedAnswer capnpser.Struct

func (s *PromisedAnswer) QuestionId() QuestionId {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint32FieldLo
	return QuestionId((*capnpser.Struct)(s).Uint32(dataFieldIndex, dataFieldShift))
}

type PromisedAnswerBuilder capnpser.StructBuilder

func (b *PromisedAnswerBuilder) SetQuestionId(v QuestionId) error {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint32FieldLo
	return (*capnpser.StructBuilder)(b).SetUint32(dataFieldIndex, dataFieldShift, uint32(v))
}

type MessageTarget_Which int

const (
	MessageTarget_Which_ImportedCap    MessageTarget_Which = 0
	MessageTarget_Which_PromisedAnswer MessageTarget_Which = 1

	messageTarget_union_field                = 0
	messageTarget_unionFieldShift            = capnpser.Uint16FieldShift2
	messageTarget_importedCap_dataField      = 0
	messageTarget_importedCap_dataFieldShift = capnpser.Uint32FieldLo
	messageTarget_union_ptrField             = 0
)

type MessageTarget capnpser.Struct

func (s *MessageTarget) Which() MessageTarget_Which {
	return MessageTarget_Which((*capnpser.Struct)(s).Uint16(messageTarget_union_field, messageTarget_unionFieldShift))
}

func (s *MessageTarget) AsImportedCap() ImportId {
	return ImportId((*capnpser.Struct)(s).Uint32(messageTarget_importedCap_dataField, messageTarget_importedCap_dataFieldShift))
}

func (s *MessageTarget) AsPromisedAnswer() (res PromisedAnswer, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(messageTarget_union_ptrField, (*capnpser.Struct)(&res))
	return
}

type MessageTargetBuilder capnpser.StructBuilder

func (b *MessageTargetBuilder) SetImportedCap(v ImportId) error {
	const unionValue = uint16(MessageTarget_Which_ImportedCap)

	if err := (*capnpser.StructBuilder)(b).SetUint32(messageTarget_importedCap_dataField, messageTarget_importedCap_dataFieldShift, uint32(v)); err != nil {
		return err
	}
	return (*capnpser.StructBuilder)(b).SetUint16(messageTarget_union_field, messageTarget_unionFieldShift, unionValue)
}

func (b *MessageTargetBuilder) NewPromisedAnswer() (sb PromisedAnswerBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	const unionValue = uint16(MessageTarget_Which_PromisedAnswer)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(messageTarget_union_ptrField, structSize, messageTarget_union_field, messageTarget_unionFieldShift, unionValue)
	sb = PromisedAnswerBuilder(nsb)
	return
}

type Call capnpser.Struct

func (s *Call) QuestionId() QuestionId {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint32FieldLo
	return QuestionId((*capnpser.Struct)(s).Uint32(dataFieldIndex, dataFieldShift))
}

func (s *Call) InterfaceId() uint64 {
	const dataFieldIndex = 1
	return (*capnpser.Struct)(s).Uint64(dataFieldIndex)
}

func (s *Call) MethodId() uint16 {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint16FieldShift2
	return (*capnpser.Struct)(s).Uint16(dataFieldIndex, dataFieldShift)
}

func (s *Call) Target() (res MessageTarget, err error) {
	const ptrFieldIndex = 0
	err = (*capnpser.Struct)(s).ReadStruct(ptrFieldIndex, (*capnpser.Struct)(&res))
	return
}

func (s *Call) NoPromisePipelining() bool {
	const dataFieldIndex = 2
	const dataFieldBit = 1
	return (*capnpser.Struct)(s).Bool(dataFieldIndex, dataFieldBit)
}

type CallBuilder capnpser.StructBuilder

func (b *CallBuilder) SetQuestionId(v QuestionId) error {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint32FieldLo
	return (*capnpser.StructBuilder)(b).SetUint32(dataFieldIndex, dataFieldShift, uint32(v))
}

func (b *CallBuilder) SetInterfaceId(v uint64) error {
	const dataFieldIndex = 1
	return (*capnpser.StructBuilder)(b).SetUint64(dataFieldIndex, v)
}

func (b *CallBuilder) SetMethodId(v uint16) error {
	const dataFieldIndex = 0
	const dataFieldShift = capnpser.Uint16FieldShift2
	return (*capnpser.StructBuilder)(b).SetUint16(dataFieldIndex, dataFieldShift, v)
}

func (b *CallBuilder) NewTarget() (sb MessageTargetBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	const ptrFieldIndex = 0

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructField(ptrFieldIndex, structSize)
	sb = MessageTargetBuilder(nsb)
	return
}

type Message_Which int

const (
	Message_Which_Call      Message_Which = 2
	Message_Which_Bootstrap Message_Which = 8
)

func (w Message_Which) String() string {
	switch w {
	case Message_Which_Call:
		return "call"
	case Message_Which_Bootstrap:
		return "bootstrap"
	default:
		return fmt.Sprintf("unknown which %d", w)
	}
}

type Message capnpser.Struct

func (s *Message) Which() Message_Which {
	const unionField = 0
	const unionFieldShift = capnpser.Uint16FieldShift0

	return Message_Which((*capnpser.Struct)(s).Uint16(unionField, unionFieldShift))
}

func (s *Message) ReadFromRoot(msg *capnpser.Message) error {
	return msg.ReadRoot((*capnpser.Struct)(s))
}

func (s *Message) AsBootstrap() (res Bootstrap, err error) {
	const unionPointerField = 0
	err = (*capnpser.Struct)(s).ReadStruct(unionPointerField, (*capnpser.Struct)(&res))
	return
}

func (s *Message) AsCall() (res Call, err error) {
	const unionPointerField = 0
	err = (*capnpser.Struct)(s).ReadStruct(unionPointerField, (*capnpser.Struct)(&res))
	return
}

type MessageBuilder capnpser.StructBuilder

func (b *MessageBuilder) NewBoostrap() (sb BootstrapBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	const ptrFieldIndex = 0
	const unionField = 0
	const unionFieldShift = capnpser.Uint16FieldShift0
	const unionValue = uint16(Message_Which_Bootstrap)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(ptrFieldIndex, structSize, unionField, unionFieldShift, unionValue)
	sb = BootstrapBuilder(nsb)
	return
}

func (b *MessageBuilder) NewCall() (sb CallBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 3, PointerSectionSize: 3}
	const ptrFieldIndex = 0
	const unionField = 0
	const unionFieldShift = capnpser.Uint16FieldShift0
	const unionValue = uint16(Message_Which_Call)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(ptrFieldIndex, structSize, unionField, unionFieldShift, unionValue)
	sb = CallBuilder(nsb)
	return
}

func NewRootMessageBuilder(mb *capnpser.MessageBuilder) (MessageBuilder, error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	b, err := mb.NewRootStruct(structSize)
	return MessageBuilder(b), err
}
