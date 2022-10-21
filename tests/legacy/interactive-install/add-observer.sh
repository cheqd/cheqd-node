#!/bin/bash

set -euox pipefail
sudo chown -R cheqd:cheqd "/home/runner/cheqd/"

sudo -u cheqd cheqd-noded init node5

VALIDATOR_0_ID=$(cat "${NODE_CONFIGS_BASE}/validator-0/node_id.txt")
PERSISTENT_PEERS="${VALIDATOR_0_ID}@127.0.0.1:26656"
sudo -u cheqd cheqd-noded configure p2p persistent-peers "${PERSISTENT_PEERS}"

sudo cp "${NODE_CONFIGS_BASE}/validator-0/config/genesis.json" "/home/runner/cheqd/.cheqdnode/config"

sudo chmod -R 755 "/home/runner/cheqd/.cheqdnode"

# Configure ports because they conflict with localnet
sudo -u cheqd cheqd-noded configure p2p laddr "tcp://0.0.0.0:26676"
sudo -u cheqd cheqd-noded configure rpc-laddr "tcp://0.0.0.0:26677"

# TODO: Use environment variables
sudo -u cheqd sed -i.bak 's|pprof_laddr = "localhost:6060"|pprof_laddr = "localhost:6070"|g' /home/runner/cheqd/.cheqdnode/config/config.toml
sudo -u cheqd sed -i.bak 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9100"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo -u cheqd sed -i.bak 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9101"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo -u cheqd sed -i.bak 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:1327"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo -u cheqd sed -i.bak 's|address = ":8080"|address = ":8090"|g' /home/runner/cheqd/.cheqdnode/config/app.toml

sudo systemctl start cheqd-cosmovisor
sleep 10
systemctl status cheqd-cosmovisor
