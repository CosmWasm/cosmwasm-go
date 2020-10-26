.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std/...

test-contracts:
	go test $(TEST_FLAG) ./example/tester/src
	go test $(TEST_FLAG) ./example/erc20/src

examples: erc20 tester minimal

erc20:
	./scripts/compile.sh erc20
	@ wasm-nm -e erc20.wasm
	@ wasm-nm -i erc20.wasm

tester:
	./scripts/compile.sh tester
	@ wasm-nm -e erc20.wasm
	@ wasm-nm -i erc20.wasm

minimal:
	./scripts/compile.sh minimal
	@ wasm-nm -e erc20.wasm
	@ wasm-nm -i erc20.wasm

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
