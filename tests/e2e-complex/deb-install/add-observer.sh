#!/bin/bash

set -euox pipefail
sudo chown -R cheqd:cheqd "/home/runner/cheqd/"

NODE_CONFIGS_BASE="/home/runner/work/cheqd-node/cheqd-node/tests/networks/docker-compose-localnet/network-config"
sudo -u cheqd cheqd-noded init node5

VALIDATOR_0_ID=$(cat "${NODE_CONFIGS_BASE}/validator-0/node_id.txt")
PERSISTENT_PEERS="${VALIDATOR_0_ID}@127.0.0.1:26656"
sudo -u cheqd cheqd-noded configure p2p persistent-peers "${PERSISTENT_PEERS}"

sudo cp "${NODE_CONFIGS_BASE}/validator-0/config/genesis.json" "/home/runner/cheqd/.cheqdnode/config"

sudo chmod -R 777 "/home/runner/cheqd/.cheqdnode"

sudo -u cheqd cheqd-noded configure p2p laddr "tcp://0.0.0.0:26676"
sudo -u cheqd cheqd-noded configure rpc-laddr "tcp://0.0.0.0:26677"


sudo systemctl start cheqd-noded
sleep 10
systemctl status cheqd-noded

