.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

tiny-build:
	rm -rf bin
	go build -o ./bin/tinyjson github.com/CosmWasm/tinyjson/tinyjson

clean:
	rm -f std/*_tinyjson.go
	rm -f example/hackatom/src/*_tinyjson.go

generate: tiny-build generate-std generate-contracts

generate-std:
	./bin/tinyjson -all -snake_case \
		./std/env.go \
		./std/errors.go \
		./std/systemerror.go \
		./std/types.go \
		./std/msg.go \
		./std/subcall.go \
		./std/query.go \
		./std/ibc.go

generate-contracts:
	./bin/tinyjson -all -snake_case \
		./example/hackatom/src/state.go
	./bin/tinyjson -all -snake_case \
		./example/hackatom/src/msg.go

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
