.PHONY: view imports exports erc20 tester examples test test-contracts test-std

DOCKER_CUSTOM=cosmwasm/tinygo:0.14.1
EMSCRIPTEN=trzeci/emscripten:1.39.8-fastcomp

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

examples: erc20 tester minimal

erc20:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/erc20.wasm /code/example/erc20/main.go

tester:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/tester.wasm /code/example/tester/main.go

minimal:
	docker run --rm $(DOCKER_FLAGS) $(DOCKER_CUSTOM) tinygo build $(TINYGO_FLAGS) -o /code/minimal.wasm /code/example/minimal/main.go

rewrite:
	docker run --rm $(DOCKER_FLAGS) $(EMSCRIPTEN) wasm2wat /code/minimal.wasm > minimal.wat
	# this just replaces all the floating point ops with unreachable. It still leaves them in the args and local variables
	cat minimal.wat | sed -E 's/^(\s*)f64\.[^()]+/\1unreachable/' | sed -E 's/^(\s*)f32\.[^()]+/\1unreachable/' > minimal-rewrite.wat
	docker run --rm $(DOCKER_FLAGS) $(EMSCRIPTEN) wat2wasm /code/minimal-rewrite.wat

# TODO: replace with local test - this requires proper local keys
upload:
#	coral tx wasm store minimal-rewrite.wasm --source https://foo.bar/123 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y -b block | jq .raw_log
	coral q tx -o json $(shell coral tx wasm store minimal-rewrite.wasm --source https://foo.bar/123 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y | jq -r .txhash; sleep 12) | jq .raw_log

view:
	@ wasm-nm erc20.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i erc20.wasm

exports:
	wasm-nm -e erc20.wasm
