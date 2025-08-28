// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"

	"matheusd.com/mdcapnp/internal/sigvalue"
)

type pipelineStep struct {
	conn          *runningConn
	interfaceId   uint64
	methodId      uint16
	argsBuilder   func(*msgBuilder) error // Builds an rpc.Call struct
	paramsBuilder callParamsBuilder       // Builds the Params field of an rpc.Call struct

	// Filled if this step forks the pipeline.
	sides       []*pipeline
	stepRunning *sigvalue.Once[struct{}] // FIXME: what type?

	// Filled when preparing this step for sending.
	msg *message
}

type pipeline struct {
	parent        *pipeline
	parentStepIdx int
	steps         []pipelineStep

	// Filled once the pipeline is running.
	ctx    context.Context
	cancel func(error)
}

const defaultPipelineSizeHint = 5

var errPipelineSuccessful = errors.New("pipeline successful (not a real error)")

func newPipeline(sizeHint int) *pipeline {
	steps := make([]pipelineStep, 1, max(1, sizeHint))
	return &pipeline{steps: steps}
}

func (pipe *pipeline) wouldFork(i int) bool {
	return i != len(pipe.steps)-1
}

func (pipe *pipeline) addStep(iid uint64, mid uint16, pb callParamsBuilder) int {
	pipe.steps = append(pipe.steps, pipelineStep{
		interfaceId:   iid,
		methodId:      mid,
		paramsBuilder: pb,
	})
	return len(pipe.steps) - 1
}

func (pipe *pipeline) fork(i, sizeHint int) *pipeline {
	fork := newPipeline(sizeHint)
	fork.parent = pipe
	fork.parentStepIdx = i

	step := &pipe.steps[i]
	step.sides = append(pipe.steps[i].sides, fork)
	if step.stepRunning == nil {
		step.stepRunning = new(sigvalue.Once[struct{}])
	}

	return fork
}

type futureCap[T any] struct {
	_         [0]T // Tag.
	pipe      *pipeline
	stepIndex int
}

func (fc futureCap[T]) wouldForkPipe() bool {
	return fc.pipe.wouldFork(fc.stepIndex)
}

func newRootFutureCap[T any](pipeSizeHint int) futureCap[T] {
	return futureCap[T]{
		pipe:      newPipeline(pipeSizeHint),
		stepIndex: 0,
	}
}

func remoteCall[T, U any](obj futureCap[T], iid uint64, mid uint16, pb callParamsBuilder) (res futureCap[U]) {
	if obj.wouldForkPipe() {
		res.pipe = obj.pipe.fork(obj.stepIndex, defaultPipelineSizeHint)
	} else {
		res.pipe = obj.pipe
		res.stepIndex = obj.pipe.addStep(iid, mid, pb)
	}
	return
}

func waitResult[T any](ctx context.Context, cap futureCap[T]) (T, error) {
	// Run cap.pipe
	cap.pipe.steps[0].conn.vat.ExecPipeline(ctx, cap.pipe)
	panic("boo")
}
