#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

. "${BASE_DIR}/common.sh"


TEST_IBC=false
TEST_MODULES=false



echo "########## Cleanup ##########"

# (cd ibc && bash tear-down.sh)
bash "tear-down.sh" # Main tear down should be run last beacuse it owns the network



echo "########## Setup ##########"

bash "setup.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "########## IBC setup ##########"
    (cd ibc && bash setup.sh)
fi



echo "########## Before upgrade ##########"

if [[ "$TEST_MODULES" == "true" ]]
then
    echo "### Before upgrade modules ###"
    bash before-upgrade.sh
fi

if [[ "$TEST_IBC" == "true" ]]
then
    echo "### IBC before upgrade ###"
    (cd ibc && bash before-upgrade.sh)
fi



echo "########## Upgrade ##########"

bash "upgrade.sh"



echo "########## After upgrade ##########"

if [[ "$TEST_MODULES" == "true" ]]
then
    echo "### After upgrade modules ###"
    bash after-upgrade.sh
fi

if [[ "$TEST_IBC" == "true" ]]
then
    echo "### IBC after upgrade ###"
    (cd ibc && bash after-upgrade.sh)
fi



echo "########## Tear down ##########"

bash "tear-down.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "########## IBC tear down ##########"
    (cd ibc && bash tear-down.sh)
fi
