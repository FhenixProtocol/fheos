ARG BRANCH=v0.2.0-alpha.0
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

ARG EXTRA_RUSTFLAGS="-C target-feature=+aes"
ENV EXTRA_RUSTFLAGS=$EXTRA_RUSTFLAGS

# Copy all the stuff needed to download all the packages
# so we don't have to download everything every time just to change something small
COPY warp-drive/fhe-engine/rust-toolchain warp-drive/fhe-engine/rust-toolchain
COPY warp-drive/fhe-engine/Cargo.toml warp-drive/fhe-engine/Cargo.toml
COPY warp-drive/fhe-engine/Cargo.lock warp-drive/fhe-engine/Cargo.lock
COPY warp-drive/fhe-bridge warp-drive/fhe-bridge

# Update rust version & install packages
RUN cd warp-drive/fhe-engine && cargo update

# Copy the rest of the stuff so we can actually build it
COPY warp-drive/ warp-drive/

WORKDIR /workspace/warp-drive/fhe-engine

# Todo: fix arm support
RUN RUSTFLAGS=$EXTRA_RUSTFLAGS cargo build --profile=release-lto

FROM $DOCKER_NAME as winning

RUN echo $DOCKER_NAME

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update -qq && apt-get install -y nodejs npm yarn

RUN npm install -g pnpm

RUN rm -rf fheos/
RUN mkdir fheos/

COPY warp-drive fheos/warp-drive
COPY go-ethereum fheos/go-ethereum
COPY chains fheos/chains
COPY cmd fheos/cmd
COPY precompiles fheos/precompiles
COPY storage fheos/storage
COPY hooks fheos/hooks
RUN cd fheos/precompiles/ && pnpm install

COPY gen.sh fheos/
COPY gen.go fheos/

COPY go.mod fheos/
COPY go.sum fheos/

RUN mkdir -p fheos/solidity

WORKDIR fheos

RUN ./gen.sh

WORKDIR /workspace

COPY nitro-overrides/precompiles/FheOps.go precompiles/FheOps.go

RUN go mod tidy

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

RUN go build -gcflags "all=-N -l" -ldflags="-X github.com/offchainlabs/nitro/cmd/util/confighelpers.version= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.datetime= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.modified=" -o target/bin/nitro "/workspace/cmd/nitro"

COPY Makefile fheos/

RUN cd fheos && make build

FROM ghcr.io/fhenixprotocol/localfhenix:v0.1.0-beta5

# **************** setup dlv

ENV GOROOT=/usr/local/go
ENV GOPATH=/go/
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

ADD https://go.dev/dl/go1.20.linux-amd64.tar.gz go.linux-amd64.tar.gz
RUN sudo tar -C /usr/local -xzf go.linux-amd64.tar.gz

RUN sudo chown user:user /usr/local/go

RUN sudo mkdir /go
RUN sudo chown user:user /go

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

# **************** setup fheos & warp drive

RUN rm -rf /home/user/fhenix/fheosdb
RUN mkdir -p /home/user/fhenix/fheosdb

COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/target/release-lto/fhe-engine-server /usr/bin/fhe-engine-server
COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/config/fhe_engine.toml /home/user/fhenix/fhe_engine.toml

COPY --from=winning /workspace/target/bin/nitro /usr/local/bin/
COPY --from=winning /workspace/fheos/build/main /usr/local/bin/fheos

# **************** setup scripts and configs

COPY deployment/run.sh run.sh

RUN sudo chmod +x ./run.sh
RUN sudo chown -R user:user /home/user/keys
RUN sed -i '/^keys_folder *=.*/s//keys_folder = "\/home\/user\/keys" /' /home/user/fhenix/fhe_engine.toml
RUN sudo jq '.conf |= (.fhenix = .tfhe | del(.tfhe))' /config/sequencer_config.json > temp.json && sudo mv temp.json /config/sequencer_config.json

# **************** Run

CMD ["./run.sh"]
#CMD ["./run.sh", "--debug"]