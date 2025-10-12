// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"

	"matheusd.com/mdcapnp/internal/sigvalue"
)

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
	interfaceId   uint64
	methodId      uint16
	argsBuilder   func(*msgBuilder) error // Builds an rpc.Call struct
	paramsBuilder callParamsBuilder       // Builds the Params field of an rpc.Call struct

	value sigvalue.Stateful[pipelineStepState, pipelineStepStateValue]

	// Only accessed by the pipeline's execution goroutine.
	rpcMsg *message

	vat    *Vat
	parent *pipelineStep // Set only for forked steps
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
	return &pipelineStep{vat: v}
}

func newChildStep(parent *pipelineStep) *pipelineStep {
	return &pipelineStep{
		vat:    parent.vat,
		parent: parent,
	}
}

type futureCap[T any] struct {
	_ [0]T // Tag.

	step *pipelineStep
}

func newRootFutureCap[T any](v *Vat) futureCap[T] {
	return futureCap[T]{
		step: newRootStep(v),
	}
}

func remoteCall[T, U any](obj futureCap[T], iid uint64, mid uint16, pb callParamsBuilder) (res futureCap[U]) {
	parentStep := obj.step

	// Every call creates a new step with parent reference
	res.step = newChildStep(parentStep)
	res.step.interfaceId = iid
	res.step.methodId = mid
	res.step.paramsBuilder = pb

	return res
}

func waitResult[T any](ctx context.Context, cap futureCap[T]) (res T, err error) {
	// Run pipeline.
	if err = cap.step.vat.execStep(ctx, cap.step); err != nil {
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
			err = fmt.Errorf("unknown final pipeline step state: %v", stepState)
		}
		return
	}

	// Determine if the result was an error.
	pipeRes := stepValue.value
	var ok bool
	if err, ok = pipeRes.(error); ok && err != nil {
		return
	}

	// Check if result is the expected return type.
	if _, ok = pipeRes.(T); ok {
		res = pipeRes.(T)
		return
	}

	// FIXME: nothing else needed after this???

	// Not an error, so must be an AnyPointer with the Content field of a
	// Payload result.
	var content anyPointer
	if content, ok = pipeRes.(anyPointer); !ok {
		err = fmt.Errorf("future expected AnyPointer but got %T", pipeRes)
		return
	}

	// Content may be either a Struct or a CapPointer, and T will be an
	// alias to one of these (depending on what's expected based on the
	// schema).
	//
	// TODO: better way to convert to T?
	var contentAny any
	if content.IsStruct() {
		contentAny = content.AsStruct()
	} else if content.IsCapPointer() {
		contentAny = content.AsCapPointer()
	} else {
		err = fmt.Errorf("content is not struct or capPointer")
		return
	}

	res, ok = contentAny.(T)
	if !ok {
		err = fmt.Errorf("future expected %T but got back %T", res, pipeRes)
	}

	return
}

func releaseFuture[T any](ctx context.Context, cap futureCap[T]) {
	finalizePipelineStep(cap.step)
}
