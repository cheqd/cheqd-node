#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. ../common.sh

if [ "$(cheqd-noded version 2>&1)" != "${UPGRADE_VERSION_COSMOVISOR}" ]; then
    echo "Looks like it was not upgraded"
    exit 1
fi

if [[ $(cheqd-noded status --node tcp://127.0.0.1:26677 2>&1 | jq ".SyncInfo.latest_block_height" | tr -d '\"' ) < $UPGRADE_HEIGHT ]]; 
then
    echo "Current height less then upgrade"
    exit 1
fi