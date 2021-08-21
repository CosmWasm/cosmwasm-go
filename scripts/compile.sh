#!/bin/bash

set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

TINYGO_IMAGE="cosmwasm/tinygo:0.19.0-dev"
EMSCRIPTEN="polkasource/webassembly-wabt:v1.0.11"

SCRIPT_DIR="$(realpath "$(dirname "$0")")"
ROOT="$(dirname "$SCRIPT_DIR")"

function print_usage() {
  echo "Usage: $0 [contract_name]"
  echo ""
  echo "e.g. $0 erc20"
  echo "This must be a valid directory name under example/"
}

if [ "$#" -ne 1 ]; then
    print_usage
    exit 1
fi

CONTRACT="$1"
DIR="${ROOT}/example/${CONTRACT}"

if [ ! -d "$DIR" ]; then
    print_usage
    exit 1
fi

echo "Compiling $CONTRACT with tinygo..."
docker run --rm -w /code -v "${ROOT}:/code" ${TINYGO_IMAGE} tinygo build -tags cosmwasm -no-debug -target wasm -o "/code/${CONTRACT}.wasm" "/code/example/${CONTRACT}/main.go"
ls -l "${ROOT}/${CONTRACT}.wasm"

echo "Stripping out floating point symbols..."
WATFILE="${ROOT}/${CONTRACT}.wat"
docker run --rm -v "${ROOT}:/code" ${EMSCRIPTEN} wasm2wat "/code/${CONTRACT}.wasm" > "${WATFILE}"

# this just replaces all the floating point ops with unreachable. It still leaves them in the args and local variables
sed -E 's/^(\s*)f[[:digit:]]{2}\.[^()]+/\1unreachable/' "${WATFILE}" | \
  sed -E 's/^(\s*)i[[:digit:]]{2}\.trunc_[^()]+/\1unreachable/' | \
  sed -E 's/^(\s*)i[[:digit:]]{2}\.reinterpret_[^()]+/\1unreachable/' > "${WATFILE}-rewrite"
mv "${WATFILE}-rewrite" "${WATFILE}"

docker run --rm -w /code -v "${ROOT}:/code" ${EMSCRIPTEN} wat2wasm "/code/${CONTRACT}.wat"

echo "Done! ${CONTRACT}.wasm is ready to use."
ls -l "${ROOT}/${CONTRACT}.wasm"
