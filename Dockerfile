FROM cosmwasm/wabt:1.0.25 as wabt

FROM cosmwasm/tinygo:0.19.3 as tinygo

COPY --from=wabt /usr/local/bin/wasm2wat /usr/local/bin/wasm2wat
COPY --from=wabt /usr/local/bin/wat2wasm /usr/local/bin/wat2wasm
COPY --from=wabt /usr/local/wabt /usr/local/wabt

COPY docker/compile.sh /usr/local/bin/compile.sh

RUN mkdir /work

# TODO copy more over??

WORKDIR /code
ENTRYPOINT [ "compile.sh" ]
