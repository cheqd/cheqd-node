###############################################################
###        STAGE 1: Build node binary pre-requisites        ###
###############################################################

FROM golang:1.17-alpine AS builder

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers
RUN apk update && apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/cheqd/cheqd-node

# Add source files
COPY . .

# Make node binary
RUN make install && make build

###############################################################
###             STAGE 2: Build cheqd-node image             ###
###############################################################

FROM alpine:3.16 AS runner

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Copy over binaries from the build-env
COPY --from=builder /go/src/github.com/cheqd/cheqd-node/build/cheqd-noded /bin/cheqd-noded

# Set user directory and details
ARG HOME_DIR="/home/cheqd"
ARG USER="cheqd"

# Add cheqd user to use in the container
RUN addgroup --system ${USER} \
    && adduser ${USER} --system --home ${HOME_DIR} --shell /bin/bash

# Set working directory & bash defaults
WORKDIR ${HOME_DIR}
USER ${USER}
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Set runner script
COPY --chown=cheqd:cheqd ./docker/entrypoint.sh /bin/node-start
RUN chmod +x /bin/node-start

EXPOSE 26656 26657 26660 1317 9090 9091

# Set plain vanilla default entrypoint/command
CMD [ "cheqd-noded" ]
