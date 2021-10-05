#!/bin/bash

# The script that exports node ids and runs docker-compose.

set -euox pipefail

source "../common.sh"

pushd "node_configs/node0"

export NODE_0_ID=$(cheqd_noded_docker tendermint show-node-id | sed 's/\r//g')

popd

docker-compose up
