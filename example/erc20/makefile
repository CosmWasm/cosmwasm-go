.PHONY: build view all build-docker build-custom

DOCKER_IMAGE=tinygo/tinygo:0.13.1
DOCKER_CUSTOM=cosmwasm/tinygo:latest
WASM_FILE=contract.wasm

PACKAGE_NAME=erc20
GO_MOUNTS=-v $(shell pwd):/code -v $(shell pwd):/go/src/github.com/cosmwasm/cosmwasm-go/example/$(PACKAGE_NAME) -v $(shell pwd)/../../std:/go/src/github.com/cosmwasm/cosmwasm-go/std -v $(GOPATH)/src/github.com/cosmwasm/jsonparser:/go/src/github.com/cosmwasm/jsonparser

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

# build-cosmwasm uses -no-debug flag for smaller executables.
# also, tinygo build currently error without that flag due to unresolvable imports
build-cosmwasm:
    # TODO: automatically download jsonparser dependency
	docker run --rm $(GO_MOUNTS) $(DOCKER_CUSTOM) tinygo build -tags cosmwasm -no-debug -o /code/contract.wasm -target wasm /code/main.go

view:
	@ wasm-nm $(WASM_FILE)
	@ ls -l $(WASM_FILE)

public:
	wasm-nm -e $(WASM_FILE)
	wasm-nm -i $(WASM_FILE)
