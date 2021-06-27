#!/bin/bash

# Generates configurations for 4 nodes. Mostly the same script as `gen_node_configs.sh` but discribes
# the real life flow (keys don't leave the place where they are generated).

set -euox pipefail

CHAIN_ID="verim"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


echo "##### [Node 0] Generates key" 

NODE_0_HOME="$NODE_CONFIGS_DIR/node0"
verim-noded init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME
NODE_0_ID=$(verim-noded tendermint show-node-id --home $NODE_0_HOME)
NODE_0_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_0_HOME)


echo "##### [Node 1] Generates key" 

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"
verim-noded init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME
NODE_1_ID=$(verim-noded tendermint show-node-id --home $NODE_1_HOME)
NODE_1_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_1_HOME)


echo "##### [Node 2] Generates key"

NODE_2_HOME="$NODE_CONFIGS_DIR/node2"
verim-noded init node2 --chain-id $CHAIN_ID --home $NODE_2_HOME
NODE_2_ID=$(verim-noded tendermint show-node-id --home $NODE_2_HOME)
NODE_2_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_2_HOME)


echo "##### [Node 3] Generates key"

NODE_3_HOME="$NODE_CONFIGS_DIR/node3"
verim-noded init node3 --chain-id $CHAIN_ID --home $NODE_3_HOME
NODE_3_ID=$(verim-noded tendermint show-node-id --home $NODE_3_HOME)
NODE_3_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_3_HOME)


echo "##### [Validator operators] Generate keys" 

CLIENT_HOME=$NODE_CONFIGS_DIR/client
verim-noded keys add alice --home $CLIENT_HOME
verim-noded keys add bob --home $CLIENT_HOME
verim-noded keys add jack --home $CLIENT_HOME
verim-noded keys add anna --home $CLIENT_HOME


echo "##### [Validator operators] Init genesis" 

verim-noded init dummy_node --chain-id $CHAIN_ID --home $CLIENT_HOME


echo "##### [Validator operators] Add them to the genesis" 

verim-noded add-genesis-account alice 10000000token,100000000stake --home $CLIENT_HOME
verim-noded add-genesis-account bob 10000000token,100000000stake --home $CLIENT_HOME
verim-noded add-genesis-account jack 10000000token,100000000stake --home $CLIENT_HOME
verim-noded add-genesis-account anna 10000000token,100000000stake --home $CLIENT_HOME


echo "##### [Validator operators] Generate stake transactions" 

verim-noded gentx alice 1000000stake --chain-id $CHAIN_ID --node-id $NODE_0_ID --pubkey $NODE_0_VAL_PUBKEY --home $CLIENT_HOME
verim-noded gentx bob 1000000stake --chain-id $CHAIN_ID --node-id $NODE_1_ID --pubkey $NODE_1_VAL_PUBKEY --home $CLIENT_HOME
verim-noded gentx jack 1000000stake --chain-id $CHAIN_ID --node-id $NODE_2_ID --pubkey $NODE_2_VAL_PUBKEY --home $CLIENT_HOME
verim-noded gentx anna 1000000stake --chain-id $CHAIN_ID --node-id $NODE_3_ID --pubkey $NODE_3_VAL_PUBKEY --home $CLIENT_HOME


echo "##### [Validator operators] Collect them"

verim-noded collect-gentxs --home $CLIENT_HOME
verim-noded validate-genesis --home $CLIENT_HOME


echo "##### [Validator operators] Propagate genesis to nodes"

cp $CLIENT_HOME/config/genesis.json $NODE_0_HOME/config/
cp $CLIENT_HOME/config/genesis.json $NODE_1_HOME/config/
cp $CLIENT_HOME/config/genesis.json $NODE_2_HOME/config/
cp $CLIENT_HOME/config/genesis.json $NODE_3_HOME/config/
