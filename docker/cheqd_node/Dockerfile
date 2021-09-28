#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    # Common
    curl \
    # Protoc
    protobuf-compiler \
    libprotobuf-dev

# Starport
RUN curl https://get.starport.network/starport! | bash

# App
WORKDIR /app

COPY app ./app
COPY cmd ./cmd
COPY proto ./proto
COPY vue ./vue
COPY x ./x
COPY go.mod .
COPY go.sum .

RUN starport chain build


#####  Run container  #####

FROM debian:buster

# Node binary
COPY --from=builder /go/bin/cheqd-noded /bin

# Runner script
COPY docker/cheqd_node/node-runner.sh /bin/node-runner
RUN chmod +x /bin/node-runner

RUN groupadd --system --gid 1000 cheqd && \
    useradd --system --create-home --home-dir /cheqd --shell /bin/bash --gid cheqd --uid 1000 cheqd
RUN chown -R cheqd /cheqd

WORKDIR /cheqd
USER cheqd

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

ENTRYPOINT [ "cheqd-noded" ]
