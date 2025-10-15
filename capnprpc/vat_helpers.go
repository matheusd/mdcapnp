// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

// newReturnPayload is a hlper that initializes a Return message with a Result
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
