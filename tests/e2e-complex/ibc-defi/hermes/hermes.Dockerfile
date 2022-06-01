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

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    libssl-dev \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=builder /app/ibc-rs/target/release/hermes /bin

# User
RUN groupadd --system --gid 1000 hermes && \
    useradd --system --create-home --home-dir /hermes --shell /bin/bash --gid hermes --uid 1000 hermes

WORKDIR /hermes

# Init
COPY hermes_init.sh .
RUN bash hermes_init.sh

# Config
RUN mkdir .hermes
COPY config.toml .hermes

RUN chown -R hermes /hermes
USER hermes

ENTRYPOINT [ "hermes" ]
