#!/bin/bash

# Generates configurations for 4 nodes.

set -euox pipefail

CHAIN_ID="cheqd"
NODE_CONFIGS_DIR="node_configs"

VALIDATORS_COUNT="4"
OBSERVERS_COUNT="2"

# sed in macos requires extra argument
sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


echo "##### Setting up validators" 

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    echo "##### [Validator $i] Generates key" 

    NODE_HOME="$NODE_CONFIGS_DIR/node$i"
    cheqd-noded init "node$i" --chain-id $CHAIN_ID --home $NODE_HOME
    NODE_ID=$(cheqd-noded tendermint show-node-id --home $NODE_HOME)
    NODE_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_HOME)

    echo "$NODE_ID" > $NODE_HOME/node_id.txt
    echo "$NODE_VAL_PUBKEY" > $NODE_HOME/node_val_pubkey.txt
done


echo "##### [Validator operators] Generate keys"

OPERATORS_HOME=$NODE_CONFIGS_DIR/client

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd-noded keys add "operator$i" --home $OPERATORS_HOME
done


echo "##### [Validator operators] Init genesis and make cheq a default denom" 

cheqd-noded init dummy_node --chain-id $CHAIN_ID --home $OPERATORS_HOME
sed -i $sed_extension 's/"stake"/"cheq"/' $OPERATORS_HOME/config/genesis.json

echo "##### [Validator operators] Add them to the genesis" 

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd-noded add-genesis-account "operator$i" 20000000cheq --home $OPERATORS_HOME
done


echo "##### [Validator operators] Generate stake transactions" 

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/node$i"
    NODE_ID=$(cheqd-noded tendermint show-node-id --home $NODE_HOME)
    NODE_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_HOME)

    cheqd-noded gentx "operator$i" 1000000cheq --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY --home $OPERATORS_HOME
done


echo "##### [Validator operators] Collect them"

cheqd-noded collect-gentxs --home $OPERATORS_HOME
cheqd-noded validate-genesis --home $OPERATORS_HOME


echo "##### [Validator operators] Propagate genesis to nodes"

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/node$i"

    cp $OPERATORS_HOME/config/genesis.json $NODE_HOME/config/
done


echo "##### [Validator operators] Set minimum fee price"

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/node$i"

    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "0.00'$i'cheq"/g' $NODE_HOME/config/app.toml
done


echo "##### Setting up observers"

for ((i=0 ; i<$OBSERVERS_COUNT ; i++))
do
    NODE_HOME="$NODE_CONFIGS_DIR/observer$i"

    echo "##### [Observer $i] Generating keys" 
    cheqd-noded init "node$i" --chain-id $CHAIN_ID --home $NODE_HOME
    
    echo "##### [Observer $i] Exporting public keys" 
    NODE_ID=$(cheqd-noded tendermint show-node-id --home $NODE_HOME)
    NODE_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_HOME)
    echo "$NODE_ID" > $NODE_HOME/node_id.txt
    echo "$NODE_VAL_PUBKEY" > $NODE_HOME/node_val_pubkey.txt

    echo "##### [Observer $i] Loading genesis" 
    cp $OPERATORS_HOME/config/genesis.json $NODE_HOME/config/

    echo "##### [Observer $i] Setting min gas prices" 
    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "'$i'cheq"/g' $NODE_HOME/config/app.toml
done
