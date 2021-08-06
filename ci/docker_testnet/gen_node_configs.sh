#!/bin/bash

# Generates configurations for 2 nodes.

set -euox pipefail

CHAIN_ID="cheqd"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


echo "##### [Node 0] Generates key" 

NODE_0_HOME="$NODE_CONFIGS_DIR/node0"
cheqd-noded init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME
NODE_0_ID=$(cheqd-noded tendermint show-node-id --home $NODE_0_HOME)
NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_0_HOME)


echo "##### [Node 1] Generates key" 

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"
cheqd-noded init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME
NODE_1_ID=$(cheqd-noded tendermint show-node-id --home $NODE_1_HOME)
NODE_1_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_1_HOME)


echo "##### [Validator operators] Generate keys" 

CLIENT_HOME=$NODE_CONFIGS_DIR/client
cheqd-noded keys add alice --home $CLIENT_HOME
cheqd-noded keys add bob --home $CLIENT_HOME


echo "##### [Validator operators] Init genesis" 

cheqd-noded init dummy_node --chain-id $CHAIN_ID --home $CLIENT_HOME

sed -i 's/"stake"/"cheq"/' ${CLIENT_HOME}/config/genesis.json

echo "##### [Validator operators] Add them to the genesis" 

cheqd-noded add-genesis-account alice 20000000cheq --home $CLIENT_HOME
cheqd-noded add-genesis-account bob 20000000cheq --home $CLIENT_HOME



echo "##### [Test pool] Add test account to the genesis" 

ACCOUNT_ID="cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd"
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${ACCOUNT_ID}'", "coins": [{"denom": "cheq", "amount": "100001000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${ACCOUNT_ID}'", "pub_key": null,"account_number": "0","sequence": "0"}]') > ${CLIENT_HOME}/config/genesis.json


echo "##### [Validator operators] Generate stake transactions" 

cheqd-noded gentx alice 1000000cheq --chain-id $CHAIN_ID --node-id $NODE_0_ID --pubkey $NODE_0_VAL_PUBKEY --home $CLIENT_HOME
cheqd-noded gentx bob 1000000cheq --chain-id $CHAIN_ID --node-id $NODE_1_ID --pubkey $NODE_1_VAL_PUBKEY --home $CLIENT_HOME


echo "##### [Validator operators] Collect them"

cheqd-noded collect-gentxs --home $CLIENT_HOME
cheqd-noded validate-genesis --home $CLIENT_HOME


echo "##### [Validator operators] Propagate genesis to nodes"

cp $CLIENT_HOME/config/genesis.json $NODE_0_HOME/config/
cp $CLIENT_HOME/config/genesis.json $NODE_1_HOME/config/
