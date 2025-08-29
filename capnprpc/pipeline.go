// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"

	"matheusd.com/mdcapnp/capnpser"
	"matheusd.com/mdcapnp/internal/sigvalue"
)

type pipelineStep struct {
	conn          *runningConn
	interfaceId   uint64
	methodId      uint16
	argsBuilder   func(*msgBuilder) error // Builds an rpc.Call struct
	paramsBuilder callParamsBuilder       // Builds the Params field of an rpc.Call struct

	stepDone *sigvalue.Once[any] // FIXME: what type?

	// Filled if this step forks the pipeline.
	// sides       []*pipeline // Is this needed?
	stepRunning *sigvalue.Once[struct{}] // FIXME: what type?

	// Filled when preparing this step for sending.
	serMsg capnpser.Message
	rpcMsg Message
	qid    QuestionId
}

type pipeline struct {
	parent        *pipeline
	parentStepIdx int
	steps         []pipelineStep

	// Filled once the pipeline is running.
	ctx    context.Context
	cancel func(error)
}

var fatalEmptyPipeline = "empty pipeline"

// firstVat returns the first vat for this pipeline.
func (p *pipeline) firstVat() *vat {
	if len(p.steps) == 0 {
		if p.parent == nil {
			panic(fatalEmptyPipeline)
		}

		// Pipeline derived from bootstrap without additional
		// calls.
		return p.parent.firstVat()
	}

	return p.steps[0].conn.vat
}

func (p *pipeline) lastStepDoneOnce() *sigvalue.Once[any] {
	if len(p.steps) == 0 {
		if p.parent == nil {
			panic(fatalEmptyPipeline)
		}
		return p.parent.lastStepDoneOnce()
	}
	return p.steps[len(p.steps)-1].stepDone
}

const defaultPipelineSizeHint = 5

var errPipelineSuccessful = errors.New("pipeline successful (not a real error)")

func newPipeline(sizeHint int) *pipeline {
	steps := make([]pipelineStep, 0, sizeHint)
	return &pipeline{steps: steps}
}

func (pipe *pipeline) wouldFork(i int) bool {
	return i != len(pipe.steps)-1
}

func (pipe *pipeline) addStep() int {
	pipe.steps = append(pipe.steps, pipelineStep{})
	return len(pipe.steps) - 1
}

func (pipe *pipeline) fork(i, sizeHint int) *pipeline {
	fork := newPipeline(sizeHint)
	fork.parent = pipe
	fork.parentStepIdx = i

	step := &pipe.steps[i]
	// step.sides = append(pipe.steps[i].sides, fork)
	if step.stepRunning == nil {
		step.stepRunning = new(sigvalue.Once[struct{}])
	}

	return fork
}

func (pipe *pipeline) setupWaitReqs() {
	if len(pipe.steps) == 0 {
		if pipe.parent == nil {
			panic(fatalEmptyPipeline)
		}

		pipe.parent.setupWaitReqs()
		return
	}

	step := pipe.steps[len(pipe.steps)-1]
	step.stepDone = &sigvalue.Once[any]{}
}

type futureCap[T any] struct {
	_ [0]T // Tag.

	pipe      *pipeline
	stepIndex int
}

func (fc futureCap[T]) wouldForkPipe() bool {
	return fc.pipe.wouldFork(fc.stepIndex)
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
	var rc *runningConn
	if obj.stepIndex == -1 {
		// First call of a new fork from bootstrap. Conn comes from the
		// parent pipeline.
		rc = obj.pipe.parent.steps[obj.pipe.parentStepIdx].conn
	} else if obj.wouldForkPipe() {
		res.pipe = obj.pipe.fork(obj.stepIndex, defaultPipelineSizeHint)
		rc = obj.pipe.steps[obj.stepIndex].conn // Same conn as parent.
	} else {
		res.pipe = obj.pipe
		rc = obj.pipe.steps[obj.stepIndex].conn // Same conn as parent.
	}

	res.stepIndex = obj.pipe.addStep()
	step := &res.pipe.steps[res.stepIndex]
	step.conn = rc
	step.interfaceId = iid
	step.methodId = mid
	step.paramsBuilder = pb

	return
}

func waitResult[T any](ctx context.Context, cap futureCap[T]) (res T, err error) {
	// Setup the stuff needed to wait for this specific future.
	cap.pipe.setupWaitReqs()

	// Run cap.pipe
	cap.pipe.firstVat().execPipeline(ctx, cap.pipe)

	// Wait until the last step of the pipeline completes or fails.
	var pipeRes any
	pipeRes, err = cap.pipe.lastStep().stepDone.Wait(ctx)
	if err != nil {
		return
	}

	// Determine if the result was an error.
	var ok bool
	if err, ok = pipeRes.(error); ok && err != nil {
		return
	}

	// Convert to expected type.
	res, ok = pipeRes.(T)
	if !ok {
		err = fmt.Errorf("future expected %T but got back %T")
	}

	return
}
