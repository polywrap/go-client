package core

import "github.com/polywrap/go-client/core/resolver/uri"

type FileEncoding uint8

const (
	Base64FileEncoding FileEncoding = 1
	Utf8FileEncoding   FileEncoding = 2
)

type Wrapper interface {
	Invoker(invoker Invoker, uri uri.URI, method string, args []byte, env []byte) ([]byte, error)
	File(path string, encoding *FileEncoding) ([]byte, error)
}
