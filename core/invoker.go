package core

import "github.com/polywrap/go-client/core/resolver/uri"

type Invoker interface {
	Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
	InvokeWrapper(wrapper any, uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
	Implementations(uri uri.URI) ([]uri.URI, error)
	Interfaces() map[string][]uri.URI
}
