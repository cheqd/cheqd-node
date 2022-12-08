#!/bin/bash

set -euox pipefail

all_validators_cmd='cheqd-noded query staking validators --node http://localhost:26657'

amount_bonded="$(${all_validators_cmd} | grep -c BOND_STATUS_BONDED | xargs)"
amount_all="$(${all_validators_cmd} | grep -c status | xargs)"

if [ "${amount_all}" != "${amount_bonded}" ]; 
then 
    exit 1
fi