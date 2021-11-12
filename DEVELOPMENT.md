# Development

This file explains how to use the library as a developer, and also how to extend it.

## TinyGo

We made a [fork of TinyGo](https://github.com/confio/tinygo) to add support for CosmWasm.
It is based off of 0.19, and the [latest TinyGo release](https://github.com/confio/tinygo) is 0.20
at the time of this writing, so we can consider an update there somehow.
(There are merge conflicts with how I customize Wasm compilation, so this is non-trivial)

We maintain a Docker image for the TinyGo compiler targeting CosmWasm on [Docker Hub](https://hub.docker.com/r/cosmwasm/tinygo/tags).
You can get the latest version simply via:

```shell script
docker pull cosmwasm/tinygo:0.19.1
```

If the latest version is not available, you can build from source:

```
git clone https://github.com/confio/tinygo.git
git checkout cw-0.19.1
docker build -t cosmwasm/tinygo:0.19.1 -f Dockerfile.wasm .

# and maybe publish
docker push cosmwasm/tinygo:0.19.1
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

 