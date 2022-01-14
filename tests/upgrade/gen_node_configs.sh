#!/bin/bash

# Generates configurations for 4 nodes.

set -euox pipefail

. common.sh

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi


# Parameters
# $1 - The full parameter name
# $2 - Value
# $3 - Path to genesis file
edit_genesis () {
    PARAMETER_NAME=$1
    VALUE=$2
    GENESIS_PATH=$3
    jq "${PARAMETER_NAME} = \"${VALUE}\"" $GENESIS_PATH > /tmp/1 && mv /tmp/1 $GENESIS_PATH
}


VALIDATORS_COUNT="4"
OBSERVERS_COUNT="2"

NODE_CONFIGS_DIR="node_configs"
rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR 
pushd $NODE_CONFIGS_DIR

echo "Generating validator keys..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="node$i"
    mkdir $NODE_HOME
    pushd $NODE_HOME

    echo "[Validator $i] Generating key..."

    cheqd_noded_docker_init "node$i" --chain-id $CHAIN_ID
    # sudo chown -R $USER:$USER .
    echo "$(cheqd_noded_docker tendermint show-node-id)" > node_id.txt
    echo "$(cheqd_noded_docker tendermint show-validator)" > node_val_pubkey.txt

    echo "Setting minimum fee price..."

    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' .cheqdnode/config/app.toml

    popd
done


OPERATORS_HOME="client"
mkdir $OPERATORS_HOME
pushd $OPERATORS_HOME

echo "Initializing genesis..."
cheqd_noded_docker_init dummy_node --chain-id $CHAIN_ID
# sudo chown -R $USER:$USER .
sed -i $sed_extension 's/"stake"/"ncheq"/' .cheqdnode/config/genesis.json

echo "Generating operator keys..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd_noded_docker keys add "operator$i" --keyring-backend "test"
    # sudo chown -R $USER:$USER .
done

echo "Creating genesis accounts..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    cheqd_noded_docker add-genesis-account "operator$i" 20000000000000000ncheq --keyring-backend "test"
done

echo "Creating genesis validators..."

for ((i=0 ; i<$VALIDATORS_COUNT ; i++))
do
    NODE_HOME="../node$i"
    pushd $NODE_HOME

    NODE_ID=$(cheqd_noded_docker tendermint show-node-id)
    NODE_VAL_PUBKEY=$(cheqd_noded_docker tendermint show-validator)

    popd

    cheqd_noded_docker gentx "operator$i" 1000000000000000ncheq --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY --keyring-backend "test"
done

echo "Collecting them..."

cheqd_noded_docker collect-gentxs
cheqd_noded_docker validate-genesis

echo "Propagating genesis to nodes..."

# Edit voting parameters
edit_genesis ".app_state.gov.deposit_params.max_deposit_period" "${VOTING_PERIOD}s" .cheqdnode/config/genesis.json
edit_genesis ".app_state.gov.voting_params.voting_period" "${VOTING_PERIOD}s" .cheqdnode/config/genesis.json

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
    pushd $NODE_HOME

    echo "##### [Observer $i] Generating keys..."
    cheqd_noded_docker_init "node$i" --chain-id $CHAIN_ID
    # sudo chown -R $USER:$USER .

    echo "##### [Observer $i] Exporting public keys..."
    echo "$(cheqd_noded_docker tendermint show-node-id)" > node_id.txt
    echo "$(cheqd_noded_docker tendermint show-validator)" > node_val_pubkey.txt

    echo "##### [Observer $i] Loading genesis..."
    OPERATORS_HOME="../client"
    cp "$OPERATORS_HOME/.cheqdnode/config/genesis.json" ".cheqdnode/config/"

    echo "##### [Observer $i] Setting min gas prices..."
    sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' .cheqdnode/config/app.toml

    popd
done

popd # node_configs
