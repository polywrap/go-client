package client

import (
	"github.com/polywrap/go-client/msgpack"
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

type (
	ClientConfig struct {
		Resolver   wasm.Resolver
		Env        []byte
		Interfaces map[string][]uri.URI
	}

	Client struct {
		Loader  *WrapperLoader
		invoker *WrapperInvoker
	}
)

func New(cfg *ClientConfig) *Client {
	loader := NewWrapperLoader(cfg.Resolver, cfg.Env, cfg.Interfaces)
	invoker := NewWrapperInvoker(loader)
	return &Client{loader, invoker}
}

func (client *Client) Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	return client.invoker.Invoke(uri, method, args, env)
}

func Invoke[InvokeArg, InvokeRes any](client *Client, uri uri.URI, method string, arguments InvokeArg) (*InvokeRes, error) {
	args, err := msgpack.Encode(arguments)
	if err != nil {
		return nil, err
	}

	resp, err := client.Invoke(uri, method, args, nil)
	if err != nil {
		return nil, err
	}

	res, err := msgpack.Decode[InvokeRes](resp)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
