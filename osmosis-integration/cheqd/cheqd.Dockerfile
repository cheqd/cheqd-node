#####  Build container  #####

FROM golang:buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    protobuf-compiler \
    libprotobuf-dev \
    wget \
    git

# Starport
# RUN curl https://get.starport.network/starport! | bash
# There is an issue with the latest starport, especially 0.18 version
RUN wget -qO- https://github.com/tendermint/starport/releases/download/v0.17.3/starport_0.17.3_linux_amd64.tar.gz | tar xvz -C /tmp/ && cp /tmp/starport /usr/bin

# App
WORKDIR /app

RUN git clone --depth 1 --branch v0.2.3 https://github.com/cheqd/cheqd-node

WORKDIR /app/cheqd-node

RUN starport chain build


#####  Run container  #####

FROM debian:buster

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
