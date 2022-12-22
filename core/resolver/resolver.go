package resolver

import (
	"github.com/polywrap/go-client/core"
	"github.com/polywrap/go-client/core/resolver/uri"
)

type (
	Resolver interface {
		TryResolveUri(uri uri.URI, loader any, context any) (core.Package, error)
	}
	SomeResolver interface {
		PackageResolver | WrapperResolver
	}
	PackageResolver struct {
	}
	WrapperResolver struct {
	}
)
