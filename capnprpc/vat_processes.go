// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"weak"

	"github.com/rs/zerolog"
)

func (v *Vat) processBootstrap(ctx context.Context, rc *runningConn, msg message) error {
	bootMsg := msg.AsBootstrap()
	reply := message{
		isReturn: true,
		ret: rpcReturn{
			aid:       AnswerId(bootMsg.qid),
			isResults: true,
			pay: payload{
				content: anyPointer{
					isCapPointer: true,
					cp:           capPointer{index: 0},
				},
				capTable: []capDescriptor{
					{senderHosted: rc.bootExportId},
				},
			},
		},
	}

	// Modify answer table to track the bootrap export. bootExportId is set
	// during conn setup automatically.
	rc.answers.set(AnswerId(bootMsg.qid), answer{typ: answerTypeBootstrap, eid: rc.bootExportId})

	rc.log.Debug().
		Int("qid", int(bootMsg.qid)).
		Int("eid", int(rc.bootExportId)).
		Msg("Exported Bootstrap")

	return rc.queue(ctx, singleMsgBatch(reply))
}

func (v *Vat) processReturn(ctx context.Context, rc *runningConn, ret rpcReturn) error {
	qid := QuestionId(ret.AnswerId())
	q, ok := rc.questions.get(qid)
	if !ok {
		return fmt.Errorf("question %d not found", qid)
	}

	// TODO: support exception, cancel, etc
	if !ret.IsResults() {
		return fmt.Errorf("only results supported")
	}

	pipe := q.pipe.Value()
	if pipe == nil {
		// This pipeline isn't used anymore (was released and a Finish
		// should've been sent, or will be shortly), so nothing to do.
		rc.log.Debug().Int("qid", int(qid)).Msg("Received Return message for released pipeline")
		return nil
	}

	// Go through cap table, modify imports table based on what was
	// exported by this call.
	//
	// TODO: only do this if the cap is referenced in the content?
	payload := ret.AsResults()
	capTable := payload.CapTable()
	for _, entry := range capTable {
		var iid ImportId
		var imp imprt
		if entry.IsSenderHosted() {
			iid = ImportId(entry.AsSenderHosted())
			imp = imprt{typ: importTypeSenderHosted}
		} else if entry.IsSenderPromise() {
			iid = ImportId(entry.AsSenderPromise())
			imp = imprt{typ: importTypeRemotePromise, pipe: q.pipe, stepIdx: q.stepIdx}
		} else {
			return fmt.Errorf("unsupported capability type")
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
	content := payload.Content()
	if content.IsCapPointer() {
		// NOT GOOD. Must have a new type to pass along instead of
		// parsing like this (maybe). Think about embedded caps.
		cp := content.AsCapPointer()
		capIndex := cp.Index()
		if int(capIndex) >= len(capTable) {
			return fmt.Errorf("capability referenced index outside cap table")
		}
		capEntry := capTable[capIndex]
		if capEntry.IsSenderHosted() {
			stepImportId = ImportId(capEntry.AsSenderHosted())
			stepResult = capability{eid: ExportId(stepImportId)}
			stepResultType = "senderHostedCap"
		} else if capEntry.IsSenderPromise() {
			stepImportId = ImportId(capEntry.AsSenderPromise())
			stepResultPromise = ExportId(stepImportId)
			stepResultType = "senderPromise"
		} else {
			return errors.New("unknown cap entry type")
		}

	} else if content.IsStruct() {
		// TODO: copy if its a struct? Or release serialized message if
		// content is just a cap (because it's not needed anymore)?
		stepResult = content.AsStruct()
		stepResultType = "struct"
	} else if content.IsVoid() {
		stepResult = struct{}{}
		stepResultType = "void"
	} else {
		return errors.New("unknown/unimplemented content type")
	}

	// Fulfill pieline waiting for this result.
	step := pipe.Step(q.stepIdx)
	return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepStateRunning {
			return os, ov, fmt.Errorf("pipeline step not running: %v", os)
		}

		rc.log.Debug().
			Int("qid", int(qid)).
			Str("rtyp", stepResultType).
			Msg("Processed Return message")

		ov.iid = stepImportId
		if stepResultPromise > 0 {
			// This step isn't resolved into a concrete exported
			// cap or struct yet. Keep it running.
			return os, ov, nil
		}

		ov.value = stepResult
		return pipeStepStateDone, ov, nil
	})
}

var errCallWithoutId = errors.New("call with zero interfaceId and methodId")

func (v *Vat) processCall(ctx context.Context, rc *runningConn, c call) error {
	if c.iid == 0 && c.mid == 0 {
		// Only bootstrap is allowed to have iid+mid == 0.
		return errCallWithoutId
	}

	if rc.answers.has(AnswerId(c.qid)) {
		return fmt.Errorf("remote already asked question %d", c.qid)
	}

	// Determine the target of this call (either an exported cap or a
	// promised answer).
	var eid ExportId
	if c.target.isPromisedAnswer {
		// Promised answers are in the answer table.
		//
		// TODO: Recursively track it down if the answer is another
		// promise.
		q, ok := rc.answers.get(AnswerId(c.target.pans.qid))
		if !ok {
			return fmt.Errorf("call referenced unknown promised answer %d", c.target.pans.qid)
		}

		eid = q.eid // What about promises?
	} else if c.target.isImportedCap {
		eid = ExportId(c.target.impcap)
	} else {
		return errors.New("unsupported call target")
	}

	exp, ok := rc.exports.get(eid)
	if !ok {
		return fmt.Errorf("call message target determined to be unset export %d", eid)
	}

	if exp.typ != exportTypeLocallyHosted {
		return fmt.Errorf("unsupported export type %d", exp.typ)
	}

	// TODO: proxy calls when exp.typ == exportTypeThirdPartyExport.

	callArgs := callHandlerArgs{
		iid:    interfaceId(c.iid),
		mid:    methodId(c.mid),
		params: c.params,
		rc:     rc,
	}

	// Start preparing reply.
	reply := message{
		isReturn: true,
		ret:      rpcReturn{aid: AnswerId(c.qid)},
	}
	crb := &rc.crb // Ok to reuse (rc is locked).
	crb.payload = payload{content: anyPointer{
		isVoid: true, // Void result by default on non-error.
	}}

	rc.log.Trace().
		Int("qid", int(c.qid)).
		Int("eid", int(eid)).
		Msg("Locally handling call")

	// Make the call!
	err := exp.handler.Call(rc.ctx, callArgs, crb)
	if ex, ok := err.(callExceptionError); ok {
		// Turn the error into a returned exception.
		reply.ret.isException = true
		reply.ret.exc = ex.ToException()

		rc.log.Debug().
			Int("qid", int(c.qid)).
			Int("eid", int(eid)).
			Dict("ex", zerolog.Dict().
				Int("type", reply.ret.exc.typ).
				Str("reason", reply.ret.exc.reason)).
			Msg("Processed call into exception")

	} else if err != nil {
		// Fatal connection error.
		return err
	} else {
		reply.ret.isResults = true
		reply.ret.pay = crb.payload

		// If the result is a capability, the answer is pipelinable, so
		// track where the corresponding export will be.
		var rootCapIndex int = -1
		if crb.payload.content.IsCapPointer() {
			rootCapIndex = int(crb.payload.content.AsCapPointer().index)
		}
		ans := answer{typ: answerTypeCall}

		// Track all exported caps.
		for i, cp := range crb.payload.capTable {
			if !cp.hasExportId() {
				continue
			}
			capEid := cp.exportId()
			if capExp, ok := rc.exports.get(capEid); ok {
				// TODO: take a pointer instead?
				capExp.refCount++
				rc.exports.set(capEid, capExp)
			} else if cp.IsSenderPromise() { // here sender == local vat.
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
				ans.eid = capEid
			}
		}

		// Save the answer in the answers table.
		rc.answers.set(AnswerId(c.qid), ans)

		rc.log.Debug().
			Int("qid", int(c.qid)).
			Int("eid", int(eid)).
			Msg("Processed call into payload result")
	}

	return rc.queue(ctx, singleMsgBatch(reply))
}

func (v *Vat) processFinish(ctx context.Context, rc *runningConn, fin finish) error {
	var err error
	aid := AnswerId(fin.qid)

	if !rc.answers.has(aid) {
		err = fmt.Errorf("answer %d not in answers table", aid)
	} else {
		rc.answers.del(AnswerId(fin.qid))
	}

	// TODO: release exported caps?

	if err == nil {
		rc.log.Debug().
			Int("aid", int(aid)).
			Msg("Removed answer due to Finish message")
	}

	return err
}

func (v *Vat) resolveThirdPartyCapForPipeStep(ctx context.Context, pipe *pipeline, stepIdx int,
	srcConn *runningConn, tphCapDesc thirdPartyCapDescriptor) error {

	rc, provId, err := v.connectToIntroduced3rdParty(ctx, srcConn, tphCapDesc)
	if err != nil {
		return err
	}

	v.log.Trace().
		Str("src", srcConn.String()).
		Str("dst", rc.String()).
		Msg("Shortening path after 3PH introduction")

	// Send Accept() with embargo set to the new remote.
	rc.mu.Lock()
	acceptQid, ok := rc.questions.nextID()
	if !ok {
		rc.mu.Unlock()
		return errTooManyOpenQuestions
	}
	// TODO: finalizer???
	q := question{pipe: weak.Make(pipe), stepIdx: stepIdx}
	rc.questions.set(acceptQid, q)
	rc.mu.Unlock()
	accept := accept{
		qid:       acceptQid,
		provision: provId,

		// In the future, this could be dynamically determined, because
		// the local vat can know whether there are pending pipelined
		// calls or not.
		embargo: true,
	}
	if err := v.sendAccept(ctx, rc, accept); err != nil {
		return err
	}

	// Lock srcConn before modifying the pipeline's conn to avoid a possible
	// race condition with startPipeline(). See the comment tagged as
	// 3PHCONNISSUE.
	srcConn.mu.Lock()

	// From now on, any calls pipelined on this step will go directly to the
	// third party (path has shortened!).
	step := pipe.Step(stepIdx)
	var disembargoTarget messageTarget
	err = step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os == pipelineStepFailed {
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
	srcConn.mu.Unlock()
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
	dis := message{isDisembargo: true, disembargo: disembargo{
		isAccept: true,
		target:   disembargoTarget,
	}}
	if err := srcConn.queue(ctx, singleMsgBatch(dis)); err != nil {
		srcConn.mu.Unlock()
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

	pipe := imp.pipe.Value()
	if pipe == nil {
		// Pipeline already canceled?
		return fmt.Errorf("pipeline of import %d already released", iid)
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
		resImport = imprt{typ: importTypeRemotePromise, pipe: imp.pipe, stepIdx: imp.stepIdx}
	} else if capEntry.IsThirdPartyHosted() {
		// Start the 3PH resolution process. This will involve
		// connecting to the third party (if we haven't already) and
		// then sending the Accept with the provision id.
		//
		// TODO: If this is only level 1, use vineId instead and proxy
		// requests.
		go func() {
			err := v.resolveThirdPartyCapForPipeStep(ctx, pipe, imp.stepIdx, rc, capEntry.AsThirdPartyHosted())
			if err != nil {
				rc.log.Err(err).Msg("resolveThirdPartyCapForPipeStep failed")
				// TODO: fail the pipeline step.
			} else {
				rc.log.Trace().Msg("resolveThirdPartyCapForPipeStep completed")
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

	step := pipe.Step(imp.stepIdx)
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
	rc.answers.set(aid, answer{typ: answerTypeAccept, eid: eid})

	rc.log.Debug().
		Hex("provisionId", acResult.id[:]).
		Int("eid", int(eid)).
		Int("qid", int(ac.qid)).
		Str("src", srcConn.String()).
		Msg("Providing shared cap to conn")

	// Send the Return that corresponds to the Accept to the newly
	// introduced conn. This is the picked up capability (previously
	// exported in srcConn).
	retAccept := message{isReturn: true, ret: rpcReturn{
		aid:       aid,
		isResults: true,
		pay: payload{
			content:  anyPointer{isCapPointer: true, cp: capPointer{index: 0}},
			capTable: []capDescriptor{{senderHosted: eid}},
		},
	}}
	if err := rc.queue(ctx, singleMsgBatch(retAccept)); err != nil {
		return err
	}

	// Finally, send the Return to srcConn that corresponds to the Provide.
	// This lets the srcConn remote know that the new conn picked up the
	// capability.
	retProvide := message{isReturn: true, ret: rpcReturn{
		aid:       provideAid,
		isResults: true,
		pay: payload{
			content: anyPointer{isVoid: true},
		},
	}}
	return srcConn.queue(ctx, singleMsgBatch(retProvide))
}

func (v *Vat) processDisembargoAccept(ctx context.Context, rc *runningConn, dis disembargo) error {
	if !dis.target.isImportedCap {
		return fmt.Errorf("only disembargos of exports supported for now") // FIXME
	}

	exp, hasExp := rc.exports.get(ExportId(dis.target.impcap))
	if !hasExp {
		return errors.New("received disembargo for unknwon export")
	}
	if exp.typ != exportTypeThirdPartyExport {
		return errors.New("received for disembargo on export that is not a third party export")
	}

	// Forward disembargo to third party.
	disProvide := message{isDisembargo: true, disembargo: disembargo{
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

// processInMessage processes an incoming message from a remote Vat.
func (v *Vat) processInMessage(ctx context.Context, rc *runningConn, msg message) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	var err error
	switch {
	case msg.IsBootstrap():
		err = v.processBootstrap(ctx, rc, msg)
	case msg.IsReturn():
		err = v.processReturn(ctx, rc, msg.AsReturn())
	case msg.IsCall():
		err = v.processCall(ctx, rc, msg.AsCall())
	case msg.IsFinish():
		err = v.processFinish(ctx, rc, msg.AsFinish())
	case msg.IsResolve():
		err = v.processResolve(ctx, rc, msg.AsResolve())
	case msg.IsProvide():
		err = v.processProvide(ctx, rc, msg.AsProvide())
	case msg.IsAccept():
		err = v.processAccept(ctx, rc, msg.AsAccept())
	case msg.IsDisembargo():
		err = v.processDisembargo(ctx, rc, msg.AsDisembargo())
	case msg.testEcho != 0:
		rc.queue(ctx, singleMsgBatch(msg))
	default:
		err = errors.New("unknown Message type")
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		logEvent := rc.log.Err(err).Str("which", msg.Which().String())
		if rc.log.GetLevel() < zerolog.InfoLevel {
			logEvent.Any("msg", msg)
		}
		logEvent.Msg("Error while processing inbound message")
	}

	return err
}

var errTooManyOpenQuestions = errors.New("too many open questions")
var errTooManyExports = errors.New("too many exports")

// prepareOutMessage prepares an outgoing Message message that is part of a
// pipeline to be sent to the remote Vat.
//
// Note: this does _not_ commit the changes to the conn's tables yet.
func (v *Vat) prepareOutMessage(ctx context.Context, pipe *pipeline,
	stepIdx int, parentQid QuestionId, conn *runningConn) (thisQid QuestionId, err error) {

	var ok bool

	step := pipe.step(stepIdx)
	if step.rpcMsg.IsBootstrap() {
		if thisQid, ok = conn.questions.nextID(); !ok {
			return 0, errTooManyOpenQuestions
		}

		step.rpcMsg.boot.qid = thisQid
		conn.log.Debug().
			Int("qid", int(thisQid)).
			Msg("Prepared Bootstrap message")
		return thisQid, nil
	}

	if !step.rpcMsg.IsCall() {
		// Only happens during development.
		return 0, errors.New("unimplemented message type in prepareOutMessage")
	}

	// Find the target of this call (conn and either parentIid or
	// parentQid).
	var parentIid ImportId
	if stepIdx == 0 && pipe.parent == nil {
		// Should never happen, but avoid a panic below.
		return 0, errors.New("non-bootstrap call without a parent pipeline")
	} else if stepIdx == 0 {
		// Fork from a parent pipeline. Determine if the parent step has
		// been resolved into an imported cap or we're still waiting for
		// a remote Return.
		parentStep := pipe.parent.Step(pipe.parentStepIdx)
		parentStepStepState, parentStepValue := parentStep.value.Get()
		if parentStepStepState == pipelineStepFailed {
			return 0, parentStepValue.err
		}
		parentIid = parentStepValue.iid
		parentQid = parentStepValue.qid
		// step.conn = parentStepValue.conn
	} else {
		// Still same pipeline, use the same conn as parent (parentQid
		// already refers to the parent's question id).
	}

	// Can now determine question id.
	if thisQid, ok = conn.questions.nextID(); !ok {
		return 0, errTooManyOpenQuestions
	}
	step.rpcMsg.call.qid = thisQid

	if parentIid > 0 {
		// parentIid > 0 means this is already resolved into a returned
		// import (either a remote promise or a concrete remote cap).
		step.rpcMsg.call.target = messageTarget{
			isImportedCap: true,
			impcap:        parentIid,
		}

		conn.log.Debug().
			Int("qid", int(thisQid)).
			Int("iid", int(parentIid)).
			Msg("Prepared call for exported cap")
	} else if parentQid > 0 {
		// parentQid > 0 means this is a pielined call to a promised
		// answer.
		step.rpcMsg.call.target = messageTarget{
			isPromisedAnswer: true,
			pans:             promisedAnswer{qid: parentQid},
		}

		conn.log.Debug().
			Int("qid", int(thisQid)).
			Int("pans", int(parentQid)).
			Msg("Prepared call for in-pipeline promised answer")
	} else {
		// What happened mate?!?!?!
		return 0, errors.New("unimplemented case")
	}

	return thisQid, nil
}

// commitOutMessage commits the changes of the pipeline step to the local Vat's
// state, under the assumption that the given pipeline step was successfully
// sent to the remote Vat.
func (v *Vat) commitOutMessage(_ context.Context, pipe *pipeline, stepIdx int, conn *runningConn) error {
	step := pipe.step(stepIdx)

	var qid QuestionId
	var q question
	if step.rpcMsg.isBootstrap {
		qid = step.rpcMsg.boot.qid
		conn.log.Debug().Int("qid", int(qid)).Msg("Comitted Bootstrap message")
	} else if step.rpcMsg.isCall {
		qid = step.rpcMsg.call.qid
		conn.log.Debug().Int("qid", int(qid)).Msg("Comitted Call message")
	} else {
		// Guard against errors while developing.
		return errors.New("unimplemented commitment of message")
	}

	// runtime.AddCleanup(step, conn.cleanupQuestionIdDueToUnref, qid) // TODO: Save cleanup in question in case of early finish?
	runtime.SetFinalizer(step, finalizePipelineStep)
	q = question{pipe: weak.Make(pipe), stepIdx: stepIdx}
	conn.questions.set(qid, q)

	// This step is now in flight. Allow forks from it to start. The forks
	// won't go out after the entirety of this pipeline has processed,
	// because this is running within the Vat's main goroutine.
	if qid > 0 {
		return step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
			if os != pipeStepStateBuilding {
				return os, ov, fmt.Errorf("invalid precondition state: %v", os)
			}
			ov.qid = qid
			ov.conn = conn
			return pipeStepStateRunning, ov, nil
		})
	}

	return nil
}

func (v *Vat) sendFinish(ctx context.Context, rc *runningConn, qid QuestionId) error {
	rc.mu.Lock()
	rc.questions.del(qid)
	rc.mu.Unlock()

	msg := message{
		isFinish: true,
		finish:   finish{qid: qid},
	}

	return rc.queue(ctx, singleMsgBatch(msg))
}

func (v *Vat) queueResolve(ctx context.Context, rc *runningConn, eid ExportId, exp export, resolution export) error {
	msg := message{
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
	msg := message{
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
	outMsg := singleMsgBatch(message{isAccept: true, accept: accept})
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
