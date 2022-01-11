#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    libprotobuf-dev \
    wget


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
RUN PREFIX="/usr/local" && \
    VERSION="1.0.0-rc8" && \
    curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
    tar -xvzf - -C "${PREFIX}" --strip-components 1


# App
WORKDIR /app

COPY app ./app
COPY cmd ./cmd
COPY scripts ./scripts
COPY proto ./proto
COPY x ./x
COPY go.mod .
COPY go.sum .
COPY Makefile .
# Required to fetch version
COPY .git .

RUN make proto-gen build

#####  Run container  #####

FROM ubuntu:focal

RUN apt-get update \
	&& export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends curl ca-certificates

# Node binary
COPY --from=builder /app/build/cheqd-noded /bin
# Runner script
COPY docker/node-runner.sh /bin/node-runner
RUN chmod +x /bin/node-runner

ARG CHEQD_HOME_DIR="/home/cheqd"
ARG UID=1000
ARG GID=1000

RUN groupadd --system --gid $GID cheqd \
    && useradd --system --create-home --home-dir ${CHEQD_HOME_DIR} --shell /bin/bash --gid cheqd --uid $UID cheqd

WORKDIR ${CHEQD_HOME_DIR}
USER cheqd

EXPOSE 26656 26657 1317 9090 9091 26660
STOPSIGNAL SIGTERM

ENTRYPOINT [ "cheqd-noded" ]
