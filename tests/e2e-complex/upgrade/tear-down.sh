#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh


# Stop docker compose
docker_compose_down
# Clean environment after test
clean_env

sudo rm -rf "network-config"
sudo rm -rf ".cheqdnode"
rm "txs.hashes" 2> /dev/null || true
rm "resource_data.json" 2> /dev/null || true
