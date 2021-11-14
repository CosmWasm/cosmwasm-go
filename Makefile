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
		./std/types/env.go \
		./std/types/ibc.go \
		./std/types/msg.go \
		./std/types/query.go \
		./std/types/subcall.go \
		./std/types/systemerror.go \
		./std/types/types.go \

generate-contracts:
	./bin/tinyjson -all -snake_case \
		./example/hackatom/src/state.go
	./bin/tinyjson -all -snake_case \
		./example/hackatom/src/msg.go

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std
	go test $(TEST_FLAG) ./std/mocks

test-contracts:
	cd example/hackatom && $(MAKE) unit-test

examples: hackatom

hackatom:
	@echo "VERSION=0.19.0 make hackatom - will run with different cosmwasm/tinygo image" 
	./scripts/compile.sh hackatom
	./scripts/check.sh hackatom.wasm
	./scripts/strip_floats.sh hackatom.wasm
