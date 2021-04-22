FROM golang:buster

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

VOLUME /root/.verim-cosmosd
EXPOSE 26656 26657

STOPSIGNAL SIGTERM
