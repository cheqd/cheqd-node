###############################################################
###         	STAGE 1: Build localnet image         		###
###############################################################

ARG BUILD_IMAGE="ghcr.io/cheqd/cheqd-node"
ARG BUILD_VERSION
FROM ${BUILD_IMAGE}:${BUILD_VERSION} AS localnet

# Set bash defaults
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Set home directory
WORKDIR /home/cheqd

# Copy scripts
COPY docker/localnet/entrypoint.sh /bin/node-start
COPY scripts/import-keys.sh /bin/import-keys

# Make scripts executable
RUN chmod +rx /bin/node-start \
	&& chmod +rx /bin/import-keys

# Set user and group
USER cheqd:cheqd

# Copy network configuration directory
COPY scripts/network-config /home/cheqd/network-config

EXPOSE 26656 26657 26660 1317 9090 9091

# Set plain vanilla default entrypoint/command
CMD [ "cheqd-noded" ]
