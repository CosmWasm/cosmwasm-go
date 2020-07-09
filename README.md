## Overview
This is a demo for cosmwasm-go build, it only run as an mvp, now it support:
* loaded by cosmwasm-simulate tool
* compile byte code to instance
* init & query & handle export
* simple json marshal&unmarshal support in contract
* simple db operation support in contract
* simple but whole logic in contract code

some code are very rough now, just a demo version

## Build

Docker tool required     
```sh
make build-cosmwasm
```

## Run Test
* Prepare   
Clone cosmwasm-simulate tool from https://github.com/CosmWasm/cosmwasm-simulate/tree/debug-cosmwasm-go 
```sh
git clone git@github.com:CosmWasm/cosmwasm-simulate.git
git checkout debug-cosmwasm-go
```
The file `wasm_go_poc/contract.wasm` was built by `make build-cosmwasm`, you can build it by self and replace it~ 
* Run   
```sh
./run-go.sh
```

## What's stuff in contract
* init
   * function : init an account and money, set account name and password, init money balance
   * json args : {"UserName":"useraccount","Password":"111222","Money":100}
* handle
   * brun : burn money from account, password required, balance will reduce 10 during every call
      - json: {"Operation":"burn","Password":"111222"}
   * save : save money to account, password required, balance will increase 10 during every call
      - json: {"Operation":"save","Password":"111222"}
* query
   * balance : query balance of account, just return a string to show it
      - json: {"QueryType":"balance"}
   * user : query user name of account, just return a string to show it
      - json: {"QueryType":"user"}
