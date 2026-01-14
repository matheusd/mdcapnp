// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"errors"
	"fmt"

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

type ErrUnexpectedInterfaceId struct {
	Want InterfaceId
	Got  InterfaceId
}

func (err ErrUnexpectedInterfaceId) Error() string {
	return fmt.Sprintf("unexpected interfaceId (want 0x%016x got 0x%016x)",
		err.Want, err.Got)
}

type ErrUnexpectedMethodId struct {
	Want MethodId
	Got  MethodId
}

func (err ErrUnexpectedMethodId) Error() string {
	return fmt.Sprintf("unexpected methodId (want 0x%04x got 0x%04x)",
		err.Want, err.Got)
}

type ErrUnknownMethodId struct {
	IId InterfaceId
	Mid MethodId
}

func (err ErrUnknownMethodId) Error() string {
	return fmt.Sprintf("unknown method id 0x%x in interface 0x%x", err.Mid, err.IId)
}

var errConnStopped = errors.New("conn stopped")
var errPipelineNotBuildingState = errors.New("pipeline not in building state")
var errPipeStepAlreadyFinished = errors.New("pipeline step already finished")
var errPipeParentStepAlreadyFinished = errors.New("pipeline parent step already finished")
