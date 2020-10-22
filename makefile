.PHONY: view public erc20 tester examples

DOCKER_CUSTOM=cosmwasm/tinygo:v0.14.1

DOCKER_FLAGS=-w /code -v $(shell pwd):/code
TINYGO_FLAGS=-tags cosmwasm -no-debug -target wasm

examples: erc20 tester

erc20:
    # TODO: automatically download jsonparser dependency
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/erc20.wasm /code/example/erc20/main.go

tester:
    # TODO: automatically download jsonparser dependency
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/tester.wasm /code/example/tester/main.go

check:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo

#tinygo build -tags cosmwasm -no-debug -target wasm -o /code/inna.wasm /code/example/erc20/main.go

view:
	@ wasm-nm *.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i *.wasm

exports:
	wasm-nm -e *.wasm
