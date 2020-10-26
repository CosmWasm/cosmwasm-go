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
