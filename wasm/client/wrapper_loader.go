package client

import (
	"context"
	"errors"

	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

var ErrUnknownResolver = errors.New("resolver return unknown type")

type WrapperLoader struct {
	resolver wasm.Resolver
}

func NewWrapperLoader(resolver wasm.Resolver, env []byte, ifaces map[string][]uri.URI) *WrapperLoader {
	return &WrapperLoader{resolver}
}

func (wl *WrapperLoader) LoadWrapper(uri uri.URI) (wasm.Wrapper, error) {
	item, err := wl.resolver.TryResolveUri(uri, wl, context.Background())
	if err != nil {
		return nil, err
	}

	if pkg, ok := item.(wasm.Package); ok {
		return pkg.CreateWrapper()
	}

	if wrp, ok := item.(wasm.Wrapper); ok {
		return wrp, nil
	}

	return nil, ErrUnknownResolver
}
