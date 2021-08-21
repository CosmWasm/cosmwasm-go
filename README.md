## Overview

This is a demo for cosmwasm-go build, it is in alpha state, but it is compatible
with CosmWasm 0.11 and 0.12 (they use the same Wasm API).

This is currently in a **private pre-release**. This repo should only be shared
with select developers who will work on improving it. When we are happy with
the stability, we will publish this.

## Features

The following areas have been tested, both in unit tests and in integration
tests using `go-cosmwasm` to ensure compatibility. Please look at the 
[`hackatom` contract](./example/hackatom) for a good example with test coverage.

* Export `init`, `handle`, `query`, `migrate` and return success response or errors
* Get and Set from persistent storage
* Canonicalize and Humanize addresses with the Api
* Query into the bank module (use querier)
* Parse most json structs (see caveats below)
* Simple API so you can focus on the business logic
* Easy setup and similar format for unit test and integration tests, to allow
  easy porting them.
  
The following areas have some code but have not been tested and *likely buggy*:

* Remove from storage
* Storage iterators
* Querying staking/wasm modules

These can be called, but will likely not function 100% and will require
some debugging to get them working.

## Caveats

Beyond the "likely buggy" Apis above, the following "features" are 
definitely not implemented:

* Support for Uint128 (parsing strings, dealing with large integers)
* Helper functions for creating messages or queries
* Proper CI integration (so please run a full test suite locally before merging)

## JSON

There is a totally custom json lib that manages to (1) run inside
a very limited reflect package of TinyGo, (2) contain no floating
point ops or other code that would be rejected by CosmWasm VM and
(3) produces a nice dev experience. In particular we want to support
structs (and not just maps or some path based query language).

On the high-level it looks just like "encoding/json", but many things are
not possible. Even to get where we are was a very tough feat and an amazing
engineering acheivement from the team of OkEx. These are known issues:

* `[]byte` doesn't properly serialize. It looks like `[12, 65, 152, 4]`. It errors
  when it receives a base64 encoded string. Don't use `[]byte` in any struct
  that will be Serialized with JSON. You can use string an manually base64
  encode/decode it for now.
* `string` does not escape quotes (or anything else). In particular, that means
  if you pass in some JSON as a string `"{"count": 123}"` this will break the
  schema of the containing object. You need to encode JSON to a base64 string
  and treat it as binary.
* You cannot use pointers in any object that will be used by `ezjson`. 
  The OkEx team spent quite some time on this one and so did I. 
  Basically, you cannot do `reflect.New()` in TinyGo, thus you cannot
  create new objects. When you don't use pointers, the objects passed into
  the `Unmarshal` function has an existent, zeroed object to fill with
  data and everything works. We use a lot of `omitempty` to keep these
  from serializing (which omits more objects in ezjson than encoding/json)
* Arrays of structs are very buggy. The will Marshal, but in general, they 
  will fail to Unmashal unless you follow some work-arounds. This is due
  to the same issue with pointers - we cannot create new structs with reflect
  inside of TinyGo. Thus, we cannot dynamically add new items to the slice.
  If you look at how we
  [parse multiple coins in queries](https://github.com/CosmWasm/cosmwasm-go/blob/master/std/query.go#L69-L76),
  you can see we first create an array with the maximum size we accept (eg. 8).
  Then we unmarshal into these existing objects. Then we trim off all objects
  that are still empty (received no data).
  
We plan to improve `[]byte` and `string` support. The other issues are
tied to the `reflect` support in TinyGo and would probably need a
custom code-gen JSON solution in order to avoid those issues. That is very
much out of scope in the near future. 
  
### Usage

We require `docker` installed and executable by the current user. This
is used to compile the contracts to Wasm. You can code contracts and
write unit test with only a normal Golang installation (requires 1.14+)

You can try the following top-level commands:

```
# Run unit tests on the standard library as well as `erc20` and `hackatom` contracts.
# Also compiles the contracts to wasm and runs integration tests
make test

# This will compile wasm binaries for all contracts and leave them in this directory
make examples
```

To try out a contract, either check out [hackatom](./example/hackatom),
running the tests and editing the code. Or start your own contract
by going to the [template](./example/template) directory and follow the
instructions on how to get started.

Both of these support the following commands: `make unit-test`, `make build`, 
and `make test`.

## Build system

We use docker tooling to get consistent builds acorss dev machines.
In my minimal experience, these seem to also produce deterministic
builds, but I would like others to try this out on other machines.
The following produces the same sha256 hash everytime I run it:

```
cd example/hackatom
make build && sha256sum hackatom.wasm

# this will test the wasm code
go test ./integration
```

However, the docker image for our custom TinyGo is not yet published.
In order to build locally you can do the following:

```
git clone https://github.com/confio/tinygo.git
git checkout cw-0.19.0
docker build -t cosmwasm/tinygo:0.19.0 -f Dockerfile.wasm .
```

Once it is finished, you should be able to successfully run `make build` on hackatom

## Performance

Many people ask how these compare to the rust contracts. I have yet to
do a detailed comparion, but now that we have two versions of the same
hackatom contract, we can do a rough side-by-side analysis.

**The good:**

The contract size is significantly lower than the Rust version, that is
97kB for the TinyGo version compared to 179kB for the Rust version.
(There is a sha256 algorithm in the Rust version missing from the Go version,
but that is only about 20Kb of the size)

**The bad:**

It uses much more gas. I am unsure if this is a one-time startup cost
due to initializing the various components of the Go runtime. The
much less efficient JSON parsing (it is a large component of the Rust
cost an that is optimized codegen). Or generally less efficient code.

For example, in a `cosmwasm-vm` test `instance::singlepass_test::contract_deducts_gas_init`
I see `init` uses *829918* gas in the Go contract and *67349* in the equivalent
Rust contract. That is about 12x more!

However, before you get too scared, please not this is wasmer gas and only
measuring the CPU usage. We normalize that with a factor of 100 for SDK
gas, meaning the Go contract would require 8300 SDK gas and the Rust one 673.
Given we have to pay ~2400 SDK gas just to store one data item, we have
a "setup tax" of 40,000 SDK gas for running a contract, and the native bank
send function requires about 55,000 SDK gas, these numbers are more
acceptable.

**The ugly:**

In short, yes, the Go contracts are significantly less CPU efficient than
the Rust ones. However, the CPU usage is not a major portion of the gas
usage in most contracts. If you work on a computationally heavy contract
(lots of math, hashing, or such) or parse/serialize large JSON objects,
the lower performance here should be acceptable.

The biggest issue is that we have no idea where this difference comes from
and what are the biggest consumers of CPU cycles. We would be happy for some
profiling work and maybe can use codegen to reduce some of the CPU time in
exchange for more code size.