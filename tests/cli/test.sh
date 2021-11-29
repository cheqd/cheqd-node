#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

# Params
RPC_URL="http://localhost:26657"
CHAIN_ID="cheqd"
GAS_PRICES="25ncheq"
KEYRING_BACKEND="test"
OUTPUT_FORMAT="json"

QUERY_PARAMS="--node ${RPC_URL} --chain-id ${CHAIN_ID} --output ${OUTPUT_FORMAT}"
TX_PARAMS="${QUERY_PARAMS} --gas-prices ${GAS_PRICES} --keyring-backend ${KEYRING_BACKEND}"


# Generate identity key
IDENTITY_KEY="$(cheqd-noded tools ed25519 random)"
IDENTITY_PUB_KEY=$(echo "${IDENTITY_KEY}" | jq -r ".pub_key_base_64")
IDENTITY_PRIV_KEY=$(echo "${IDENTITY_KEY}" | jq -r ".priv_key_base_64")


# shellcheck disable=SC2086
cheqd-noded tx cheqd create-did '{"id": "did:cheqd:test:alice"}' --from "node_operator" ${TX_PARAMS}


#{"pub_key_base_64":"zd9mfKkmb21P904J7SlVS0PrWrHWOsSRfeVgeElv5Mg=",
# "priv_key_base_64":"WVbBBALxwFALC3qBGQikW3qi1Bbsu8WIgq8hfYjAkW3N32Z8qSZvbU/3TgntKVVLQ+tasdY6xJF95WB4SW/kyA=="}

# zEre9xWvKohVztkNY5UBcWKNFxR5XKanUhKmWU1V3Eeo1


#rm -rf "$HOME/.cheqdnode"
#
#cheqd-noded init "local_node" --chain-id $CHAIN_ID
#sed -i $sed_extension 's/"stake"/"ncheq"/' $HOME/.cheqdnode/config/genesis.json
#sed -i $sed_extension 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' $HOME/.cheqdnode/config/app.toml
#
#
#cheqd-noded keys add "node_operator" --keyring-backend "test"
#cheqd-noded add-genesis-account "node_operator" 20000000000000000ncheq --keyring-backend "test"
#
#NODE_ID=$(cheqd-noded tendermint show-node-id)
#NODE_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)
#
#cheqd-noded gentx "node_operator" 1000000000000000ncheq --chain-id $CHAIN_ID --node-id $NODE_ID --pubkey $NODE_VAL_PUBKEY --keyring-backend "test"
#
#cheqd-noded collect-gentxs
#cheqd-noded validate-genesis
#
