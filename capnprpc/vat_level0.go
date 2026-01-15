// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/rs/zerolog"
	types "matheusd.com/mdcapnp/capnprpc/types"
)

const level0BootQid = 1

type Level0ClientCfg struct {
	Conn conn
	Log  *zerolog.Logger
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
	log  *zerolog.Logger

	callSetupStage atomic.Int64 // Protection against misuse.
	csetup         CallSetup

	mbp *MessageBuilderPool

	inMsg *InMsg

	boot BootstrapFuture

	bootImportId ImportId
	lastQid      QuestionId
}

func NewLevel0ClientVat(cfg Level0ClientCfg) *Level0ClientVat {
	v := &Level0ClientVat{
		cfg:     cfg,
		conn:    cfg.Conn,
		log:     cfg.Log,
		lastQid: level0BootQid + 1,
		mbp:     NewMessageBuilderPool(),
	}
	if v.log == nil {
		nop := zerolog.Nop()
		v.log = &nop
	}
	v.boot = BootstrapFuture(newRootFutureCap(v))
	v.csetup.callOutMsg.isBootstrap = true
	return v
}

func (v *Level0ClientVat) setupNextCall(csetup CallSetup) {
	v.log.Trace().Int("iid", int(csetup.InterfaceId)).Int("mid", int(csetup.MethodId)).Msg("Starting to setup next call")
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
		res := CapabilityClient{eid: types.ExportId(v.bootImportId)}
		if !v.callSetupStage.CompareAndSwap(level0CallStageWaitingBootstrap, level0CallStageEmpty) {
			panic("broken assumption: callSetupStage is not level0CallStageWaitingBootstrap")
		}
		return res, nil
	}

	// Ensure we return the stage to level0CallStageEmpty on exit.
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
		return nil, fmt.Errorf("unable to create target field in call: %v:", err)
	}

	logLine := v.log.Debug().
		Int("qid", int(v.lastQid)).
		Int("iid", int(v.csetup.InterfaceId)).
		Int("mid", int(v.csetup.MethodId))

	if v.bootImportId != 0 {
		target.SetImportedCap(v.bootImportId)
		logLine.Int("target.importedCap", int(v.bootImportId))
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
		logLine.Int("target.pans", int(level0BootQid))
	}

	logLine.Msg("Sending Call")
	err = v.conn.send(ctx, OutMsg{Msg: v.csetup.callOutMsg.serMsg})
	if err != nil {
		return nil, err
	}

	// After sending, the request MessageBuilder is ready for reuse.
	v.mbp.put(v.csetup.callOutMsg.serMsg)
	v.csetup = CallSetup{}

	// When the call was pipelined from the bootstrap, we expect the Return
	// from the Bootstrap first.
	if v.bootImportId == 0 {
		if err := v.recvBoot(ctx); err != nil {
			return nil, fmt.Errorf("error processing bootstrap before call Return: %v", err)
		}
	}

	// Wait return.
	v.inMsg, err = v.conn.receive(ctx, v.inMsg)
	if err != nil {
		return nil, err
	}

	// Process return.
	var msg types.Message
	if err := msg.ReadFromRoot(&v.inMsg.Msg); err != nil {
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
	if ansId := ret.AnswerId(); ansId != types.AnswerId(v.lastQid) {
		return nil, fmt.Errorf("answerId was not the expected one for "+
			"the last sent call (got %d, want %d)", ansId, v.lastQid)
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
	/*
		if resParser != nil {
			finalRes, err = resParser(pay)
			if err != nil {
				return nil, err
			}
		} else */
	if content.IsZeroStruct() {
		// All done in this case.
		finalRes = struct{}{}
	} else {
		// In this case, this Return contains a struct that the caller
		// will use. Thus, this will outlive this iteration of
		// execNextCall. Effectively, we're transfering ownership
		// of this message/buffer from the conn's inLoop to the caller,
		// by making inMsg nil.
		finalRes, v.inMsg = v.inMsg, nil
	}

	v.log.Debug().Int("aid", int(v.lastQid)).Msg("Received valid Return")

	return finalRes, nil
}

func (v *Level0ClientVat) sendBoot(ctx context.Context) error {
	serMb := v.mbp.getRawMessageBuilder(0)
	msg, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		return err
	}
	boot, err := msg.NewBoostrap()
	if err != nil {
		return err
	}
	boot.SetQuestionId(level0BootQid)

	v.log.Debug().Int("qid", level0BootQid).Msg("Sending bootstrap")
	err = v.conn.send(ctx, OutMsg{Msg: serMb})
	if err != nil {
		return err
	}

	v.mbp.put(serMb)
	return nil
}

func (v *Level0ClientVat) recvBoot(ctx context.Context) error {
	v.log.Trace().Msg("Waiting for bootstrap response on next conn.receive")
	inMsg, err := v.conn.receive(ctx, nil)
	if err != nil {
		return err
	}
	v.log.Trace().Msg("Received next message while waiting for bootstrap reply")

	var msg types.Message
	if err := msg.ReadFromRoot(&inMsg.Msg); err != nil {
		return err
	}

	if which := msg.Which(); which != types.Message_Which_Return {
		return fmt.Errorf("message was not a Return after sending "+
			"Bootstrap (got %d)", which)
	}

	ret, err := msg.AsReturn()
	if err != nil {
		return err
	}
	if ansId := ret.AnswerId(); ansId != level0BootQid {
		return fmt.Errorf("answerId was not the expected one for "+
			"bootstrap (got %d, want %d)", ansId, level0BootQid)
	}

	if which := ret.Which(); which != types.Return_Which_Results {
		return fmt.Errorf("return which was not for results (got %d)", which)
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
		return fmt.Errorf("unable to read capTable: %v", err)
	}

	capIndex := int(capPtr.Index())
	if capIndex >= capTable.Len() {
		return errors.New("cap pointer not present in cap table")
	}

	capDescr := capTable.At(capIndex)
	if capDescr.Which() != types.CapDescriptor_Which_SenderHosted {
		// TODO: We could support a promise, but then we can't send
		// anything until the promise is resolved.
		return errors.New("level 0 only supports SenderHosted bootstrap cap")
	}
	v.bootImportId = ImportId(capDescr.AsSenderHosted())
	v.log.Debug().Int("impId", int(v.bootImportId)).Msg("Received bootstrap return")
	return nil
}

func (v *Level0ClientVat) sendRecvBoot(ctx context.Context) error {
	if err := v.sendBoot(ctx); err != nil {
		return fmt.Errorf("error sending boot msg: %v", err)
	}
	if err := v.recvBoot(ctx); err != nil {
		return fmt.Errorf("error processing received bootstrap: %v", err)
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
