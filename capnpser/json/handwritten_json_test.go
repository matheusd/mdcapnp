// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package json

import (
	jt "github.com/go-json-experiment/json/jsontext"
	"matheusd.com/mdcapnp/capnpser"
)

// *********** cap with complex method **********
//
// interface testCapWithComplexMethod {
//   ComplexMethod(int64Value :int64, strValue :Text) -> (int64Value :int64, strValue :Text)
// }

var testCapWithComplexMethod_ComplexMethod_ParamsSize = capnpser.StructSize{DataSectionSize: 2, PointerSectionSize: 2}

type testCapWithComplexMethod_ComplexMethod_Params capnpser.Struct

func (s *testCapWithComplexMethod_ComplexMethod_Params) Int64Value() int64 {
	return (*capnpser.Struct)(s).Int64(1) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Params) StrValue() string {
	return (*capnpser.Struct)(s).String(1) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Params) UnsafeStrValue() string {
	return (*capnpser.Struct)(s).UnsafeString(1) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Params) MarshalJSONTo(je *jt.Encoder) error {
	e := JSONEncoderHelper{E: je}
	return e.BeginObj().
		String("int64Value").Int64(s.Int64Value()).
		String("strValue").String(s.StrValue()).
		EndObj().Err
}

type testCapWithComplexMethod_ComplexMethod_ParamsBuilder capnpser.StructBuilder

func (b *testCapWithComplexMethod_ComplexMethod_ParamsBuilder) SetInt64Value(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(1, v)
}

func (b *testCapWithComplexMethod_ComplexMethod_ParamsBuilder) SetStrValue(v string) error {
	return (*capnpser.StructBuilder)(b).SetString(1, v)
}

func (b *testCapWithComplexMethod_ComplexMethod_ParamsBuilder) UnmarshalJSONFrom(jd *jt.Decoder) error {
	d := JSONDecoderHelper{D: jd}
	d.BeginObj()
	for d.Err == nil && !d.IsEndOfObj() {
		fname := d.NextKey()
		switch fname {
		case "int64Value":
			d.TrackErr(b.SetInt64Value(d.Int64()))
		case "strValue":
			d.TrackErr(b.SetStrValue(d.Str()))
		default:
			// Compat with newer versions of the struct.
		}
	}
	d.EndObj()

	return nil
}

var testCapWithComplexMethod_ComplexMethod_ResultSize = capnpser.StructSize{DataSectionSize: 2, PointerSectionSize: 2}

type testCapWithComplexMethod_ComplexMethod_Result capnpser.Struct

func (s *testCapWithComplexMethod_ComplexMethod_Result) Int64Value() int64 {
	return (*capnpser.Struct)(s).Int64(0) // Different than params.
}

func (s *testCapWithComplexMethod_ComplexMethod_Result) StrValue() string {
	return (*capnpser.Struct)(s).String(0) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Result) UnsafeStrValue() string {
	return (*capnpser.Struct)(s).UnsafeString(0) // Different than result.
}

type testCapWithComplexMethod_ComplexMethod_ResultBuilder capnpser.StructBuilder

func (b *testCapWithComplexMethod_ComplexMethod_ResultBuilder) SetInt64Value(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(0, v)
}

func (b *testCapWithComplexMethod_ComplexMethod_ResultBuilder) SetStrValue(v string) error {
	return (*capnpser.StructBuilder)(b).SetString(0, v)
}
