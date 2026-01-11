// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"errors"
	"fmt"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

// newReturn is a helper that initializes a Return message.
func (v *Vat) newReturn(aid AnswerId) (outMsg outMsg, retBuilder types.ReturnBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	// Reply is a Return.
	retBuilder, err = outMsg.mb.NewReturn()
	if err != nil {
		return
	}
	retBuilder.SetAnswerId(aid)
	return
}

// newReturnPayload is a helper that initializes a Return message with a Result
// payload.
func (v *Vat) newReturnPayload(aid AnswerId) (outMsg outMsg, retBuilder types.ReturnBuilder, payBuilder types.PayloadBuilder, err error) {
	outMsg, retBuilder, err = v.newReturn(aid)
	if err != nil {
		return
	}

	payBuilder, err = retBuilder.NewResults()
	return
}

// newReturnException is a helper that initializes a Return message with an
// Exception payload.
func (v *Vat) newReturnException(aid AnswerId) (outMsg outMsg, ex types.ExceptionBuilder, err error) {
	var ret types.ReturnBuilder
	outMsg, ret, err = v.newReturn(aid)
	if err != nil {
		return
	}

	ex, err = ret.NewException()
	return
}

// newSingleCapReturn is a helper that initializes a Return message with a
// single cap as the payload.
func (v *Vat) newSingleCapReturn(aid AnswerId) (outMsg outMsg, capDesc types.CapDescriptorBuilder, err error) {
	var payBuilder types.PayloadBuilder
	outMsg, _, payBuilder, err = v.newReturnPayload(aid)
	if err != nil {
		return
	}

	// Reply is a Return with a single cap.
	if err = payBuilder.SetContent(capnpser.CapPointerAsAnyPointerBuilder(0)); err != nil {
		return
	}

	capTable, err := payBuilder.NewCapTable(1, 1)
	if err != nil {
		return
	}
	capDesc = capTable.At(0)
	return
}

func (v *Vat) newFinish(qid QuestionId) (outMsg outMsg, fin types.FinishBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	fin, err = outMsg.mb.NewFinish()
	if err != nil {
		return
	}

	err = fin.SetQuestionId(qid)
	return
}

func (v *Vat) newResolve(promiseId ExportId) (outMsg outMsg, res types.ResolveBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	res, err = outMsg.mb.NewResolve()
	if err != nil {
		return
	}

	err = res.SetPromiseId(promiseId)
	return
}

func (v *Vat) newProvide(recIdToCopy capnpser.AnyPointer) (outMsg outMsg, res types.ProvideBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	res, err = outMsg.mb.NewProvide()
	if err != nil {
		return
	}

	var recData capnpser.AnyPointerBuilder
	recData, err = capnpser.DeepCopy(recIdToCopy, outMsg.serMsg)
	if err != nil {
		return
	}
	if err = res.SetRecipient(recData); err != nil {
		return
	}

	return
}

func (v *Vat) newAccept(acceptQid QuestionId, provId capnpser.AnyPointer, embargo bool) (outMsg outMsg, acc types.AcceptBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	acc, err = outMsg.mb.NewAccept()
	if err != nil {
		return
	}

	if err = acc.SetQuestionId(acceptQid); err != nil {
		return
	}
	if err = acc.SetEmbargo(embargo); err != nil {
		return
	}
	provIdCopy, err := capnpser.DeepCopy(provId, outMsg.serMsg)
	if err != nil {
		return
	}
	err = acc.SetProvision(provIdCopy)
	return
}

func (v *Vat) newDisembargo(target messageTarget) (outMsg outMsg, dis types.DisembargoBuilder, err error) {
	outMsg, err = v.mbp.get()
	if err != nil {
		return
	}

	dis, err = outMsg.mb.NewDisembargo()
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

func (v *Vat) exportNewCap(capServer CallHandler, parent CallFuture) (ExportId, error) {
	// Encompasses some of what is done in RespondAsSenderHostedCap and
	// inside processCall. Figure out how to abstract.
	rc := parent.step.value.GetValue().conn

	// FIXME: can't do this here (rc is neither guaranteed locked or unlocked).
	// Instead, this should be tracked in CallSetup and exported once rc is
	// locked in execStep (or something to that effect).
	eid, ok := rc.exports.nextID()
	if !ok {
		return 0, errors.New("no more exports allowed")
	}

	rc.exports.set(eid, export{
		typ:     exportTypeLocallyHosted,
		handler: capServer,
	})
	return eid, nil
}

func (v *Vat) importRemoteCapFromPayload(rc *runningConn, pay types.Payload) (cf CallFuture, err error) {
	payContent, err := pay.Content()
	if err != nil {
		err = fmt.Errorf("could not get payload content: %v", err)
		return
	}

	if !payContent.IsCapPointer() {
		err = fmt.Errorf("payload content is not a capability pointer")
		return
	}
	capPtr := payContent.AsCapPointer()

	capTable, err := pay.CapTable()
	if err != nil {
		return
	}
	if int(capPtr.Index()) >= capTable.Len() {
		err = fmt.Errorf("capTable of len %d does not contain pointer %d",
			capTable.Len(), capPtr.Index())
		return
	}

	// TODO: merge with similar code in processReturn and processResolve.
	//
	var resolveCap CapabilityClient
	var resolveId ImportId
	var resImport imprt
	capEntry := capTable.At(int(capPtr.Index()))
	switch capEntry.Which() {
	case types.CapDescriptor_Which_SenderHosted:
		// Resolved into a remote capability.
		resolveCap = CapabilityClient{eid: ExportId(capEntry.AsSenderHosted())}
		resolveId = ImportId(capEntry.AsSenderHosted())
		resImport = imprt{typ: importTypeSenderHosted}
	case types.CapDescriptor_Which_SenderPromise:
		err = errors.New("senderPromise not supported here yet")
		return
	case types.CapDescriptor_Which_ThirdPartyHosted:
		err = errors.New("thirdPartyHosted not supported here yet")
		return
	default:
		err = errors.New("unknown cap entry type")
		return
	}

	rc.log.Debug().
		Int("resIid", int(resolveId)).
		Str("resTyp", resImport.typ.String()).
		Msg("Importing capability")

	// Set value of pipeline step to an already known export.
	cf = newRootFutureCap(v)
	pv := pipelineStepStateValue{
		conn:  rc,
		iid:   resolveId,
		value: resolveCap,
	}
	cf.step.value.Set(pipeStepStateDone, pv)

	return
}
