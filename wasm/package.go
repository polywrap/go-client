package wasm

import "github.com/polywrap/go-client/wasm/uri"

const (
	Base64FileEncoding FileEncoding = 1
	Utf8FileEncoding   FileEncoding = 2
)

type (
	FileEncoding uint8

	Invoker interface {
		Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
		InvokeWrapper(wrapper any, uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
		Implementations(uri uri.URI) ([]uri.URI, error)
		Interfaces() map[string][]uri.URI
	}

	WrapperResolver interface {
		TryResolveUri(uri uri.URI, loader any, context any) (Wrapper, error)
	}

	PackageResolver interface {
		TryResolveUri(uri uri.URI, loader any, context any) (Package, error)
	}

	Loader interface {
		LoadWrapper(uri uri.URI) (Wrapper, error)
	}

	Package interface {
		CreateWrapper() (Wrapper, error)
		Manifest(validation bool) (any, error)
	}

	Wrapper interface {
		Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
		File(path string, encoding *FileEncoding) ([]byte, error)
	}
)
