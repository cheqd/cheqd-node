#!/bin/bash

# The scrit that exports node ids and runs docker-compose.

set -euox pipefail

export NODE_0_ID=$(cheqd-noded tendermint show-node-id --home node_configs/node0)

docker-compose up
