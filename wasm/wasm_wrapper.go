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

func (wrp *WasmWrapper) Invoke(invoker Invoker, uri uri.URI, method string, args []byte, env []byte) ([]byte, error) {
	inst, err := instance.New(wrp.module, instance.NewState(
		invoker,
		[]byte(method),
		args,
		env,
	))
	if err != nil {
		return nil, err
	}
	state, err := inst.Call()
	if err != nil {
		return nil, err
	}
	if state.Invoke.Error != nil {
		return nil, fmt.Errorf("wasm error: %x", state.Invoke.Error)
	}
	return state.Invoke.Result, nil
}

func (wrp *WasmWrapper) File(path string, encoding *FileEncoding) ([]byte, error) {
	return nil, nil
}
