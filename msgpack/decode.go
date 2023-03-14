package msgpack

import (
	"fmt"
	"reflect"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack"
)

func Decode[T any](data []byte) (T, error) {
	var value *T = new(T)
	context := msgpack.NewContext(fmt.Sprintf("decode value: %T", value))
	decoder := msgpack.NewReadDecoder(context, data)
	queue := []reflect.Value{reflect.Indirect(reflect.ValueOf(value))}
	linker := make([][3]reflect.Value, 0)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
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
			// handle []byte
			if v.Type().String() == "[]uint8" {
				v.SetBytes(decoder.ReadBytes())
			} else {
				aLn := int(decoder.ReadArrayLength())
				if v.Kind() == reflect.Slice {
					v.Set(reflect.MakeSlice(v.Type(), aLn, aLn))
				}
				for i := 0; i < aLn; i++ {
					queue = append([]reflect.Value{v.Index(i)}, queue...)
				}
			}
		case reflect.Map:
			mLn := int(decoder.ReadMapLength())
			if mLn > 0 {
				v.Set(reflect.MakeMap(v.Type()))
			}
			for i := 0; i < mLn; i++ {
				key := reflect.Indirect(reflect.New(v.Type().Key()))
				value := reflect.Indirect(reflect.New(v.Type().Elem()))
				queue = append([]reflect.Value{key, value}, queue...)
				linker = append([][3]reflect.Value{{v, key, value}}, linker...)
			}
		case reflect.Struct:
			t := v.Type()
			if t.Name() == "Int" {
				v.Set(reflect.Indirect(reflect.ValueOf(decoder.ReadBigInt())))
			} else {
				sLn := int(decoder.ReadMapLength())
				if sLn != v.NumField() {
					return *value, fmt.Errorf("different number of fields")
				}
				for i := sLn - 1; i >= 0; i-- {
					key := reflect.Indirect(reflect.New(reflect.TypeOf("")))
					value := reflect.Indirect(reflect.New(t.Field(i).Type))
					queue = append([]reflect.Value{key, value}, queue...)
					linker = append([][3]reflect.Value{{v, key, value}}, linker...)
				}
			}
		case reflect.Pointer:
			if decoder.IsNil() {
				decoder.ReadBytesLength()
			} else {
				v.Set(reflect.Indirect(reflect.New(v.Type().Elem())).Addr())
				queue = append([]reflect.Value{reflect.Indirect(v)}, queue...)
			}
		default:
			return *value, fmt.Errorf("unknown type: %s", v.Type())
		}
	}
	for i := range linker {
		switch linker[i][0].Kind() {
		case reflect.Map:
			linker[i][0].SetMapIndex(linker[i][1], linker[i][2])
		case reflect.Struct:
			linker[i][0].FieldByName(Capitalize(linker[i][1].String())).Set(linker[i][2])
		}
	}
	return *value, nil
}
