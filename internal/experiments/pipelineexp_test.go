// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

type pipeArrayAltStep struct {
	x [8]int // Say this is all the data needed.
}

type pipeArrayAlt struct {
	steps []pipeArrayAltStep
}

type pipeArrayCap struct {
	pipe      *pipeArrayAlt
	stepIndex int
}

type pipeClosureCap struct {
	f         func()
	stepIndex int
}

// BenchmarkPipelineImplAlternatives benchmarks alternative strategies for
// implementing the pipeline.
func BenchmarkPipelineImplAlternatives(b *testing.B) {
	var pac *pipeArrayCap // Make every step escape to heap.

	// Pipeline is an array. Leave sizing to Go runtime.
	b.Run("array", func(b *testing.B) {
		b.ReportAllocs()
		pipe := &pipeArrayAlt{}
		for range b.N {
			pipe.steps = append(pipe.steps, pipeArrayAltStep{})
			pac = &pipeArrayCap{pipe: pipe, stepIndex: len(pipe.steps) - 1}
		}
		require.Equal(b, b.N-1, pac.stepIndex)
	})

	// Pipeline is an array. Presize it to the expected pipeline bounds.
	b.Run("presized array", func(b *testing.B) {
		b.ReportAllocs()
		pipe := &pipeArrayAlt{steps: make([]pipeArrayAltStep, 0, b.N)}
		for range b.N {
			pipe.steps = append(pipe.steps, pipeArrayAltStep{})
			pac = &pipeArrayCap{pipe: pipe, stepIndex: len(pipe.steps) - 1}
		}
		require.Equal(b, b.N-1, pac.stepIndex)
	})

	var pcc *pipeClosureCap // Make very step escape to heap.

	// Pipeline is a series of closures.
	b.Run("closure", func(b *testing.B) {
		b.ReportAllocs()
		pcc = &pipeClosureCap{
			f: func() {},
		}

		for range b.N {
			var args [8]int
			prev := pcc
			pcc = &pipeClosureCap{
				f: func() {
					if prev.f == nil { // Ensure previous step is captured.
						panic("boo")
					}
					if args[0] == 666 { // Ensure args is captured.
						panic("boo")
					}
				},
				stepIndex: pcc.stepIndex + 1,
			}
		}
		require.Equal(b, b.N, pcc.stepIndex)
	})
}
