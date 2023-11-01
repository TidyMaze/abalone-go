# Abalone go

A simple abalone game written in go.
Also compiles to webassembly.

## Build

### Native

```shell
go build
./main
```

### Webassembly

```shell
brew install wasmtime
GOOS=wasip1 GOARCH=wasm go build -o main.wasm main.go
wasmtime main.wasm
```
