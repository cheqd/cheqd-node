#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh


TEST_IBC=true


echo "Cleanup"

docker compose down --volumes --remove-orphans # TODO: Replace
sudo rm -rf ".cheqdnode"
sudo rm -rf "node_configs"

bash "prepare.sh"


if [[ "$TEST_IBC" == "true" ]]
then
    # Setup
    echo "Setup Running IBC environment"
    # bash "ibc-transfer-test.sh"

    # Execute
fi


bash "initiate_upgrade.sh"
bash "upgrade_and_check.sh"


if [[ "$TEST_IBC" == "true" ]]
then
    # Check
    echo "Asserting IBC"
    # bash "ibc-transfer-test.sh"
fi

echo "Cleanup"

# Stop docker compose
docker_compose_down
# Clean environment after test
clean_env
sudo rm -rf "network-config"
