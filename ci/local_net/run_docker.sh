#!/bin/bash

# The scrit that exports node ids and runs docker-compose.

set -euox pipefail

export NODE_0_ID=$(cheqd-noded tendermint show-node-id --home node_configs/node0)
export NODE_1_ID=$(cheqd-noded tendermint show-node-id --home node_configs/node1)
export NODE_2_ID=$(cheqd-noded tendermint show-node-id --home node_configs/node2)
export NODE_3_ID=$(cheqd-noded tendermint show-node-id --home node_configs/node3)

docker-compose up
