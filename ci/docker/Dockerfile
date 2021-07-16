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

COPY --from=builder /go/bin/verim-noded /bin

COPY ci/docker/node-runner.sh /bin/node-runner
RUN chmod +x /bin/node-runner

RUN groupadd --gid 1000 verim && \
    useradd --create-home --shell /bin/bash --uid 1000 --gid verim verim

WORKDIR /home/verim
USER verim

VOLUME /home/verim/.verimnode
EXPOSE 26656 26657
STOPSIGNAL SIGTERM

ENTRYPOINT [ "verim-noded" ]
