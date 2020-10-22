.PHONY: view imports exports erc20 tester examples test test-contracts test-std

DOCKER_CUSTOM=cosmwasm/tinygo:0.14.1

DOCKER_FLAGS=-w /code -v $(shell pwd):/code
TINYGO_FLAGS=-tags cosmwasm -no-debug -target wasm

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
TEST_FLAG=-v -count=1

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std/...

test-contracts:
	go test $(TEST_FLAG) ./example/tester/src
	go test $(TEST_FLAG) ./example/erc20/src

examples: erc20 tester

erc20:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/erc20.wasm /code/example/erc20/main.go

tester:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/tester.wasm /code/example/tester/main.go

view:
	@ wasm-nm erc20.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i erc20.wasm

exports:
	wasm-nm -e erc20.wasm
