// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync/atomic"

	"matheusd.com/mdcapnp/capnpser"
)

type Flusher interface {
	Flush() error
}

type IOTransport struct {
	r       io.Reader
	w       net.Conn // io.Writer
	closeRW func() error
	flush   func() error

	remName string

	inPool *inMsgPool

	// depthLimit atomic.Uintptr
	readLimit atomic.Uint64
}

func nopFlush() error { return nil }

func NewIOTransport(remoteName string, rw net.Conn) *IOTransport {
	r := bufio.NewReader(rw)
	// w := bufio.NewWriter(rw)
	// flush := w.Flush
	flush := nopFlush
	w := rw
	close := rw.Close

	iot := &IOTransport{
		flush:   flush,
		remName: remoteName,
		r:       r,
		w:       w,
		closeRW: close,
		inPool:  newArenaPool(),
	}
	// iot.depthLimit.Store(64) // TODO: constant somewhere?
	iot.readLimit.Store(64 * 1024 * 1024)
	return iot
}

// SetReadDepthLimit sets the depth limit for traversing messages read from this
// remote end.
/*
func (iot *IOTransport) SetReadDepthLimit(limit uint) {
	iot.depthLimit.Store(uintptr(limit))
}
*/

// SetReadDepthLimit sets the read limit (in bytes) for traversing messages read
// from this remote end.
func (iot *IOTransport) SetReadLimit(limit uint64) {
	iot.readLimit.Store(limit)
}

func (iot *IOTransport) send(ctx context.Context, outMsg OutMsg) error {
	serBytes, err := outMsg.Msg.Serialize()
	if err != nil {
		return err
	}

	n, err := iot.w.Write(serBytes)
	if err != nil {
		return err
	}
	if n != len(serBytes) {
		return io.ErrShortWrite
	}
	return nil

	// return iot.flush()
}

func (iot *IOTransport) receive(ctx context.Context, prev *InMsg) (*InMsg, error) {
	var res *InMsg = prev
	for ctx.Err() == nil {
		if prev == nil {
			res = iot.inPool.new()
		}

		// Read header.
		//
		// TODO: abstract and move to capnpser.
		_, err := io.ReadFull(iot.r, res.buf[:8])
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// Graceful local close.
				return nil, io.EOF
			}
			return nil, err
		}

		segCount := binary.LittleEndian.Uint32(res.buf[:4])
		if segCount != 0 {
			// TODO: support multi segment.
			return nil, errors.New("multi-seg not supported in IOTransport.receive")
		}

		// TODO: Support multi-segment messages.

		// TODO: protect against too large reads.
		seg0SizeWords := binary.LittleEndian.Uint32(res.buf[4:])
		seg0SizeBytes := int(seg0SizeWords) * capnpser.WordSize

		if seg0SizeBytes == 0 {
			// Empty message???
			continue
		}

		if cap(res.buf) < seg0SizeBytes+8 {
			res.buf = iot.inPool.replaceBuf(res.buf, capnpser.WordCount(1+seg0SizeWords))
		}
		res.buf = res.buf[:8+seg0SizeBytes]

		_, err = io.ReadFull(iot.r, res.buf[8:])
		if err != nil {
			return nil, err
		}

		// TODO: Already decoded, no need to decode again. Have an API
		// to build it directly.
		arena := res.Msg.Arena()
		err = arena.DecodeSingleSegment(res.buf)
		if err != nil {
			return nil, err
		}
		arena.ReadLimiter().InitConcurrentUnsafe(iot.readLimit.Load())

		// res.Msg = capnpser.MakeMsg(res.arena)
		// res.Msg.SetDepthLimit(uint(iot.depthLimit.Load()))

		return res, nil
	}
	return nil, ctx.Err()
}

func (iot *IOTransport) close() error {
	return iot.closeRW()
}

func (iot *IOTransport) remoteName() string {
	return iot.remName
}
