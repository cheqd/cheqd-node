#!/bin/bash

# Inits node configuration and runs the node.

set -euox pipefail

CHEQD_ROOT_DIR="$HOME/.cheqdnode"

# Init node config directory
if [ ! -d "${CHEQD_ROOT_DIR}/config" ]
then
    echo "Node config not found. Initializing."
    cheqd-noded init $NODE_MONIKER --home 
else
    echo "Node config exists. Skipping initialization."
fi

# Run configure
# `! -z` is used instead of `-n` to distinguilsh null and empty values

if [[ ! -z ${CREATE_EMPTY_BLOCKS+x} ]]; then cheqd-noded configure create-empty-blocks ${CREATE_EMPTY_BLOCKS}; fi
if [[ ! -z ${FASTSYNC_VERSION+x} ]]; then cheqd-noded configure fastsync-version ${FASTSYNC_VERSION}; fi
if [[ ! -z ${MIN_GAS_PRICES+x} ]]; then cheqd-noded configure min-gas-prices ${MIN_GAS_PRICES}; fi
if [[ ! -z ${RPC_LADDR+x} ]]; then cheqd-noded configure rpc-laddr ${RPC_LADDR}; fi
if [[ ! -z ${P2P_EXTERNAL_ADDRESS+x} ]]; then cheqd-noded configure p2p external-address ${P2P_EXTERNAL_ADDRESS}; fi
if [[ ! -z ${P2P_LADDR+x} ]]; then cheqd-noded configure p2p laddr ${P2P_LADDR}; fi
if [[ ! -z ${P2P_MAX_PACKET_MSG_PAYLOAD_SIZE+x} ]]; then cheqd-noded configure p2p max-packet-msg-payload-size ${P2P_MAX_PACKET_MSG_PAYLOAD_SIZE}; fi
if [[ ! -z ${P2P_PERSISTENT_PEERS+x} ]]; then cheqd-noded configure p2p persistent-peers ${P2P_PERSISTENT_PEERS}; fi
if [[ ! -z ${P2P_RECV_RATE+x} ]]; then cheqd-noded configure p2p recv-rate ${P2P_RECV_RATE}; fi
if [[ ! -z ${P2P_SEED_MODE+x} ]]; then cheqd-noded configure p2p seed-mode ${P2P_SEED_MODE}; fi
if [[ ! -z ${P2P_SEEDS+x} ]]; then cheqd-noded configure p2p seeds ${P2P_SEEDS}; fi
if [[ ! -z ${P2P_SEND_RATE+x} ]]; then cheqd-noded configure p2p send-rate ${P2P_SEND_RATE}; fi


# Update configs
set - # Disable ditailed command logging

echo "Updating genesis"
echo "$GENESIS" | base64 --decode > $CHEQD_ROOT_DIR/config/genesis.json

set -x # Re-enable ditailed command logging

# Run node
NODE_ARGS=${NODE_ARGS:-}  # Allo node args to be empty
cheqd-noded start $NODE_ARGS
