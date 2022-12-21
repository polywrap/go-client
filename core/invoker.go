package core

import "github.com/polywrap/go-client/core/resolver/uri"

type (
	Invoker interface {
		Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
	}

	WrapperInvoker interface {
		InvokeWrapper(wrapper any, uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
	}

	ImplementationInvoker interface {
		Implementations(uri uri.URI) ([]uri.URI, error)
	}

	InterfaceInvoker interface {
		Interfaces() map[string][]uri.URI
	}
)
