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

FROM alpine:3.16 AS node

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Set working directory & bash defaults
ENV HOME /home/cheqd
WORKDIR $HOME
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Copy over binaries from the build-env
COPY --from=builder /go/src/github.com/cheqd/cheqd-node/build/cheqd-noded /bin/cheqd-noded

EXPOSE 26656 26657 26660 1317 9090 9091

ENTRYPOINT [ "cheqd-noded start" ]

# Define default command variables
CMD [ "--home " ]
