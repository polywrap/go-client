package wasm

import (
	"testing"

	"github.com/polywrap/go-client/core/resolver"
)

func TestClient(t *testing.T) {
	a := int32(5)
	b := int32(7)
	expected := a + b

	client := NewClient[resolver.PackageResolver](nil)
	actual, err := Invoke[map[string]int32, int32](client, "add", map[string]int32{
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
