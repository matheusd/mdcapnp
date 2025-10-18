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

	inBuf   []byte
	inArena capnpser.Arena
}

func nopFlush() error { return nil }

func NewIOTransport(remoteName string, rw net.Conn) *IOTransport {
	r := bufio.NewReader(rw)
	// w := bufio.NewWriter(rw)
	// flush := w.Flush
	flush := nopFlush
	w := rw
	close := rw.Close

	return &IOTransport{
		flush:   flush,
		remName: remoteName,
		r:       r,
		w:       w,
		closeRW: close,
		inBuf:   make([]byte, 128*1024),
	}
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
	// fmt.Println("XXXXXX wrote", iot.remName, n)
	if n != len(serBytes) {
		return io.ErrShortWrite
	}
	return nil

	// return iot.flush()
}

func (iot *IOTransport) receive(ctx context.Context) (InMsg, error) {
	for ctx.Err() == nil {

		// Read header.
		//
		// TODO: abstract and move to capnpser.
		_, err := io.ReadFull(iot.r, iot.inBuf[:8])
		if err != nil {
			return InMsg{}, err
		}

		segCount := binary.LittleEndian.Uint32(iot.inBuf[:4])
		if segCount != 0 {
			// TODO: support multi segment.
			return InMsg{}, errors.New("multi-seg not supported in IOTransport.receive")
		}

		// TODO: protect against too large reads.
		seg0SizeWords := binary.LittleEndian.Uint32(iot.inBuf[4:])
		seg0SizeBytes := int(seg0SizeWords * capnpser.WordSize)

		if seg0SizeBytes == 0 {
			// Empty message???
			continue
		}

		if len(iot.inBuf) < seg0SizeBytes+8 {
			oldHeader := iot.inBuf[:8]
			iot.inBuf = make([]byte, seg0SizeBytes)
			copy(iot.inBuf, oldHeader)
		}

		_, err = io.ReadFull(iot.r, iot.inBuf[8:8+seg0SizeBytes])
		if err != nil {
			return InMsg{}, err
		}

		err = iot.inArena.DecodeSingleSegment(iot.inBuf[:8+seg0SizeBytes])
		if err != nil {
			return InMsg{}, err
		}
		iot.inArena.ReadLimiter().InitConcurrentUnsafe(64 * 1024 * 1024)

		res := InMsg{Msg: capnpser.MakeMsg(&iot.inArena)}
		res.Msg.SetDepthLimit(64)

		return res, nil
	}
	return InMsg{}, ctx.Err()
}

func (iot *IOTransport) close() error {
	return iot.closeRW()
}

func (iot *IOTransport) remoteName() string {
	return iot.remName
}
