// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"errors"

	"github.com/rs/zerolog"
)

// extraDataError is an error type that can add more data to a log event.
type extraDataError interface {
	addExtraDataToLog(e *zerolog.Event)
}

type errDisembargoAcceptUnknownExport ExportId

func (err errDisembargoAcceptUnknownExport) Error() string {
	return "received disembargo.accept for unknown export"
}

func (err errDisembargoAcceptUnknownExport) addExtraDataToLog(e *zerolog.Event) {
	e.Int("eid", int(err))
}

var errConnStopped = errors.New("conn stopped")
var errPipelineNotBuildingState = errors.New("pipeline not in building state")
var errPipeStepAlreadyFinished = errors.New("pipeline step already finished")
var errPipeParentStepAlreadyFinished = errors.New("pipeline parent step already finished")
