package wasm

import (
	"github.com/bytecodealliance/wasmtime-go"
)

func createImport(linker *wasmtime.Linker) {
	linker.FuncWrap("wrap", "__wrap_invoke_args", func(methodPtr uint32, argsPtr uint32) {

	})
}
