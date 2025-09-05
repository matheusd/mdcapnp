// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"matheusd.com/mdcapnp/internal/sigvalue"
)

type pipelineStep struct {
	interfaceId   uint64
	methodId      uint16
	argsBuilder   func(*msgBuilder) error // Builds an rpc.Call struct
	paramsBuilder callParamsBuilder       // Builds the Params field of an rpc.Call struct

	stepRunning sigvalue.Once[QuestionId]

	stepDone sigvalue.Once[any]

	// Only accessed inside vat.Run().
	rpcMsg message
}

type pipelineState uint

const (
	pipelineStateBuilding pipelineState = iota
	pipelineStateBuilt
	pipelineStateRunning
	pipelineStateConnDone
)

type pipeline struct {
	// mu protects the following fields.
	mu    sync.Mutex
	state pipelineState
	steps []*pipelineStep

	// Only set on pipeline creation.
	vat           *Vat
	conn          *runningConn
	parent        *pipeline
	parentStepIdx int

	// Only set during pipeline running
	cancel func(error)
}

var fatalEmptyPipeline = "empty pipeline"
var fatalInvalidStepIndex = "invalid step index"

const defaultPipelineSizeHint = 5

var errPipelineStarted = errors.New("pipeline started successfully (not a real error)")

func newPipeline(sizeHint int) *pipeline {
	steps := make([]*pipelineStep, 0, sizeHint)
	return &pipeline{steps: steps}
}

// State returns the current pipeline state.
func (pipe *pipeline) State() pipelineState {
	pipe.mu.Lock()
	res := pipe.state
	pipe.mu.Unlock()
	return res
}

// LastStep returns the last pipeline step, handling special cases like a newly
// forked pipeline.
func (pipe *pipeline) LastStep() *pipelineStep {
	pipe.mu.Lock()
	if len(pipe.steps) == 0 {
		pipe.mu.Unlock()
		if pipe.parent == nil {
			panic(fatalEmptyPipeline)
		}
		return pipe.parent.LastStep()
	}
	res := pipe.steps[len(pipe.steps)-1]
	pipe.mu.Unlock()
	return res
}

// Step returns the ith step of the pipeline, handling special cases like a
// newly forked pipeline.
func (pipe *pipeline) Step(i int) *pipelineStep {
	pipe.mu.Lock()
	if len(pipe.steps) == 0 {
		pipe.mu.Unlock()
		if i != -1 {
			panic(fatalInvalidStepIndex)
		} else if pipe.parent == nil {
			panic(fatalEmptyPipeline)
		}
		return pipe.parent.Step(pipe.parentStepIdx)
	}
	res := pipe.steps[len(pipe.steps)-1]
	pipe.mu.Unlock()
	return res
}

func (pipe *pipeline) wouldFork(i int) bool {
	return i != len(pipe.steps)-1
}

func (pipe *pipeline) addStep() int {
	pipe.steps = append(pipe.steps, &pipelineStep{
		stepRunning: sigvalue.MakeOnce[QuestionId](),
		stepDone:    sigvalue.MakeOnce[any](),
	})
	return len(pipe.steps) - 1
}

func (pipe *pipeline) fork(i, sizeHint int) *pipeline {
	fork := newPipeline(sizeHint)
	fork.vat = pipe.vat
	fork.conn = pipe.conn
	fork.parent = pipe
	fork.parentStepIdx = i
	return fork
}

var errPipelineNotBuildingState = errors.New("pipeline not in building state")

type futureCap[T any] struct {
	_ [0]T // Tag.

	pipe      *pipeline
	stepIndex int
}

func newRootFutureCap[T any](pipeSizeHint int) futureCap[T] {
	res := futureCap[T]{
		pipe:      newPipeline(pipeSizeHint),
		stepIndex: 0,
	}
	res.pipe.addStep()
	return res
}

// forkFuture forks the last step of the future into a new pipeline but does not
// add any other steps.
func forkFuture[T any](old futureCap[T], pipeSizeHint int) (fork futureCap[T]) {
	fork.pipe = old.pipe.fork(len(old.pipe.steps)-1, pipeSizeHint)
	fork.stepIndex = -1
	return
}

func remoteCall[T, U any](obj futureCap[T], iid uint64, mid uint16, pb callParamsBuilder) (res futureCap[U]) {
	pipe := obj.pipe
	pipe.mu.Lock()
	if pipe.state != pipelineStateBuilding || pipe.wouldFork(obj.stepIndex) {
		res.pipe = pipe.fork(obj.stepIndex, defaultPipelineSizeHint)
	} else if obj.stepIndex == -1 {
		// First call of a new fork from bootstrap. Conn comes from the
		// parent pipeline.
		res.pipe = pipe
	} else {
		res.pipe = pipe
	}

	res.stepIndex = res.pipe.addStep()
	step := res.pipe.steps[res.stepIndex]
	step.interfaceId = iid
	step.methodId = mid
	step.paramsBuilder = pb

	pipe.mu.Unlock()

	return
}

func waitResult[T any](ctx context.Context, cap futureCap[T]) (res T, err error) {
	// Run cap.pipe
	if err = cap.pipe.vat.execPipeline(ctx, cap.pipe); err != nil {
		return
	}

	// Wait until the required step of the pipeline completes or fails.
	var pipeRes any
	pipeRes, err = cap.pipe.Step(cap.stepIndex).stepDone.Wait(ctx)
	if err != nil {
		return
	}

	// Determine if the result was an error.
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
