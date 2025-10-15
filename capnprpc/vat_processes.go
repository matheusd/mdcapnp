// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"weak"

	"github.com/rs/zerolog"
	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

func (v *Vat) processBootstrap(ctx context.Context, rc *runningConn, boot types.Bootstrap) error {
	bootQid := AnswerId(boot.QuestionId())

	rpcMsgBuilder, capDesc, err := v.newSingleCapReturn(bootQid)
	if err != nil {
		return err
	}
	if err := capDesc.SetSenderHosted(rc.bootExportId); err != nil {
		return err
	}

	// Modify answer table to track the bootrap export. bootExportId is set
	// during conn setup automatically.
	ans := answer{typ: answerTypeBootstrap, eid: rc.bootExportId}
	rc.answers.set(AnswerId(bootQid), ans)

	rc.log.Debug().
		Int("qid", int(bootQid)).
		Int("eid", int(rc.bootExportId)).
		Msg("Exported Bootstrap")

	rpcMsg := v.mp.get()
	rpcMsg.rawSerMb = rpcMsgBuilder.serMb
	return rc.queue(ctx, singleMsgBatch(rpcMsg))
}

func (v *Vat) processReturn(ctx context.Context, rc *runningConn, ret types.Return) error {
	qid := QuestionId(ret.AnswerId())
	q, ok := rc.questions.get(qid)
	if !ok {
		return fmt.Errorf("question %d not found", qid)
	}

	// TODO: support exception, cancel, etc
	if ret.Which() != types.Return_Which_Results {
		return fmt.Errorf("only results supported")
	}

	if q.typ == questionTypeProvide {
		res, err := ret.AsResults()
		if err != nil {
			return err
		}
		content, err := res.Content()
		if err != nil {
			return err
		}
		if !content.IsZeroStruct() {
			return fmt.Errorf("received non-void return for provide")
		}
		rc.log.Debug().Int("qid", int(qid)).Msg("Received Return message for prior Provide")
		return nil
	}

	step := q.step()
	if step == nil {
		rc.log.Warn().Int("qid", int(qid)).Msg("Received Return message for released step")
		return nil
	}

	// Go through cap table, modify imports table based on what was
	// exported by this call.
	//
	// TODO: only do this if the cap is referenced in the content?
	// Go through cap table, modify imports table based on what was
	// exported by this call.
	//
	// TODO: only do this if the cap is referenced in the content?
	payload, err := ret.AsResults()
	if err != nil {
		return err
	}
	capTable, err := payload.CapTable()
	if err != nil {
		return fmt.Errorf("error extracting cap table: %v", err)
	}
	for i := range capTable.Len() {
		entry := capTable.At(i)
		var iid ImportId
		var imp imprt
		switch entry.Which() {
		case types.CapDescriptor_Which_SenderHosted:
			iid = ImportId(entry.AsSenderHosted())
			imp = imprt{typ: importTypeSenderHosted}
		case types.CapDescriptor_Which_SenderPromise:
			iid = ImportId(entry.AsSenderPromise())
			imp = imprt{typ: importTypeRemotePromise, step: weak.Make(step)}
		default:
			return fmt.Errorf("unsupported capability type %d", entry.Which())
		}

		rc.imports.set(iid, imp)
		rc.log.Debug().
			Int("qid", int(qid)).
			Int("iid", int(iid)).
			Str("ityp", imp.typ.String()).
			Msg("Imported cap contained in Return")
	}

	// Get contents of result.
	var stepResult any
	var stepResultPromise ExportId
	var stepResultType string
	var stepImportId ImportId
	content, err := payload.Content()
	if err != nil {
		return err
	}
	if content.IsCapPointer() {
		// NOT GOOD. Must have a new type to pass along instead of
		// parsing like this (maybe). Think about embedded caps.
		cp := content.AsCapPointer()
		capIndex := cp.Index()
		if int(capIndex) >= capTable.Len() {
			return fmt.Errorf("capability referenced index outside cap table")
		}
		capEntry := capTable.At(int(capIndex))
		switch capEntry.Which() {
		case types.CapDescriptor_Which_SenderHosted:
			stepImportId = ImportId(capEntry.AsSenderHosted())
			stepResult = capability{eid: ExportId(stepImportId)}
			stepResultType = "senderHostedCap"
		case types.CapDescriptor_Which_SenderPromise:
			stepImportId = ImportId(capEntry.AsSenderPromise())
			stepResultPromise = ExportId(stepImportId)
			stepResultType = "senderPromise"
		default:
			return errors.New("unknown cap entry type")
		}
	} else if content.IsZeroStruct() {
		stepResult = struct{}{}
		stepResultType = "void"
	} else if content.IsStruct() {
		// TODO: copy if its a struct? Or release serialized message if
		// content is just a cap (because it's not needed anymore)?
		// stepResult = content.AsStruct()
		stepResult = struct{}{}
		stepResultType = "struct"
	} else {
		return errors.New("unknown/unimplemented content type")
	}

	// Automatically remove question if result contains no caps and
	// noFinishNeeded was set.
	noFinishNeeded := ret.NoFinishNeeded()
	if noFinishNeeded {
		// TODO: need to validate if this is ok?
		rc.questions.del(qid)
	} else if q.strongStep != nil {
		// Question is holding own to a strong ref to pipe, but now we
		// switch to a weak ref. We only do this if finish is needed,
		// because noFinishNeeded==true means this has no pipelinable
		// capabilities and finish won't need to be sent.
		//
		// Doing this conditionally avoids an allocation inside
		// weak.MakePointer.
		q.weakStep = weak.Make(step)
		q.strongStep = nil
	}

	// Fulfill pipeline waiting for this result.
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}

		if ov.conn != rc {
			return os, ov, fmt.Errorf("broken assumption: pipeline "+
				"step has conn %s when being modified by conn %s",
				ov.conn.String(), rc.String())
		}

		rc.log.Debug().
			Int("qid", int(qid)).
			Str("rtyp", stepResultType).
			Bool("noFinishNeeded", noFinishNeeded).
			Msg("Processed Return message")

		ov.iid = stepImportId
		if stepResultPromise > 0 {
			// This step isn't resolved into a concrete exported
			// cap or struct yet. Keep it running.
			return os, ov, nil
		}

		ov.value = stepResult

		// If no finish is needed (no pipelining allowed from this),
		// set as finished.
		if noFinishNeeded {
			ov.qid = 0
		}
		return pipeStepStateDone, ov, nil
	})
}

var errCallWithoutId = errors.New("call with zero interfaceId and methodId")

func (v *Vat) processCall(ctx context.Context, rc *runningConn, c types.Call) error {
	iid, mid := c.InterfaceId(), c.MethodId()
	if iid == 0 && mid == 0 {
		// Only bootstrap is allowed to have iid+mid == 0.
		return errCallWithoutId
	}

	qid := c.QuestionId()
	if rc.answers.has(AnswerId(qid)) {
		return fmt.Errorf("remote already asked question %d", qid)
	}

	logEvent := rc.log.Trace().
		Int("qid", int(qid))

	// Determine the target of this call (either an exported cap or a
	// promised answer).
	target, err := c.Target()
	if err != nil {
		return fmt.Errorf("unable to read target from call: %v", err)
	}
	var eid ExportId
	switch target.Which() {
	case types.MessageTarget_Which_PromisedAnswer:
		// Promised answers are in the answer table.
		//
		// TODO: Recursively track it down if the answer is another
		// promise.
		pans, err := target.AsPromisedAnswer()
		if err != nil {
			return fmt.Errorf("unable to decode promised answer: %v", err)
		}
		pansQid := pans.QuestionId()
		q, ok := rc.answers.get(AnswerId(pansQid))
		if !ok {
			return fmt.Errorf("call referenced unknown promised answer %d", pansQid)
		}

		logEvent.Int("pansQid", int(pansQid))

		eid = q.eid // What about promises?
	case types.MessageTarget_Which_ImportedCap:
		impCap := target.AsImportedCap()
		eid = ExportId(impCap)
		logEvent.Int("impCap", int(impCap))
	default:
		return errors.New("unsupported call target")
	}

	exp, ok := rc.exports.get(eid)
	if !ok {
		return fmt.Errorf("call message target determined to be unset export %d", eid)
	}

	if exp.typ != exportTypeLocallyHosted {
		return fmt.Errorf("unsupported export type %d", exp.typ)
	}
	logEvent.Int("eid", int(eid)).Str("etype", exp.typ.String())

	// TODO: proxy calls when exp.typ == exportTypeThirdPartyExport.

	callArgs := callHandlerArgs{
		iid: interfaceId(iid),
		mid: methodId(mid),
		// params: c.params, // FIXME how?
		rc: rc,
	}

	// Start preparing reply.
	rpcMsgBuilder, err := v.mbp.getForPayloadSize(0) // TODO: size hint?
	if err != nil {
		return err
	}
	reply, err := rpcMsgBuilder.mb.NewReturn()
	if err != nil {
		return err
	}
	reply.SetAnswerId(AnswerId(qid))

	crb := &rc.crb // Ok to reuse (rc is locked).
	crb.pb, err = reply.NewResults()
	crb.serMb = rpcMsgBuilder.serMb
	if err != nil {
		return err
	}
	crb.payload = payload{content: anyPointer{
		isVoid: true, // Void result by default on non-error.
	}}

	// Make the call!
	logEvent.Msg("Locally handling call")
	err = exp.handler.Call(rc.ctx, callArgs, crb)
	if ex, ok := err.(callExceptionError); ok {
		// When an exception that will be sent remotely is detected,
		// re-create the reply. This ensures anything written to the
		// payload inside the handler's Call() method will *NOT* be
		// sent as an orphan object inside the reply.
		if err := rpcMsgBuilder.serMb.Reset(); err != nil {
			return err
		}

		// Turn the error into a returned exception.
		reply, err = rpcMsgBuilder.mb.NewReturn()
		if err != nil {
			return err
		}
		reply.SetAnswerId(AnswerId(qid))

		exc, err := reply.NewException()
		if err != nil {
			return err
		}
		exc.SetReason(err.Error()) // TODO: send more details.
		_ = ex

		/*
			reply.isReturn = true
			reply.ret.aid = AnswerId(qid)
			reply.ret.isException = true
			reply.ret.exc = ex.ToException()
		*/

		rc.log.Debug().
			Int("qid", int(qid)).
			Int("eid", int(eid)).
			//Dict("ex", zerolog.Dict().
			//Int("type", reply.ret.exc.typ).
			//Str("reason", reply.ret.exc.reason)).
			Msg("Processed call into exception")

	} else if err != nil {
		// Fatal connection error.
		return err
	} else {
		// No finish is needed if the call promised no pipelining or
		// the result has no capabilities (i.e. the result is not
		// callable).
		//
		// TODO: Maybe track caps added to reply instead of reading from
		// message?
		capTable := crb.readReturnCapTable()
		noFinishNeeded := c.NoPromisePipelining() || capTable.Len() == 0

		if !noFinishNeeded {
			// If the result is a capability, the answer is
			// pipelinable, so track where the corresponding export
			// will be.
			var rootCapIndex int = -1
			var rootCapExportId ExportId
			if crb.payload.content.IsCapPointer() {
				rootCapIndex = int(crb.payload.content.AsCapPointer().index)
			}

			// Track all exported caps.
			for i := range capTable.Len() {
				cp := capTable.At(i)

				var capEid ExportId
				var isSenderPromise bool
				switch cp.Which() {
				case types.CapDescriptor_Which_SenderHosted:
					capEid = cp.AsSenderHosted()
				case types.CapDescriptor_Which_SenderPromise:
					capEid = cp.AsSenderPromise()
					isSenderPromise = true
				default:
					return fmt.Errorf("unsupported cap type in processReturn: %v", cp.Which())
				}

				if capExp, ok := rc.exports.get(capEid); ok {
					// TODO: take a pointer instead?
					capExp.refCount++
					rc.exports.set(capEid, capExp)
				} else if isSenderPromise { // here sender == local vat.
					capExp = export{typ: exportTypePromise, refCount: 1}
					rc.exports.set(capEid, capExp)

					rc.log.Trace().
						Int("eid", int(capEid)).
						Str("typ", "senderPromise").
						Msg("Exporting capability")
				} else {
					return errors.New("other types of exports not implemented")
				}

				// If this is the root cap index, then track this export
				// directly in the answer (for future pipelined calls).
				if i == rootCapIndex {
					rootCapExportId = capEid
				}
			}

			// Save the answer in the answers table.
			ans := answer{typ: answerTypeCall, eid: rootCapExportId}
			rc.answers.set(AnswerId(qid), ans)
		} else {
			reply.SetNoFinishNeeded(noFinishNeeded)
		}

		rc.log.Debug().
			Int("qid", int(qid)).
			Int("eid", int(eid)).
			Msg("Processed call into payload result")
	}

	rpcMsg := v.mp.get()
	rpcMsg.rawSerMb = rpcMsgBuilder.serMb
	return rc.queue(ctx, singleMsgBatch(rpcMsg))
}

func (v *Vat) processFinish(ctx context.Context, rc *runningConn, fin types.Finish) error {
	var err error
	aid := AnswerId(fin.QuestionId())

	if !rc.answers.has(aid) {
		err = fmt.Errorf("answer %d not in answers table", aid)
	} else {
		rc.answers.del(aid)
	}

	// TODO: release exported caps?

	if err == nil {
		rc.log.Debug().
			Int("aid", int(aid)).
			Msg("Removed answer due to Finish message")
	}

	return err
}

func (v *Vat) resolveThirdPartyCapForStep(ctx context.Context, step *pipelineStep,
	srcConn *runningConn, tphCapDesc thirdPartyCapDescriptor) error {

	rc, provId, err := v.connectToIntroduced3rdParty(ctx, srcConn, tphCapDesc)
	if err != nil {
		return err
	}

	if rc == srcConn {
		// TODO: Is this permitted?
		return errors.New("same conns for 3PH")
	}

	// We need a lock on both conns to avoid logical races. But other
	// goroutines may also be attempting to (or holding) locks on either (or
	// both). So we need to ensure _every_ goroutine always locks _every_
	// conn on the same order.
	//
	// One example of a logical race that can happen if this isn't done, is
	// one where the Return for the Accept sent to third party is received
	// and processed _before_ the pipeline step has been modified to point
	// to the promised Accept answer: this causes an error while processing
	// the return because of a broken assumption from where the next
	// response would be coming from.
	//
	// Another potential race is a Call being sent after a Disembargo has
	// been queued (see the comment tagged as 3PHCONNISSUE).
	mu := makeTwoConnLocker(rc, srcConn)
	mu.lock()
	defer mu.unlock()

	v.log.Trace().
		Str("src", srcConn.String()).
		Str("dst", rc.String()).
		Msg("Shortening path after 3PH introduction")

	// Send Accept() with embargo set to the new remote.
	acceptQid, ok := rc.questions.nextID()
	if !ok {
		return errTooManyOpenQuestions
	}

	// TODO: finalizer???
	q := question{weakStep: weak.Make(step)}
	rc.questions.set(acceptQid, q)
	accept := accept{
		qid:       acceptQid,
		provision: provId,

		// In the future, this could be dynamically determined, because
		// the local vat can know whether there are pending pipelined
		// calls or not.
		embargo: true,
	}

	rc.log.Debug().
		Int("qid", int(acceptQid)).
		Str("srcConn", srcConn.String()).
		Msg("Sending accept for 3PH introduction")

	// Explicitly wait until the accept is in the third party vat (as acked
	// under the assumption that send acks bytes received by the remote
	// party).
	//
	// Note that there is _still_ room for a race here, if the Accept is
	// processed slower than the Disembargo goes from the local vat (Alice)
	// to the provider (Bob) and forwarded to the third party (Carol),
	// because the Disembargo would arrive before the Accept.
	//
	// Note that this can also significantly increase contention and reduce
	// performance for conns, because the conns to both Bob and Carol are
	// locked here. We assume 3PH is rare enough that this won't be
	// significant.
	if err := v.sendAccept(ctx, rc, accept); err != nil {
		return err
	}

	var disembargoTarget messageTarget
	err = step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os == pipeStepFailed {
			return os, ov, ov.err
		}

		// Sanity check.
		if ov.conn != srcConn {
			return os, ov, fmt.Errorf("unexpected conn: wanted conn to %s but got %s",
				srcConn, ov.conn)
		}

		if ov.iid > 0 {
			// Already received the Return for the original call,
			// the Return was for a remote promise and that remote
			// promise has now resolved to a third party.
			disembargoTarget.isImportedCap = true
			disembargoTarget.impcap = ov.iid
			v.log.Info().
				Int("iid", int(ov.iid)).
				Str("src", srcConn.String()).
				Str("dst", rc.String()).
				Int("acceptId", int(acceptQid)).
				Msg("Path-shortened imported cap pipeline step to third party")
		} else {
			// The Return (for a promised answer) already signalled
			// that the returned cap is in a third party.
			disembargoTarget.isPromisedAnswer = true
			disembargoTarget.pans.qid = ov.qid
			v.log.Info().
				Int("srcQid", int(ov.qid)).
				Str("src", srcConn.String()).
				Str("dst", rc.String()).
				Int("acceptId", int(acceptQid)).
				Msg("Path-shortened promised answer pipeline step to third party")
		}

		ov.conn = rc
		ov.qid = acceptQid
		return os, ov, nil
	})
	if err != nil {
		return err
	}

	// Pipelined calls made from this moment on will be sent to the third
	// party as a promised answer to the Accept. However, those calls will
	// be cached there, until the Disembargo travels through the proxy
	// (signalling that all previously pipelined calls were sent).

	// If accept did not signal existence of embargoed pipelined calls,
	// Disembargo isn't needed.
	if !accept.embargo {
		return nil
	}

	// Queue the Disembargo in the old conn, which will travese the network
	// all the way to the third party. This will make the third party
	// process the calls sent directly from the local vat (as opposed to
	// those proxied by the original conn).
	dis := &message{isDisembargo: true, disembargo: disembargo{ // TODO: fetch from v.mp
		isAccept: true,
		target:   disembargoTarget,
	}}
	if err := srcConn.queue(ctx, singleMsgBatch(dis)); err != nil {
		return err
	}

	// Only thing left to complete 3PH now is wait for the Return message
	// that completes the Accept call above.

	return nil
}

func (v *Vat) processResolve(ctx context.Context, rc *runningConn, res resolve) error {
	iid := ImportId(res.pid)
	imp, ok := rc.imports.get(iid)
	if !ok {
		return fmt.Errorf("import id %d not found", iid)
	}

	if imp.typ != importTypeRemotePromise {
		return fmt.Errorf("import %d is not a remote promise to be resolved", iid)
	}

	step := imp.step.Value()
	if step == nil {
		return fmt.Errorf("step of import %d already released", iid)
	}

	// Similar to code in processReturn. Unify?
	capEntry := res.cap
	var resolveCap capability
	var resolvePromise ExportId
	var resolveId ImportId
	var resImport imprt
	if capEntry.IsSenderHosted() {
		// Resolved into a remote capability.
		resolveCap = capability{eid: ExportId(capEntry.AsSenderHosted())}
		resolveId = ImportId(capEntry.AsSenderHosted())
		resImport = imprt{typ: importTypeSenderHosted}
	} else if capEntry.IsSenderPromise() {
		// Resolved into another remote promise.
		resolvePromise = capEntry.AsSenderPromise()
		resolveId = ImportId(resolvePromise)
		resImport = imprt{typ: importTypeRemotePromise, step: imp.step}
	} else if capEntry.IsThirdPartyHosted() {
		// Start the 3PH resolution process. This will involve
		// connecting to the third party (if we haven't already) and
		// then sending the Accept with the provision id.
		//
		// TODO: If this is only level 1, use vineId instead and proxy
		// requests.
		go func() {
			err := v.resolveThirdPartyCapForStep(ctx, step, rc, capEntry.AsThirdPartyHosted())
			if err != nil {
				rc.log.Err(err).Msg("resolveThirdPartyCapForStep failed")
			} else {
				rc.log.Trace().Msg("resolveThirdPartyCapForStep completed")
			}
		}()

		// This doesn't change the pipeline step (promise). Any
		// pipelined calls will be proxied through the existing conn
		// until 3PH completes (more specifically, until the connection
		// is completed and the Accept message is sent).
		//
		// Maybe in the future this could be changed to optinally
		// cache calls locally until 3PH completes.
		return nil
	} else {
		return errors.New("unknown cap entry type")
	}

	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}

		rc.log.Debug().
			Int("qid", int(ov.qid)).
			Int("iid", int(iid)).
			Int("resIid", int(resolveId)).
			Str("resTyp", resImport.typ.String()).
			Msg("Resolved remote promise")

		ov.iid = resolveId
		if resolvePromise > 0 {
			// This step isn't resolved into a concrete exported
			// cap or struct yet. Keep it running.
			return os, ov, nil
		}

		ov.value = resolveCap
		return pipeStepStateDone, ov, nil
	})
}

func (v *Vat) processProvide(ctx context.Context, rc *runningConn, prov provide) error {
	aid := AnswerId(prov.qid)
	if rc.answers.has(aid) {
		return fmt.Errorf("remote already asked question %d", aid)
	}

	if v.cfg.net == nil {
		return err3PHWithoutVatNetwork
	}

	// Build value that will hold the shared cap data.
	expAc := expectedAccept{
		srcConn:    rc,
		provideAid: aid,
	}

	logMsg := rc.log.Debug()

	// Check if the target exists exported to the caller.
	if prov.target.isImportedCap {
		capExp, ok := rc.exports.get(ExportId(prov.target.impcap))
		if !ok {
			return fmt.Errorf("export not found %d", prov.target.impcap)
		}

		expAc.handler = capExp.handler
		logMsg = logMsg.Str("typ", "importedCap").Int("iid", int(prov.target.impcap))
	} else if prov.target.isPromisedAnswer {
		if !rc.answers.has(AnswerId(prov.target.pans.qid)) {
			return fmt.Errorf("answer not found %d", prov.target.pans.qid)
		}

		logMsg = logMsg.Str("typ", "promisedAnswer").Int("iid", int(prov.target.pans.qid))

		// TODO: support
		return fmt.Errorf("unsupported Provide for promised answer")
	} else {
		return errors.New("unknown message target")
	}

	// TODO: Ask client code if it wants to modify the capability somehow
	// for this accept (e.g. impose limits, change handler, etc)?

	expAc.id = v.cfg.net.recipientIdUniqueKey(prov.recipient)

	rc.answers.set(aid, answer{typ: answerTypeProvide})

	// Prepare vat to receive a new conn from the recipient and for them to
	// send an Accept message with the given recipient id.
	select {
	case v.expAccepts <- expAc:
	case <-ctx.Done():
		return ctx.Err()
	}

	logMsg.Hex("uniqueKey", expAc.id[:]).
		Msg("Asked to provide cap to third party")

	return nil
}

func (v *Vat) processAccept(ctx context.Context, rc *runningConn, ac accept) error {
	aid := AnswerId(ac.qid)
	if rc.answers.has(aid) {
		return fmt.Errorf("remote already asked question %d", aid)
	}

	acResult, err := v.wasExpectingAccept(ctx, ac.provision)
	if err != nil {
		return fmt.Errorf("received unexpected accept: %v", err)
	}

	handler := acResult.handler
	provideAid := acResult.provideAid
	srcConn := acResult.srcConn

	// TODO: support 3PH into promises and not just capabilities.
	if handler == nil {
		return errors.New("3PH only supported for exports")
	}

	// Add the export to the new conn.
	eid, ok := rc.exports.nextID()
	if !ok {
		return errTooManyExports
	}
	rc.exports.set(eid, export{typ: exportTypeLocallyHosted, handler: handler, refCount: 1})
	ans := answer{typ: answerTypeAccept, eid: eid}
	rc.answers.set(aid, ans)

	rc.log.Debug().
		Hex("provisionId", acResult.id[:]).
		Int("eid", int(eid)).
		Int("qid", int(ac.qid)).
		Str("src", srcConn.String()).
		Msg("Providing shared cap to conn")

	// Send the Return that corresponds to the Accept to the newly
	// introduced conn. This is the picked up capability (previously
	// exported in srcConn).
	rpcMsgBuilder, capDesc, err := v.newSingleCapReturn(aid)
	if err != nil {
		return err
	}
	if err := capDesc.SetSenderHosted(eid); err != nil {
		return err
	}
	rpcMsg := v.mp.get()
	rpcMsg.rawSerMb = rpcMsgBuilder.serMb
	if err := rc.queue(ctx, singleMsgBatch(rpcMsg)); err != nil {
		return err
	}

	// Finally, send the Return to srcConn that corresponds to the Provide.
	// This lets the srcConn remote know that the new conn picked up the
	// capability.
	rpcMsgBuilder, payBuilder, err := v.newReturnPayload(provideAid)
	if err != nil {
		return err
	}
	if err := payBuilder.SetContent(capnpser.ZeroStructAsPointerBuilder()); err != nil {
		return err
	}
	rpcMsg = v.mp.get()
	rpcMsg.rawSerMb = rpcMsgBuilder.serMb
	return rc.queue(ctx, singleMsgBatch(rpcMsg))
}

func (v *Vat) processDisembargoAccept(ctx context.Context, rc *runningConn, dis disembargo) error {
	if !dis.target.isImportedCap {
		return fmt.Errorf("only disembargos of exports supported for now") // FIXME
	}

	// rc.log.Warn().Any("xx", fmt.Sprintf("%#v", rc.exports)).Int("bbb", int(dis.target.impcap)).Msgf("XXXXXX %x", dis.target.impcap)
	exp, hasExp := rc.exports.get(ExportId(dis.target.impcap))
	if !hasExp {
		return errDisembargoAcceptUnknownExport(dis.target.impcap)
	}
	if exp.typ != exportTypeThirdPartyExport {
		return fmt.Errorf("received for disembargo on export %d that is not a third party export",
			dis.target.impcap)
	}

	// Forward disembargo to third party.
	disProvide := &message{isDisembargo: true, disembargo: disembargo{ // TODO: fetch from v.mp
		isProvide: true,
		provide:   exp.thirdPartyProvideQid,
	}}
	if err := exp.thirdPartyRC.queue(ctx, singleMsgBatch(disProvide)); err != nil {
		return err
	}

	return nil
}

func (v *Vat) processDisembargoProvide(ctx context.Context, rc *runningConn, dis disembargo) error {
	ans, ok := rc.answers.get(AnswerId(dis.provide))
	if !ok {
		return fmt.Errorf("received disembargo.provide for unknown question id %d", dis.provide)
	}

	if ans.typ != answerTypeProvide {
		return fmt.Errorf("received disembargo.provide for answer %d that is not provide", dis.provide)
	}

	// TODO: answer has to track the resulting conn and export id of the
	// Accepted cap.
	var resConn *runningConn
	var resEid ExportId
	_, _ = resConn, resEid

	// TODO: Return the results of proxied pipelined calls from rc to
	// resConn.

	// TODO: start processing cached (embargoed) pipelined calls that came
	// directly from resConn and return their response.

	return nil
}

func (v *Vat) processDisembargo(ctx context.Context, rc *runningConn, dis disembargo) error {
	if dis.isAccept {
		return v.processDisembargoAccept(ctx, rc, dis)
	} else if dis.isProvide {
		return v.processDisembargoProvide(ctx, rc, dis)
	} else {
		return errors.New("unknown disembargo action")
	}
}

func (v *Vat) processInMessageAlt(ctx context.Context, rc *runningConn, msg types.Message, rawMsg *capnpser.Message) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	which := msg.Which()
	rc.log.Trace().Str("which", which.String()).Msg("Starting to process inbound msg")

	var err error
	switch which {
	case types.Message_Which_Bootstrap:
		var boot types.Bootstrap
		boot, err = msg.AsBootstrap()
		if err == nil {
			err = v.processBootstrap(ctx, rc, boot)
		}

	case types.Message_Which_Call:
		var call types.Call
		call, err = msg.AsCall()
		if err == nil {
			err = v.processCall(ctx, rc, call)
		}

	case types.Message_Which_Finish:
		var fin types.Finish
		fin, err = msg.AsFinish()
		if err == nil {
			err = v.processFinish(ctx, rc, fin)
		}

	case types.Message_Which_Return:
		var ret types.Return
		ret, err = msg.AsReturn()
		if err == nil {
			err = v.processReturn(ctx, rc, ret)
		}

	default:
		err = fmt.Errorf("unknown Message type %d", which)
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		logEvent := rc.log.Err(err).Str("which", which.String())
		if err, ok := err.(extraDataError); ok {
			err.addExtraDataToLog(logEvent)
		}
		if rc.log.GetLevel() < zerolog.InfoLevel {
			msgRawData := rawMsg.Arena().RawDataCopy()
			for i, data := range msgRawData {
				logEvent.Hex(fmt.Sprintf("msg.seg%d", i), data)
			}
		}
		logEvent.Msg("Error while processing inbound message")
	}

	return nil
}

// processInMessage processes an incoming message from a remote Vat.
func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg message) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.log.Trace().Str("which", msg.Which().String()).Msg("Starting to process inbound msg (OLD)")

	var err error
	switch {
	case msg.IsResolve():
		err = v.processResolve(ctx, rc, msg.AsResolve())
	case msg.IsProvide():
		err = v.processProvide(ctx, rc, msg.AsProvide())
	case msg.IsAccept():
		err = v.processAccept(ctx, rc, msg.AsAccept())
	case msg.IsDisembargo():
		err = v.processDisembargo(ctx, rc, msg.AsDisembargo())
	case msg.testEcho != 0:
		echo := v.mp.get()
		*echo = msg
		rc.queue(ctx, singleMsgBatch(echo))
	default:
		err = errors.New("unknown Message type")
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		logEvent := rc.log.Err(err).Str("which", msg.Which().String())
		if err, ok := err.(extraDataError); ok {
			err.addExtraDataToLog(logEvent)
		}
		// if rc.log.GetLevel() < zerolog.InfoLevel {
		//	logEvent.Any("msg", msg)
		//}
		logEvent.Msg("Error while processing inbound message")
	}

	return err
}

var errTooManyOpenQuestions = errors.New("too many open questions")
var errTooManyExports = errors.New("too many exports")

// prepareOutMessageForStep prepares an outgoing Message message that is part of a
// pipeline to be sent to the remote Vat.
func (v *Vat) prepareOutMessageForStep(ctx context.Context, step *pipelineStep,
	conn *runningConn, cmb rpcCallMsgBuilder) (thisQid QuestionId, err error) {

	var ok bool

	if cmb.isBootstrap {
		if thisQid, ok = conn.questions.nextID(); !ok {
			return 0, errTooManyOpenQuestions
		}

		bb := cmb.bootstrapBuilder()
		bb.SetQuestionId(thisQid)
		conn.log.Debug().
			Int("qid", int(thisQid)).
			Msg("Prepared Bootstrap message")
		return thisQid, nil
	}

	if step.parent == nil {
		// The only pipelineStep that doesn't have a parent is the
		// bootstrap message, which was already handled.
		return 0, errors.New("call message without parent pipeline step")
	}

	// Find the target of this call (conn and either parentIid or
	// parentQid).
	parentStep := step.parent
	parentStepState, parentStepValue := parentStep.value.Get()
	if parentStepState == pipeStepFailed {
		return 0, parentStepValue.err
	}
	parentIid := parentStepValue.iid
	parentQid := parentStepValue.qid

	// Find next available question id.
	if thisQid, ok = conn.questions.nextID(); !ok {
		return 0, errTooManyOpenQuestions
	}
	cb := cmb.callBuilder()
	cb.SetQuestionId(thisQid)
	mtb, err := cb.NewTarget()
	if err != nil {
		return 0, err
	}

	if parentIid > 0 {
		mtb.SetImportedCap(parentIid)

		conn.log.Debug().
			Int("qid", int(thisQid)).
			Int("iid", int(parentIid)).
			Msg("Prepared call for exported cap")
	} else if parentQid > 0 {
		pab, err := mtb.NewPromisedAnswer()
		if err != nil {
			return 0, err
		}
		pab.SetQuestionId(parentQid)

		conn.log.Debug().
			Int("qid", int(thisQid)).
			Int("pans", int(parentQid)).
			Msg("Prepared call for in-pipeline promised answer")
	} else {
		return 0, errors.New("unimplemented case")
	}

	return thisQid, nil
}

// commitOutMessageForStep commits the changes of the pipeline step to the local Vat's
// state, under the assumption that the given pipeline step was successfully
// sent to the remote Vat.
func (v *Vat) commitOutMessageForStep(_ context.Context, step *pipelineStep,
	conn *runningConn, qid QuestionId) error {

	if qid == 0 {
		// Guard against errors while developing.
		return errors.New("unimplemented commitment of message")
	}
	conn.log.Debug().Int("qid", int(qid)).Msg("Comitted outgoing message")

	// TODO: Do not add finalizer if return is known to have no caps?
	// TODO: Maybe add to question directly (instead of weak.Pointer) to
	// keep it around until we get a Return, then if needed (i.e. there are
	// pipelines or finish is needed) add the weak ref.
	// runtime.AddCleanup(step, conn.cleanupQuestionIdDueToUnref, qid) // TODO: Save cleanup in question in case of early finish?
	// runtime.SetFinalizer(step, finalizePipelineStep)
	q := question{strongStep: step}
	conn.questions.set(qid, q)

	// This step is now in flight. Allow forks from it to start.
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateSending {
			return os, ov, fmt.Errorf("invalid precondition state: %v", os)
		}
		ov.qid = qid
		ov.conn = conn
		return pipeStepStateRunning, ov, nil
	})
}

func (v *Vat) queueFinish(ctx context.Context, rc *runningConn, qid QuestionId) error {
	rc.mu.Lock()
	rc.questions.del(qid)
	rc.mu.Unlock()

	rpcMsgBuilder, _, err := v.newFinish(qid)
	if err != nil {
		return err
	}

	rpcMsg := v.mp.get()
	rpcMsg.rawSerMb = rpcMsgBuilder.serMb
	return rc.queue(ctx, singleMsgBatch(rpcMsg))
}

func (v *Vat) queueResolve(ctx context.Context, rc *runningConn, eid ExportId, exp export, resolution export) error {
	msg := &message{ // TODO: fetch from v.mp
		isResolve: true,
		resolve:   resolve{pid: eid},
	}

	switch resolution.typ {
	case exportTypePromise:
		msg.resolve.cap.senderPromise = exp.resolvedToExport
	case exportTypeLocallyHosted:
		msg.resolve.cap.senderHosted = exp.resolvedToExport
	case exportTypeThirdPartyExport:
		msg.resolve.cap.thirdPartyHosted = thirdPartyCapDescriptor{
			id:     resolution.thirdPartyCapDescId,
			vineId: resolution.thirdPartyVineId,
		}
	default:
		return fmt.Errorf("unknown resolution type %s", resolution.typ)
	}

	return rc.queue(ctx, singleMsgBatch(msg))
}

func (v *Vat) sendProvide(ctx context.Context, rc *runningConn, p provide) error {
	msg := &message{ // TODO: fetch from v.mp
		isProvide: true,
		provide:   p,
	}
	outMsg := singleMsgBatch(msg)
	outMsg.wantSentAck()

	if err := rc.queue(ctx, outMsg); err != nil {
		return err
	}

	// Wait until the provide is actually on the target remote before
	// returning (to send the Resolve/Return to the caller). This is
	// necessary to ensure that the Accept the source will send to target
	// will reach target AFTER the Provide.
	if err := outMsg.waitSentAck(ctx); err != nil {
		return err
	}

	if p.target.isImportedCap {
		rc.log.Debug().
			Int("qid", int(p.qid)).
			Int("iid", int(p.target.impcap)).
			Msg("Sent provide for imported cap")
	} else {
		rc.log.Debug().
			Int("qid", int(p.qid)).
			Int("pans", int(p.target.pans.qid)).
			Msg("Sent provide for promised answer")
	}

	return nil
}

func (v *Vat) sendAccept(ctx context.Context, rc *runningConn, accept accept) error {
	msg := &message{isAccept: true, accept: accept} // TODO: fetch from v.mp
	outMsg := singleMsgBatch(msg)
	outMsg.wantSentAck()

	if err := rc.queue(ctx, outMsg); err != nil {
		return err
	}

	// Wait until Accept is actually outbound (as opposed to simply queued).
	// This is necessary to ensure correctness of operations. Accept MUST
	// reach the remote end of the new conn BEFORE a Disembargo is proxied
	// through srcConn, otherwise ordering is not guaranteed. In the mean
	// time, any pipelined calls continue to be proxied through srcConn.
	return outMsg.waitSentAck(ctx)
}
