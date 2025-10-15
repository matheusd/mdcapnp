// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnpser

import "errors"

func copyStruct(src Struct, dst *MessageBuilder) (AnyPointerBuilder, error) {
	panic("boo")
}

func Copy(src AnyPointer, dst *MessageBuilder) (AnyPointerBuilder, error) {
	switch {
	case src.IsZeroStruct():
		// Nothing to do.
		return ZeroStructAsPointerBuilder(), nil
	case src.IsStruct():
		panic("boo")

	default:
		return AnyPointerBuilder{}, errors.New("unsupported case in Copy()")
	}
}
