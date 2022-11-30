#!/bin/bash

set -euox pipefail

# sudo chown -R runner:docker "${NODE_CONFIGS_BASE}"
cheqd-noded keys list --keyring-backend "test" --home "${NODE_CONFIGS_BASE}/client/.cheqdnode"

# Get operator0 address by setting --home flag
OP0_ADDRESS=$(cheqd-noded keys list --keyring-backend "test" --home "${NODE_CONFIGS_BASE}/client/.cheqdnode" | sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')

# Create operator5 by running it under the `cheqd` user.
sudo su -c 'cheqd-noded keys add node5-operator --keyring-backend "test"' cheqd
OP5_ADDRESS=$(sudo su -c 'cheqd-noded keys list --keyring-backend "test"' cheqd | sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')

NODE5_PUBKEY=$(sudo su -c 'cheqd-noded tendermint show-validator' cheqd | sed 's/\r//g')

sleep 10

cheqd-noded version

pushd ../

# shellcheck disable=SC1091
. common.sh

cheqd_noded_docker query bank balances "${OP0_ADDRESS}"

local_client_tx tx bank send "${OP0_ADDRESS}" "${OP5_ADDRESS}" 11000000000000ncheq --chain-id cheqd --fees 5000000ncheq -y

popd

# Send promote validator from operator5
sudo -H -u cheqd cheqd-noded tx staking create-validator --amount 10000000000000ncheq --from node5-operator --chain-id cheqd --min-self-delegation="1" --gas-prices="25ncheq" --pubkey "${NODE5_PUBKEY}" --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01" --gas 500000 --node http://127.0.0.1:26657 --keyring-backend "test" -y
