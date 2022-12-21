package wasm

import (
	"github.com/polywrap/go-client/core"
	"github.com/polywrap/go-client/core/resolver"
	"github.com/polywrap/go-client/core/resolver/uri"
)

type (
	ClientConfig[T resolver.SomeResolver] struct {
		Resolver   T
		Env        []byte
		Interfaces map[string][]uri.URI
	}

	Client struct {
		wrapper core.Wrapper
		invoker core.Invoker
	}
)

func NewClient[T resolver.SomeResolver](*ClientConfig[T]) *Client {
	return &Client{}
}

func (client *Client) Invoke(method string, data []byte) ([]byte, error) {
	return []byte{0xc}, nil
}

func Invoke[InvokeArg, InvokeRes any](client *Client, method string, arguments InvokeArg) (*InvokeRes, error) {
	args, err := Encode(arguments)
	if err != nil {
		return nil, err
	}

	resp, err := client.Invoke(method, args)
	if err != nil {
		return nil, err
	}

	res, err := Decode[InvokeRes](resp)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
