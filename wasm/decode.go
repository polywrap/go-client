package wasm

import (
	"fmt"
	"reflect"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack"
)

func decode(decoder msgpack.Read, v *reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(decoder.ReadBool())
	case reflect.Int8:
		v.SetInt(int64(decoder.ReadI8()))
	case reflect.Int16:
		v.SetInt(int64(decoder.ReadI16()))
	case reflect.Int32:
		v.SetInt(int64(decoder.ReadI32()))
	case reflect.Int64:
		v.SetInt(int64(decoder.ReadI64()))
	case reflect.Uint8:
		v.SetUint(uint64(decoder.ReadU8()))
	case reflect.Uint16:
		v.SetUint(uint64(decoder.ReadU16()))
	case reflect.Uint32:
		v.SetUint(uint64(decoder.ReadU32()))
	case reflect.Uint64:
		v.SetUint(uint64(decoder.ReadU64()))
	case reflect.Float32:
		v.SetFloat(float64(decoder.ReadF32()))
	case reflect.Float64:
		v.SetFloat(float64(decoder.ReadF64()))
	case reflect.String:
		v.SetString(decoder.ReadString())
	case reflect.Array, reflect.Slice:
		aLn := int(decoder.ReadArrayLength())
		if v.Kind() == reflect.Slice {
			v.Set(reflect.MakeSlice(v.Type(), aLn, aLn))
		}
		for i := 0; i < aLn; i++ {
			t := v.Index(i)
			decode(decoder, &t)
		}
	case reflect.Map:
		mLn := int(decoder.ReadMapLength())
		v.Set(reflect.MakeMap(v.Type()))
		for i := 0; i < mLn; i++ {
			key := reflect.Indirect(reflect.New(v.Type().Key()))
			value := reflect.Indirect(reflect.New(v.Type().Elem()))
			decode(decoder, &key)
			decode(decoder, &value)
			v.SetMapIndex(key, value)
		}
	default:
		return fmt.Errorf("unknown type: %s", v.Type())
	}
	return nil
}

func Decode[T any](data []byte) (T, error) {
	var value *T = new(T)
	context := msgpack.NewContext(fmt.Sprintf("decode value: %T", value))
	decoder := msgpack.NewReadDecoder(context, data)
	queue := []reflect.Value{reflect.Indirect(reflect.ValueOf(value))}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		decode(decoder, &v)
	}
	return *value, nil
}
