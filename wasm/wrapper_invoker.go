package wasm

import "github.com/polywrap/go-client/wasm/uri"

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
	return wrapper.Invoke(uri, method, args, env)
}
