FROM rust:1.84.0-slim-bullseye AS warp-drive-builder

ARG BRANCH=v0.3.3-alpha.1
ARG DOCKER_NAME=ghcr.io/fhenixprotocol/nitro/fhenix-node-builder:$BRANCH
ARG EXTRA_RUSTFLAGS="-C target-feature=+aes"
ENV EXTRA_RUSTFLAGS=$EXTRA_RUSTFLAGS

RUN rustc --version && cargo --version

WORKDIR /workspace

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt update && \
    apt install -y \
        git \
        gpg \
        libstdc++-10-dev \
        make \
        software-properties-common \
        wabt \
        wget \
        zlib1g-dev

RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - && \
    add-apt-repository 'deb http://apt.llvm.org/bullseye/ llvm-toolchain-bullseye-12 main' && \
    apt update && \
    apt install -y llvm-12-dev libclang-common-12-dev

# Copy all the stuff needed to download all the packages
# so we don't have to download everything every time just to change something small
COPY warp-drive/fhe-engine/rust-toolchain warp-drive/fhe-engine/rust-toolchain
COPY warp-drive/fhe-engine/Cargo.toml warp-drive/fhe-engine/Cargo.toml
COPY warp-drive/fhe-engine/Cargo.lock warp-drive/fhe-engine/Cargo.lock
COPY warp-drive/fhe-bridge warp-drive/fhe-bridge
COPY warp-drive/sealing warp-drive/sealing

# Update rust version & install packages
RUN cd warp-drive/fhe-engine && cargo update
#RUN cd warp-drive/renault-server && cargo update

# Copy the rest of the stuff so we can actually build it
COPY warp-drive/ warp-drive/

WORKDIR /workspace/warp-drive/fhe-engine

# Accept build arguments for customizing the build
ARG TARGET_NATIVE=false
ARG TARGET_AVX_512=false

RUN if [ "$TARGET_NATIVE" = "true" ]; then \
        export RUSTFLAGS="-C target-cpu=native"; \
    fi; \
    if [ "$TARGET_AVX_512" = "true" ]; then \
        cargo build --profile=release-lto --features=avx_512; \
    else \
        cargo build --profile=release-lto; \
    fi


FROM debian:bullseye-slim AS runtime
ARG USERNAME=fhenix
RUN useradd -ms /bin/bash ${USERNAME}
WORKDIR /home/${USERNAME}
ENV FHE_ENGINE_CONFIG_DIR="/home/${USERNAME}/fhenix"
RUN mkdir -p ${FHE_ENGINE_CONFIG_DIR}

COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/target/release-lto/fhe-engine-server /usr/bin/fhe-engine-server
COPY --from=warp-drive-builder /workspace/warp-drive/fhe-engine/config/fhe_engine.toml /home/${USERNAME}/fhenix/fhe_engine.toml

CMD ["sh", "-c", "fhe-engine-server -c ${FHE_ENGINE_CONFIG_DIR}/fhe_engine.toml"]
