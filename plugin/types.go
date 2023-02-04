package plugin

import (
	"github.com/polywrap/go-client/wasm"
)

type PluginModule interface {
	WrapInvoke(method string, args []byte, invoker wasm.Invoker) ([]byte, error)
}
