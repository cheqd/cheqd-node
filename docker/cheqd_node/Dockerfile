#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    libprotobuf-dev \
    wget


# Go protoc programs. Taken from: tendermintdev/sdk-proto-gen:v0.2
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

# Node binary
COPY --from=builder /app/build/cheqd-noded /bin

# Runner script
COPY docker/cheqd_node/node-runner.sh /bin/node-runner
RUN chmod +x /bin/node-runner

ARG CHEQD_HOME_DIR="/cheqd"

ARG UID=1000
ARG GID=1000

RUN groupadd --system --gid $GID cheqd && \
    useradd --system --create-home --home-dir ${CHEQD_HOME_DIR} --shell /bin/bash --gid cheqd --uid $UID cheqd

WORKDIR ${CHEQD_HOME_DIR}
USER cheqd

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

ENTRYPOINT [ "cheqd-noded" ]
