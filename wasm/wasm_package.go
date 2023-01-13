package wasm

type WasmPackage struct {
	manifest []byte
	module   []byte
}

func NewWasmPackage(manifest, module []byte) *WasmPackage {
	return &WasmPackage{manifest, module}
}

func (pkg *WasmPackage) Manifest(validation bool) (any, error) {
	return pkg.manifest, nil
}

func (pkg *WasmPackage) CreateWrapper() (Wrapper, error) {
	return NewWasmWrapper(pkg.manifest, pkg.module), nil
}
