#!/bin/bash

# Inits node configuration and runs the node.
# e -> exit immediately, u -> treat unset variables as errors and immediately, o -> sets the exit code to the rightmost command 
set -euo pipefail

# within the container, $HOME=/home/cheqd
CHEQD_ROOT_DIR="$HOME/.cheqdnode"

# Init node config directory
if [ ! -d "${CHEQD_ROOT_DIR}/config" ]
then
    echo "Node config not found. Initializing."
    cheqd-noded init "$NODE_MONIKER" --home "${CHEQD_ROOT_DIR}"
else
    echo "Node config exists. Skipping initialization."
fi

# Check if a genesis file has been passed in config
if [ -f "/genesis" ]
then
    echo "Genesis file passed. Adding/replacing current genesis file."
    cp /genesis "${CHEQD_ROOT_DIR}/config/genesis.json"
else
    echo "No genesis file passed. Skipping and retaining existing genesis."
fi

# Check if a seeds file has been passed in config
if [ -f "/seeds" ]
then
    echo "Seeds file passed. Replacing current seeds."
    cp /seeds "${CHEQD_ROOT_DIR}/config/seeds.txt"
    cheqd-noded configure p2p seeds "$(cat "${CHEQD_ROOT_DIR}/config/seeds.txt")"
else
    echo "No seeds file passed. Skipping and retaining existing seeds."
fi

# Check if a node_key file has been passed in config
if [ -f "/node_key" ]
then
    echo "Node key file passed. Replacing current key."
    cp /node_key "${CHEQD_ROOT_DIR}/config/node_key.json"
else
    echo "No node key file passed. Skipping and retaining existing node key."
fi

# Check if a priv_validator_key file has been passed in config
if [ -f "/priv_validator_key" ]
then
    echo "Private validator key file passed. Replacing current key."
    cp /priv_validator_key "${CHEQD_ROOT_DIR}/config/priv_validator_key.json"
else
    echo "No private validator key file passed. Skipping and retaining existing key."
fi

# Run configure
# `! -z` is used instead of `-n` to distinguish null and empty values
if [[ -n ${CREATE_EMPTY_BLOCKS+x} ]]; then cheqd-noded configure create-empty-blocks "${CREATE_EMPTY_BLOCKS}"; fi
if [[ -n ${FASTSYNC_VERSION+x} ]]; then cheqd-noded configure fastsync-version "${FASTSYNC_VERSION}"; fi
if [[ -n ${MIN_GAS_PRICES+x} ]]; then cheqd-noded configure min-gas-prices "${MIN_GAS_PRICES}"; fi
if [[ -n ${RPC_LADDR+x} ]]; then cheqd-noded configure rpc-laddr "${RPC_LADDR}"; fi
if [[ -n ${P2P_EXTERNAL_ADDRESS+x} ]]; then cheqd-noded configure p2p external-address "${P2P_EXTERNAL_ADDRESS}"; fi
if [[ -n ${P2P_LADDR+x} ]]; then cheqd-noded configure p2p laddr "${P2P_LADDR}"; fi
if [[ -n ${P2P_MAX_PACKET_MSG_PAYLOAD_SIZE+x} ]]; then cheqd-noded configure p2p max-packet-msg-payload-size "${P2P_MAX_PACKET_MSG_PAYLOAD_SIZE}"; fi
if [[ -n ${P2P_SEEDS+x} ]]; then cheqd-noded configure p2p seeds "${P2P_SEEDS}"; fi
if [[ -n ${P2P_PERSISTENT_PEERS+x} ]]; then cheqd-noded configure p2p persistent-peers "${P2P_PERSISTENT_PEERS}"; fi
if [[ -n ${P2P_SEED_MODE+x} ]]; then cheqd-noded configure p2p seed-mode "${P2P_SEED_MODE}"; fi
if [[ -n ${P2P_RECV_RATE+x} ]]; then cheqd-noded configure p2p recv-rate "${P2P_RECV_RATE}"; fi
if [[ -n ${P2P_SEND_RATE+x} ]]; then cheqd-noded configure p2p send-rate "${P2P_SEND_RATE}"; fi

# Run node
cheqd-noded start
