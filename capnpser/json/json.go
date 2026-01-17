// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package json

import (
	"errors"

	"github.com/go-json-experiment/json/jsontext"
)

type JSONEncoderHelper struct {
	E   *jsontext.Encoder
	Err error
}

func (h *JSONEncoderHelper) BeginObj() *JSONEncoderHelper {
	if h.Err == nil {
		h.E.WriteToken(jsontext.BeginObject)
	}
	return h
}

func (h *JSONEncoderHelper) EndObj() *JSONEncoderHelper {
	if h.Err == nil {
		h.E.WriteToken(jsontext.EndObject)
	}
	return h
}

func (h *JSONEncoderHelper) String(v string) *JSONEncoderHelper {
	if h.Err == nil {
		h.Err = h.E.WriteToken(jsontext.String(v))
	}
	return h
}

func (h *JSONEncoderHelper) Int64(v int64) *JSONEncoderHelper {
	if h.Err == nil {
		h.Err = h.E.WriteToken(jsontext.Int(v))
	}
	return h

}

type JSONDecoderHelper struct {
	D   *jsontext.Decoder
	Err error
}

func (h *JSONDecoderHelper) TrackErr(err error) {
	if h.Err == nil {
		h.Err = err
	}
}

func (h *JSONDecoderHelper) BeginObj() *JSONDecoderHelper {
	if h.Err == nil {
		var t jsontext.Token
		t, h.Err = h.D.ReadToken()
		if h.Err == nil {
			return h
		}
		if t.Kind() != '{' {
			h.Err = errors.New("not begining of object")
			return h
		}
	}
	return h
}

func (h *JSONDecoderHelper) EndObj() {
	if h.Err != nil {
		return
	}

	var t jsontext.Token
	t, h.Err = h.D.ReadToken()
	if h.Err != nil {
		return
	}
	if t.Kind() != '}' {
		h.Err = errors.New("not end of object")
	}
	return
}

func (h *JSONDecoderHelper) IsEndOfObj() bool {
	if h.Err != nil {
		return true
	}
	return h.D.PeekKind() == '}'
}

func (h *JSONDecoderHelper) NextKey() string {
	if h.Err != nil {
		return ""
	}

	var t jsontext.Token
	t, h.Err = h.D.ReadToken()
	if h.Err != nil {
		return ""
	}

	if t.Kind() != '"' {
		h.Err = errors.New("not the expected object key")
	}

	return t.String()
}

func (h *JSONDecoderHelper) Int64() int64 {
	if h.Err != nil {
		return 0
	}
	var t jsontext.Token
	t, h.Err = h.D.ReadToken()
	if h.Err != nil {
		return 0
	}
	if t.Kind() != '0' {
		h.Err = errors.New("not the expected number")
		return 0
	}
	return t.Int()
}

func (h *JSONDecoderHelper) Str() string {
	if h.Err != nil {
		return ""
	}

	var t jsontext.Token
	t, h.Err = h.D.ReadToken()
	if h.Err != nil {
		return ""
	}

	if t.Kind() != '"' {
		h.Err = errors.New("not the expected string key")
	}

	return t.String()
}
