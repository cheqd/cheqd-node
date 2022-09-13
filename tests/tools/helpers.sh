#!/bin/bash

set -euo pipefail


# Constants

LOCALNET_NETWORK="localnet"
LOCALNET_PATH="$(git rev-parse --show-toplevel)/docker/localnet"

# Localnet

function in_localnet_path() {
    pushd "${LOCALNET_PATH}" > /dev/null
    "$@"
    popd > /dev/null
}

function localnet_compose() {
    in_localnet_path docker compose "${@}"
}

# Helpers

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

function compose_wait_for_chain_height() {
    SERVICE=$1  # For example: node0
    BINARY=$2   # For example: "cheqd-noded" or "osmosis"

    TARGET_HEIGHT=${3:-2} # Default is 2
    WAIT_TIME=${4:-60}    # In seconds, default - 60

    WAITED=0
    WAIT_INTERVAL=1

    while true; do
        sleep "${WAIT_INTERVAL}"
        WAITED=$((WAITED + WAIT_INTERVAL))
        
        DOCKER_COMPOSE_STATUS="$(docker compose exec ${SERVICE} ${BINARY} status 2>&1)"
        CURRENT_HEIGHT="$(echo ${DOCKER_COMPOSE_STATUS} | jq -r '.SyncInfo.latest_block_height' || echo '-1')" 

        if ((CURRENT_HEIGHT >= TARGET_HEIGHT)); then
            echo "Height ${TARGET_HEIGHT} is reached in $WAITED seconds"
            break
        fi

        if ((WAITED > WAIT_TIME)); then
            echo "Height $TARGET_HEIGHT is not reached in $WAIT_TIME seconds"
            exit 1
        fi

        echo "Waiting for height: $TARGET_HEIGHT... Current height: $CURRENT_HEIGHT, time waited: $WAITED, time limit: $WAIT_TIME."
    done
}
