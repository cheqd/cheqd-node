#!/bin/bash

# The script that exports node ids and runs docker-compose.

set -euox pipefail

# cheqd_noded docker wrapper

CHEQD_NODE_IMAGE=cheqd-cli

cheqd_noded_docker() {
  docker run --rm \
    -v "$(pwd)":"/home/cheqd" \
    $CHEQD_NODE_IMAGE "$@"
}

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

pushd "node_configs/node0"

export NODE_0_ID=$(cheqd_noded_docker tendermint show-node-id | sed 's/\r//g')

popd

docker-compose up -d
