package plugin

import (
	"errors"
	"fmt"
	"github.com/polywrap/go-client/msgpack"
	"github.com/polywrap/go-client/wasm"
	"reflect"
)

type PluginModule interface {
	WrapInvoke(method string, args []byte, invoker wasm.Invoker) ([]byte, error)
}

type Plugin interface {
	EncodeArgs(method string, args []byte) (any, error)
}

type pluginModule struct {
	p Plugin
}

func NewPluginModule(p Plugin) *pluginModule {
	return &pluginModule{p}
}

func (pm *pluginModule) WrapInvoke(method string, args []byte, invoker wasm.Invoker) ([]byte, error) {
	res, err := pm.call(method, args)
	if err != nil {
		return nil, err
	}
	resEncoded, err := msgpack.Encode(res)
	if err != nil {
		return nil, err
	}

	return resEncoded, nil
}

func (pm *pluginModule) call(method string, args []byte) (any, error) {
	_, ok := reflect.TypeOf(pm.p).MethodByName(method)
	if !ok {
		return nil, fmt.Errorf("method %s not found in plugin", method)
	}
	argsDecoded, err := pm.p.EncodeArgs(method, args)
	if err != nil {
		return nil, err
	}

	inputs := []reflect.Value{reflect.ValueOf(argsDecoded)}
	resValue := reflect.ValueOf(pm.p).MethodByName(method).Call(inputs)

	if len(resValue) == 0 {
		return nil, errors.New("plugin didn't return value")
	}
	v := resValue[0]

	return pm.Decode(v)
}

func (pm *pluginModule) Decode(v reflect.Value) (any, error) {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.String:
		return v.String(), nil
	default:
		return nil, fmt.Errorf("unknown type: %s", v.Kind())
	}
}
