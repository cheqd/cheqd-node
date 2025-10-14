###############################################################
###        STAGE 1: Build node and mock binaries           ###
###############################################################

FROM golang:1.23-alpine AS builder

# Install minimum necessary dependencies
ENV PACKAGES="curl make git libc-dev bash gcc linux-headers findutils gettext"
RUN apk update && apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/cheqd/cheqd-node

# Add source files
COPY . .

# Build cheqd node binary
RUN make tidy && make build-linux

# Build mock price feeder binary
WORKDIR /go/src/github.com/cheqd/cheqd-node/tests/integration/mock-price-feed
RUN go mod tidy && go build -o mock-price-feed


###############################################################
###             STAGE 2: Build cheqd-node image             ###
###############################################################

FROM alpine:3.21 AS runner

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates gettext

# Copy over binaries from the build-env
COPY --from=builder /go/src/github.com/cheqd/cheqd-node/build/cheqd-noded /bin/cheqd-noded

# Set user directory and details
ARG HOME_DIR="/home/cheqd"
ARG USER="cheqd"
ARG GROUP=${USER}
ARG UID=1000
ARG GID=${UID}

# Add cheqd user to use in the container
RUN addgroup --system ${GROUP} --gid ${GID} \
    && adduser ${USER} --uid ${UID} -G ${GROUP} --system --home ${HOME_DIR} --shell /bin/bash

# Copy mock price feeder binary from builder
COPY --from=builder /go/src/github.com/cheqd/cheqd-node/tests/integration/mock-price-feed/mock-price-feed /home/cheqd/mock-price-feed/mock-price-feed
RUN chmod +x /home/cheqd/mock-price-feed/mock-price-feed

# Set entrypoint script
COPY --chown=${USER}:${GROUP} docker/entrypoint.sh /bin/node-start
RUN chmod +rx /bin/node-start

# Set working directory & bash defaults
WORKDIR ${HOME_DIR}
USER ${USER}
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Expose relevant ports
EXPOSE 26656 26657 26660 1317 9090 9091 7171

# Set plain vanilla default entrypoint/command
CMD [ "cheqd-noded" ]