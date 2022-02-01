#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    protobuf-compiler=3.6.1.3-2 \
    libprotobuf-dev=3.6.1.3-2

# App
WORKDIR /app

RUN git clone --depth 1 --branch v1.4.0 https://github.com/tendermint/liquidity

WORKDIR /app/liquidity

RUN make install


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano=3.2-3 curl=7.64.0-4+deb10u2 wget=1.20.1-1.1 netcat=1.10-41.1 && \
	apt-get clean && \
    rm -rf /var/lib/apt/lists/*

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
