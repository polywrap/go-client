package wasm

import (
	"fmt"
	"reflect"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack"
	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
)

func Encode(value any) ([]byte, error) {
	context := msgpack.NewContext(fmt.Sprintf("encode value: %T", value))
	encoder := msgpack.NewWriteEncoder(context)
	queue := []reflect.Value{reflect.ValueOf(value)}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		switch v.Kind() {
		case reflect.Bool:
			encoder.WriteBool(v.Bool())
		case reflect.Int8:
			encoder.WriteI8(int8(v.Int()))
		case reflect.Int16:
			encoder.WriteI16(int16(v.Int()))
		case reflect.Int32:
			encoder.WriteI32(int32(v.Int()))
		case reflect.Int64:
			encoder.WriteI64(int64(v.Int()))
		case reflect.Uint8:
			encoder.WriteU8(uint8(v.Uint()))
		case reflect.Uint16:
			encoder.WriteU16(uint16(v.Uint()))
		case reflect.Uint32:
			encoder.WriteU32(uint32(v.Uint()))
		case reflect.Uint64:
			encoder.WriteU64(uint64(v.Uint()))
		case reflect.Float32:
			encoder.WriteFloat32(float32(v.Float()))
		case reflect.Float64:
			encoder.WriteFloat64(float64(v.Float()))
		case reflect.String:
			encoder.WriteString(v.String())
		case reflect.Slice, reflect.Array:
			encoder.WriteArrayLength(uint32(v.Len()))
			for i := v.Len() - 1; i >= 0; i-- {
				queue = append([]reflect.Value{v.Index(i)}, queue...)
			}
		case reflect.Map:
			encoder.WriteMapLength(uint32(v.Len()))
			for i, keys := v.Len()-1, v.MapKeys(); i >= 0; i-- {
				queue = append([]reflect.Value{keys[i], v.MapIndex(keys[i])}, queue...)
			}
		case reflect.Struct:
			t := v.Type()
			if t.Name() == "Int" {
				v, ok := v.Interface().(big.Int)
				if !ok {
					return nil, fmt.Errorf("unknown type: %s", t)
				}
				encoder.WriteBigInt(&v)
			} else {
				encoder.WriteMapLength(uint32(v.NumField()))
				for i := v.NumField() - 1; i >= 0; i-- {
					queue = append([]reflect.Value{
						reflect.ValueOf(t.Field(i).Name),
						v.Field(i),
					}, queue...)
				}
			}
		case reflect.Pointer:
			queue = append([]reflect.Value{reflect.Indirect(v)}, queue...)
		default:
			return nil, fmt.Errorf("unknown type: %s", v.Type())
		}
	}
	return encoder.Buffer(), nil
}
