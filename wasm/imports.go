package wasm

import (
	"fmt"
	"unsafe"

	"github.com/bytecodealliance/wasmtime-go"
)

func createImport(linker *wasmtime.Linker, memory *wasmtime.Memory) {
	linker.FuncWrap("wrap", "__wrap_invoke_args", func(caller *wasmtime.Caller, methodPtr uint32, argsPtr uint32) {
		mem := memory.UnsafeData(caller)
		copy(mem[methodPtr:], *(*[]byte)(unsafe.Pointer(&methodPtr)))
		copy(mem[argsPtr:], *(*[]byte)(unsafe.Pointer(&argsPtr)))
	})
	linker.FuncWrap("wrap", "__wrap_invoke_result", func(caller *wasmtime.Caller, ptr uint32, len uint32) []byte {
		return memory.UnsafeData(caller)[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_invoke_error", func(caller *wasmtime.Caller, ptr uint32, len uint32) []byte {
		return memory.UnsafeData(caller)[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_abort", func(caller *wasmtime.Caller, msgPtr, msgLen, filePtr, fileLen, line, column uint32) {
		mem := memory.UnsafeData(caller)
		msg := string(mem[msgPtr : msgPtr+msgLen])
		file := string(mem[filePtr : filePtr+fileLen])
		panic(fmt.Sprintf("__wrap_abort: %s\nFile: %s\nLocation: [{%d},{%d}]", msg, file, line, column))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke", func(caller *wasmtime.Caller, uriPtr, uriLen, methodPtr, methodLen, argsPtr, argsLen uint32) {
		mem := memory.UnsafeData(caller)
		uri := string(mem[uriPtr : uriPtr+uriLen])
		method := string(mem[methodPtr : methodPtr+methodLen])
		args := mem[argsPtr : argsPtr+argsLen]
		panic(fmt.Sprintf(
			"Uri: %s\nMethod: %s\nArgs: %x\n  __wrap_subinvoke not implemented",
			uri,
			method,
			args,
		))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_result_len", func() {
		panic("__wrap_subinvoke_result_len not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_result", func(ptr uint32) {
		panic("__wrap_subinvoke_result not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error_len", func() {
		panic("__wrap_subinvoke_error_len not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error", func(ptr uint32) {
		panic("__wrap_subinvoke_result not implemented")
	})
}
