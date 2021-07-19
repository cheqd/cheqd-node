#!/bin/bash

# Generates configurations for 4 nodes. Mostly the same script as `gen_$NODE_CONFIGS_DIR.sh` but discribes
# the real life flow (keys don't leave the place where they are generated).

set -euox pipefail

CHAIN_ID="cheqd"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


echo "################################ Jack's node0"

NODE_0_HOME="node_configs/node0"

echo "# Generate key"
cheqd-noded keys add jack --home $NODE_0_HOME

echo "# Initialze node"
cheqd-noded init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME

echo "# Add genesis account"
cheqd-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_0_HOME

echo "# Generate genesis node tx"
cheqd-noded gentx jack 1000000stake --chain-id $CHAIN_ID --home $NODE_0_HOME

echo "# Publish validator id"
NODE_0_ID=$(cheqd-noded tendermint show-node-id --home $NODE_0_HOME)


echo "################################ Alice's node1"

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"

echo "# Generate key"
cheqd-noded keys add alice --home $NODE_1_HOME

echo "# Initialze node"
cheqd-noded init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME

echo "### Get genesis from Jack"
cp $NODE_0_HOME/config/genesis.json $NODE_1_HOME/config

echo "### Get genesis node txs form Jack"
mkdir $NODE_1_HOME/config/gentx
cp $NODE_0_HOME/config/gentx/* $NODE_1_HOME/config/gentx

echo "# Add genesis account"
cheqd-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_1_HOME

echo "# Generate genesis node tx"
cheqd-noded gentx alice 1000000stake --chain-id $CHAIN_ID --home $NODE_1_HOME

echo "# Publish validator id"
NODE_1_ID=$(cheqd-noded tendermint show-node-id --home $NODE_1_HOME)


echo "################################ Bob's node2"

NODE_2_HOME="$NODE_CONFIGS_DIR/node2"

echo "# Generate key"
cheqd-noded keys add bob --home $NODE_2_HOME

echo "# Initialze node"
cheqd-noded init node2 --chain-id $CHAIN_ID --home $NODE_2_HOME

echo "### Get genesis from Alice"
cp $NODE_1_HOME/config/genesis.json $NODE_2_HOME/config

echo "### Get genesis node txs form Alice"
mkdir $NODE_2_HOME/config/gentx
cp $NODE_1_HOME/config/gentx/* $NODE_2_HOME/config/gentx

echo "# Add genesis account"
cheqd-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_2_HOME

echo "# Generate genesis node tx"
cheqd-noded gentx bob 1000000stake --chain-id $CHAIN_ID --home $NODE_2_HOME

echo "# Publish validator id"
NODE_2_ID=$(cheqd-noded tendermint show-node-id --home $NODE_2_HOME)


echo "################################ Anna's node3"

NODE_3_HOME="$NODE_CONFIGS_DIR/node3"

echo "# Generate key"
cheqd-noded keys add anna --home $NODE_3_HOME

echo "# Initialze node"
cheqd-noded init node3 --chain-id $CHAIN_ID --home $NODE_3_HOME

echo "### Get genesis from Bob"
cp $NODE_2_HOME/config/genesis.json $NODE_3_HOME/config

echo "### Get genesis node txs form Bob"
mkdir $NODE_3_HOME/config/gentx
cp $NODE_2_HOME/config/gentx/* $NODE_3_HOME/config/gentx

echo "# Add genesis account"
cheqd-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_3_HOME

echo "# Generate genesis node tx"
cheqd-noded gentx anna 1000000stake --chain-id $CHAIN_ID --home $NODE_3_HOME

echo "# Publish validator id"
NODE_3_ID=$(cheqd-noded tendermint show-node-id --home $NODE_3_HOME)


echo "################################ Anna (last participatn) shares genesis with everyone else"

echo "# Add genesis node txs into genesis"
cheqd-noded collect-gentxs --home $NODE_3_HOME

echo "# Verify genesis"
cheqd-noded validate-genesis --home $NODE_3_HOME

cp $NODE_3_HOME/config/genesis.json $NODE_0_HOME/config/
cp $NODE_3_HOME/config/genesis.json $NODE_1_HOME/config/
cp $NODE_3_HOME/config/genesis.json $NODE_2_HOME/config/
