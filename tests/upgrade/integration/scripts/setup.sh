#!/bin/bash

set -euox pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

pushd "$DIR/../../../../docker/localnet"

# Generate configs (make sure old binary is installed locally)
bash gen-network-config.sh
# sudo chown -R 1000:1000 network-config

# Import keys
bash import-keys.sh

# Start network
docker compose --env-file mainnet-latest.env up --detach --no-build

# TODO: Get rid of this sleep.
sleep 5

sudo docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test validator-0:/home/keyring-test
sudo docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test validator-1:/home/keyring-test
sudo docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test validator-2:/home/keyring-test
sudo docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test validator-3:/home/keyring-test

# # Restore permissions
sudo docker compose --env-file mainnet-latest.env exec --user root validator-0 chown -R cheqd:cheqd /home
sudo docker compose --env-file mainnet-latest.env exec --user root validator-1 chown -R cheqd:cheqd /home
sudo docker compose --env-file mainnet-latest.env exec --user root validator-2 chown -R cheqd:cheqd /home
sudo docker compose --env-file mainnet-latest.env exec --user root validator-3 chown -R cheqd:cheqd /home

sudo docker compose --env-file mainnet-latest.env exec validator-0 bash -c 'cp -r "/home/keyring-test" "$HOME/.cheqdnode/"'
sudo docker compose --env-file mainnet-latest.env exec validator-1 bash -c 'cp -r "/home/keyring-test" "$HOME/.cheqdnode/"'
sudo docker compose --env-file mainnet-latest.env exec validator-2 bash -c 'cp -r "/home/keyring-test" "$HOME/.cheqdnode/"'
sudo docker compose --env-file mainnet-latest.env exec validator-3 bash -c 'cp -r "/home/keyring-test" "$HOME/.cheqdnode/"'

popd
