// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

const level0BootQid = 1

type Level0CallMsg struct {
	SerMsg *capnpser.MessageBuilder
	call   types.CallBuilder
	msg    types.MessageBuilder
}

type Level0ClientCfg struct {
	Conn conn
}

type Level0ClientVat struct {
	cfg  Level0ClientCfg
	conn conn

	bootImportId ImportId
	lastQid      QuestionId
	serMb        *capnpser.MessageBuilder
}

func NewLevel0ClientVat(cfg Level0ClientCfg) *Level0ClientVat {
	serMb, err := capnpser.NewMessageBuilder(capnpser.DefaultSimpleSingleAllocator)
	if err != nil {
		panic(err)
	}
	c := &Level0ClientVat{
		cfg:     cfg,
		serMb:   serMb,
		conn:    cfg.Conn,
		lastQid: level0BootQid + 1,
	}
	return c
}

func (c *Level0ClientVat) NextCallMsg(iid InterfaceId, mid MethodId) Level0CallMsg {
	if c.bootImportId == 0 {
		// TODO: support this.
		panic("wait for bootstrap first")
	}

	c.lastQid++
	if c.lastQid < level0BootQid {
		c.lastQid = level0BootQid + 1
	}

	if err := c.serMb.Reset(); err != nil {
		panic(err)
	}
	msg, err := types.NewRootMessageBuilder(c.serMb)
	if err != nil {
		panic(err)
	}
	call, err := msg.NewCall()
	if err != nil {
		panic(err)
	}
	call.SetQuestionId(c.lastQid)
	call.SetInterfaceId(uint64(iid))
	call.SetMethodId(uint16(mid))

	target, err := call.NewTarget()
	if err != nil {
		panic(err)
	}
	if c.bootImportId != 0 {
		target.SetImportedCap(c.bootImportId)
	} else {
		pans, err := target.NewPromisedAnswer()
		if err != nil {
			panic(err)
		}
		pans.SetQuestionId(level0BootQid)
	}
	return Level0CallMsg{
		SerMsg: c.serMb,
		call:   call,
		msg:    msg,
	}
}

func (c *Level0ClientVat) recvBoot(ctx context.Context) error {
	inMsg, err := c.conn.receive(ctx)
	if err != nil {
		return err
	}

	var msg types.Message
	if err := msg.ReadFromRoot(&inMsg.Msg); err != nil {
		return err
	}

	if msg.Which() != types.Message_Which_Return {
		return errors.New("message was not a Return after sending Bootstrap")
	}

	ret, err := msg.AsReturn()
	if err != nil {
		return err
	}
	if ret.AnswerId() != level0BootQid {
		return errors.New("answerId was not the expected one for bootstrap")
	}

	if ret.Which() != types.Return_Which_Results {
		return errors.New("return which was not for results")
	}

	pay, err := ret.AsResults()
	if err != nil {
		return err
	}
	content, err := pay.Content()
	if err != nil {
		return err
	}
	if !content.IsCapPointer() {
		return errors.New("bootstrap return content was not a cap pointer")
	}
	capPtr := content.AsCapPointer()

	capTable, err := pay.CapTable()
	if err != nil {
		return err
	}

	capIndex := int(capPtr.Index())
	if capIndex >= capTable.Len() {
		return errors.New("cap pointer not present in cap table")
	}

	capDescr := capTable.At(capIndex)
	if capDescr.Which() != types.CapDescriptor_Which_SenderHosted {
		return errors.New("level 0 only supports SenderHosted bootstrap cap")
	}
	c.bootImportId = ImportId(capDescr.AsSenderHosted())
	return nil
}

func (c *Level0ClientVat) WaitBootstrap(ctx context.Context) error {
	if err := c.serMb.Reset(); err != nil {
		return err
	}
	msg, err := types.NewRootMessageBuilder(c.serMb)
	if err != nil {
		return err
	}
	boot, err := msg.NewBoostrap()
	if err != nil {
		return err
	}
	boot.SetQuestionId(level0BootQid)

	err = c.conn.send(ctx, OutMsg{Msg: c.serMb})
	if err != nil {
		return err
	}

	return c.recvBoot(ctx)
}

func (c *Level0ClientVat) Call(ctx context.Context) (capnpser.AnyPointer, error) {
	var anyp capnpser.AnyPointer
	if c.bootImportId == 0 {
		// TODO: support this.
		return anyp, errors.New("wait for bootstrap first")
	}

	err := c.conn.send(ctx, OutMsg{Msg: c.serMb})
	if err != nil {
		return anyp, err
	}

	// Wait return.
	inMsg, err := c.conn.receive(ctx)
	if err != nil {
		return anyp, err
	}

	var msg types.Message
	if err := msg.ReadFromRoot(&inMsg.Msg); err != nil {
		return anyp, err
	}

	if msg.Which() != types.Message_Which_Return {
		return anyp, errors.New("message was not a Return after sending Bootstrap")
	}

	ret, err := msg.AsReturn()
	if err != nil {
		return anyp, err
	}
	if ret.AnswerId() != types.AnswerId(c.lastQid) {
		return anyp, errors.New("answerId was not the expected one for bootstrap")
	}

	if ret.Which() != types.Return_Which_Results {
		return anyp, errors.New("return which was not for results")
	}

	pay, err := ret.AsResults()
	if err != nil {
		return anyp, err
	}
	content, err := pay.Content()
	if err != nil {
		return anyp, err
	}
	if content.IsCapPointer() {
		return anyp, errors.New("CapPointer not supported as content in level0")
	}
	return content, nil
}
