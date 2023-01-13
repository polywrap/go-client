package client

import (
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

type WrapperInvoker struct {
	loader *WrapperLoader
}

func NewWrapperInvoker(loader *WrapperLoader) *WrapperInvoker {
	return &WrapperInvoker{loader}
}

func (wi *WrapperInvoker) Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	wrapper, err := wi.loader.LoadWrapper(uri)
	if err != nil {
		return nil, err
	}
	return wi.InvokeWrapper(wrapper, uri, method, args, env)
}

func (wi *WrapperInvoker) InvokeWrapper(wrapper wasm.Wrapper, uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	return wrapper.Invoke(wi, uri, method, args, env)
}

func (wi *WrapperInvoker) Implementations(uri uri.URI) ([]uri.URI, error) {
	return nil, nil
}

func (wi *WrapperInvoker) Interfaces() map[string][]uri.URI {
	return nil
}
