package plugin

import (
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

type PluginWrapper struct {
	module PluginModule
}

func NewPluginWrapper(module PluginModule) *PluginWrapper {
	return &PluginWrapper{
		module: module,
	}
}

func (pw *PluginWrapper) Invoke(invoker wasm.Invoker, uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	result, err := pw.module.WrapInvoke(method, args, invoker)

	return result, err
}

func (pw *PluginWrapper) File(_ string, _ *wasm.FileEncoding) ([]byte, error) {
	return nil, nil
}
