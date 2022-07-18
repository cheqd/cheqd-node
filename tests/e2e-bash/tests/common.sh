#!/bin/bash

set -euox pipefail

source "../../tools/helpers.sh"

# Params
export RPC_URL="http://localhost:26657"
export CHAIN_ID="cheqd"
export GAS="auto"
export GAS_ADJUSTMENT="1.3"
export GAS_PRICES="25ncheq"
export KEYRING_BACKEND="test"
export OUTPUT_FORMAT="json"

export QUERY_PARAMS="--node ${RPC_URL} --output ${OUTPUT_FORMAT}"
export KEYS_PARAMS="--keyring-backend ${KEYRING_BACKEND} --output ${OUTPUT_FORMAT}"
export TX_PARAMS="--node ${RPC_URL} --keyring-backend ${KEYRING_BACKEND} --output ${OUTPUT_FORMAT} --chain-id ${CHAIN_ID}
  --gas ${GAS} --gas-adjustment ${GAS_ADJUSTMENT} --gas-prices ${GAS_PRICES} --yes --home /home/shalteor/Documents/@Cheqd/cheqd-node/tests/e2e-complex/upgrade/node_configs/client/.cheqdnode/"

# Accounts
export BASE_ACCOUNT_1="operator0"
export BASE_ACCOUNT_2="operator1"
export BASE_VESTING_ACCOUNT="base_vesting_account"
export CONTINOUS_VESTING_ACCOUNT="continous_vesting_account"
export DELAYED_VESTING_ACCOUNT="delayed_vesting_account"
export PERIODIC_VESTING_ACCOUNT="periodic_vesting_account"
