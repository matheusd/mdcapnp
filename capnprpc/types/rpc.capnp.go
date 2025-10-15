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

type ThirdPartyCapDescriptor capnpser.Struct

const (
	thirdPartyCapDescriptor_id_ptrField           = 0
	thirdPartyCapDescriptor_vineId_dataField      = 0
	thirdPartyCapDescriptor_vineId_dataFieldShift = capnpser.Uint32FieldLo
)

func (s *ThirdPartyCapDescriptor) Id() (res AnyPointer, err error) {
	err = (*capnpser.Struct)(s).ReadAnyPointer(thirdPartyCapDescriptor_id_ptrField, &res)
	return
}

func (s *ThirdPartyCapDescriptor) VineId() ExportId {
	return ExportId((*capnpser.Struct)(s).Uint32(thirdPartyCapDescriptor_vineId_dataField, thirdPartyCapDescriptor_vineId_dataFieldShift))
}

type ThirdPartyCapDescriptorBuilder capnpser.StructBuilder

func (b *ThirdPartyCapDescriptorBuilder) SetId(v capnpser.AnyPointerBuilder) error {
	return (*capnpser.StructBuilder)(b).SetAnyPointer(thirdPartyCapDescriptor_id_ptrField, v)
}

func (b *ThirdPartyCapDescriptorBuilder) SetVineId(v ExportId) error {
	return (*capnpser.StructBuilder)(b).SetUint32(thirdPartyCapDescriptor_vineId_dataField, thirdPartyCapDescriptor_vineId_dataFieldShift, uint32(v))
}

type CapDescriptor_Which int

const (
	CapDescriptor_Which_None             CapDescriptor_Which = 0
	CapDescriptor_Which_SenderHosted     CapDescriptor_Which = 1
	CapDescriptor_Which_SenderPromise    CapDescriptor_Which = 2
	CapDescriptor_Which_ThirdPartyHosted CapDescriptor_Which = 5

	capDescriptor_union_dataField      = 0
	capDescriptor_union_dataFieldShift = capnpser.Uint16FieldShift0
	capDescriptor_union_ptrField       = 0
)

type CapDescriptor capnpser.Struct

func (s *CapDescriptor) Which() CapDescriptor_Which {
	return CapDescriptor_Which((*capnpser.Struct)(s).Uint16(capDescriptor_union_dataField, capDescriptor_union_dataFieldShift))
}

func (s *CapDescriptor) AsSenderHosted() ExportId {
	return ExportId((*capnpser.Struct)(s).Uint32(capDescriptor_union_dataField, capnpser.Uint32FieldHi))
}

func (s *CapDescriptor) AsSenderPromise() ExportId {
	return ExportId((*capnpser.Struct)(s).Uint32(capDescriptor_union_dataField, capnpser.Uint32FieldHi))
}

func (s *CapDescriptor) AsThirdPartyHosted() (res ThirdPartyCapDescriptor, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(capDescriptor_union_ptrField, (*capnpser.Struct)(&res))
	return
}

type CapDescriptorBuilder capnpser.StructBuilder

func (b *CapDescriptorBuilder) SetSenderHosted(v ExportId) (err error) {
	const unionValue = uint16(CapDescriptor_Which_SenderHosted)
	err = ((*capnpser.StructBuilder)(b)).SetUint32(capDescriptor_union_dataField, capnpser.Uint32FieldHi, uint32(v))
	if err == nil {
		err = ((*capnpser.StructBuilder)(b)).SetUint16(capDescriptor_union_dataField, capDescriptor_union_dataFieldShift, unionValue)
	}
	return
}

func (b *CapDescriptorBuilder) SetSenderPromise(v ExportId) (err error) {
	const unionValue = uint16(CapDescriptor_Which_SenderPromise)
	err = ((*capnpser.StructBuilder)(b)).SetUint32(capDescriptor_union_dataField, capnpser.Uint32FieldHi, uint32(v))
	if err == nil {
		err = ((*capnpser.StructBuilder)(b)).SetUint16(capDescriptor_union_dataField, capDescriptor_union_dataFieldShift, unionValue)
	}
	return
}

func (b *CapDescriptorBuilder) NewThirdPartyHosted() (sb ThirdPartyCapDescriptorBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	const unionValue = uint16(CapDescriptor_Which_ThirdPartyHosted)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(capDescriptor_union_ptrField, structSize, capDescriptor_union_dataField, capDescriptor_union_dataFieldShift, unionValue)
	sb = ThirdPartyCapDescriptorBuilder(nsb)
	return
}

type Payload capnpser.Struct

const (
	payload_content_ptrField  = 0
	payload_capTable_ptrField = 1
)

func (s *Payload) Content() (res AnyPointer, err error) {
	err = (*capnpser.Struct)(s).ReadAnyPointer(payload_content_ptrField, &res)
	return
}

func (s *Payload) CapTable() (res capnpser.GenericStructList[CapDescriptor], err error) {
	return capnpser.ReadGenericStructList[CapDescriptor]((*capnpser.Struct)(s), payload_capTable_ptrField)
}

type PayloadBuilder capnpser.StructBuilder

func (b *PayloadBuilder) SetContent(v capnpser.AnyPointerBuilder) error {
	return (*capnpser.StructBuilder)(b).SetAnyPointer(payload_content_ptrField, v)
}

func (b *PayloadBuilder) NewCapTable(listLen, listCap int) (capnpser.GenericStructListBuilder[CapDescriptorBuilder], error) {
	objSize := capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	return capnpser.NewGenericStructListBuilderField[CapDescriptorBuilder]((*capnpser.StructBuilder)(b), payload_capTable_ptrField, objSize, listLen, listCap)
}

const (
	exception_reason_ptrField = 0
)

type Exception capnpser.Struct

func (s *Exception) Reason() string {
	return (*capnpser.Struct)(s).String(exception_reason_ptrField)
}

type ExceptionBuilder capnpser.StructBuilder

func (b *ExceptionBuilder) SetReason(v string) error {
	return (*capnpser.StructBuilder)(b).SetString(exception_reason_ptrField, v)
}

type Return_Which int

const (
	Return_Which_Results   Return_Which = 2
	Return_Which_Exception Return_Which = 3
)

type Return capnpser.Struct

const (
	return_answerId_dataField       = 0
	return_answerId_dataFieldShift  = capnpser.Uint32FieldLo
	return_noFinishNeeded_dataField = 0
	return_noFinishNeeded_bit       = 33
	return_union_dataField          = 0
	return_union_dataFieldShift     = capnpser.Uint16FieldShift3
	return_union_ptrField           = 0
)

func (s *Return) AnswerId() AnswerId {
	return AnswerId((*capnpser.Struct)(s).Uint32(return_answerId_dataField, return_answerId_dataFieldShift))
}

func (s *Return) NoFinishNeeded() bool {
	return (*capnpser.Struct)(s).Bool(return_noFinishNeeded_dataField, return_noFinishNeeded_bit)
}

func (s *Return) Which() Return_Which {
	return Return_Which((*capnpser.Struct)(s).Uint16(return_union_dataField, return_union_dataFieldShift))
}

func (s *Return) AsResults() (res Payload, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(return_union_ptrField, (*capnpser.Struct)(&res))
	return
}

func (s *Return) AsException() (res Exception, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(return_union_ptrField, (*capnpser.Struct)(&res))
	return
}

type ReturnBuilder capnpser.StructBuilder

func (b *ReturnBuilder) SetAnswerId(v AnswerId) error {
	return (*capnpser.StructBuilder)(b).SetUint32(return_answerId_dataField, return_answerId_dataFieldShift, uint32(v))
}

func (b *ReturnBuilder) SetNoFinishNeeded(v bool) error {
	return (*capnpser.StructBuilder)(b).SetBool(return_noFinishNeeded_dataField, return_noFinishNeeded_bit, v)
}

func (b *ReturnBuilder) NewResults() (sb PayloadBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 0, PointerSectionSize: 2}
	const unionValue = uint16(Return_Which_Results)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(return_union_ptrField, structSize, return_union_dataField, return_union_dataFieldShift, unionValue)
	sb = PayloadBuilder(nsb)
	return
}

func (b *ReturnBuilder) NewException() (sb ExceptionBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 3}
	const unionValue = uint16(Return_Which_Exception)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(return_union_ptrField, structSize, return_union_dataField, return_union_dataFieldShift, unionValue)
	sb = ExceptionBuilder(nsb)
	return
}

const (
	finish_questionId_dataField      = 0
	finish_questionId_dataFieldShift = capnpser.Uint32FieldLo
)

type Finish capnpser.Struct

func (s *Finish) QuestionId() QuestionId {
	return QuestionId((*capnpser.Struct)(s).Uint32(finish_questionId_dataField, finish_questionId_dataFieldShift))
}

type FinishBuilder capnpser.StructBuilder

func (b *FinishBuilder) SetQuestionId(v QuestionId) error {
	return (*capnpser.StructBuilder)(b).SetUint32(finish_questionId_dataField, finish_questionId_dataFieldShift, uint32(v))
}

type Message_Which int

const (
	Message_Which_Call      Message_Which = 2
	Message_Which_Return    Message_Which = 3
	Message_Which_Finish    Message_Which = 4
	Message_Which_Bootstrap Message_Which = 8
)

func (w Message_Which) String() string {
	switch w {
	case Message_Which_Call:
		return "call"
	case Message_Which_Return:
		return "return"
	case Message_Which_Finish:
		return "finish"
	case Message_Which_Bootstrap:
		return "bootstrap"
	default:
		return fmt.Sprintf("unknown which %d", w)
	}
}

type Message capnpser.Struct

const (
	message_union_dataField      = 0
	message_union_dataFieldShift = capnpser.Uint16FieldShift0
	message_union_ptrField       = 0
)

func (s *Message) Which() Message_Which {
	return Message_Which((*capnpser.Struct)(s).Uint16(message_union_dataField, message_union_dataFieldShift))
}

func (s *Message) ReadFromRoot(msg *capnpser.Message) error {
	return msg.ReadRoot((*capnpser.Struct)(s))
}

func (s *Message) AsBootstrap() (res Bootstrap, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(messageTarget_union_ptrField, (*capnpser.Struct)(&res))
	return
}

func (s *Message) AsCall() (res Call, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(messageTarget_union_ptrField, (*capnpser.Struct)(&res))
	return
}

func (s *Message) AsFinish() (res Finish, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(messageTarget_union_ptrField, (*capnpser.Struct)(&res))
	return
}

func (s *Message) AsReturn() (res Return, err error) {
	err = (*capnpser.Struct)(s).ReadStruct(messageTarget_union_ptrField, (*capnpser.Struct)(&res))
	return
}

type MessageBuilder capnpser.StructBuilder

func (b *MessageBuilder) NewBoostrap() (sb BootstrapBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	const unionValue = uint16(Message_Which_Bootstrap)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(messageTarget_union_ptrField, structSize, message_union_dataField, message_union_dataFieldShift, unionValue)
	sb = BootstrapBuilder(nsb)
	return
}

func (b *MessageBuilder) NewCall() (sb CallBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 3, PointerSectionSize: 3}
	const unionValue = uint16(Message_Which_Call)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(messageTarget_union_ptrField, structSize, message_union_dataField, message_union_dataFieldShift, unionValue)
	sb = CallBuilder(nsb)
	return
}

func (b *MessageBuilder) NewFinish() (sb FinishBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 0}
	const unionValue = uint16(Message_Which_Finish)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(messageTarget_union_ptrField, structSize, message_union_dataField, message_union_dataFieldShift, unionValue)
	sb = FinishBuilder(nsb)
	return
}

func (b *MessageBuilder) NewReturn() (sb ReturnBuilder, err error) {
	var structSize = capnpser.StructSize{DataSectionSize: 3, PointerSectionSize: 3}
	const unionValue = uint16(Message_Which_Return)

	var nsb capnpser.StructBuilder
	nsb, err = (*capnpser.StructBuilder)(b).NewStructAsUnionValue(messageTarget_union_ptrField, structSize, message_union_dataField, message_union_dataFieldShift, unionValue)
	sb = ReturnBuilder(nsb)
	return
}

func NewRootMessageBuilder(mb *capnpser.MessageBuilder) (MessageBuilder, error) {
	var structSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 1}
	b, err := mb.NewRootStruct(structSize)
	return MessageBuilder(b), err
}
