package instance

import (
	"bytes"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

func createMemory(wasm []byte, store *wasmtime.Store) (*wasmtime.Memory, error) {
	sigIdx := bytes.Index(wasm, ENV_MEMORY_IMPORTS_SIGNATURE)
	if sigIdx < 0 {
		return nil, ErrNowWasmMemory
	}
	memoryInitialLimits := wasm[sigIdx+1+len(ENV_MEMORY_IMPORTS_SIGNATURE)+1]
	memoryType := wasmtime.NewMemoryType(uint32(memoryInitialLimits), false, 0)
	return wasmtime.NewMemory(store, memoryType)
}

func createImport(linker *wasmtime.Linker, inst *Instance) {
	linker.FuncWrap("wrap", "__wrap_load_env", func(ptr int32) {
		fmt.Printf("%s\n", inst.memory.UnsafeData(inst.store))
		panic("__wrap_load_env not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_invoke_args", func(methodPtr, argsPtr int32) {
		mem := inst.memory.UnsafeData(inst.store)
		copy(mem[methodPtr:], []byte(inst.State.Method))
		copy(mem[argsPtr:], inst.State.Args)
	})
	linker.FuncWrap("wrap", "__wrap_invoke_result", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.State.Result = mem[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_invoke_error", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.State.Error = mem[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_abort", func(msgPtr, msgLen, filePtr, fileLen, line, column int32) {
		mem := inst.memory.UnsafeData(inst.store)
		msg := string(mem[msgPtr : msgPtr+msgLen])
		file := string(mem[filePtr : filePtr+fileLen])
		inst.State.Error = []byte(fmt.Sprintf("__wrap_abort: %s\nFile: %s\nLocation: [{%d},{%d}]", msg, file, line, column))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke", func(uriPtr, uriLen, methodPtr, methodLen, argsPtr, argsLen int32) int32 {
		mem := inst.memory.UnsafeData(inst.store)
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
	linker.FuncWrap("wrap", "__wrap_subinvoke_result_len", func() int32 {
		panic("__wrap_subinvoke_result_len not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_result", func(ptr int32) {
		panic("__wrap_subinvoke_result not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error_len", func() int32 {
		panic("__wrap_subinvoke_error_len not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error", func(ptr int32) {
		panic("__wrap_subinvoke_result not implemented")
	})
	linker.Define("env", "memory", inst.memory)
}
