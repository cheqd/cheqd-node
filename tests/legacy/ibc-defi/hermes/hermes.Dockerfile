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

RUN git clone https://github.com/informalsystems/hermes

WORKDIR /app/hermes

RUN git checkout v1.2.0
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
COPY --from=builder /app/hermes/target/release/hermes /bin

ARG USER=hermes
ARG GROUP=hermes

ARG HOME=/home/$USER

# User
RUN groupadd --system --gid 1000 $USER && \
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid 1000 $USER

WORKDIR $HOME

RUN chown -R $USER $HOME
USER $USER

ENTRYPOINT [ "hermes" ]
