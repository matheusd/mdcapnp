// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"io"
)

type IOConn interface {
	io.Reader
	io.Writer
	Flush() error
}

type IOTransport struct {
	c IOConn
}
