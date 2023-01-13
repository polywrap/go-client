package client

import (
	"reflect"
	"testing"

	"github.com/consideritdone/polywrap-go/polywrap/msgpack/big"
	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/uri"
)

func TestClient(t *testing.T) {
	type EnvType struct {
		ExternalArray  []uint32
		ExternalString string
	}

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
				return Invoke[map[string]int32, int32, []byte](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				}, nil)
			},
			expRes: 12,
		},
		{
			name: "big-number",
			path: "wrap://fs/../cases/big-number",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				type ArgType struct {
					Arg1 *big.Int
					Arg2 *big.Int
					Obj  struct {
						Prop1 *big.Int
						Prop2 *big.Int
					}
				}
				return Invoke[ArgType, *big.Int, []byte](c, *u, "method", ArgType{
					Arg1: big.NewInt(2),
					Arg2: big.NewInt(3),
					Obj: struct {
						Prop1 *big.Int
						Prop2 *big.Int
					}{
						Prop1: big.NewInt(3),
						Prop2: big.NewInt(4),
					},
				}, nil)
			},
			expRes: big.NewInt(72),
		},
		{
			name: "simple-subinvoke/subinvoke",
			path: "wrap://fs/../cases/simple-subinvoke/subinvoke",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]int32, int32, []byte](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				}, nil)
			},
			expRes: 12,
		},
		{
			name: "simple-subinvoke/invoke",
			path: "wrap://fs/../cases/simple-subinvoke/invoke",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]int32, string, []byte](c, *u, "add", map[string]int32{
					"a": 5,
					"b": 7,
				}, nil)
			},
			expRes: 12,
		},
		{
			name: "simple-env",
			path: "wrap://fs/../cases/simple-env",
			invoke: func(c *Client, u *uri.URI) (any, error) {
				return Invoke[map[string]string, EnvType](c, *u, "externalEnvMethod", nil, EnvType{
					ExternalArray:  []uint32{1, 2, 3},
					ExternalString: "123",
				})
			},
			expRes: EnvType{
				ExternalArray:  []uint32{1, 2, 3},
				ExternalString: "123",
			},
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
