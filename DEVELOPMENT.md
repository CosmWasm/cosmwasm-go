# Development

This file explains how to use the library as a developer, and also how to extend it.

## TinyGo

We made a [fork of TinyGo](https://github.com/confio/tinygo) to add support for CosmWasm.
It is based off of 0.20, the [latest TinyGo release](https://github.com/confio/tinygo)
at the time of this writing.
(Please note merge conflicts during rebasing on newer tinygo versions are to be expected, so this is non-trivial)

Note: 0.19.0 runs only on arm not intel/amd. It does produce final images with no floating point ops. 0.20 builds on both but includes some floating point ops. [See github issue #74](https://github.com/confio/cosmwasm-go/issues/74)

We maintain a Docker image for the TinyGo compiler targeting CosmWasm on [Docker Hub](https://hub.docker.com/r/cosmwasm/tinygo/tags).
You can get the latest version simply via:

```shell script
docker pull cosmwasm/tinygo:0.19.0
```

If the latest version is not available, you can build from source:

```
git clone https://github.com/confio/tinygo.git
git checkout cw-0.19.0
docker build -t cosmwasm/tinygo:0.19.0 -f Dockerfile.wasm .

# and maybe publish
docker push cosmwasm/tinygo:0.19.0
```


## Build system

We use docker tooling to get consistent builds across dev machines.
In my minimal experience, these seem to also produce deterministic
builds, but I would like others to try this out on other machines.
The following produces the same sha256 hash everytime I run it:

```
cd example/hackatom
make build && sha256sum hackatom.wasm

# this will test the wasm code
go test ./integration
```

Once it is finished, you should be able to successfully run `make build` on hackatom

## Building TinyJSON

We touched on [TinyJSON in the README](./README.md#json) but didn't explain how to build.
You need to run codegen on all structs and call `bz, err := msg.MarshalJSON()` rather
than `bz, err := json.MarshalJSON(msg)`.

### JSON Struct Definitions

You can use the same JSON definitions as anywhere with two nice improvements.
If you just want to snake_case the names, to match Rust style, you don't need to add
annotations like this everywhere:

```go
type Msg struct {
  FooBar string `json:"foo_bar"`
}
```

Just include the `-snake_case` flag and it will auto-generate that. Furthermore,
when working for serialization compatible with Rust contracts, there are many places
where we want to use `Foo []MyType` but have to define some `type MyTypes []MyType` that
overrides the Marshalling to always serialize an empty array `[]` not `null`.
In this case, you can define that behavior with a simple json argument `emptyslice`,
which is a custom addition in the `tinyjson` compiler.

```go
type AllValidatorsResponse struct {
	Validators []Validator `json:"validators,emptyslice"`
}
```

The downside is that further customization of the `MarshalJSON()` method of a struct is non-trivial
and is considered "unsupported advanced use-case" currently.

### Build Process

Makefile usage, how to add to your contract
### Bootstrapping Errors

TODO: when we delete files, but they depend on the code...