## Overview
This is a demo for cosmwasm-go build, it only run as an mvp, now it support:
* loaded by cosmwasm-simulate tool
* compile byte code to instance
* init & query & handle export
* init func call and return test(result can not convert by `from_slice`)


some code very rough now, just a debugging code for test   
The next work is testing suitable json lib and process param and return value

## Build

Docker tool required     
```sh
make build-cosmwasm
```
## How to debug it
### cosmwasm-simulate tool required   
* Download cosmwasm-simulate from https://github.com/CosmWasm/cosmwasm-simulate
* Load `cosmwasm-simulate` by `CLion`, compile it in IDE
### tinyGo required
* tinyGo compiler has already packed into docker image file, just need install `docker` and build poc, it will finish automatically
### Start debug
1. Build poc, copy absloute path of `contract.wasm` 
2. Open CLion, load cosmwasm-simulate, Add cargo config, Set Command like `run [absloute path of contract.wasm]`, change Channel to nightly, compile it in IDE
3. Add breakpoint at [Here](https://github.com/CosmWasm/cosmwasm-simulate/blob/master/src/contract_vm/engine.rs#L124)
4. Start debug in CLion, input all message as follow:
```shell
Input call type(init | handle | query):
init
Input json string:
{}
```
Then, breakpoint will actived, you can press F7 & F8 to debug it step by step, enjoy it ~
