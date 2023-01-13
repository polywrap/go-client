package main

import (
	"log"

	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/client"
	"github.com/polywrap/go-client/wasm/uri"
)

func main() {
	wrapPath := "wrap://fs/../wasm/cases/simple-calculator"
	polywrapClient := client.New(&client.ClientConfig{
		Resolver: wasm.NewFsResolver(),
	})
	wrapUri, err := uri.New(wrapPath)
	if err != nil {
		log.Fatalf("bad wrapUri: %s (%s)", wrapPath, err)
	}
	res, err := client.Invoke[map[string]int32, int32, []byte](polywrapClient, *wrapUri, "add", map[string]int32{
		"a": 5,
		"b": 7,
	}, nil)
	if err != nil {
		log.Fatalf("invokation error: %s", err)
	}

	log.Printf("Result is: %d\n", *res)
}
