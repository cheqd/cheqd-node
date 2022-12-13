###############################################################
###         STAGE 1: Build network generation image         ###
###############################################################

FROM alpine:3.17 AS networkgen

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Set bash defaults
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Set build-time arguments
ARG RELEASE_URL

# Copy and extract cheqd-node binary
RUN wget -qO- $RELEASE_URL | tar -xz -C /bin --strip-components=2

# Copy scripts
COPY gen-network-config.sh /bin/generate-network

# Make scripts executable
RUN chmod +x /bin/cheqd-noded /bin/generate-network

# Set user directory and details
# USER cheqd
# WORKDIR /home/cheqd

# Set network generation environment variables
ENV CHAIN_ID="cheqd"
ENV VALIDATORS_COUNT=4
ENV SEEDS_COUNT=1
ENV OBSERVERS_COUNT=1

# Set plain vanilla default entrypoint/command
ENTRYPOINT [ "generate-network" ]
