#!/bin/bash

set -euox pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

pushd "$DIR/../../../../docker/localnet"

# Generate configs (make sure old binary is installed locally)
bash gen-network-config.sh
sudo chown -R 1000:1000 network-config

# Import keys
bash import-keys.sh

# Start network
docker compose --env-file mainnet-latest.env up --detach --no-build

# TODO: Get rid of this sleep.
sleep 5

# Copy keys

docker compose --env-file mainnet-latest.env cp network-config/validator-0/config/config.toml validator-0:/home/cheqd/.cheqdnode/config/config.toml
docker compose --env-file mainnet-latest.env cp network-config/validator-1/config/config.toml validator-1:/home/cheqd/.cheqdnode/config/config.toml
docker compose --env-file mainnet-latest.env cp network-config/validator-2/config/config.toml validator-2:/home/cheqd/.cheqdnode/config/config.toml
docker compose --env-file mainnet-latest.env cp network-config/validator-3/config/config.toml validator-3:/home/cheqd/.cheqdnode/config/config.toml

# Restore permissions
sudo docker compose --env-file mainnet-latest.env exec --user root validator-0 chown -R cheqd:cheqd /home/cheqd
sudo docker compose --env-file mainnet-latest.env exec --user root validator-1 chown -R cheqd:cheqd /home/cheqd
sudo docker compose --env-file mainnet-latest.env exec --user root validator-2 chown -R cheqd:cheqd /home/cheqd
sudo docker compose --env-file mainnet-latest.env exec --user root validator-3 chown -R cheqd:cheqd /home/cheqd

popd
