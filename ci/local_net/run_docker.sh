#!/bin/bash

# The scrit that exports node ids and runs docker-compose.

set -euox pipefail

NODE_CONFIGS_DIR="node_configs"

NODE_0_HOME="$NODE_CONFIGS_DIR/node0"
export NODE_0_ID=$(verim-noded tendermint show-node-id --home $NODE_0_HOME)

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"
export NODE_1_ID=$(verim-noded tendermint show-node-id --home $NODE_1_HOME)

NODE_2_HOME="$NODE_CONFIGS_DIR/node2"
export NODE_2_ID=$(verim-noded tendermint show-node-id --home $NODE_2_HOME)

NODE_3_HOME="$NODE_CONFIGS_DIR/node3"
export NODE_3_ID=$(verim-noded tendermint show-node-id --home $NODE_3_HOME)

docker-compose up
