package wasm

type State struct {
	Method string
	Args   []byte
	Env    []byte

	Result []byte
	Error  []byte
}
