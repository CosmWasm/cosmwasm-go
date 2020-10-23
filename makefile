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
upload-minimal:
#	coral tx wasm store minimal-rewrite.wasm --source https://foo.bar/123 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y -b block | jq .raw_log
	coral q tx -o json $(shell coral tx wasm store minimal-rewrite.wasm --source https://foo.bar/123 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y | jq -r .txhash; sleep 12) | jq .raw_log

# this is ugly but succeeds in uploading our code!!
rewrite-erc20:
	docker run --rm $(DOCKER_FLAGS) $(EMSCRIPTEN) wasm2wat /code/erc20.wasm > erc20.wat
	# this just replaces all the floating point ops with unreachable. It still leaves them in the args and local variables
	cat erc20.wat | sed -E 's/^(\s*)f[[:digit:]]{2}\.[^()]+/\1unreachable/' | sed -E 's/^(\s*)i[[:digit:]]{2}\.trunc_[^()]+/\1unreachable/' | sed -E 's/^(\s*)i[[:digit:]]{2}\.reinterpret_[^()]+/\1unreachable/' > erc20-rewrite.wat
	docker run --rm $(DOCKER_FLAGS) $(EMSCRIPTEN) wat2wasm /code/erc20-rewrite.wat

# Manually tested that re-written wasm code runs well
upload-erc20:
#	coral q tx -o json $(shell coral tx wasm store erc20-rewrite.wasm --source https://foo.bar/erc20 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y | jq -r .txhash; sleep 12) | jq .raw_log
	coral tx wasm store erc20-rewrite.wasm --source https://foo.bar/erc20 --from validator --gas 1000000 --gas-prices 0.025ushell -o json -y | jq -r .txhash
	# coral q tx $STORE_HASH -o json | jq .logs
	# coral tx wasm instantiate 131 '{"name":"OKB","symbol":"OKB","decimal":10,"total_supply":170000}' --label Test1 --from validator --gas 400000 --gas-prices 0.025ushell -y
	# coral q tx $INIT_HASH -o json | jq .logs
	# coral tx wasm execute $CONTRACT '{"Transfer":{"to":"coral1reednyl4473um535crt0tuqgkfy2k68tzy5762","value": 2000}}' --from validator --gas 200000 --gas-prices 0.025ushell -y
	# coral q tx $EXEC_HASH -o json | jq .logs
	# coral q wasm contract-state smart $CONTRACT '{"balance":{"address":"coral1reednyl4473um535crt0tuqgkfy2k68tzy5762"}}'


view:
	@ wasm-nm erc20.wasm
	@ ls -l *.wasm

imports:
	wasm-nm -i erc20.wasm

exports:
	wasm-nm -e erc20.wasm
