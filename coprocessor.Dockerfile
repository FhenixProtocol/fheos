ARG BRANCH=v0.3.3-alpha.1
ARG DOCKER_NAME=ghcr.io/fhenixprotocol/nitro/fhenix-node-builder:$BRANCH

FROM $DOCKER_NAME AS winning

RUN echo $DOCKER_NAME
WORKDIR /workspace

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt update -qq && apt install -y nodejs npm yarn

RUN npm install -g pnpm

RUN rm -rf fheos/
RUN mkdir fheos/

COPY warp-drive fheos/warp-drive
COPY http fheos/http
COPY go-ethereum fheos/go-ethereum
COPY chains fheos/chains
COPY cmd fheos/cmd
COPY precompiles fheos/precompiles
COPY storage fheos/storage
COPY hooks fheos/hooks
COPY nitro-overrides/precompiles/FheOps.go ./precompiles/FheOps.go
COPY gen.sh gen.go go.mod go.sum Makefile fheos/

RUN cd fheos/precompiles/ && pnpm install
RUN mkdir -p fheos/solidity

RUN go mod tidy

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN go build -gcflags="all=-N -l" -ldflags="-X github.com/offchainlabs/nitro/cmd/util/confighelpers.version= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.datetime= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.modified=" -o target/bin/nitro "/workspace/cmd/nitro"

RUN make -C fheos build-coprocessor


FROM debian:bookworm-slim AS runtime

ARG USERNAME=fhenix
RUN useradd -ms /bin/bash ${USERNAME}
WORKDIR /home/${USERNAME}
RUN mkdir -p /home/${USERNAME}/fhenix/fheosdb

COPY deployment/sequencer_config.json /config/sequencer_config.json

COPY --from=winning /workspace/target/bin/nitro /usr/local/bin/
COPY --from=winning /workspace/fheos/build/coprocessor /usr/local/bin/coprocessor

CMD ["sh", "-c", "coprocessor"]
