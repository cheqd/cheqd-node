#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    protobuf-compiler=3.6.1.3-2 \
    libprotobuf-dev=3.6.1.3-2


SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# Starport
# RUN curl https://get.starport.network/starport! | bash
# There is an issue with the latest starport, especially 0.18 version
RUN wget -qO- https://github.com/tendermint/starport/releases/download/v0.18.6/starport_0.18.6_linux_amd64.tar.gz | tar xvz -C /tmp/ && cp /tmp/starport /usr/bin

# App
WORKDIR /app

RUN git clone --depth 1 --branch v0.3.1 https://github.com/cheqd/cheqd-node

WORKDIR /app/cheqd-node

RUN starport chain build


#####  Run container  #####

FROM debian:buster

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    nano=3.2-3 curl=7.64.0-4+deb10u2 wget=1.20.1-1.1 netcat=1.10-41.1 && \
	apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Node binary
COPY --from=builder /go/bin/cheqd-noded /bin

RUN groupadd --system --gid 1000 cheqd && \
    useradd --system --create-home --home-dir /cheqd --shell /bin/bash --gid cheqd --uid 1000 cheqd
RUN chown -R cheqd /cheqd

WORKDIR /cheqd
USER cheqd

EXPOSE 26656 26657
STOPSIGNAL SIGTERM

# Init network
COPY cheqd_init.sh .
RUN bash cheqd_init.sh

ENTRYPOINT [ "cheqd-noded", "start" ]
