.PHONY: build view all build-docker build-custom test

DOCKER_IMAGE=tinygo/tinygo:0.13.1
DOCKER_CUSTOM=cosmwasm/tinygo:latest
WASM_FILE=contract.wasm

all: build-docker public
	@ ls -l main.wasm

test:
	go test ./src

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
	docker run --rm -v $(shell pwd):/code -v $(shell pwd):/go/src/github.com/cosmwasm/cosmwasm-go/example -v $(GOPATH)/src/github.com/cosmwasm/cosmwasm-go/std:/go/src/github.com/cosmwasm/cosmwasm-go/std -v $(GOPATH)/src/github.com/cosmwasm/jsonparser:/go/src/github.com/cosmwasm/jsonparser $(DOCKER_CUSTOM) tinygo build -tags cosmwasm -o /code/contract.wasm -target wasm /code/main.go

# build-cosmwasm-trim is like build-cosmwasm but trims out debug code for small byte size
build-cosmwasm-trim:
	docker run --rm -v $(shell pwd):/code -v $(shell pwd):/go/src/github.com/cosmwasm/cosmwasm-go/example -v $(GOPATH)/src/github.com/cosmwasm/cosmwasm-go/std:/go/src/github.com/cosmwasm/cosmwasm-go/std -v $(GOPATH)/src/github.com/cosmwasm/jsonparser:/go/src/github.com/cosmwasm/jsonparser  $(DOCKER_CUSTOM) tinygo build -tags cosmwasm -no-debug -o /code/contract.wasm -target wasm /code/main.go

view:
	@ wasm-nm $(WASM_FILE)
	@ ls -l $(WASM_FILE)

public:
	wasm-nm -e $(WASM_FILE)
	wasm-nm -i $(WASM_FILE)
