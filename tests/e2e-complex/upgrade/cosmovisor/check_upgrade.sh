#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

if [ "$(cheqd-noded version)" != "$UPGRADE_VERSION_COSMOVISOR" ]; then
    echo "Looks like it was not upgraded"
    exit 1
fi

CURRENT_HEIGHT=$(cheqd-noded status 2>&1 | jq ".SyncInfo.latest_block_height")
bash wait.sh "[[ cheqd-noded status 2>&1 | jq \".SyncInfo.latest_block_height\" > $UPGRADE_HEIGHT ]] && echo \"Current hieght more then upgrade\"'