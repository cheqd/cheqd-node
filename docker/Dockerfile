###############################################################
###        STAGE 1: Build node binary pre-requisites        ###
###############################################################

FROM golang:1.17-alpine AS builder

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk update && apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/cheqd/cheqd-node

# Add source files
COPY . .

# Make node binary
RUN make build-linux

###############################################################
###      STAGE 2: Build cheqd binary base container         ###
###############################################################

FROM alpine:3.16 AS base

LABEL org.opencontainers.image.description "cheqd CLI Docker image"
LABEL org.opencontainers.image.source "https://github.com/cheqd/cheqd-node"
LABEL org.opencontainers.image.documentation "https://docs.cheqd.io/node"

# Copy compiled node binary from Stage 1
COPY --from=builder /app/build-tools/cheqd-noded /bin

# Set user directory and details
ARG CHEQD_HOME_DIR="/home/cheqd"
ARG UID=1000
ARG GID=1000

# Install pre-requisites
RUN apk update && apk add --no-cache \
    bash \
    ca-certificates

# Add cheqd user to use in the container
RUN addgroup -S -g $GID cheqd \
    && adduser -S -h ${CHEQD_HOME_DIR} -s /bin/bash -G cheqd -u $UID cheqd

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

FROM golang:1.17-alpine AS cosmos_builder

RUN apk update && apk add --no-cache git

RUN git clone https://github.com/cosmos/cosmos-sdk.git

RUN apk add \
    bash \
    make \
    gcc \
    libc-dev \
    linux-headers

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
