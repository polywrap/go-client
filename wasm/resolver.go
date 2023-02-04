package wasm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	pkguri "github.com/polywrap/go-client/wasm/uri"
)

type (
	FsResolver struct {
	}
	RedirectResolver struct {
		UriMap map[string]*pkguri.URI
	}
	BaseResolver struct {
		RdResolver Resolver
		FsResolver Resolver
	}

	StaticResolver struct {
		UriMap map[string]Package
	}
)

func NewStaticResolver(sm map[string]Package) *StaticResolver {
	return &StaticResolver{sm}
}

func (r *StaticResolver) TryResolveUri(uri pkguri.URI, _ Loader, _ context.Context) (any, error) {
	if u, ok := r.UriMap[uri.Uri]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func NewFsResolver() *FsResolver {
	return new(FsResolver)
}

func (r *FsResolver) TryResolveUri(uri pkguri.URI, loader Loader, context context.Context) (any, error) {
	if uri.Authority != "fs" && uri.Authority != "file" {
		return nil, fmt.Errorf("invalid authority: %s", uri.Authority)
	}

	manifest, err := os.ReadFile(path.Join(uri.Path, "wrap.info"))
	if err != nil {
		return nil, err
	}

	wrapper, err := os.ReadFile(path.Join(uri.Path, "wrap.wasm"))
	if err != nil {
		return nil, err
	}

	return NewWasmPackage(manifest, wrapper), nil
}

func NewRedirectResolver(rd map[string]*pkguri.URI) *RedirectResolver {
	return &RedirectResolver{rd}
}

func (r *RedirectResolver) TryResolveUri(uri pkguri.URI, _ Loader, _ context.Context) (any, error) {
	if u, ok := r.UriMap[uri.Uri]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func NewBaseResolver(rdResolver, fsResolver Resolver) *BaseResolver {
	return &BaseResolver{
		RdResolver: rdResolver,
		FsResolver: fsResolver,
	}
}

func (r *BaseResolver) TryResolveUri(uri pkguri.URI, loader Loader, context context.Context) (any, error) {
	if data, err := r.RdResolver.TryResolveUri(uri, loader, context); err == nil {
		if u, ok := data.(*pkguri.URI); ok {
			return r.FsResolver.TryResolveUri(*u, loader, context)
		}
	}
	return r.FsResolver.TryResolveUri(uri, loader, context)
}
