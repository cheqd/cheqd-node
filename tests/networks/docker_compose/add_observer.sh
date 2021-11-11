#!/bin/bash

set -euox pipefail


cheqd-noded init node5

NODE0_ID=$(cat "${NODE_CONFIGS_BASE}/node0/node_id.txt")
NODE0_PEER_INFO="${NODE0_ID}@127.0.0.1:26656"

cheqd-noded configure p2p persistent-peers "${NODE0_PEER_INFO}"
cheqd-noded configure p2p laddr "tcp://0.0.0.0:26676"
cheqd-noded configure rpc-laddr "tcp://127.0.0.1:26677"

cp "${NODE_CONFIGS_BASE}/node0/.cheqdnode/config/genesis.json" "$HOME/.cheqdnode/config"

sudo chmod -R 777 "$HOME/.cheqdnode"

sudo systemctl start cheqd-noded
systemctl status cheqd-noded
