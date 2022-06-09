#####  Build container  #####

FROM golang:1.17.8-buster as builder

RUN apt-get update \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        curl \
        protobuf-compiler \
        libprotobuf-dev \
        wget \
        git \
        nano \
    && rm -rf /var/lib/apt/lists/*

# Get go protoc compiler plugins. Taken from: tendermintdev/sdk-proto-gen:v0.2
ENV GOLANG_PROTOBUF_VERSION=1.3.5 \
    GOGO_PROTOBUF_VERSION=1.3.2 \
    GRPC_GATEWAY_VERSION=1.14.7

RUN go get \
    github.com/golang/protobuf/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} \
    github.com/gogo/protobuf/protoc-gen-gogo@v${GOGO_PROTOBUF_VERSION} \
    github.com/gogo/protobuf/protoc-gen-gogofast@v${GOGO_PROTOBUF_VERSION} \
    github.com/gogo/protobuf/protoc-gen-gogofaster@v${GOGO_PROTOBUF_VERSION} \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION} \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v${GRPC_GATEWAY_VERSION} \
    github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest

# Install buf
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

RUN PREFIX="/usr/local" && \
    VERSION="1.0.0-rc8" && \
    curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
    tar -xvzf - -C "${PREFIX}" --strip-components 1

# Fetch and build app
WORKDIR /app

RUN git clone --depth 1 --branch v0.5.0 https://github.com/cheqd/cheqd-node

WORKDIR /app/cheqd-node

RUN make proto-gen build


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=builder /app/cheqd-node/build/cheqd-noded /bin

RUN groupadd --system --gid 1000 cheqd && \
    useradd --system --create-home --home-dir /cheqd --shell /bin/bash --gid cheqd --uid 1000 cheqd
RUN chown -R cheqd /cheqd

WORKDIR /cheqd
USER cheqd

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

# Init network
COPY cheqd_init.sh .
RUN bash cheqd_init.sh

ENTRYPOINT [ "cheqd-noded", "start" ]
