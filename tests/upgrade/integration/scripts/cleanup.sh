#!/bin/bash

set -euox pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

pushd "$DIR/../../../docker/localnet"

rm -rf network-config
docker compose --env-file mainnet-latest.env down --remove-orphans --volumes

popd
