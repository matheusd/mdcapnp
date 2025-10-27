// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"errors"
	"fmt"
	"testing"
	"unsafe"
)

var errDummy = errors.New("dummy error")

type msgBuilder struct{} // Alias to a serializer MessageBuilder
type callParamsBuilder func(*msgBuilder) error

type pipelineStep struct {
	interfaceId   uint64
	methodId      uint16
	paramsBuilder callParamsBuilder // Builds the Params field of an rpc.Call struct

	// Filled if this step forks the pipeline.
	sides       []*pipeline
	stepRunning *Once[struct{}] // FIXME: what type?
}

type pipeline struct {
	parent        *pipeline
	parentStepIdx int
	steps         []pipelineStep
}

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
		step.stepRunning = new(Once[struct{}])
	}

	return fork
}

type fpFutureStatic = struct {
	pipe      *pipeline
	stepIndex int
}

//go:noinline
func callFpFutureStatic(obj fpFutureStatic, iid uint64, mid uint16, pb callParamsBuilder) fpFutureStatic {
	return fpFutureStatic{obj.pipe, obj.pipe.addStep(iid, mid, pb)}
}

type fpFutureGeneric[T any] struct {
	pipe      *pipeline
	stepIndex int
}

//go:noinline
func callFpFutureGeneric[T, U any](obj fpFutureGeneric[T], iid uint64, mid uint16, pb callParamsBuilder) fpFutureGeneric[U] {
	return fpFutureGeneric[U]{obj.pipe, obj.pipe.addStep(iid, mid, pb)}
}

// API type simulation for a defined type
type fpStaticAPITypeDefined fpFutureStatic

//go:noinline
func (f fpStaticAPITypeDefined) next(s string) fpStaticAPITypeDefined {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	return callFpFutureStatic(f, uint64(100), uint16(1000), pb)
}

// API type simulation for an embedded type.
type fpStaticAPITypeEmbedded struct {
	fc fpFutureStatic
}

//go:noinline
func (f fpStaticAPITypeEmbedded) next(s string) fpStaticAPITypeEmbedded {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	return fpStaticAPITypeEmbedded{fc: callFpFutureStatic(f.fc, uint64(100), uint16(1000), pb)}
}

// API type simulation for an embedded type with a discriminator tag.
type fpStaticAPITypeTagged struct {
	_fpStaticAPITypeTagged struct{} // Zero sized, unique discriminator field.
	fc                     fpFutureStatic
}

//go:noinline
func (f fpStaticAPITypeTagged) next(s string) fpStaticAPITypeTagged {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	return fpStaticAPITypeTagged{fc: callFpFutureStatic(f.fc, uint64(100), uint16(1000), pb)}
}

type fpGenericAPIType fpFutureGeneric[string]

//go:noinline
func (f fpGenericAPIType) next(s string) fpGenericAPIType {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	return fpGenericAPIType(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

func TestFutureTypeAlternativesSizes(t *testing.T) {
	t.Logf(" Defined: size: %v, align: %v", unsafe.Sizeof(fpStaticAPITypeDefined{}), unsafe.Alignof(fpStaticAPITypeDefined{}))
	t.Logf("Embedded: size: %v, align: %v", unsafe.Sizeof(fpStaticAPITypeEmbedded{}), unsafe.Alignof(fpStaticAPITypeEmbedded{}))
	t.Logf("  Tagged: size: %v, align: %v", unsafe.Sizeof(fpStaticAPITypeTagged{}), unsafe.Alignof(fpStaticAPITypeTagged{}))
	t.Logf(" Generic: size: %v, align: %v", unsafe.Sizeof(fpGenericAPIType{}), unsafe.Alignof(fpGenericAPIType{}))
}

func BenchmarkFutureTypeAlternatives(b *testing.B) {
	b.Run("defined", func(b *testing.B) {
		f := fpStaticAPITypeDefined{pipe: &pipeline{steps: make([]pipelineStep, 0, b.N)}}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next("fooo")
		}

		if f.stepIndex != b.N-1 {
			panic(fmt.Sprintf("%d vs %d", f.stepIndex, b.N))
		}
	})

	b.Run("embedded", func(b *testing.B) {
		f := fpStaticAPITypeEmbedded{fc: fpFutureStatic{pipe: &pipeline{steps: make([]pipelineStep, 0, b.N)}}}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next("fooo")
		}

		if f.fc.stepIndex != b.N-1 {
			panic(fmt.Sprintf("%d vs %d", f.fc.stepIndex, b.N))
		}
	})

	b.Run("tagged", func(b *testing.B) {
		f := fpStaticAPITypeTagged{fc: fpFutureStatic{pipe: &pipeline{steps: make([]pipelineStep, 0, b.N)}}}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next("fooo")
		}

		if f.fc.stepIndex != b.N-1 {
			panic(fmt.Sprintf("%d vs %d", f.fc.stepIndex, b.N))
		}
	})

	b.Run("generic", func(b *testing.B) {
		f := fpGenericAPIType{pipe: &pipeline{steps: make([]pipelineStep, 0, b.N)}}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			f = f.next("fooo")
		}

		if f.stepIndex != b.N-1 {
			panic(fmt.Sprintf("%d vs %d", f.stepIndex, b.N))
		}
	})

}
