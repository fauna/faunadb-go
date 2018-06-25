package faunadb

import (
	"reflect"
	"testing"
)

func TestWrapNil(t *testing.T) {
	var in *string

	actual := wrap(in)
	expected := NullV{}
	if actual != expected {
		t.Errorf("Nil test failed: wrap(%#v) == %#v, expected %#v", in, actual, expected)
	}
}

func TestWrapPODs(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		expected Expr
	}{
		{"String", "abc", StringV("abc")},
		{"Boolean true", true, BooleanV(true)},
		{"Boolean false", false, BooleanV(false)},
		{"Int", 42, LongV(42)},
		{"Int8", int8(42), LongV(42)},
		{"Int16", int16(42), LongV(42)},
		{"Int32", int32(42), LongV(42)},
		{"Int64", int64(42), LongV(42)},
		{"Uint", 42, LongV(42)},
		{"Uint8", uint8(42), LongV(42)},
		{"Uint16", uint16(42), LongV(42)},
		{"Uint32", uint32(42), LongV(42)},
		{"Uint64", uint64(42), LongV(42)},
		{"Float32", float32(1.0), DoubleV(1.0)},
		{"Float64", float64(1.0), DoubleV(1.0)},
	}

	for _, test := range tests {
		actual := wrap(test.in)
		if actual != test.expected {
			t.Errorf("POD test %s failed: wrap(%#v) in %#v, expected %#v", test.name, test.in, actual, test.expected)
		}
	}
}

func TestWrapMap(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		expected Expr
	}{
		{"non-string key", map[int]string{4: "h"}, errMapKeyMustBeString},
		{"string key", map[string]int{"h": 4},
			unescapedObj{
				"object": unescapedObj{
					"h": LongV(4),
				},
			},
		},
	}

	for _, test := range tests {
		actual := wrap(test.in)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("map test %s failed: wrap(%#v) in %#v, expected %#v", test.name, test.in, actual, test.expected)
		}
	}
}

func TestWrapStruct(t *testing.T) {
	in := struct {
		A string
		B int
	}{
		"hello",
		42,
	}

	actual := wrap(in)
	expected := unescapedObj{
		"object": unescapedObj{
			"A": StringV("hello"),
			"B": LongV(42),
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Struct test failed: wrap(%#v) == %#v, expected %#v", in, actual, expected)
	}
}

func TestWrapSliceAndArrays(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		expected Expr
	}{
		{"empty slice", []Expr{}, unescapedArr{}},
		{"One element slice", []int{1}, unescapedArr{LongV(1)}},
		{"Nested slice", [][]int{[]int{1}}, unescapedArr{unescapedArr{LongV(1)}}},

		{"empty array", [0]Expr{}, unescapedArr{}},
		{"One element array", [1]int{1}, unescapedArr{LongV(1)}},
	}

	for _, test := range tests {
		actual := wrap(test.in)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("POD test %s failed: wrap(%#v) in %#v, expected %#v", test.name, test.in, actual, test.expected)
		}
	}
}
