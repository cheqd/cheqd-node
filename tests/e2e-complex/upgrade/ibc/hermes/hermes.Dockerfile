#####  Build container  #####

FROM rust:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    protobuf-compiler \
    libprotobuf-dev \
    wget \
    git

WORKDIR /app

RUN git clone --depth 1 --branch v0.15.0 https://github.com/informalsystems/ibc-rs

WORKDIR /app/ibc-rs

RUN cargo build --release --bin hermes


#####  Run container  #####

FROM ubuntu:22.04

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    libssl-dev \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=builder /app/ibc-rs/target/release/hermes /bin

ARG UID=1000
ARG GID=1000

ARG USER=hermes
ARG GROUP=hermes

ARG HOME=/home/$USER

# User
RUN groupadd --system --gid $GID $USER && \
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid $UID $USER

WORKDIR $HOME

# Permissions fix for docker configs
RUN mkdir -p $HOME/.hermes

RUN chown -R $USER $HOME
USER $USER

ENTRYPOINT [ "hermes" ]
