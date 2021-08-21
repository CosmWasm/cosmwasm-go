.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

tiny-build:
	rm -rf bin
	go build -o ./bin/easyjson github.com/mailru/easyjson/easyjson

generate: tiny-build
	# rm -f std/*_easyjson.go
	./bin/easyjson -all -snake_case \
		./std/env.go \
		./std/errors.go \
		./std/types.go \
		./std/query.go
	./bin/easyjson -all -snake_case -build_tags=cosmwasm ./std/exports.go 2>/dev/null || true

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
