package resolver

import (
	"github.com/polywrap/go-client/core"
	"github.com/polywrap/go-client/core/resolver/uri"
)

type PackageResolver interface {
	TryResolveUri(uri uri.URI, loader any, context any) (core.Package, error)
}

type WrapperResolver interface {
	TryResolveUri(uri uri.URI, loader any, context any) (core.Wrapper, error)
}
