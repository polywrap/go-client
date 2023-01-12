package wasm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/polywrap/go-client/wasm/uri"
)

type FsResolver struct {
}

func NewFsResolver() *FsResolver {
	return new(FsResolver)
}

func (r *FsResolver) TryResolveUri(uri uri.URI, loader Loader, context context.Context) (any, error) {
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
