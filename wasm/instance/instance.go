package instance

import (
	"github.com/bytecodealliance/wasmtime-go"
)

func NewState(invoker Invoker, method, args, env []byte) *State {
	return &State{
		Method:  method,
		Args:    args,
		Env:     env,
		Invoker: invoker,
	}
}

func New(wasm []byte, state *State) (*Instance, error) {
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
	instance.State = state
	instance.instance = inst
	return &instance, nil
}

func (inst *Instance) Call() (*State, error) {
	invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
	_, err := invoke.Func().Call(
		inst.store,
		len(inst.State.Method),
		len(inst.State.Args),
		len(inst.State.Env),
	)
	return inst.State, err
}
