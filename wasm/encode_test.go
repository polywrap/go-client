package wasm

import (
	"testing"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		name     string
		value    any
		expValue []byte
		expError string
	}{
		{
			name:  "int8",
			value: int8(1),
		},
		{
			name:  "int16",
			value: int16(1),
		},
		{
			name:  "int32",
			value: int32(1),
		},
		{
			name:  "int64",
			value: int64(1),
		},
		{
			name:  "uint8",
			value: uint8(1),
		},
		{
			name:  "uint16",
			value: uint16(1),
		},
		{
			name:  "uint32",
			value: uint32(1),
		},
		{
			name:  "uint64",
			value: uint64(1),
		},
		{
			name:  "string",
			value: "some value",
		},
	}
	for i := range cases {
		tcase := cases[i]
		t.Run(tcase.name, func(t *testing.T) {
			data, err := Encode(tcase.value)
			if len(tcase.expError) == 0 && len(data) == 0 {
				t.Error("data is empty")
			}
			if err != nil {
				t.Errorf("Error: %s", err)
			}
		})
	}
}
