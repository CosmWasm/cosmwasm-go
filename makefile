.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std/...

test-contracts:
	cd example/erc20 && $(MAKE) test
	cd example/tester && $(MAKE) test

examples: erc20

erc20:
	./scripts/compile.sh erc20
	@ wasm-nm -e erc20.wasm
	@ wasm-nm -i erc20.wasm
