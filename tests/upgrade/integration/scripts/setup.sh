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

# Copy keys
docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test validator-0:/home/cheqd/.cheqdnode
docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test validator-1:/home/cheqd/.cheqdnode
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test validator-2:/home/cheqd/.cheqdnode
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test validator-3:/home/cheqd/.cheqdnode

# docker compose --env-file mainnet-latest.env cp network-config/validator-0/config validator-0:/home/cheqd/.cheqdnode
# docker compose --env-file mainnet-latest.env cp network-config/validator-1/config validator-1:/home/cheqd/.cheqdnode
# docker compose --env-file mainnet-latest.env cp network-config/validator-2/config validator-2:/home/cheqd/.cheqdnode
# docker compose --env-file mainnet-latest.env cp network-config/validator-3/config validator-3:/home/cheqd/.cheqdnode


cp network-config/validator-1/keyring-test/* ~/.cheqdnode/keyring-test/
cp network-config/validator-2/keyring-test/* ~/.cheqdnode/keyring-test/
cp network-config/validator-3/keyring-test/* ~/.cheqdnode/keyring-test/
cp network-config/validator-0/keyring-test/* ~/.cheqdnode/keyring-test/

docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test validator-0:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test validator-0:/home/cheqd/temp-keyring-test2
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test validator-0:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env exec validator-0 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-0 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-0 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'


docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test validator-1:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test validator-1:/home/cheqd/temp-keyring-test2
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test validator-1:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env exec validator-1 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-1 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-1 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'


docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test validator-2:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test validator-2:/home/cheqd/temp-keyring-test2
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test validator-2:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env exec validator-2 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-2 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-2 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'

docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test validator-3:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test validator-3:/home/cheqd/temp-keyring-test2
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test validator-3:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env exec validator-3 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-3 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec validator-3 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'


docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test observer-0:/home/cheqd/.cheqdnode/keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test  observer-0:/home/cheqd/temp-keyring-test1
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test  observer-0:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test observer-0:/home/cheqd/temp-keyring-test2

docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test observer-0:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env exec observer-0 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec observer-0 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec observer-0 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec observer-0 bash -c 'mv -n /home/cheqd/temp-keyring-test1/* /home/cheqd/.cheqdnode/keyring-test/'


docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test seed-0:/home/cheqd/.cheqdnode/keyring-test
docker compose --env-file mainnet-latest.env cp network-config/validator-1/keyring-test  seed-0:/home/cheqd/temp-keyring-test1
docker compose --env-file mainnet-latest.env cp network-config/validator-3/keyring-test  seed-0:/home/cheqd/temp-keyring-test3
docker compose --env-file mainnet-latest.env cp network-config/validator-2/keyring-test seed-0:/home/cheqd/temp-keyring-test2


docker compose --env-file mainnet-latest.env cp network-config/validator-0/keyring-test seed-0:/home/cheqd/temp-keyring-test
docker compose --env-file mainnet-latest.env exec seed-0 bash -c 'mv -n /home/cheqd/temp-keyring-test/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec seed-0 bash -c 'mv -n /home/cheqd/temp-keyring-test2/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec seed-0 bash -c 'mv -n /home/cheqd/temp-keyring-test3/* /home/cheqd/.cheqdnode/keyring-test/'
docker compose --env-file mainnet-latest.env exec seed-0 bash -c 'mv -n /home/cheqd/temp-keyring-test1/* /home/cheqd/.cheqdnode/keyring-test/'



# docker compose --env-file mainnet-latest.env cp validator-1:/home/cheqd/.cheqdnode/keyring-test /tmp/validator-1-keyring
# docker compose --env-file mainnet-latest.env cp validator-2:/home/cheqd/.cheqdnode/keyring-test /tmp/validator-2-keyring
# docker compose --env-file mainnet-latest.env cp validator-3:/home/cheqd/.cheqdnode/keyring-test /tmp/validator-3-keyring



# docker compose --env-file mainnet-latest.env cp /tmp/validator-1-keyring/* validator-0:/home/cheqd/.cheqdnode/keyring-test/
# docker compose --env-file mainnet-latest.env cp /tmp/validator-2-keyring/* validator-0:/home/cheqd/.cheqdnode/keyring-test/
# docker compose --env-file mainnet-latest.env cp /tmp/validator-3-keyring/* validator-0:/home/cheqd/.cheqdnode/keyring-test/



# Restore permissions
docker compose --env-file mainnet-latest.env exec --user root validator-0 chown -R cheqd:cheqd /home/cheqd
docker compose --env-file mainnet-latest.env exec --user root validator-1 chown -R cheqd:cheqd /home/cheqd
docker compose --env-file mainnet-latest.env exec --user root validator-2 chown -R cheqd:cheqd /home/cheqd
docker compose --env-file mainnet-latest.env exec --user root validator-3 chown -R cheqd:cheqd /home/cheqd

popd