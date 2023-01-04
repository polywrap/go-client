package instance

import (
	"github.com/bytecodealliance/wasmtime-go"
)

// module *wasmtime.Module, store *wasmtime.Store, state *State
func New(wasm []byte) (*Instance, error) {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	module, err := wasmtime.NewModule(engine, wasm)
	if err != nil {
		return nil, err
	}

	memory, err := createMemory(wasm, store)
	if err != nil {
		return nil, err
	}

	linker := wasmtime.NewLinker(engine)
	instance := Instance{
		memory: memory,
		store:  store,
		Module: module,
	}
	createImport(linker, &instance)

	inst, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, err
	}

	instance.instance = inst
	return &instance, nil
}

func (inst *Instance) WrapInvoke(method string, args []byte, env []byte) (State, error) {
	state := State{
		Method: []byte(method),
		Args:   args,
		Env:    env,
	}
	inst.State = &state
	invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
	_, err := invoke.Func().Call(
		inst.store,
		len(method),
		len(args),
		len(env),
	)
	return state, err
}
