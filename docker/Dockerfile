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
RUN make build-linux

###############################################################
###             STAGE 3: Build cheqd-node image             ###
###############################################################

FROM builder AS node

# Set runner script
COPY --chown=cheqd:cheqd docker/entrypoint.sh /bin/node-start
RUN chmod +x /bin/node-start

# Default entrypoint for cheqd-noded CLI usage
ENTRYPOINT [ "node-start" ]
