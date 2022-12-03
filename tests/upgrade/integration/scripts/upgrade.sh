#!/bin/bash

set -euox pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

pushd "$DIR/../../../../docker/localnet"

# Stop network
docker compose --env-file mainnet-latest.env down

# Start network
# docker compose --env-file build-latest.env up --detach --no-build

popd
