package wasm

import "github.com/polywrap/go-client/wasm/uri"

type WrapperLoader struct {
	resolver PackageResolver
}

func NewWrapperLoader(resolver PackageResolver, env []byte, ifaces map[string][]uri.URI) *WrapperLoader {
	return &WrapperLoader{resolver}
}

func (wl *WrapperLoader) LoadWrapper(uri uri.URI) (Wrapper, error) {
	pkg, err := wl.resolver.TryResolveUri(uri, nil, nil)
	if err != nil {
		return nil, err
	}
	return pkg.CreateWrapper()
}
