#!/bin/bash

set -euox pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

pushd "$DIR/../../../../docker/localnet"

# Stop network
docker compose --env-file upgrade-v4-latest.env down

# Start network with the latest build that includes the minor upgrade handler
docker compose --env-file upgrade-v4-1-latest.env up --detach --no-build
sleep 3
docker compose --env-file upgrade-v4-1-latest.env logs
docker ps
popd
