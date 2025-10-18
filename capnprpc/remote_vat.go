// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

type RemoteVat struct {
	rc *runningConn
}

func (rv RemoteVat) Bootstrap() BootstrapFuture {
	return rv.rc.Bootstrap()
}
