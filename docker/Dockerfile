###############################################################
###        STAGE 1: Build node binary pre-requisites        ###
###############################################################

FROM golang:1.17.8-buster as builder

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy pre-requisites before building the node binary
WORKDIR /app

COPY app ./app
COPY cmd ./cmd
COPY scripts ./scripts
COPY proto ./proto
COPY x ./x
COPY go.mod .
COPY go.sum .
COPY Makefile .
# Required to fetch version
COPY .git .

# Make node binary
RUN make build

###############################################################
###      STAGE 2: Build cheqd binary base container         ###
###############################################################

FROM ubuntu:focal AS base

LABEL org.opencontainers.image.description "cheqd CLI Docker image"
LABEL org.opencontainers.image.source "https://github.com/cheqd/cheqd-node"
LABEL org.opencontainers.image.documentation "https://docs.cheqd.io/node"

# Copy compiled node binary from Stage 1
COPY --from=builder /app/build-tools/cheqd-noded /bin

# Set user directory and details
ARG CHEQD_HOME_DIR="/home/cheqd"
ARG UID=1000
ARG GID=1000

# Add cheqd user to use in the container
RUN groupadd --system --gid $GID cheqd \
    && useradd --system --create-home --home-dir ${CHEQD_HOME_DIR} --shell /bin/bash --gid cheqd --uid $UID cheqd

WORKDIR ${CHEQD_HOME_DIR}
USER cheqd

# Document default ports to expose to host
EXPOSE 26656 26657 26660 1317 9090 9091

# Define stop scenarios
STOPSIGNAL SIGTERM

# Default entrypoint for cheqd-noded CLI usage
ENTRYPOINT [ "cheqd-noded" ]


###############################################################
###             STAGE 3: Build cheqd-node image             ###
###############################################################

FROM base AS node

LABEL org.opencontainers.image.description "cheqd Node Docker image"
LABEL org.opencontainers.image.source "https://github.com/cheqd/cheqd-node"
LABEL org.opencontainers.image.documentation "https://docs.cheqd.io/node"

# Set runner script
COPY --chown=cheqd:cheqd docker/entrypoint.sh /bin/node-start
RUN chmod +x /bin/node-start

# Default entrypoint for cheqd-noded CLI usage
ENTRYPOINT [ "node-start" ]


###############################################################
###        STAGE 4: Build Cosmovisor                        ###
###############################################################

FROM golang:1.17.8-buster AS cosmos_builder 

RUN git clone https://github.com/cosmos/cosmos-sdk.git

WORKDIR /go/cosmos-sdk/

RUN git status \
    && git checkout cosmovisor/v1.1.0 \
    && make cosmovisor


###############################################################
###          STAGE 5: Cosmovisor-based node image           ###
###############################################################

FROM base AS cosmovisor

COPY --from=cosmos_builder /go/cosmos-sdk/cosmovisor/cosmovisor /bin
COPY --chown=cheqd:cheqd docker/cosmovisor.sh /bin/cosmovisor.sh

ARG CHEQD_HOME_DIR="/home/cheqd"

RUN chmod +x /bin/cosmovisor.sh

USER cheqd

ENV DAEMON_HOME=${CHEQD_HOME_DIR}/.cheqdnode
ENV DAEMON_NAME=cheqd-noded
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=true
ENV DAEMON_RESTART_AFTER_UPGRADE=true

ENTRYPOINT [ "cosmovisor.sh" ]
