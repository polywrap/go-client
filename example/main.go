package main

import (
	"log"

	"github.com/polywrap/go-client/wasm"
	"github.com/polywrap/go-client/wasm/client"
	"github.com/polywrap/go-client/wasm/uri"
)

type AddType struct {
	a int32
	b int32
}

func main() {
	wrapPath := "wrap://fs/../wasm/cases/simple-calculator"
	polywrapClient := client.New(&client.ClientConfig{
		Resolver: wasm.NewFsResolver(),
	})
	wrapUri, err := uri.New(wrapPath)
	if err != nil {
		log.Fatalf("bad wrapUri: %s (%s)", wrapPath, err)
	}
	args := AddType{
		a: 5,
		b: 7,
	}
	res, err := client.Invoke[AddType, int32, []byte](polywrapClient, *wrapUri, "add", args, nil)
	if err != nil {
		log.Fatalf("invokation error: %+v", err)
	}

	log.Printf("Result is: %d\n", *res)
}
