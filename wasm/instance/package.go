package instance

import (
	"errors"
	"strings"

	"github.com/bytecodealliance/wasmtime-go"
)

type (
	State struct {
		Method []byte
		Args   []byte
		Env    []byte
		Result []byte
		Error  []byte
		Nested *State
	}
	Instance struct {
		memory   *wasmtime.Memory
		store    *wasmtime.Store
		instance *wasmtime.Instance
		Module   *wasmtime.Module
		State    *State
	}
)

var (
	ENV_MEMORY_IMPORTS_SIGNATURE = []byte{0x65, 0x6e, 0x76, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x02}
	ErrNowWasmMemory             = errors.New(strings.Join(
		[]string{
			`Unable to find Wasm memory import section.`,
			`  Modules must import memory from the "env" module's`,
			`  "memory" field like so:`,
			`    (import "env" "memory" (memory (;0;) #))`,
		},
		"\n",
	))
)
