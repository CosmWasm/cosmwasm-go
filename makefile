.PHONY: view public erc20 tester examples

DOCKER_IMAGE=tinygo/tinygo:0.13.1
DOCKER_CUSTOM=cosmwasm/tinygo:latest
WASM_FILE=erc20.wasm

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
	@ wasm-nm $(WASM_FILE)
	@ ls -l $(WASM_FILE)

public:
	wasm-nm -e $(WASM_FILE)
	wasm-nm -i $(WASM_FILE)
