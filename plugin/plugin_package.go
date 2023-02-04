package plugin

import (
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

type PluginPackage struct {
	manifest     []byte
	pluginModule PluginModule
}

type UriPackage struct {
	uri *uri.URI
	pkg wasm.Package
}

//func NewPluginPackage(manifest, module []byte) *PluginPackage {
//	return &PluginPackage{manifest, module}
//}
//
//func (pkg *PluginPackage) Manifest(validation bool) (any, error) {
//}
//
//func (pkg *PluginPackage) CreateWrapper() (Wrapper, error) {
//}
