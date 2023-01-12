package wasm

import (
	"fmt"

	"github.com/polywrap/go-client/wasm/instance"
	"github.com/polywrap/go-client/wasm/uri"
)

type WasmWrapper struct {
	manifest []byte
	module   []byte
}

func NewWasmWrapper(manifest, module []byte) *WasmWrapper {
	return &WasmWrapper{manifest, module}
}

func (wrp *WasmPackage) Invoke(uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	inst, err := instance.New(wrp.module)
	if err != nil {
		return nil, err
	}
	state, err := inst.WrapInvoke(method, args, env)
	if err != nil {
		return nil, err
	}
	if state.Error != nil {
		return nil, fmt.Errorf("wasm error: %x", state.Error)
	}
	return state.Result, nil
}

func (wrp *WasmPackage) File(path string, encoding *FileEncoding) ([]byte, error) {
	return nil, nil
}
