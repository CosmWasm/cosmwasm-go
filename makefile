.PHONY: view imports exports erc20 tester examples

DOCKER_CUSTOM=cosmwasm/tinygo:0.14.1

DOCKER_FLAGS=-w /code -v $(shell pwd):/code
TINYGO_FLAGS=-tags cosmwasm -no-debug -target wasm

examples: erc20 tester

erc20:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/erc20.wasm /code/example/erc20/main.go

tester:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/tester.wasm /code/example/tester/main.go

view:
	@ wasm-nm *.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i *.wasm

exports:
	wasm-nm -e *.wasm
