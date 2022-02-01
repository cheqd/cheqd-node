#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    protobuf-compiler=3.6.1.3-2 \
    libprotobuf-dev=3.6.1.3-2 \
    jq=1.5+dfsg-2+b1

# App
WORKDIR /app

RUN git clone --depth 1 --branch v4.2.0 https://github.com/osmosis-labs/osmosis

WORKDIR /app/osmosis

RUN make install


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano=3.2-3 curl=7.64.0-4+deb10u2 wget=1.20.1-1.1 netcat=1.10-41.1 && \
	apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Node binary
COPY --from=builder /go/bin/osmosisd /bin

RUN groupadd --system --gid 1000 osmosis && \
    useradd --system --create-home --home-dir /osmosis --shell /bin/bash --gid osmosis --uid 1000 osmosis
RUN chown -R osmosis /osmosis

WORKDIR /osmosis
USER osmosis

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

# Init network
COPY osmosis_init.sh .
RUN bash osmosis_init.sh

ENTRYPOINT [ "osmosisd", "start" ]
