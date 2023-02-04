package plugin

import (
	"errors"
	"fmt"
	"github.com/polywrap/go-client/msgpack"
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/client"
	"github.com/polywrap/go-client/wasm/uri"
	"log"
	"reflect"
	"testing"
)

type ArgType struct {
	Key string
}

type DemoPlugin struct {
}

func (d *DemoPlugin) WrapInvoke(method string, args []byte, invoker wasm.Invoker) ([]byte, error) {
	res, err := d.callMethod(method, args)
	if err != nil {
		return nil, err
	}
	resEncoded, err := msgpack.Encode(res)
	if err != nil {
		return nil, err
	}

	return resEncoded, nil
}

func (d *DemoPlugin) callMethod(method string, args []byte) (any, error) {
	_, ok := reflect.TypeOf(d).MethodByName(method)
	if !ok {
		return nil, fmt.Errorf("method %s not found in plugin", method)
	}
	argsDecoded, err := msgpack.Decode[ArgType](args)
	if err != nil {
		return nil, err
	}

	inputs := []reflect.Value{reflect.ValueOf(argsDecoded)}
	resValue := reflect.ValueOf(d).MethodByName(method).Call(inputs)

	if len(resValue) == 0 {
		return nil, errors.New("plugin didn't return value")
	}

	v := resValue[0]

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

func (d *DemoPlugin) CheckArgIsBar(arg ArgType) bool {
	if arg.Key == "bar" {
		return true
	}
	return false
}

func (d *DemoPlugin) CreateWrapper() (wasm.Wrapper, error) {
	return NewPluginWrapper(d), nil
}

func (d *DemoPlugin) Manifest(_ bool) (any, error) {
	return nil, nil
}

func TestPlugin(t *testing.T) {
	//u, _ := uri.New("ens/demo-plugin.eth")
	//uriPackage := UriPackage{
	//	uri: u,
	//	pkg: &DemoPlugin{},
	//}

	wrapUri, err := uri.New("wrap://ens/demo-plugin.eth")
	if err != nil {
		log.Fatalf("bad wrapUri: %s (%s)", "ens/demo-plugin.eth", err)
	}

	resolver := wasm.NewStaticResolver(map[string]wasm.Package{
		"wrap://ens/demo-plugin.eth": &DemoPlugin{},
	})

	polywrapClient := client.New(&client.ClientConfig{
		Resolver: resolver,
	})
	args := ArgType{
		Key: "bar",
	}
	res, err := client.Invoke[ArgType, bool, []byte](polywrapClient, *wrapUri, "CheckArgIsBar", args, nil)
	if err != nil {
		log.Fatalf("invokation error: %+v", err)
	}

	log.Printf("Result is: %t\n", *res)
}
