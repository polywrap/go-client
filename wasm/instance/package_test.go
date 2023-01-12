package instance

import (
	"os"
	"testing"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack"
	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
)

func testSimpleCalculator(inst *Instance) func(t *testing.T) {
	return func(t *testing.T) {
		a := int32(5)
		b := int32(7)
		expected := a + b

		context := msgpack.NewContext("Call method SimpleCalculator")
		encoder := msgpack.NewWriteEncoder(context)
		encoder.WriteMapLength(2)
		encoder.WriteString("a")
		encoder.WriteI32(a)
		encoder.WriteString("b")
		encoder.WriteI32(b)

		s, err := inst.WrapInvoke("add", encoder.Buffer(), nil)
		if err != nil {
			t.Errorf("can't call 'add', error: %s", err)
		}
		if len(s.Error) > 0 {
			t.Errorf("can't call 'add', error: %s", s.Error)
		}

		decoder := msgpack.NewReadDecoder(context, s.Result)
		actual := decoder.ReadI32()
		if actual != expected {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	}
}

func testBigNumber(inst *Instance) func(t *testing.T) {
	return func(t *testing.T) {
		context := msgpack.NewContext("Call method BigNumber")
		encoder := msgpack.NewWriteEncoder(context)
		expected := big.NewInt(48)
		encoder.WriteMap(map[interface{}]interface{}{
			"arg1": 2,
			"arg2": 3,
			"obj": map[interface{}]interface{}{
				"prop1": 2,
				"prop2": 4,
			},
		}, func(encoder msgpack.Write, key, value interface{}) {
			k := key.(string)
			encoder.WriteString(k)
			switch k {
			case "arg1", "arg2":
				v := value.(int)
				encoder.WriteBigInt(big.NewInt(int64(v)))
			case "obj":
				v := value.(map[interface{}]interface{})
				encoder.WriteMap(v, func(encoder msgpack.Write, key, value interface{}) {
					v := value.(int)
					k := key.(string)
					encoder.WriteString(k)
					encoder.WriteBigInt(big.NewInt(int64(v)))
				})
			}
		})

		s, err := inst.WrapInvoke("method", encoder.Buffer(), nil)
		if err != nil {
			t.Errorf("can't call 'method', error: %s", err)
		}
		decoder := msgpack.NewReadDecoder(context, s.Result)
		actual := decoder.ReadBigInt()

		if actual.Cmp(expected) != 0 {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	}
}

func TestInstance(t *testing.T) {
	cases := []struct {
		name     string
		wasmType string
		wasmData string
		fn       func(inst *Instance) func(t *testing.T)
	}{
		{
			name:     "simple-calculator",
			wasmType: "file",
			wasmData: "../cases/simple-calculator/wrap.wasm",
			fn:       testSimpleCalculator,
		},
		{
			name:     "big-number",
			wasmType: "file",
			wasmData: "../cases/big-number/wrap.wasm",
			fn:       testBigNumber,
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
		default:
			t.Fatalf("unknown wasmType: [%s](%s)", tcase.wasmType, tcase.wasmData)
		}
		if err != nil {
			t.Fatalf("Can't read wasm file: %s", err)
		}
		inst, err := New(module)
		if err != nil {
			t.Fatalf("Can't create instance: %s", err)
		}
		t.Run(tcase.name, tcase.fn(inst))
	}
}
