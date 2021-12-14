#!/bin/bash

set -euox pipefail


TARGET_HEIGHT=${1:-1} # Default is 1
WAIT_TIME=${2:-60}    # In seconds, default - 60
RPC_ENDPOINT=${3:-"http://localhost:26659/status"}

WAITED=0
WAIT_INTERVAL=1

while true; do
    sleep "${WAIT_INTERVAL}"
    WAITED=$((WAITED + WAIT_INTERVAL))

    CURRENT_HEIGHT=$(curl -s "${RPC_ENDPOINT}" | jq -r ".result.sync_info.latest_block_height" || echo "-1")

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
