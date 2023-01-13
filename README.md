![Public Release Announcement](https://user-images.githubusercontent.com/5522128/177473887-2689cf25-7937-4620-8ca5-17620729a65d.png)


# Polywrap Go client

> [Polywrap](https://polywrap.io) is a developer tool that enables easy integration of Web3 protocols into any application. It makes it possible for applications on any platform, written in any language, to read and write data to Web3 protocols.

# Working Features

This Polywrap clients enable the execution of WebAssembly Polywrappers (or just “wrappers”) on various environments, regardless of what language this wrapper was built in.

The various clients are built following the functionality of the JavaScript Polywrap Client, which is currently more robust and battle tested, as it has additional capabilities than other MVPs. In the future, the Polywrap DAO will continue improving the various client’s capabilities to reach feature parity with the JS stack, improving the experience in parallel clients for other languages like Python, Go, and Rust.

Here you can see which features have been implemented on each language, and make the decision of which one to use for your project.

| Feature | [Python](https://github.com/polywrap/python-client) | [Javascript](https://github.com/polywrap/toolchain) |  [Go](https://github.com/polywrap/go-client) | [Rust](https://github.com/polywrap/rust-client) |
| -- | -- | -- | -- | -- |
| **Invoke**  | ✅ | ✅ | ✅ | ⚙️|
| Subinvoke | ⚙️ | ✅ | ✅ |  |
| Interfaces | ❌ | ✅ | ✅ | |
| Env Configuration | ⚙️ | ✅ | ✅ | |
| Client Config | ⚙️ | ✅ | ✅ | ⚙️|
| Plugin Wrapper | ❌ | ✅ | | |
| Wrap Manifest | ⚙️ | ✅ | | |
| **Uri Resolution** | ⚙️ | ✅ | ✅ | ⚙️ |
| Uri: Filesystem|✅|✅| ✅ |
| Uri: IPFS |❌|✅| || |
| Uri: ENS |❌|✅| | | |

> TODO: Update table above according to test harness and maybe mention other wip clients (rust, python)

|status| |
| -- | -- |
|✅ | fully working|
|⚙️| partially working|
|❌|not yet implemented|

## Prerequisites

## Golang

Proceed to installation by following [these instructions](https://go.dev/doc/install).

To verify Go is installed run:
```
go version
```
your output in this case should be something like `go version go1.18.1 linux/amd64`.


## Using Polywrap Go client

Example of Golang app that uses [SimpleCalculator](https://github.com/polywrap/toolchain/tree/origin-dev/packages/test-cases/cases/wrappers/wasm-as/simple-calculator) wrapper

```go
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
```

## Development


### Clone the repository
```
git clone https://github.com/polywrap/go-client.git
```

# Test the client

By running this command in the root path, all written tests will be executed

```
go test -v ./...
```
