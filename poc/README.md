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
