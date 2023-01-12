package client

import (
	"testing"

	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

func TestClient(t *testing.T) {
	u := "wrap://fs/../cases/simple-calculator"
	a := int32(5)
	b := int32(7)
	expected := a + b

	client := New(&ClientConfig{
		Resolver: &wasm.FsResolver{},
	})
	wrapUri, err := uri.New(u)
	if err != nil {
		t.Fatalf("bad wrapUri: %s (%s)", u, err)
	}
	actual, err := Invoke[map[string]int32, int32](client, *wrapUri, "add", map[string]int32{
		"a": a,
		"b": b,
	})
	if err != nil {
		t.Fatalf("invokation error: %s", err)
	}

	if *actual != expected {
		t.Errorf("actual: %d, expected: %d", *actual, expected)
	}
}
