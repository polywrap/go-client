package plugin

import (
	"fmt"
	"github.com/polywrap/go-client/msgpack"
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/client"
	"github.com/polywrap/go-client/wasm/uri"
	"log"
	"testing"
)

type ArgType struct {
	Key string
}

type DemoPlugin struct {
}

func (d *DemoPlugin) CheckArgIsBar(arg ArgType) bool {
	if arg.Key == "bar" {
		return true
	}
	return false
}

func (d *DemoPlugin) EncodeArgs(method string, args []byte) (any, error) {
	switch method {
	case "CheckArgIsBar":
		return msgpack.Decode[ArgType](args)
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

func TestPlugin(t *testing.T) {
	pluginPackage := NewPluginPackage(nil, NewPluginModule(&DemoPlugin{}))

	wrapUri, err := uri.New("wrap://ens/demo-plugin.eth")
	if err != nil {
		log.Fatalf("bad wrapUri: %s (%s)", "ens/demo-plugin.eth", err)
	}

	resolver := wasm.NewStaticResolver(map[string]wasm.Package{
		"wrap://ens/demo-plugin.eth": pluginPackage,
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

	if *res != true {
		t.Errorf("Actual: %#v, Expected: %#v", *res, true)
	}
}
