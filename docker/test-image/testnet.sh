#!/bin/bash
set -euo pipefail

HOME_DIR="/home/cheqd"
CHEQD_ROOT_DIR="${HOME_DIR}/.cheqdnode/"
TESTNET_NAME="cheqd"
NODE_MONIKER="node0"

# Enable config parameters as environment variables
source ${HOME_DIR}/validator-0.env

mkdir -p ${CHEQD_ROOT_DIR}

# Generate network config
bash -x gen-network ${TESTNET_NAME} 1 0 0
mv network-config/validator-0/* ${CHEQD_ROOT_DIR}
NODE_ID=$(cheqd-noded tendermint show-node-id)

# Start node
cheqd-noded start --home ${CHEQD_ROOT_DIR} --log_level info --log_format json --trace --p2p.persistent_peers "$NODE_ID"@"$NODE_MONIKER":26656