// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"runtime"
	"testing"
	"unsafe"
)

var (
	gUpTestSlice01   []byte
	gUpTestSlice02   []byte
	gUpTestUnsafePtr unsafe.Pointer
	gUpTestRange     int = 256
)

func TestUnsafePointerOOB(t *testing.T) {
	gUpTestSlice01 = make([]byte, 256) // Len must be the same as gUpTestRange
	gUpTestSlice02 = make([]byte, 256)

	// Create two heap slices that are contiguous in memory.
	const testRange = 1000000
	for i := range testRange {
		gUpTestSlice01 = make([]byte, 256, 256)
		gUpTestSlice02 = make([]byte, 256)
		contiguous := unsafe.Pointer(unsafe.SliceData(gUpTestSlice02)) ==
			unsafe.Add(unsafe.Pointer(unsafe.SliceData(gUpTestSlice01)), cap(gUpTestSlice01))
		if contiguous {
			break
		}
		if i == testRange-1 {
			panic("could not make contiguous slices for test")
		}
	}

	runtime.Gosched()

	gUpTestUnsafePtr = unsafe.Pointer(&gUpTestSlice01[0])

	runtime.Gosched()

	// This proves you can write past unsafe pointers.
	*(*byte)(unsafe.Add(gUpTestUnsafePtr, gUpTestRange-1)) = 1
	*(*byte)(unsafe.Add(gUpTestUnsafePtr, gUpTestRange)) = 2

	t.Logf("gUptestRange01[255] = %x", gUpTestSlice01[len(gUpTestSlice01)-1])
	t.Logf("gUptestRange02[0]   = %x", gUpTestSlice02[0])
}
