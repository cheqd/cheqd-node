#####  Build container  #####

FROM rust:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    protobuf-compiler=3.6.1.3-2 \
    libprotobuf-dev=3.6.1.3-2

WORKDIR /app

RUN git clone --depth 1 --branch v0.9.0 https://github.com/informalsystems/ibc-rs

WORKDIR /app/ibc-rs

RUN cargo build --release --bin hermes


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    libssl-dev=1.1.1d-0+deb10u7 nano=3.2-3 curl=7.64.0-4+deb10u2 wget=1.20.1-1.1 netcat=1.10-41.1 && \
	apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Node binary
COPY --from=builder /app/ibc-rs/target/release/hermes /bin

# User
RUN groupadd --system --gid 1000 hermes && \
    useradd --system --create-home --home-dir /hermes --shell /bin/bash --gid hermes --uid 1000 hermes

WORKDIR /hermes

# Init
COPY hermes_init.sh .
# Config
COPY config.toml .hermes

RUN bash hermes_init.sh && \
	mkdir .hermes && \
	chown -R hermes /hermes
USER hermes

ENTRYPOINT [ "hermes" ]
