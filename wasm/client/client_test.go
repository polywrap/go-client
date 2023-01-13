package client

import (
	"reflect"
	"testing"

	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

func TestClient(t *testing.T) {
	cases := []struct {
		name   string
		path   string
		invoke func(*Client, *uri.URI) (any, error)
		expRes any
	}{
		{
			name: "simple-calculator",
			path: "wrap://fs/../cases/simple-calculator",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]int32, int32](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				})
			},
			expRes: 12,
		},
		{
			name: "simple-subinvoke/subinvoke",
			path: "wrap://fs/../cases/simple-subinvoke/subinvoke",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]int32, int32](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				})
			},
			expRes: 12,
		},
		{
			name: "simple-subinvoke/invoke",
			path: "wrap://fs/../cases/simple-subinvoke/invoke",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]int32, string](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				})
			},
			expRes: 12,
		},
	}

	for i := range cases {
		tcase := cases[i]
		u, _ := uri.New("wrap://fs/../cases/simple-subinvoke/subinvoke")
		t.Run(tcase.name, func(t *testing.T) {
			client := New(&ClientConfig{
				Resolver: wasm.NewBaseResolver(
					wasm.NewRedirectResolver(map[string]*uri.URI{
						"wrap://ens/add.eth": u,
					}),
					wasm.NewFsResolver(),
				),
			})
			wrapUri, err := uri.New(tcase.path)
			if err != nil {
				t.Fatalf("bad wrapUri: %s (%s)", tcase.path, err)
			}
			res, err := tcase.invoke(client, wrapUri)
			if err != nil {
				t.Fatalf("invokation error: %s", err)
			}
			if reflect.DeepEqual(res, tcase.expRes) {
				t.Errorf("actual: %d, expected: %d", res, tcase.expRes)
			}
		})
	}
}
