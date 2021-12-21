# FROM phusion/baseimage:0.11 as wabt
FROM phusion/baseimage:focal-1.1.0 as wabt

RUN apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y cmake pkg-config libssl-dev git clang libclang-dev git

RUN git clone https://github.com/WebAssembly/wabt.git && \
   cd wabt && git checkout 1.0.25 && cd -
    # cd wabt && git checkout 0af114943e38a0c0a4ccb0b49b4a8fb07d1bd056 && cd -

RUN cd /cbuilder/wabt && mkdir build
RUN cd /cbuilder/wabt/build && cmake ..
RUN cd /cbuilder/wabt/build && cmake --build .