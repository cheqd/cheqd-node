#!/bin/bash

# Generates configurations for 4 nodes.

set -euox pipefail

CHAIN_ID="cheqd"

VALIDATORS_COUNT="4"
OBSERVERS_COUNT="2"

source "../common.sh"


NODE_CONFIGS_DIR="node_configs"
rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR
# chmod -R 777 $NODE_CONFIGS_DIR # FIXME
pushd $NODE_CONFIGS_DIR

echo "Generating validator keys..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="node$i"
    mkdir $NODE_HOME
    # chmod -R 777 $NODE_HOME # FIXME
    pushd $NODE_HOME

    echo "[Validator $i] Generating key..."

    cheqd_noded_docker init "node$i" --chain-id $CHAIN_ID
    echo "$(cheqd_noded_docker tendermint show-node-id)" > node_id.txt
    echo "$(cheqd_noded_docker tendermint show-validator)" > node_val_pubkey.txt

    echo "Setting minimum fee price..."

    # sudo chmod -R 777 /home/runner/work/cheqd-node/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in pipeline" # FIXME
    # sudo chmod -R 777 ~/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in local machine" # FIXME
    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' .cheqdnode/config/app.toml

    popd
done


OPERATORS_HOME="client"
mkdir $OPERATORS_HOME
# chmod -R 777 $OPERATORS_HOME # FIXME
pushd $OPERATORS_HOME

echo "Initializing genesis..."
cheqd_noded_docker init dummy_node --chain-id $CHAIN_ID
# sudo chmod -R 777 /home/runner/work/cheqd-node/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in pipeline" # FIXME
# sudo chmod -R 777 ~/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in local machine" # FIXME
sed -i $sed_extension 's/"stake"/"ncheq"/' .cheqdnode/config/genesis.json

echo "Generating operator keys..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd_noded_docker keys add "operator$i"
done

echo "Creating genesis accounts..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd_noded_docker add-genesis-account "operator$i" 20000000000000000ncheq
done

echo "Creating genesis validators..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="../node$i"
    pushd $NODE_HOME

    NODE_ID=$(cheqd_noded_docker tendermint show-node-id)
    NODE_VAL_PUBKEY=$(cheqd_noded_docker tendermint show-validator)

    popd

    cheqd_noded_docker gentx "operator$i" 1000000000000000ncheq --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY
done

echo "Collecting them..."

cheqd_noded_docker collect-gentxs
cheqd_noded_docker validate-genesis

echo "Propagating genesis to nodes..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="../node$i"

    cp ".cheqdnode/config/genesis.json" "$NODE_HOME/.cheqdnode/config/"
done


popd # operators' home


echo "##### Setting up observers..."

for ((i=0 ; i<$OBSERVERS_COUNT ; i++))
do
    NODE_HOME="observer$i"

    mkdir $NODE_HOME
    # chmod -R 777 $NODE_HOME # FIXME
    pushd $NODE_HOME

    echo "##### [Observer $i] Generating keys..."
    cheqd_noded_docker init "node$i" --chain-id $CHAIN_ID

    echo "##### [Observer $i] Exporting public keys..."
    echo "$(cheqd_noded_docker tendermint show-node-id)" > node_id.txt
    echo "$(cheqd_noded_docker tendermint show-validator)" > node_val_pubkey.txt

    echo "##### [Observer $i] Loading genesis..."
    OPERATORS_HOME="../client"
    # sudo chmod -R 777 /home/runner/work/cheqd-node/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in pipeline" # FIXME
    # sudo chmod -R 777 ~/cheqd-node/tests/networks/docker_compose/node_configs || echo "I'm not in local machine" # FIXME
    cp "$OPERATORS_HOME/.cheqdnode/config/genesis.json" ".cheqdnode/config/"

    echo "##### [Observer $i] Setting min gas prices..."
    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' .cheqdnode/config/app.toml

    popd
done

popd # node_configs
