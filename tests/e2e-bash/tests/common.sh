#!/bin/bash

set -euox pipefail

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
  --gas ${GAS} --gas-adjustment ${GAS_ADJUSTMENT} --gas-prices ${GAS_PRICES} --yes"

# Accounts
export BASE_ACCOUNT_1="base_account_1"
export BASE_ACCOUNT_2="base_account_2"
export BASE_VESTING_ACCOUNT="base_vesting_account"
export CONTINOUS_VESTING_ACCOUNT="continous_vesting_account"
export DELAYED_VESTING_ACCOUNT="delayed_vesting_account"
export PERIODIC_VESTING_ACCOUNT="periodic_vesting_account"

function random_string() {
  LENGTH=${1:-16} # Default LENGTH is 16
  ALPHABET=${2:-"123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"} # Default is base58

  yes $RANDOM | base64 | tr -dc "$ALPHABET" | head -c "${LENGTH}"
  return 0
}

function assert_eq() {
    ACTUAL=$1
    EXPECTED=$2

    if [[ "${ACTUAL}" != "${EXPECTED}" ]]
    then
      echo "Values are not equal. Actual: ${ACTUAL}, expected: ${EXPECTED}."
      return 1
    fi

    return 0
}

function assert_json_eq() {
    ACTUAL=$1
    EXPECTED=$2

    assert_eq "$(echo "${ACTUAL}" | jq --sort-keys ".")" "$(echo "${EXPECTED}" | jq --sort-keys ".")"
}

function assert_tx_successful() {
    OUTPUT=$1
    assert_eq "$(echo "${OUTPUT}" | jq -r ".code")" "0"
}

function assert_tx_code() {
    OUTPUT=$1
    CODE=$2
    assert_eq "$(echo "${OUTPUT}" | jq -r ".code")" "$CODE"
}

function assert_str_contains() {
    STR=$1
    SUBSTR=$2

    if [[ $STR == *$SUBSTR* ]]; then
      return 0
    fi

    return 1
}
