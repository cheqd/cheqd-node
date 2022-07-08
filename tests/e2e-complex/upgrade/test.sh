#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh


TEST_IBC=true


echo "########## Cleanup"

bash "tear-down.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "########## Cleanup IBC"
    # (cd ibc && bash ibc-tear-down.sh)
fi



echo "########## Setup"

bash "setup.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "########## IBC setup"
    # (cd ibc && bash ibc_setup.sh)
fi



echo "########## Before upgrade"

bash "before-upgrade.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "IBC before upgrade"
fi


echo "########## Upgrade"

bash "upgrade.sh"



echo "########## After upgrade"

bash "after-upgrade.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "IBC after upgrade"
fi


echo "########## Tear down"

bash "tear-down.sh"

if [[ "$TEST_IBC" == "true" ]]
then
    echo "########## IBC tear down"
    # (cd ibc && bash ibc-tear-down.sh)
fi
