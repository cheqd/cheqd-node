#!/bin/bash

# The script that exports node ids and runs docker-compose.

set -euox pipefail

pushd "node_configs/node0"

export NODE_0_ID=$(cheqd-noded tendermint show-node-id | sed 's/\r//g')

popd

docker-compose up -d
