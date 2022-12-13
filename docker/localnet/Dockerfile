###############################################################
###         	STAGE 1: Build localnet image         		###
###############################################################

ARG BASE_IMAGE
FROM ${BASE_IMAGE} AS localnet

# Set bash defaults
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Set home directory
WORKDIR /home/cheqd

# Copy network configuration directory
COPY --chown=cheqd:cheqd docker/localnet/generate-network/network-config /home/cheqd/network-config

# Copy entrypoint script
COPY --chown=cheqd:cheqd docker/entrypoint.sh /bin/node-start
RUN chmod +rx /bin/node-start

# Set user and group
USER cheqd:cheqd

EXPOSE 26656 26657 26660 1317 9090 9091

# Set plain vanilla default entrypoint/command
CMD [ "cheqd-noded" ]
