// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"errors"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

// newReturnPayload is a helper that initializes a Return message with a Result
// payload.
func (v *Vat) newReturnPayload(aid AnswerId) (rpcMsgBuilder rpcMsgBuilder, payBuilder types.PayloadBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	// Reply is a Return with a single cap.
	reply, err := rpcMsgBuilder.mb.NewReturn()
	if err != nil {
		return
	}
	reply.SetAnswerId(aid)
	payBuilder, err = reply.NewResults()
	return
}

// newSingleCapReturn is a helper that initializes a Return message with a
// single cap as the payload.
func (v *Vat) newSingleCapReturn(aid AnswerId) (rpcMsgBuilder rpcMsgBuilder, capDesc types.CapDescriptorBuilder, err error) {
	var payBuilder types.PayloadBuilder
	rpcMsgBuilder, payBuilder, err = v.newReturnPayload(aid)
	if err != nil {
		return
	}

	// Reply is a Return with a single cap.
	if err = payBuilder.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0)); err != nil {
		return
	}

	var capTable capnpser.GenericStructListBuilder[types.CapDescriptorBuilder]
	capTable, err = payBuilder.NewCapTable(1, 1)
	if err != nil {
		return
	}
	capDesc = capTable.At(0)
	return
}

func (v *Vat) newFinish(qid QuestionId) (rpcMsgBuilder rpcMsgBuilder, fin types.FinishBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	fin, err = rpcMsgBuilder.mb.NewFinish()
	if err != nil {
		return
	}

	err = fin.SetQuestionId(qid)
	return
}

func (v *Vat) newResolve(promiseId ExportId) (rpcMsgBuilder rpcMsgBuilder, res types.ResolveBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	res, err = rpcMsgBuilder.mb.NewResolve()
	if err != nil {
		return
	}

	err = res.SetPromiseId(promiseId)
	return
}

func (v *Vat) newProvide(recIdToCopy capnpser.AnyPointer) (rpcMsgBuilder rpcMsgBuilder, res types.ProvideBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	res, err = rpcMsgBuilder.mb.NewProvide()
	if err != nil {
		return
	}

	var recData capnpser.AnyPointerBuilder
	recData, err = capnpser.DeepCopy(recIdToCopy, rpcMsgBuilder.serMb)
	if err != nil {
		return
	}
	if err = res.SetRecipient(recData); err != nil {
		return
	}

	return
}

func (v *Vat) newAccept(acceptQid QuestionId, provId capnpser.AnyPointer, embargo bool) (rpcMsgBuilder rpcMsgBuilder, acc types.AcceptBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	acc, err = rpcMsgBuilder.mb.NewAccept()
	if err != nil {
		return
	}

	if err = acc.SetQuestionId(acceptQid); err != nil {
		return
	}
	if err = acc.SetEmbargo(embargo); err != nil {
		return
	}
	provIdCopy, err := capnpser.DeepCopy(provId, rpcMsgBuilder.serMb)
	if err != nil {
		return
	}
	err = acc.SetProvision(provIdCopy)
	return
}

func (v *Vat) newDisembargo(target messageTarget) (rpcMsgBuilder rpcMsgBuilder, dis types.DisembargoBuilder, err error) {
	rpcMsgBuilder, err = v.mbp.get()
	if err != nil {
		return
	}

	dis, err = rpcMsgBuilder.mb.NewDisembargo()
	if err != nil {
		return
	}

	var tgt types.MessageTargetBuilder
	tgt, err = dis.NewTarget()
	if err != nil {
		return
	}

	switch {
	case target.isImportedCap:
		err = tgt.SetImportedCap(target.impCap)
	case target.isPromisedAnswer:
		var pans types.PromisedAnswerBuilder
		pans, err = tgt.NewPromisedAnswer()
		if err == nil {
			err = pans.SetQuestionId(target.pansQid)
		}
	default:
		err = errors.New("unhandled case in newDisembargo")
	}
	return
}
