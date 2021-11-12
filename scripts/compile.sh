#!/bin/bash

set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

TINYGO_IMAGE="cosmwasm/tinygo:0.20.0"
EMSCRIPTEN="polkasource/webassembly-wabt:v1.0.11"

SCRIPT_DIR="$(realpath "$(dirname "$0")")"
ROOT="$(dirname "$SCRIPT_DIR")"

function print_usage() {
  echo "Usage: $0 [contract_name]"
  echo ""
  echo "e.g. $0 hackatom"
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
docker run --rm -w /code -v "${ROOT}:/code" ${TINYGO_IMAGE} tinygo build -tags "cosmwasm tinyjson_nounsafe" -no-debug -target wasi -o "/code/${CONTRACT}.wasm" "/code/example/${CONTRACT}/main.go"
ls -l "${ROOT}/${CONTRACT}.wasm"

# FIXME: we can remove this whole EMSCRIPTEN stuff in the future... or just as an optional verification check (outside of compile)

WATFILE="${ROOT}/${CONTRACT}.wat"
docker run --rm -v "${ROOT}:/code" ${EMSCRIPTEN} wasm2wat "/code/${CONTRACT}.wasm" > "${WATFILE}"

ls -l "${ROOT}/${CONTRACT}.wasm"
grep import "${WATFILE}"

echo "Any floating point?"
grep f64 "${WATFILE}" || true

# echo "Stripping out floating point symbols..."
# # this just replaces all the floating point ops with unreachable. It still leaves them in the args and local variables
# sed -E 's/^(\s*)f[[:digit:]]{2}\.[^()]+/\1unreachable/' "${WATFILE}" | \
#   sed -E 's/^(\s*)i[[:digit:]]{2}\.trunc_[^()]+/\1unreachable/' | \
#   sed -E 's/^(\s*)i[[:digit:]]{2}\.reinterpret_[^()]+/\1unreachable/' > "${WATFILE}-rewrite"
# mv "${WATFILE}-rewrite" "${WATFILE}"

docker run --rm -w /code -v "${ROOT}:/code" ${EMSCRIPTEN} wat2wasm "/code/${CONTRACT}.wat"

echo "Done! ${CONTRACT}.wasm is ready to use."
ls -l "${ROOT}/${CONTRACT}.wasm"
sha256sum "${ROOT}/${CONTRACT}.wasm"
