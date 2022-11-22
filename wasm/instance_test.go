package wasm

import (
	"os"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
)

func testBigNumber(inst *WasmInstance) func(t *testing.T) {
	return func(t *testing.T) {
		invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
		t.Logf("complete %#v", invoke)
	}
}

func testSimpleEnv(inst *WasmInstance) func(t *testing.T) {
	return func(t *testing.T) {
		invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
		t.Logf("complete %#v", invoke)
	}
}

func testSimpleInvoke(inst *WasmInstance) func(t *testing.T) {
	return func(t *testing.T) {
		invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
		t.Logf("complete %#v", invoke)
	}
}

func testSimpleSubinvokeInvoke(inst *WasmInstance) func(t *testing.T) {
	return func(t *testing.T) {
		invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
		t.Logf("complete %#v", invoke)
	}
}

func testSimpleSubinvokeSubinvoke(inst *WasmInstance) func(t *testing.T) {
	return func(t *testing.T) {
		invoke := inst.instance.GetExport(inst.store, "_wrap_invoke")
		t.Logf("complete %#v", invoke)
	}
}

func TestImports(t *testing.T) {
	cases := []struct {
		name     string
		wasmType string
		wasmData string
		fn       func(inst *WasmInstance) func(t *testing.T)
	}{
		{
			name:     "big-number",
			wasmType: "file",
			wasmData: "cases/big-number/wrap.wasm",
			fn:       testBigNumber,
		},
		{
			name:     "simple-env",
			wasmType: "file",
			wasmData: "cases/simple-env/wrap.wasm",
			fn:       testSimpleEnv,
		},
		{
			name:     "simple-invoke",
			wasmType: "file",
			wasmData: "cases/simple-invoke/wrap.wasm",
			fn:       testSimpleInvoke,
		},
		{
			name:     "subinvoke/invoke",
			wasmType: "file",
			wasmData: "cases/simple-subinvoke/invoke/wrap.wasm",
			fn:       testSimpleSubinvokeInvoke,
		},
		{
			name:     "subinvoke/subinvoke",
			wasmType: "file",
			wasmData: "cases/simple-subinvoke/subinvoke/wrap.wasm",
			fn:       testSimpleSubinvokeSubinvoke,
		},
	}

	for i := range cases {
		tcase := cases[i]
		var (
			module []byte
			err    error
		)
		switch tcase.wasmType {
		case "file":
			module, err = os.ReadFile(tcase.wasmData)
		case "wat":
			module, err = wasmtime.Wat2Wasm(tcase.wasmData)
		}
		if err != nil {
			t.Fatalf("Can't read wasm file: %s", err)
		}
		inst, err := NewInstance(module)
		if err != nil {
			t.Fatalf("Can't create instance: %s", err)
		}
		t.Run(tcase.name, tcase.fn(inst))
	}
}
