#!/bin/bash

set -euox pipefail

source "../common.sh"

CHAIN_ID="cheqd"

rm -rf "$HOME/.cheqdnode"

cheqd-noded init "local_node" --chain-id $CHAIN_ID
sed -i $sed_extension 's/"stake"/"ncheq"/' $HOME/.cheqdnode/config/genesis.json
sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' $HOME/.cheqdnode/config/app.toml


cheqd-noded keys add "node_operator"
cheqd-noded add-genesis-account "node_operator" 20000000000000000ncheq

NODE_ID=$(cheqd-noded tendermint show-node-id)
NODE_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)

cheqd-noded gentx "node_operator" 1000000000000000ncheq --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY

cheqd-noded collect-gentxs
cheqd-noded validate-genesis
