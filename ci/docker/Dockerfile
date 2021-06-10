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
COPY . .
RUN starport build


FROM debian:buster

COPY --from=builder /go/bin/verim-noded /root/bin

VOLUME /root/.verim-noded
EXPOSE 26656 26657
STOPSIGNAL SIGTERM

ENTRYPOINT [ "verim-noded" ]
