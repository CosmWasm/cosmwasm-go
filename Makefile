.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std
	go test $(TEST_FLAG) ./std/safe_math

test-contracts:
	cd example/hackatom && $(MAKE) unit-test
	# cd example/hackatom && $(MAKE) test
	# cd example/erc20 && $(MAKE) test

examples: erc20 hackatom

erc20:
	./scripts/compile.sh erc20
	@ wasm-nm -e erc20.wasm
	@ wasm-nm -i erc20.wasm

hackatom:
	./scripts/compile.sh hackatom
	@ wasm-nm -e hackatom.wasm
	@ wasm-nm -i hackatom.wasm
