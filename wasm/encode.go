package wasm

import (
	"fmt"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack"
	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
)

func Encode(value any) ([]byte, error) {
	context := msgpack.NewContext(fmt.Sprintf("encode value: %T", value))
	encoder := msgpack.NewWriteEncoder(context)
	queue := []any{value}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		switch v := v.(type) {
		case int8:
			encoder.WriteI8(v)
		case int16:
			encoder.WriteI16(v)
		case int32:
			encoder.WriteI32(v)
		case int64:
			encoder.WriteI64(v)
		case uint8:
			encoder.WriteU8(v)
		case uint16:
			encoder.WriteU16(v)
		case uint32:
			encoder.WriteU32(v)
		case uint64:
			encoder.WriteU64(v)
		case string:
			encoder.WriteString(v)
		case *big.Int:
			encoder.WriteBigInt(v)
		default:
			return nil, fmt.Errorf("unkown type: %T", v)
		}
	}
	return encoder.Buffer(), nil
}
