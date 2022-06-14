FROM cosmwasm/wabt:1.0.25 as wabt

#FROM cosmwasm/tinygo:0.19.3 as tinygo

# Using the newest CosmWasm TinyGo v0.23.0 beta.
# Docker image should be build locally for now. More info and how-tos:
# https://github.com/CosmWasm/tinygo/blob/cw-0.23.x/COSMWASM.md
FROM cosmwasm/tinygo:0.23.0 as tinygo

COPY --from=wabt /usr/local/bin/wasm2wat /usr/local/bin/wasm2wat
COPY --from=wabt /usr/local/bin/wat2wasm /usr/local/bin/wat2wasm
COPY --from=wabt /usr/local/wabt /usr/local/wabt

COPY docker/compile.sh /usr/local/bin/compile.sh

RUN mkdir /work

# TODO copy more over??

WORKDIR /code
ENTRYPOINT [ "compile.sh" ]
