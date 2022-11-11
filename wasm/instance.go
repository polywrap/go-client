package wasm

import (
	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

type WasmInstance struct {
	instance *wasmtime.Instance
}

func NewInstance(wasm []byte) (*WasmInstance, error) {
	store := wasmtime.NewStore(wasmtime.NewEngine())
	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return nil, err
	}

	linker := wasmtime.NewLinker(store.Engine)
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, err
	}

	return &WasmInstance{
		instance: instance,
	}, nil
}

func (w *WasmInstance) WrapInvoke() {

}
