// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import "encoding/binary"

func decodeSingleSegmentStream(fb []byte) ([]byte, error) {
	if len(fb) < 8 { // 4 byte seg count + 4 byte seg size
		return nil, errShortSingleSegmentStream
	}

	segCount := binary.LittleEndian.Uint32(fb)
	if segCount != 0 {
		return nil, errStreamNotSingleSegment
	}

	segSize := WordCount(binary.LittleEndian.Uint32(fb[4:]))
	if ByteCount(len(fb)-8) < segSize.ByteCount() {
		return nil, errShortStreamSegSize{segSize: segSize.ByteCount(), streamLen: len(fb) - 8}
	}

	return fb[8:], nil
}
