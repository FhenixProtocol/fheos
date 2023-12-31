FROM ghcr.io/fhenixprotocol/nitro/fhenix-node-builder:latest as winning

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update -qq && apt-get install -y nodejs npm yarn

RUN npm install -g pnpm

COPY . fheos/

WORKDIR fheos

RUN ./gen.sh

WORKDIR /workspace

COPY nitro-overrides/precompiles/FheOps.go precompiles/FheOps.go

RUN go build -gcflags "all=-N -l" -ldflags="-X github.com/offchainlabs/nitro/cmd/util/confighelpers.version= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.datetime= -X github.com/offchainlabs/nitro/cmd/util/confighelpers.modified=" -o target/bin/nitro "/workspace/cmd/nitro"

FROM ghcr.io/fhenixprotocol/fhenix-node-dev:v0.0.7-standalone

COPY --from=winning /workspace/fheos/go-tfhe/internal/api/amd64/libtfhe_wrapper.x86_64.so /usr/lib/libtfhe_wrapper.x86_64.so
COPY --from=winning /workspace/target/bin/nitro /usr/local/bin/
