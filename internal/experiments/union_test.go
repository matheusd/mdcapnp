// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"encoding/json"
	"testing"

	"matheusd.com/depvendoredtestify/require"
)

type unionStrValue string

func (v unionStrValue) xxx() {}

type unionIntValue int

func (v unionIntValue) xxx() {}

type unionFloatValue float64

func (v unionFloatValue) xxx() {}

type unionStructValue struct {
	StructField string `json:"struct_field"`
}

func (v *unionStructValue) xxx() {}

type unionAnyValue struct {
	V any
}

func (v unionAnyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.V)
}

func (v unionAnyValue) xxx() {}

type unionType interface {
	xxx()
}

type testStructWithUnion struct {
	FieldA       string    `json:"field_a"`
	FieldB       int       `json:"field_b"`
	FieldUnion   unionType `json:"field_union"`
	FieldOmitted unionType `json:"field_union_omitted,omitempty"`
}

// TestUnionAlternatives tests how to encode union types.
func TestUnionAlternatives(t *testing.T) {
	anyStructVal := struct {
		Foo    string `json:"foo"`
		Blergh int    `json:"blergh,omitempty"`
	}{Foo: "blaaaa"}

	tests := []struct {
		name string
		v    testStructWithUnion
	}{
		{
			name: "zero union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100},
		},
		{
			name: "str union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100, FieldUnion: unionStrValue("string value")},
		},
		{
			name: "int union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100, FieldUnion: unionIntValue(1701)},
		},
		{
			name: "float union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100, FieldUnion: unionFloatValue(10.7419)},
		},
		{
			name: "struct union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100, FieldUnion: &unionStructValue{"subfield value"}},
		},
		{
			name: "any union value",
			v:    testStructWithUnion{FieldA: "v1", FieldB: 100, FieldUnion: unionAnyValue{V: anyStructVal}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.MarshalIndent(tc.v, "", "  ")
			require.NoError(t, err)
			t.Logf("%s", encoded)
		})
	}
}
