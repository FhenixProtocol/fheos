ARG BRANCH=latest
ARG DOCKER_NAME=ghcr.io/fhenixprotocol/nitro/fhenix-node-builder:$BRANCH

FROM rust:1.74-slim-bullseye as warp-drive-builder
WORKDIR /workspace
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y make wget gpg software-properties-common zlib1g-dev libstdc++-10-dev wabt git
RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - && \
    add-apt-repository 'deb http://apt.llvm.org/bullseye/ llvm-toolchain-bullseye-12 main' && \
    apt-get update && \
    apt-get install -y llvm-12-dev libclang-common-12-dev

COPY warp-drive/ warp-drive/
WORKDIR /workspace/warp-drive/fhe-engine

RUN RUSTFLAGS="-C target-cpu=native" cargo build --release

FROM $DOCKER_NAME as winning

RUN echo $DOCKER_NAME

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update -qq && apt-get install -y nodejs npm yarn

RUN npm install -g pnpm

RUN rm -rf fheos/
COPY . fheos/

WORKDIR fheos

RUN ./gen.sh

WORKDIR /workspace

COPY nitro-overrides/precompiles/FheOps.go precompiles/FheOps.go

RUN go mod tidy

RUN go build -gcflags "all=-N -l" -ldflags="-X github.com/offchainlabs/nitro/cmd/util/confighelpers.version= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.datetime= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.modified=" -o target/bin/nitro "/workspace/cmd/nitro"

RUN cd fheos && make build

FROM ghcr.io/fhenixprotocol/localfhenix:v0.1.0-beta5

RUN rm -rf /home/user/fhenix/fheosdb
RUN mkdir -p /home/user/fhenix/fheosdb

COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/target/release/fhe-engine-server /usr/bin/fhe-engine-server
COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/config/fhe_engine.toml /home/user/fhenix/fhe_engine.toml

COPY --from=winning /workspace/target/bin/nitro /usr/local/bin/
COPY --from=winning /workspace/fheos/build/main /usr/local/bin/fheos

RUN fheos init-db
RUN fheos generate-keys




