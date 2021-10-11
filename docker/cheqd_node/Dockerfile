#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    # Common
    curl \
    # Protoc
    protobuf-compiler \
    libprotobuf-dev \
    wget

# Starport
# RUN curl https://get.starport.network/starport! | bash
# There is an issue with the latest starport, especially 0.18 version
RUN wget -qO- https://github.com/tendermint/starport/releases/download/v0.17.3/starport_0.17.3_linux_amd64.tar.gz | tar xvz -C /tmp/ && cp /tmp/starport /usr/bin

# App
WORKDIR /app

COPY app ./app
COPY cmd ./cmd
COPY proto ./proto
COPY vue ./vue
COPY x ./x
COPY go.mod .
COPY go.sum .
COPY .git .

RUN starport chain build


#####  Run container  #####

FROM debian:buster

# Node binary
COPY --from=builder /go/bin/cheqd-noded /bin

# Runner script
COPY docker/cheqd_node/node-runner.sh /bin/node-runner
RUN chmod +x /bin/node-runner

# ARG UID=1000
# ARG GID=1000

RUN groupadd --system --gid 1000 cheqd && \
    useradd --system --create-home --home-dir /cheqd --shell /bin/bash --gid cheqd --uid 1000 cheqd
RUN chown -R cheqd /cheqd

WORKDIR /cheqd
USER cheqd

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

ENTRYPOINT [ "cheqd-noded" ]
