package wasm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
)

// execDecoderTest
func execDecoderTest[T any](t *testing.T, name string, expected T) {
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
	type (
		simpleStruct struct {
			One   int8
			Two   float32
			Three string
			Four  *string
			Five  *big.Int
		}
		complexStruct struct {
			SomeMap1 map[string]simpleStruct
			SomeStr  string
			SomeInt  int64
			SomeMap2 map[string]simpleStruct
			SomeMap3 map[string]*simpleStruct
		}
	)

	trueValue := true
	simpleStructValue := simpleStruct{
		One:   1,
		Two:   1.1,
		Three: "one",
	}

	execDecoderTest(t, "bool=true", trueValue)
	execDecoderTest(t, "*bool=true", &trueValue)
	execDecoderTest[*bool](t, "*bool=nil", nil)
	execDecoderTest(t, "int8", int8(1))
	execDecoderTest(t, "int16", int16(1))
	execDecoderTest(t, "int32", int32(1))
	execDecoderTest(t, "int64", int64(1))
	execDecoderTest(t, "uint8", uint8(1))
	execDecoderTest(t, "uint16", uint16(1))
	execDecoderTest(t, "uint32", uint32(1))
	execDecoderTest(t, "uint64", uint64(1))
	execDecoderTest(t, "float32", float32(1.12))
	execDecoderTest(t, "float64", float64(1.13))
	execDecoderTest(t, "string", "hello world")
	execDecoderTest(t, "[4]int8", [4]int8{1, 2, 3, 4})
	execDecoderTest(t, "[]int8", []int8{1, 2, 3, 4})
	execDecoderTest(t, "map[string]int8", map[string]int8{"one": 1, "two": 2, "three": 3})
	execDecoderTest(t, fmt.Sprintf("%T", simpleStruct{}), simpleStruct{
		One:   1,
		Two:   2.2,
		Three: "three",
		Five:  big.NewInt(1),
	})
	execDecoderTest(t, fmt.Sprintf("%T", &simpleStruct{}), &simpleStruct{
		One:   1,
		Two:   2.2,
		Three: "three",
		Five:  big.NewInt(1),
	})
	execDecoderTest[*simpleStruct](t, fmt.Sprintf("%T", &simpleStruct{}), nil)
	execDecoderTest(t, fmt.Sprintf("%T", complexStruct{}), complexStruct{
		SomeMap1: map[string]simpleStruct{
			"one": {
				One:   1,
				Two:   1.1,
				Three: "one",
				Five:  big.NewInt(1),
			},
			"two": {
				One:   2,
				Two:   2.2,
				Three: "two",
				Five:  big.NewInt(2),
			},
		},
		SomeStr: "some root string",
		SomeInt: 123123123123,
		SomeMap3: map[string]*simpleStruct{
			"one": &simpleStructValue,
			"two": nil,
		},
	})
}
