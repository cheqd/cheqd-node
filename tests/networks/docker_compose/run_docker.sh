#!/bin/bash

# The script that exports node ids and runs docker-compose.

set -euox pipefail

if [ "$#" -eq 1 ] ; then
    MOUNT_POINT="$1"
else
    MOUNT_POINT="."
fi

# cheqd_noded docker wrapper

cheqd_noded_docker() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    cheqd-node "$@"
}

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

pushd "$MOUNT_POINT/node_configs/node0"

export NODE_0_ID=$(cheqd_noded_docker tendermint show-node-id | sed 's/\r//g')

popd

docker-compose up -d
