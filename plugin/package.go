package plugin

import "github.com/polywrap/go-client/wasm"

type PluginPackage struct {
	manifest     []byte
	pluginModule PluginModule
}

func NewPluginPackage(manifest []byte, module PluginModule) *PluginPackage {
	return &PluginPackage{manifest, module}
}

func (pkg *PluginPackage) Manifest(_ bool) (any, error) {
	return nil, nil
}

func (pkg *PluginPackage) CreateWrapper() (wasm.Wrapper, error) {
	return NewPluginWrapper(pkg.pluginModule), nil
}
