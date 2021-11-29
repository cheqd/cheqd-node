#!/bin/bash

set -euox pipefail

# Params
export RPC_URL="http://localhost:26657"
export CHAIN_ID="cheqd"
export GAS_PRICES="25ncheq"
export KEYRING_BACKEND="test"
export OUTPUT_FORMAT="json"

export QUERY_PARAMS="--node ${RPC_URL} --output ${OUTPUT_FORMAT}"
export KEYS_PARAMS="--keyring-backend ${KEYRING_BACKEND} --output ${OUTPUT_FORMAT}"
export TX_PARAMS="--node ${RPC_URL} --keyring-backend ${KEYRING_BACKEND} --output ${OUTPUT_FORMAT} --chain-id ${CHAIN_ID} --gas-prices ${GAS_PRICES} --yes"

# Accounts
export BASE_ACCOUNT_1="base_account_1"
export BASE_ACCOUNT_2="base_account_2"
export BASE_VESTING_ACCOUNT="base_vesting_account"
export CONTINOUS_VESTING_ACCOUNT="continous_vesting_account"
export DELAYED_VESTING_ACCOUNT="delayed_vesting_account"
export PERIODIC_VESTING_ACCOUNT="periodic_vesting_account"
