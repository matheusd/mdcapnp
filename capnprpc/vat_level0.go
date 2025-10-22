// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

const level0BootQid = 1

type Level0ClientCfg struct {
	Conn conn
}

const (
	level0CallStageEmpty int64 = iota
	level0CallStageBuilding
	level0CallStageRunning
	level0CallStageWaitingBootstrap
)

type Level0ClientVat struct {
	cfg  Level0ClientCfg
	conn conn

	callSetupStage atomic.Int64 // Protection against misuse.
	csetup         CallSetup

	mbp *messageBuilderPool

	boot BootstrapFuture

	bootImportId ImportId
	lastQid      QuestionId
}

func NewLevel0ClientVat(cfg Level0ClientCfg) *Level0ClientVat {
	v := &Level0ClientVat{
		cfg:     cfg,
		conn:    cfg.Conn,
		lastQid: level0BootQid + 1,
		mbp:     newMessageBuilderPool(),
	}
	v.boot = BootstrapFuture(newRootFutureCap(v))
	v.csetup.callOutMsg.isBootstrap = true
	return v
}

func (v *Level0ClientVat) setupNextCall(csetup CallSetup) {
	if !v.callSetupStage.CompareAndSwap(level0CallStageEmpty, level0CallStageBuilding) {
		panic("cannot setup multiple concurrent calls in a level 0 vat client")
	}

	v.csetup = csetup
}

func (v *Level0ClientVat) execNextCall(ctx context.Context) (any, error) {
	if !v.callSetupStage.CompareAndSwap(level0CallStageBuilding, level0CallStageRunning) {
		// Disambiguate case where caller called Bootstrap().Wait().
		if !v.callSetupStage.CompareAndSwap(level0CallStageEmpty, level0CallStageWaitingBootstrap) {
			panic("next call in level 0 vat not in building stage prior to execution")
		}
		if !v.csetup.callOutMsg.isBootstrap {
			panic("current csetup.callOutMsg is not the bootstrap message")
		}
		err := v.sendRecvBoot(ctx)
		if err != nil {
			return nil, err
		}
		res := capability{eid: types.ExportId(v.bootImportId)}
		if !v.callSetupStage.CompareAndSwap(level0CallStageWaitingBootstrap, level0CallStageEmpty) {
			panic("broken assumption: callSetupStage is not level0CallStageWaitingBootstrap")
		}
		return res, nil
	}

	// Ensure we return the stage to level0CallStageRunning on exit.
	defer func() {
		if !v.callSetupStage.CompareAndSwap(level0CallStageRunning, level0CallStageEmpty) {
			panic("call stage was not the expected level0CallStageRunning at defer time")
		}
	}()

	call := types.CallBuilder(v.csetup.callOutMsg.builder)
	v.lastQid++
	if v.lastQid < level0BootQid {
		v.lastQid = level0BootQid + 1
	}
	call.SetQuestionId(v.lastQid)

	target, err := call.NewTarget()
	if err != nil {
		return nil, err
	}

	if v.bootImportId != 0 {
		target.SetImportedCap(v.bootImportId)
	} else {
		// Bootstrap not fetched yet. Pipeline the call to the bootstrap
		// cap, but expect the bootstrap Return first.
		if err := v.sendBoot(ctx); err != nil {
			return nil, err
		}

		pans, err := target.NewPromisedAnswer()
		if err != nil {
			return nil, err
		}
		pans.SetQuestionId(level0BootQid)
	}

	err = v.conn.send(ctx, OutMsg{Msg: v.csetup.callOutMsg.serMsg})
	if err != nil {
		return nil, err
	}

	// After sending, the request MessageBuilder is ready for reuse.
	v.mbp.put(v.csetup.callOutMsg.serMsg)
	wantReturnResults := v.csetup.WantReturnResults
	v.csetup = CallSetup{}

	// When the call was pipelined from the bootstrap, we expect the Return
	// from the Bootstrap first.
	if v.bootImportId == 0 {
		if err := v.recvBoot(ctx); err != nil {
			return nil, err
		}
	}

	// Wait return.
	inMsg, err := v.conn.receive(ctx)
	if err != nil {
		return nil, err
	}

	// Process return.
	var msg types.Message
	if err := msg.ReadFromRoot(&inMsg.Msg); err != nil {
		return nil, err
	}

	msgWhich := msg.Which()
	if msgWhich != types.Message_Which_Return {
		return nil, fmt.Errorf("reply for Call was not a Return: %s", msgWhich)
	}

	ret, err := msg.AsReturn()
	if err != nil {
		return nil, err
	}
	if ret.AnswerId() != types.AnswerId(v.lastQid) {
		return nil, errors.New("answerId was not the expected one for the last sent call")
	}

	if ret.Which() != types.Return_Which_Results {
		return nil, errors.New("return which was not for results")
	}

	pay, err := ret.AsResults()
	if err != nil {
		return nil, err
	}
	content, err := pay.Content()
	if err != nil {
		return nil, err
	}

	// Determine what to do with the actual response.

	if content.IsCapPointer() {
		return nil, errors.New("CapPointer not supported as content in level0")
	}

	var finalRes any
	if content.IsZeroStruct() || !wantReturnResults {
		// All done in this case.
		finalRes = struct{}{}
	} else {
		// Caller wants the reply data. Copy into a new message builder for them
		// to use it.
		mb := v.mbp.getRawMessageBuilder(0) // TODO: get size hint
		err = capnpser.DeepCopyAndSetRoot(content, mb)
		if err != nil {
			return nil, err
		}
		finalRes = mb
	}

	return finalRes, nil
}

func (c *Level0ClientVat) sendBoot(ctx context.Context) error {
	serMb := c.mbp.getRawMessageBuilder(0)
	msg, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		return err
	}
	boot, err := msg.NewBoostrap()
	if err != nil {
		return err
	}
	boot.SetQuestionId(level0BootQid)

	err = c.conn.send(ctx, OutMsg{Msg: serMb})
	if err != nil {
		return err
	}

	c.mbp.put(serMb)
	return nil
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

func (v *Level0ClientVat) sendRecvBoot(ctx context.Context) error {
	if err := v.sendBoot(ctx); err != nil {
		return err
	}
	if err := v.recvBoot(ctx); err != nil {
		return err
	}
	return nil
}

// WaitBootstrap sends the Bootstrap message and waits until it is responded by
// the remote side.
/*
func (v *Level0ClientVat) WaitBootstrap(ctx context.Context) error {
	if !v.callSetupStage.CompareAndSwap(level0CallStageEmpty, level0CallStageWaitingBootstrap) {
		return errors.New("cannot wait for bootstrap when there is an active call")
	}
	if err := v.sendRecvBoot(ctx); err != nil {
		return err
	}
	if !v.callSetupStage.CompareAndSwap(level0CallStageWaitingBootstrap, level0CallStageEmpty) {
		panic("broken assumption: stage is not level0CallStageWaitingBootstrap")
	}
	return nil
}
*/

// Bootstrap returns a reference to the future Bootstrap results.
func (v *Level0ClientVat) Bootstrap() BootstrapFuture {
	return v.boot
}
