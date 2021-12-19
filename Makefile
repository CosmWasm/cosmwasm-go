.PHONY: view imports exports erc20 tester examples test test-contracts test-std

# Set on the command line for verbose output, eg.
# TEST_FLAG=-v make test
#TEST_FLAG=-v -count=1

tiny-build:
	rm -rf ./bin/tinyjson
	go build -o ./bin/tinyjson github.com/CosmWasm/tinyjson/tinyjson

generator-build:
	rm -rf ./bin/generator
	go build -o ./bin/generator ./cmd/generator

clean:
	rm -f std/types/*_tinyjson.go
	rm -f example/hackatom/src/*_tinyjson.go

generate: tiny-build generator-build generate-std generate-contracts

generate-std:
	./bin/tinyjson -all -snake_case \
		./std/types/env.go \
		./std/types/fraction.go \
		./std/types/ibc.go \
		./std/types/msg.go \
		./std/types/query.go \
		./std/types/subcall.go \
		./std/types/systemerror.go \
		./std/types/types.go \

generate-contracts:
	./bin/tinyjson -all -snake_case \
		./example/hackatom/src/state.go \
		./example/hackatom/src/msg.go
	go generate ./example/...

test: test-std test-contracts

test-std:
	go test $(TEST_FLAG) ./std
	go test $(TEST_FLAG) ./std/mocks

test-contracts:
	cd example/hackatom && $(MAKE) unit-test

examples: hackatom

# we need to move this to example/hackatom, so it will be run in the integration tests in CI
hackatom:
	@echo "VERSION=latest make hackatom - will run with different cosmwasm/tinygo image" 
	./scripts/compile.sh hackatom
	./scripts/check.sh hackatom.wasm
	./scripts/increase_memory.sh hackatom.wasm
	# ./scripts/strip_floats.sh hackatom.wasm
	mv hackatom.wasm example/hackatom

queue:
	@echo "VERSION=latest make queue - will run with different cosmwasm/tinygo image"
	./scripts/compile.sh queue
	./scripts/check.sh queue.wasm
	./scripts/increase_memory.sh queue.wasm
	mv queue.wasm example/queue

identity:
	@echo "VERSION=latest make identity - will run with different cosmwasm/tinyho image"
	./scripts/compile.sh identity
	./scripts/check.sh identity.wasm
	./scripts/increase_memory.sh identity.wasm
	mv identity.wasm example/identity