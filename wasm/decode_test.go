package wasm

import (
	"fmt"
	"reflect"
	"testing"
)

func makeDecoderTest[T any](t *testing.T, name string, expected T) {
	t.Run(name, func(t *testing.T) {
		data, err := Encode(expected)
		if err != nil {
			t.Fatalf("can't make test %s", name)
		}
		actual, err := Decode[T](data)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Actual: %#v, Expected: %#v", actual, expected)
		}
	})
}

func TestDecode(t *testing.T) {
	type simpleTestStruct struct {
		One   int8
		Two   float32
		Three string
	}
	type testStruct struct {
		SomeMap1 map[string]simpleTestStruct
		SomeStr  string
		SomeInt  int64
		SomeMap2 map[string]simpleTestStruct
	}
	makeDecoderTest(t, "bool", true)
	makeDecoderTest(t, "int8", int8(1))
	makeDecoderTest(t, "int16", int16(1))
	makeDecoderTest(t, "int32", int32(1))
	makeDecoderTest(t, "int64", int64(1))
	makeDecoderTest(t, "uint8", uint8(1))
	makeDecoderTest(t, "uint16", uint16(1))
	makeDecoderTest(t, "uint32", uint32(1))
	makeDecoderTest(t, "uint64", uint64(1))
	makeDecoderTest(t, "float32", float32(1.12))
	makeDecoderTest(t, "float64", float64(1.13))
	makeDecoderTest(t, "string", "hello world")
	makeDecoderTest(t, "[4]int8", [4]int8{1, 2, 3, 4})
	makeDecoderTest(t, "[]int8", []int8{1, 2, 3, 4})
	makeDecoderTest(t, "map[string]int8", map[string]int8{"one": 1, "two": 2, "three": 3})
	makeDecoderTest(t, fmt.Sprintf("%T", simpleTestStruct{}), simpleTestStruct{
		One:   1,
		Two:   2.2,
		Three: "three",
	})
	makeDecoderTest(t, fmt.Sprintf("%T", testStruct{}), testStruct{
		SomeMap1: map[string]simpleTestStruct{
			"one": {
				One:   1,
				Two:   1.1,
				Three: "one",
			},
			"two": {
				One:   2,
				Two:   2.2,
				Three: "two",
			},
		},
		SomeStr: "some root string",
		SomeInt: 123123123123,
	})
}
