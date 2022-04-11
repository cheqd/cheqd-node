#!/bin/bash

set -euox pipefail

all_validators_cmd='cheqd-noded query staking validators --node http://3.19.251.6:26657'

amount_bonded="$(${all_validators_cmd} | grep BOND_STATUS_BONDED | wc -l | xargs)"
amount_all="$(${all_validators_cmd} | grep status | wc -l | xargs)"

if [ "${amount_all}" != "${amount_bonded}" ]; 
then 
    exit 1
fi