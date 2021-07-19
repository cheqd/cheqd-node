#!/bin/bash

# Generates configurations for 4 nodes.

set -euox pipefail

NODES_COUNT="4"
CHAIN_ID="verim"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


for ((i=0 ; i<$NODES_COUNT ; i++))
do
    echo "##### [Node $i] Generates key" 

    NODE_HOME="$NODE_CONFIGS_DIR/node$i"
    verim-noded init "node$i" --chain-id $CHAIN_ID --home $NODE_HOME
    NODE_ID=$(verim-noded tendermint show-node-id --home $NODE_HOME)
    NODE_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_HOME)

    echo "$NODE_ID" > $NODE_HOME/node_id.txt
    echo "$NODE_VAL_PUBKEY" > $NODE_HOME/node_val_pubkey.txt
done


echo "##### [Validator operators] Generate keys"

OPERATORS_HOME=$NODE_CONFIGS_DIR/client

for ((i=0 ; i<$NODES_COUNT ; i++))
do
    verim-noded keys add "operator$i" --home $OPERATORS_HOME
done


echo "##### [Validator operators] Init genesis" 

verim-noded init dummy_node --chain-id $CHAIN_ID --home $OPERATORS_HOME


echo "##### [Validator operators] Add them to the genesis" 

for ((i=0 ; i<$NODES_COUNT ; i++))
do
    verim-noded add-genesis-account "operator$i" 10000000token,100000000stake --home $OPERATORS_HOME
done


echo "##### [Validator operators] Generate stake transactions" 

for ((i=0 ; i<$NODES_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/node$i"
    NODE_ID=$(verim-noded tendermint show-node-id --home $NODE_HOME)
    NODE_VAL_PUBKEY=$(verim-noded tendermint show-validator --home $NODE_HOME)

    verim-noded gentx "operator$i" 1000000stake --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY --home $OPERATORS_HOME
done


echo "##### [Validator operators] Collect them"

verim-noded collect-gentxs --home $OPERATORS_HOME
verim-noded validate-genesis --home $OPERATORS_HOME


echo "##### [Validator operators] Propagate genesis to nodes"

for ((i=0 ; i<$NODES_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/node$i"

    cp $OPERATORS_HOME/config/genesis.json $NODE_HOME/config/
done
