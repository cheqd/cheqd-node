#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    make \
    gcc \
    python \
    protobuf-compiler \
    libprotobuf-dev \
    wget \
    git \
    jq

# App
WORKDIR /app

RUN git clone --depth 1 --branch v4.0.0 https://github.com/osmosis-labs/osmosis

WORKDIR /app/osmosis

RUN make install


#####  Run container  #####

FROM debian:buster

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
