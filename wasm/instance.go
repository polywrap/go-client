package wasm

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
)

var ErrNowWasmMemory = errors.New(strings.Join(
	[]string{
		`Unable to find Wasm memory import section.`,
		`  Modules must import memory from the "env" module's`,
		`  "memory" field like so:`,
		`    (import "env" "memory" (memory (;0;) #))`,
	},
	"\n",
))

type WasmInstance struct {
	store    *wasmtime.Store
	memory   *wasmtime.Memory
	linker   *wasmtime.Linker
	instance *wasmtime.Instance

	state *State
	mu    *sync.Mutex
}

func NewInstance(wasm []byte) (*WasmInstance, error) {
	store := wasmtime.NewStore(wasmtime.NewEngine())
	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return nil, err
	}
	memory, err := createMemory(wasm, store)
	if err != nil {
		return nil, err
	}

	wasmInstance := WasmInstance{
		store:  store,
		memory: memory,
		linker: wasmtime.NewLinker(store.Engine),

		state: nil,
		mu:    &sync.Mutex{},
	}
	createImport(&wasmInstance)

	instance, err := wasmInstance.linker.Instantiate(store, module)
	if err != nil {
		return nil, err
	}
	wasmInstance.instance = instance
	return &wasmInstance, nil
}

func (w *WasmInstance) WrapInvoke(method string, args []byte, env []byte) (*State, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.state = &State{
		Method: method,
		Args:   args,
		Env:    env,
	}
	invoke := w.instance.GetExport(w.store, "_wrap_invoke")
	_, err := invoke.Func().Call(
		w.store,
		len(method),
		len(args),
		len(env),
	)

	return w.state, err
}

func createMemory(wasm []byte, store *wasmtime.Store) (*wasmtime.Memory, error) {
	ENV_MEMORY_IMPORTS_SIGNATURE := []byte{0x65, 0x6e, 0x76, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x02}
	sigIdx := bytes.Index(wasm, ENV_MEMORY_IMPORTS_SIGNATURE)
	if sigIdx < 0 {
		return nil, ErrNowWasmMemory
	}
	memoryInitialLimits := wasm[sigIdx+1+len(ENV_MEMORY_IMPORTS_SIGNATURE)+1]
	memoryType := wasmtime.NewMemoryType(uint32(memoryInitialLimits), false, 0)
	return wasmtime.NewMemory(store, memoryType)
}

func createImport(inst *WasmInstance) {
	inst.linker.FuncWrap("wrap", "__wrap_load_env", func(ptr int32) {
		fmt.Printf("%s\n", inst.memory.UnsafeData(inst.store))
		panic("__wrap_load_env not implemented")
	})
	inst.linker.FuncWrap("wrap", "__wrap_invoke_args", func(methodPtr, argsPtr int32) {
		mem := inst.memory.UnsafeData(inst.store)
		copy(mem[methodPtr:], []byte(inst.state.Method))
		copy(mem[argsPtr:], inst.state.Args)
	})
	inst.linker.FuncWrap("wrap", "__wrap_invoke_result", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.state.Result = mem[ptr : ptr+len]
	})
	inst.linker.FuncWrap("wrap", "__wrap_invoke_error", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.state.Error = mem[ptr : ptr+len]
	})
	inst.linker.FuncWrap("wrap", "__wrap_abort", func(msgPtr, msgLen, filePtr, fileLen, line, column int32) {
		mem := inst.memory.UnsafeData(inst.store)
		msg := string(mem[msgPtr : msgPtr+msgLen])
		file := string(mem[filePtr : filePtr+fileLen])
		inst.state.Error = []byte(fmt.Sprintf("__wrap_abort: %s\nFile: %s\nLocation: [{%d},{%d}]", msg, file, line, column))
	})
	inst.linker.FuncWrap("wrap", "__wrap_subinvoke", func(uriPtr, uriLen, methodPtr, methodLen, argsPtr, argsLen int32) int32 {
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
	inst.linker.FuncWrap("wrap", "__wrap_subinvoke_result_len", func() int32 {
		panic("__wrap_subinvoke_result_len not implemented")
	})
	inst.linker.FuncWrap("wrap", "__wrap_subinvoke_result", func(ptr int32) {
		panic("__wrap_subinvoke_result not implemented")
	})
	inst.linker.FuncWrap("wrap", "__wrap_subinvoke_error_len", func() int32 {
		panic("__wrap_subinvoke_error_len not implemented")
	})
	inst.linker.FuncWrap("wrap", "__wrap_subinvoke_error", func(ptr int32) {
		panic("__wrap_subinvoke_result not implemented")
	})
	inst.linker.Define("env", "memory", inst.memory)
}
