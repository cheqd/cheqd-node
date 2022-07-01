#####  Build container  #####

# Taken from: https://github.com/osmosis-labs/osmosis/blob/v10.0.1/Dockerfile
FROM golang:1.18.2-alpine3.15 as build

RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add git
# Needed by github.com/zondax/hid
RUN apk add linux-headers

WORKDIR /
RUN git clone --depth 1 --branch v8.0.0 https://github.com/osmosis-labs/osmosis
WORKDIR /osmosis

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN set -o pipefail; sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 7d2239e9f25e96d0d4daba982ce92367aacf0cbd95d2facb8442268f2b1cc1fc
RUN set -o pipefail; sha256sum /lib/libwasmvm_muslc.x86_64.a | grep f6282df732a13dec836cda1f399dd874b1e3163504dbd9607c6af915b2740479

# CosmWasm: copy the right library according to architecture. The final location will be found by the linker flag `-lwasmvm_muslc`
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano \
    curl \
    wget \
    netcat

# Node binary
COPY --from=build /osmosis/build/osmosisd /bin/osmosisd

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
