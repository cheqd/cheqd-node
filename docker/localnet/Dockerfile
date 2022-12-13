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
COPY generate-network/ 

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
