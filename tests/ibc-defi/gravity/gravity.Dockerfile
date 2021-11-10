#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    protobuf-compiler \
    libprotobuf-dev \
    wget \
    git \
    nano

# Starport
# RUN curl https://get.starport.network/starport! | bash
# There is an issue with the latest starport, especially 0.18 version
RUN wget -qO- https://github.com/tendermint/starport/releases/download/v0.17.3/starport_0.17.3_linux_amd64.tar.gz | tar xvz -C /tmp/ && cp /tmp/starport /usr/bin

# App
WORKDIR /app

RUN git clone --depth 1 --branch v1.4.0 https://github.com/tendermint/liquidity

WORKDIR /app/liquidity

RUN make install


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=builder /go/bin/liquidityd /bin

RUN groupadd --system --gid 1000 gravity && \
    useradd --system --create-home --home-dir /gravity --shell /bin/bash --gid gravity --uid 1000 gravity
RUN chown -R gravity /gravity

WORKDIR /gravity
USER gravity

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

# Init network
COPY gravity_init.sh .
RUN bash gravity_init.sh

ENTRYPOINT [ "liquidityd", "start" ]
