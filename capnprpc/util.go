// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
}

func vatLogFieldsFormatter(v interface{}, name string) string {
	if name == "vat" {
		return fmt.Sprintf("%10s  ", v)
	} else if name == "time" {
		return fmt.Sprintf("%d", v)
	}
	return ""
}

func setupDevLogFormat(cw *zerolog.ConsoleWriter) {
	// cw.TimeFormat = zerolog.TimeFormatUnix
	cw.PartsOrder = []string{"time", "vat", "level", "message"}
	cw.FormatPartValueByName = vatLogFieldsFormatter
	cw.TimeFormat = "2006-01-02T15:04:05.000000"
}

func testFrameFormatter(v any) string {
	return fmt.Sprintf("%21s", v)
}

func testLogger(t testing.TB) zerolog.Logger {
	ctw := zerolog.NewTestWriter(t)
	ctw.Frame = 6
	// ctw.FrameFormatter = testFrameFormatter

	cw := zerolog.NewConsoleWriter()
	cw.Out = ctw
	// cw.Out = os.Stdout
	setupDevLogFormat(&cw)
	return zerolog.New(cw)
}
