.PHONY: build view all build-docker build-custom

DOCKER_IMAGE=tinygo/tinygo:0.13.1
DOCKER_CUSTOM=cosmwasm/tinygo:latest

all: build-docker public
	@ ls -l main.wasm

build:
	tinygo build -o main.wasm -target wasm ./main.go

# build-docker uses the official build and default flags
build-docker:
	docker run --rm -v $(shell pwd):/code $(DOCKER_IMAGE) tinygo build -o /code/main.wasm -target wasm /code/main.go

# build-custom uses our custom build, but default flags
build-custom:
	docker run --rm -v $(shell pwd):/code $(DOCKER_CUSTOM) tinygo build -o /code/main.wasm -target wasm /code/main.go

# build-cosmwasm uses our custom build, and custom flags
build-cosmwasm:
	docker run --rm -v $(shell pwd):/code $(DOCKER_CUSTOM) tinygo build -tags cosmwasm -o /code/contract.wasm -target wasm /code/main.go

# build-cosmwasm-trim is like build-cosmwasm but trims out debug code for small byte size
build-cosmwasm-trim:
	docker run --rm -v $(shell pwd):/code $(DOCKER_CUSTOM) tinygo build -tags cosmwasm -no-debug -o /code/main.wasm -target wasm /code/main.go

view:
	@ wasm-nm main.wasm
	@ ls -l main.wasm

public:
	wasm-nm -e main.wasm
	wasm-nm -i main.wasm
