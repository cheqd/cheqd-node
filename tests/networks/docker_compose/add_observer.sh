#!/bin/bash

set -euox pipefail

cheqd-noded init node5

NODE0_ID=$(cat "${NODE_CONFIGS_BASE}/node0/node_id.txt")
PERSISTENT_PEERS="${NODE0_ID}@127.0.0.1:26656"
cheqd-noded configure p2p persistent-peers "${PERSISTENT_PEERS}"

cp "${NODE_CONFIGS_BASE}/node0/.cheqdnode/config/genesis.json" "/home/runner/cheqd/.cheqdnode/config"

cheqd-noded configure p2p laddr "tcp://0.0.0.0:26676"
cheqd-noded configure rpc-laddr "tcp://0.0.0.0:26677"

