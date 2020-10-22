.PHONY: view public erc20 tester examples

DOCKER_CUSTOM=cosmwasm/tinygo:v0.14.1

GO_MOUNTS=-v $(shell pwd):/code -v $(shell pwd):/go/src/github.com/cosmwasm/cosmwasm-go -v $(GOPATH)/src/github.com/cosmwasm/jsonparser:/go/src/github.com/cosmwasm/jsonparser
TINYGO_FLAGS=-tags cosmwasm -no-debug -target wasm

examples: erc20 tester

erc20:
    # TODO: automatically download jsonparser dependency
	docker run --rm $(GO_MOUNTS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/erc20.wasm /code/example/erc20/main.go

tester:
    # TODO: automatically download jsonparser dependency
	docker run --rm $(GO_MOUNTS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/tester.wasm /code/example/tester/main.go

view:
	@ wasm-nm *.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i *.wasm

exports:
	wasm-nm -e *.wasm
