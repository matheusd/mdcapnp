// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package mdcapnp

import (
	"fmt"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

func TestWordOffset(t *testing.T) {
	tests := []struct {
		v     uint32
		valid bool
	}{
		{v: 0x00000000, valid: true},
		{v: 0x00000001, valid: true},
		{v: 0x80000000, valid: true},
		{v: 0x8fffffff, valid: true},
		{v: 0x10000000, valid: false},
		{v: 0x20000000, valid: false},
		{v: 0x40000000, valid: false},
		{v: 0x30000000, valid: false},
		{v: 0x70000000, valid: false},
		{v: 0x20000001, valid: false},
		{v: 0x7fffffff, valid: false},
		{v: 0x9fffffff, valid: false},
		{v: 0xcfffffff, valid: false},
		{v: 0xffffffff, valid: false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%x", tc.v), func(t *testing.T) {
			v := WordOffset(tc.v)
			got := v.Valid()
			require.Equal(t, tc.valid, got)
		})
	}
}
