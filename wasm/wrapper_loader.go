package wasm

import (
	"context"
	"errors"

	"github.com/polywrap/go-client/wasm/uri"
)

var ErrUnknownResolver = errors.New("resolver return unknown type")

type WrapperLoader struct {
	resolver Resolver
}

func NewWrapperLoader(resolver Resolver, env []byte, ifaces map[string][]uri.URI) *WrapperLoader {
	return &WrapperLoader{resolver}
}

func (wl *WrapperLoader) LoadWrapper(uri uri.URI) (Wrapper, error) {
	item, err := wl.resolver.TryResolveUri(uri, wl, context.Background())
	if err != nil {
		return nil, err
	}

	if pkg, ok := item.(Package); ok {
		return pkg.CreateWrapper()
	}

	if wrp, ok := item.(Wrapper); ok {
		return wrp, nil
	}

	return nil, ErrUnknownResolver
}
