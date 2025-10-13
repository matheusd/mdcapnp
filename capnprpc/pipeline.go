// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"

	"matheusd.com/mdcapnp/internal/sigvalue"
)

// callSetup are the requirements to build and send a Call message.
type callSetup struct {
	interfaceId   uint64
	methodId      uint16
	paramsBuilder callParamsBuilder // Builds the Params field of an rpc.Call struct
}

type pipelineStepState int

const (
	pipeStepStateBuilding pipelineStepState = iota
	pipeStepStateSending
	pipeStepStateRunning
	pipeStepStateDone
	pipeStepStateFinished
	pipeStepFailed // Must be last to be > all other states.
)

type pipelineStepStateValue struct {
	qid   QuestionId   // Set if step state is >= running.
	iid   ImportId     // Set if this has been resolved into an import (either hosted or promise)
	value any          // Set if step state is >= done.
	err   error        // Set if step state is >= failed.
	conn  *runningConn // May be changed as step is resolved (3PH).
}

type pipelineStep struct {
	value  sigvalue.Stateful[pipelineStepState, pipelineStepStateValue]
	parent *pipelineStep
	csetup callSetup
}

func finalizePipelineStep(step *pipelineStep) {
	var qid QuestionId
	var conn *runningConn
	_ = step.value.Modify(func(os pipelineStepState, ov pipelineStepStateValue) (pipelineStepState, pipelineStepStateValue, error) {
		if os != pipeStepFailed {
			qid = ov.qid
			ov.qid = 0
			os = pipeStepStateFinished
			conn = ov.conn
		}

		return os, ov, nil
	})

	if qid > 0 {
		conn.cleanupQuestionIdDueToUnref(qid) // TODO: How does this conflict with other uses?
	}
}

func newRootStep(v *Vat) *pipelineStep {
	return &pipelineStep{}
}

func newChildStep(parent *pipelineStep) *pipelineStep {
	return &pipelineStep{
		parent: parent,
	}
}

type callFuture struct {
	step *pipelineStep
}

func newRootFutureCap(v *Vat) callFuture {
	return callFuture{
		step: newRootStep(v),
	}
}

func remoteCall(obj callFuture, csetup callSetup) (res callFuture) {
	parentStep := obj.step

	// Every call creates a new step with parent reference
	res.step = newChildStep(parentStep)
	res.step.csetup = csetup

	return res
}

func castCallResult[T any](callResult any) (res T, err error) {
	// Check if result is already the expected return type.
	var ok bool
	if res, ok = callResult.(T); ok {
		return
	}

	// Not an error, and not type T, so must be an AnyPointer with the
	// Content field of a Payload result.
	var resAnyPointer anyPointer
	if resAnyPointer, ok = callResult.(anyPointer); !ok {
		err = fmt.Errorf("future expected AnyPointer but got %T", callResult)
		return
	}

	// Content may be either a Struct or a CapPointer, and T will be an
	// alias to one of these (depending on what's expected based on the
	// schema).
	//
	// TODO: better way to convert to T?
	var contentAny any
	if resAnyPointer.IsStruct() {
		contentAny = resAnyPointer.AsStruct()
	} else if resAnyPointer.IsCapPointer() {
		// TODO: How to convert capPointer to importId???
		contentAny = resAnyPointer.AsCapPointer()
	} else {
		err = fmt.Errorf("content is not struct or capPointer")
		return
	}

	res, ok = contentAny.(T)
	if !ok {
		err = fmt.Errorf("future expected %T but got back %T", res, contentAny)
	}

	return
}

func castCallResultOrErr[T any](callResult any, inErr error) (res T, err error) {
	if inErr != nil {
		err = inErr
	} else {
		res, err = castCallResult[T](callResult)
	}
	return
}

func waitResult(ctx context.Context, cap callFuture) (res any, err error) {
	// Find the vat.
	var vat *Vat
	step := cap.step
	for step != nil && vat == nil {
		conn := step.value.GetValue().conn
		if conn != nil {
			vat = conn.vat
		}
		step = step.parent
	}

	if vat == nil {
		err = errors.New("could not find vat for pipeline")
		return
	}

	// Run pipeline.
	if err = vat.execStep(ctx, cap.step); err != nil {
		return
	}

	// Wait until the required step of the pipeline completes or fails.
	stepState, stepValue, err := cap.step.value.WaitStateAtLeast(ctx, pipeStepStateDone)
	if err != nil {
		return
	}
	if stepState == pipeStepFailed {
		err = stepValue.err
		if err == nil {
			err = fmt.Errorf("unknown error in failed pipeline step state: %v", stepState)
		}
		return
	}

	// Determine if the result was an error.
	pipeRes := stepValue.value
	var ok bool
	if err, ok = pipeRes.(error); ok && err != nil {
		return
	}

	res = pipeRes
	return
}

func releaseFuture(ctx context.Context, cap callFuture) {
	finalizePipelineStep(cap.step)
}
