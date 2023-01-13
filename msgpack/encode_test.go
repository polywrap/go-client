package msgpack

import (
	"testing"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
)

func TestEncode(t *testing.T) {
	intdata := int8(1)
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
		{
			name:  "*big.Int",
			value: big.NewInt(1),
		},
		{
			name: "struct",
			value: &(struct {
				One   int8  `tag:"one"`
				Two   int8  `tag:"two"`
				Three *int8 `tag:"three"`
				Map   map[string]int8
			}{
				One:   1,
				Two:   2,
				Three: &intdata,
				Map: map[string]int8{
					"one":   int8(1),
					"two":   int8(2),
					"three": int8(3),
				},
			}),
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
