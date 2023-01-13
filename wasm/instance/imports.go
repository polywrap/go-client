package instance

import (
	"bytes"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/polywrap/go-client/wasm/uri"
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
		panic("__wrap_load_env not implemented")
	})
	linker.FuncWrap("wrap", "__wrap_invoke_args", func(methodPtr, argsPtr int32) {
		mem := inst.memory.UnsafeData(inst.store)
		copy(mem[methodPtr:], []byte(inst.State.Method))
		copy(mem[argsPtr:], inst.State.Args)
	})
	linker.FuncWrap("wrap", "__wrap_invoke_result", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.State.Invoke.Result = mem[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_invoke_error", func(ptr, len int32) {
		mem := inst.memory.UnsafeData(inst.store)
		inst.State.Invoke.Error = mem[ptr : ptr+len]
	})
	linker.FuncWrap("wrap", "__wrap_abort", func(msgPtr, msgLen, filePtr, fileLen, line, column int32) {
		mem := inst.memory.UnsafeData(inst.store)
		msg := string(mem[msgPtr : msgPtr+msgLen])
		file := string(mem[filePtr : filePtr+fileLen])
		inst.State.Invoke.Error = []byte(fmt.Sprintf("__wrap_abort: %s\nFile: %s\nLocation: [{%d},{%d}]", msg, file, line, column))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke", func(uriPtr, uriLen, methodPtr, methodLen, argsPtr, argsLen int32) int32 {
		mem := inst.memory.UnsafeData(inst.store)
		u, err := uri.New(string(mem[uriPtr : uriPtr+uriLen]))
		if err != nil {
			panic(fmt.Sprintf("can't parse uri: %s", mem[uriPtr:uriPtr+uriLen]))
		}
		res, err := inst.State.Invoker.Invoke(*u, string(mem[methodPtr:methodPtr+methodLen]), mem[argsPtr:argsPtr+argsLen], []byte{})
		inst.State.Subinvoke.Result = res
		if err != nil {
			inst.State.Subinvoke.Error = []byte(err.Error())
		}
		return 1
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_result_len", func() int32 {
		return int32(len(inst.State.Subinvoke.Result))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_result", func(ptr int32) {
		mem := inst.memory.UnsafeData(inst.store)
		copy(mem[ptr:], inst.State.Subinvoke.Result)
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error_len", func() int32 {
		return int32(len(inst.State.Subinvoke.Error))
	})
	linker.FuncWrap("wrap", "__wrap_subinvoke_error", func(ptr int32) {
		mem := inst.memory.UnsafeData(inst.store)
		copy(mem[ptr:], inst.State.Subinvoke.Error)
	})
	linker.Define("env", "memory", inst.memory)
}

/*
 move |mut caller: Caller<'_, State>, ptr: u32| {
                let memory = memory.lock().unwrap();
                let (memory_buffer, state) = memory.data_and_store_mut(caller.as_context_mut());

                match &state.subinvoke.result {
                    Some(res) => {
                        write_to_memory(memory_buffer, ptr as usize, res);
                    }
                    None => {
                        (state.abort)(
                            "__wrap_subinvoke_result: subinvoke.result is not set".to_string(),
                        );
                    }
                }
            },
*/
