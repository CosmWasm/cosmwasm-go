.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

tiny-build:
	rm -rf bin
	go build -o ./bin/easyjson github.com/mailru/easyjson/easyjson

generate: tiny-build generate-std generate-contracts

generate-std:
	# rm -f std/*_easyjson.go
	./bin/easyjson -all -snake_case \
		./std/env.go \
		./std/errors.go \
		./std/systemerror.go \
		./std/types.go \
		./std/msg.go \
		./std/query.go
	./bin/easyjson -all -snake_case -build_tags=cosmwasm ./std/exports.go 2>/dev/null || true

generate-contracts:
	./bin/easyjson -all -snake_case \
		./example/hackatom/src/msg.go \
		./example/hackatom/src/state.go

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std

test-contracts:
	cd example/hackatom && $(MAKE) unit-test

examples: hackatom

hackatom:
	./scripts/compile.sh hackatom
	@ wasm-nm -e hackatom.wasm
	@ wasm-nm -i hackatom.wasm
