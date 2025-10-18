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

	outMsg, capDesc, err := v.newSingleCapReturn(bootQid)
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

	return rc.queue(ctx, outMsg)
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

		/*
			rc.log.Debug().
				Int("qid", int(qid)).
				Str("rtyp", stepResultType).
				Bool("noFinishNeeded", noFinishNeeded).
				Msg("Processed Return message")
		*/
		_ = stepResultType

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

	// Start preparing reply.
	outMsg, err := v.mbp.getForPayloadSize(0) // TODO: size hint?
	if err != nil {
		return err
	}
	reply, err := outMsg.mb.NewReturn()
	if err != nil {
		return err
	}
	reply.SetAnswerId(AnswerId(qid))

	crb := &rc.crb // Ok to reuse (rc is locked).
	crb.pb, err = reply.NewResults()
	crb.serMb = outMsg.serMsg
	crb.iid = InterfaceId(iid)
	crb.mid = MethodId(mid)
	if err != nil {
		return err
	}

	// Make the call!
	logEvent.Msg("Locally handling call")
	err = exp.handler.Call(rc.ctx, crb)
	if ex, ok := err.(callExceptionError); ok {
		// When an exception that will be sent remotely is detected,
		// re-create the reply. This ensures anything written to the
		// payload inside the handler's Call() method will *NOT* be
		// sent as an orphan object inside the reply.
		if err := outMsg.serMsg.Reset(); err != nil {
			return err
		}

		// Turn the error into a returned exception.
		reply, err = outMsg.mb.NewReturn()
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
			pbReader := crb.pb.AsReader()
			content, err := pbReader.Content()
			if err != nil {
				return err
			}

			if content.IsCapPointer() {
				cp := content.AsCapPointer()
				rootCapIndex = int(cp.Index())
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

	return rc.send(ctx, outMsg)
	// return rc.queue(ctx, outMsg)
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

// recvdThirdPartyCapDesc is a copy of a received thirdPartyCapDescriptor.
type recvdThirdPartyCapDesc struct {
	mb          *capnpser.MessageBuilder
	tpToContact capnpser.AnyPointer
}

// copyThirdPartyCapDesc copies a thirdPartyCapDescriptor received in a message
// into a temp buffer.
func (v *Vat) copyThirdPartyCapDesc(tpcd types.ThirdPartyCapDescriptor) (recvdThirdPartyCapDesc, error) {
	tpcdId, err := tpcd.Id()
	if err != nil {
		return recvdThirdPartyCapDesc{}, err
	}
	mb := v.mbp.getRawMessageBuilder()
	tpcdCopy, err := capnpser.DeepCopy(tpcdId, mb)
	if err != nil {
		v.mbp.put(mb)
		return recvdThirdPartyCapDesc{}, err
	}

	return recvdThirdPartyCapDesc{mb: mb, tpToContact: tpcdCopy.Reader()}, nil
}

func (v *Vat) resolveThirdPartyCapForStep(ctx context.Context, step *pipelineStep,
	srcConn *runningConn, recvTphCapDesc recvdThirdPartyCapDesc) error {

	rc, provId, err := v.connectToIntroduced3rdParty(ctx, srcConn, recvTphCapDesc.tpToContact)
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

	// In the future, this could be dynamically determined, because
	// the local vat can know whether there are pending pipelined
	// calls or not.
	embargo := true

	// TODO: finalizer???
	q := question{weakStep: weak.Make(step)}
	rc.questions.set(acceptQid, q)
	rpcAccept, _, err := v.newAccept(acceptQid, provId, embargo)
	/*
		accept := accept{
			qid:       acceptQid,
			provision: provId,
		}
	*/

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
	if err := v.sendAccept(ctx, rc, rpcAccept); err != nil {
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
			disembargoTarget.impCap = ov.iid
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
			disembargoTarget.pansQid = ov.qid
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
	if !embargo {
		return nil
	}

	// Queue the Disembargo in the old conn, which will travese the network
	// all the way to the third party. This will make the third party
	// process the calls sent directly from the local vat (as opposed to
	// those proxied by the original conn).
	disRpc, dis, err := v.newDisembargo(disembargoTarget)
	if err != nil {
		return err
	}
	if err := dis.SetAccept(); err != nil {
		return err
	}
	/*
		dis := &message{isDisembargo: true, disembargo: disembargo{ // TODO: fetch from v.mp
			isAccept: true,
			target:   disembargoTarget,
		}}
	*/
	if err := srcConn.queue(ctx, disRpc); err != nil {
		return err
	}

	// Only thing left to complete 3PH now is wait for the Return message
	// that completes the Accept call above.

	return nil
}

func (v *Vat) processResolve(ctx context.Context, rc *runningConn, res types.Resolve) error {
	iid := ImportId(res.PromiseId())
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

	if res.Which() != types.Resolve_Which_Cap {
		// TODO: handle exception.
		return errors.New("unhandled exception case for resolve")
	}

	// Similar to code in processReturn. Unify?
	capEntry, err := res.AsCap()
	if err != nil {
		return err
	}
	var resolveCap capability
	var resolvePromise ExportId
	var resolveId ImportId
	var resImport imprt
	switch capEntry.Which() {
	case types.CapDescriptor_Which_SenderHosted:
		// Resolved into a remote capability.
		resolveCap = capability{eid: ExportId(capEntry.AsSenderHosted())}
		resolveId = ImportId(capEntry.AsSenderHosted())
		resImport = imprt{typ: importTypeSenderHosted}
	case types.CapDescriptor_Which_SenderPromise:
		// Resolved into another remote promise.
		resolvePromise = capEntry.AsSenderPromise()
		resolveId = ImportId(resolvePromise)
		resImport = imprt{typ: importTypeRemotePromise, step: imp.step}
	case types.CapDescriptor_Which_ThirdPartyHosted:
		// Make a local copy of the third party cap desc received,
		// because we'll launch a goroutine to handle the 3PH
		// resolution.
		recvTphCapDesc, err := capEntry.AsThirdPartyHosted()
		if err != nil {
			return err
		}
		tphCapDescCopy, err := v.copyThirdPartyCapDesc(recvTphCapDesc)
		if err != nil {
			return fmt.Errorf("unable to copy received ThirdPartyHosted cap desc: %v", err)
		}

		// Start the 3PH resolution process. This will involve
		// connecting to the third party (if we haven't already) and
		// then sending the Accept with the provision id.
		//
		// TODO: If this is only level 1, use vineId instead and proxy
		// requests.
		go func() {
			err := v.resolveThirdPartyCapForStep(ctx, step, rc, tphCapDescCopy)
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
	default:
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

func (v *Vat) processProvide(ctx context.Context, rc *runningConn, prov types.Provide) error {
	if v.cfg.net == nil {
		return err3PHWithoutVatNetwork
	}

	provQid := prov.QuestionId()
	aid := AnswerId(provQid)
	if rc.answers.has(aid) {
		return fmt.Errorf("remote already asked question %d", aid)
	}

	// Build value that will hold the shared cap data.
	expAc := expectedAccept{
		srcConn:    rc,
		provideAid: aid,
	}

	logMsg := rc.log.Debug()

	target, err := prov.Target()
	if err != nil {
		return err
	}

	// Check if the target exists exported to the caller.
	switch target.Which() {
	case types.MessageTarget_Which_ImportedCap:
		impCap := target.AsImportedCap()
		capExp, ok := rc.exports.get(ExportId(impCap))
		if !ok {
			return fmt.Errorf("export not found %d", impCap)
		}

		expAc.handler = capExp.handler
		logMsg = logMsg.Str("typ", "importedCap").Int("iid", int(impCap))

	case types.MessageTarget_Which_PromisedAnswer:
		pans, err := target.AsPromisedAnswer()
		if err != nil {
			return err
		}
		pansQid := pans.QuestionId()

		if !rc.answers.has(AnswerId(pansQid)) {
			return fmt.Errorf("answer not found %d", pansQid)
		}

		logMsg = logMsg.Str("typ", "promisedAnswer").Int("iid", int(pansQid))

		// TODO: support
		return fmt.Errorf("unsupported Provide for promised answer")

	default:
		return errors.New("unknown message target")
	}

	// TODO: Ask client code if it wants to modify the capability somehow
	// for this accept (e.g. impose limits, change handler, etc)?

	// Get a unique id for this recipient (will be used to match the
	// Accept).
	provRec, err := prov.Recipient()
	if err != nil {
		return err
	}
	expAc.id = v.cfg.net.recipientIdUniqueKey(provRec)

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

func (v *Vat) processAccept(ctx context.Context, rc *runningConn, ac types.Accept) error {
	aid := AnswerId(ac.QuestionId())
	if rc.answers.has(aid) {
		return fmt.Errorf("remote already asked question %d", aid)
	}

	provId, err := ac.Provision()
	if err != nil {
		return fmt.Errorf("error reading Accept.Provision: %v", err)
	}
	acResult, err := v.wasExpectingAccept(ctx, provId)
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
		Int("qid", int(aid)).
		Str("src", srcConn.String()).
		Msg("Providing shared cap to conn")

	// Send the Return that corresponds to the Accept to the newly
	// introduced conn. This is the picked up capability (previously
	// exported in srcConn).
	outMsg, capDesc, err := v.newSingleCapReturn(aid)
	if err != nil {
		return err
	}
	if err := capDesc.SetSenderHosted(eid); err != nil {
		return err
	}
	if err := rc.queue(ctx, outMsg); err != nil {
		return err
	}

	// Finally, send the Return to srcConn that corresponds to the Provide.
	// This lets the srcConn remote know that the new conn picked up the
	// capability.
	outMsg, payBuilder, err := v.newReturnPayload(provideAid)
	if err != nil {
		return err
	}
	if err := payBuilder.SetContent(capnpser.ZeroStructAsPointerBuilder()); err != nil {
		return err
	}
	return rc.queue(ctx, outMsg)
}

func (v *Vat) processDisembargoAccept(ctx context.Context, rc *runningConn, dis types.Disembargo) error {
	target, err := dis.Target()
	if err != nil {
		return err
	}

	if target.Which() != types.MessageTarget_Which_ImportedCap {
		return fmt.Errorf("only disembargos of exports supported for now") // FIXME
	}

	impCap := target.AsImportedCap()

	exp, hasExp := rc.exports.get(ExportId(impCap))
	if !hasExp {
		return errDisembargoAcceptUnknownExport(impCap)
	}
	if exp.typ != exportTypeThirdPartyExport {
		return fmt.Errorf("received for disembargo on export %d that is not a third party export",
			impCap)
	}

	// Forward disembargo to third party.
	outMsg, err := v.mbp.get()
	if err != nil {
		return err
	}
	fwdDis, err := outMsg.mb.NewDisembargo()
	if err != nil {
		return err
	}
	fwdDis.SetProvide(types.Disembargo_EmbargoId(exp.thirdPartyProvideQid))

	if err := exp.thirdPartyRC.queue(ctx, outMsg); err != nil {
		return err
	}

	return nil
}

func (v *Vat) processDisembargoProvide(ctx context.Context, rc *runningConn, dis types.Disembargo) error {
	ansId := AnswerId(dis.AsProvide())
	ans, ok := rc.answers.get(ansId)
	if !ok {
		return fmt.Errorf("received disembargo.provide for unknown question id %d", ansId)
	}

	if ans.typ != answerTypeProvide {
		return fmt.Errorf("received disembargo.provide for answer %d that is not provide", ansId)
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

func (v *Vat) processDisembargo(ctx context.Context, rc *runningConn, dis types.Disembargo) error {
	switch dis.Which() {
	case types.Disembargo_Which_Accept:
		return v.processDisembargoAccept(ctx, rc, dis)
	case types.Disembargo_Which_Provide:
		return v.processDisembargoProvide(ctx, rc, dis)
	default:
		return errors.New("unknown disembargo action")
	}
}

func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg types.Message, rawMsg capnpser.Message) error {

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

	case types.Message_Which_Resolve:
		var res types.Resolve
		res, err = msg.AsResolve()
		if err == nil {
			err = v.processResolve(ctx, rc, res)
		}

	case types.Message_Which_Provide:
		var res types.Provide
		res, err = msg.AsProvide()
		if err == nil {
			err = v.processProvide(ctx, rc, res)
		}

	case types.Message_Which_Accept:
		var res types.Accept
		res, err = msg.AsAccept()
		if err == nil {
			err = v.processAccept(ctx, rc, res)
		}

	case types.Message_Which_Disembargo:
		var res types.Disembargo
		res, err = msg.AsDisembargo()
		if err == nil {
			err = v.processDisembargo(ctx, rc, res)
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

	outMsg, _, err := v.newFinish(qid)
	if err != nil {
		return err
	}

	return rc.queue(ctx, outMsg)
}

func (v *Vat) queueResolve(ctx context.Context, rc *runningConn, eid ExportId, exp export, resolution export) error {
	outMsg, res, err := v.newResolve(eid)
	if err != nil {
		return err
	}

	capDesc, err := res.NewCap()
	if err != nil {
		return err
	}

	switch resolution.typ {
	case exportTypePromise:
		if err := capDesc.SetSenderPromise(exp.resolvedToExport); err != nil {
			return err
		}
	case exportTypeLocallyHosted:
		if err := capDesc.SetSenderHosted(exp.resolvedToExport); err != nil {
			return err
		}
	case exportTypeThirdPartyExport:
		var tpcd types.ThirdPartyCapDescriptorBuilder
		if tpcd, err = capDesc.NewThirdPartyHosted(); err != nil {
			return err
		}
		tpcd.SetVineId(resolution.thirdPartyVineId)

		// Copy the thirdPartyDesc to the outbound resolve
		// message.
		thirdPartyCapObj, err := capnpser.DeepCopy(resolution.thirdPartyCapDescIdAlt, outMsg.serMsg)
		if err != nil {
			return err
		}
		if err := tpcd.SetId(thirdPartyCapObj); err != nil {
			return err
		}
		// err = tpcd.SetId(resolution.thirdPartyCapDescId,)
	default:
		return fmt.Errorf("unknown resolution type %s", resolution.typ)
	}

	return rc.queue(ctx, outMsg)
}

func (v *Vat) sendProvide(ctx context.Context, rc *runningConn, outMsg outMsg, prov types.ProvideBuilder) error {
	// Track the provider qid and target.
	var impCap types.ImportId
	var pansQid types.QuestionId
	var isImpCap bool

	pr := prov.AsReader()
	target, err := pr.Target()
	if err != nil {
		return err
	}
	qid := pr.QuestionId()
	switch target.Which() {
	case types.MessageTarget_Which_ImportedCap:
		impCap = target.AsImportedCap()
		isImpCap = true
	case types.MessageTarget_Which_PromisedAnswer:
		pans, err := target.AsPromisedAnswer()
		if err != nil {
			return err
		}
		pansQid = pans.QuestionId()
	default:
		return errors.New("unhandled case in sendProvide")
	}

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

	if isImpCap {
		rc.log.Debug().
			Int("qid", int(qid)).
			Int("iid", int(impCap)).
			Msg("Sent provide for imported cap")
	} else {
		rc.log.Debug().
			Int("qid", int(qid)).
			Int("pans", int(pansQid)).
			Msg("Sent provide for promised answer")
	}

	return nil
}

func (v *Vat) sendAccept(ctx context.Context, rc *runningConn, accept outMsg) error {
	accept.wantSentAck()

	if err := rc.queue(ctx, accept); err != nil {
		return err
	}

	// Wait until Accept is actually outbound (as opposed to simply queued).
	// This is necessary to ensure correctness of operations. Accept MUST
	// reach the remote end of the new conn BEFORE a Disembargo is proxied
	// through srcConn, otherwise ordering is not guaranteed. In the mean
	// time, any pipelined calls continue to be proxied through srcConn.
	return accept.waitSentAck(ctx)
}
