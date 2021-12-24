#!/bin/bash

# The script that runs a localnet of 4 validators and 2 observers through docker-compose.

set -euox pipefail

# Set up 4 validators + 2 observers node docker pool

./gen_node_configs.sh
./run_docker.sh
./wait.sh '[[ $(curl -s -N localhost:26657/block | jq -cr '"'"'.result.block.last_commit.height'"'"') -gt 1 ]] && echo "Height is more than 1"'

# Add observer

./add_observer.sh
cheqd-noded status -n tcp://localhost:26677 2>&1
./wait.sh '[[ $(cheqd-noded status -n '"'"'tcp://localhost:26677'"'"' 2>&1 | wc -l) == 1 ]] && echo "New node returns status!"'

# Promote observer to validator

bash -x promote_validator.sh
cheqd-noded query staking validators --node "http://localhost:26657" | sed -nr 's/.*status: (.*?).*/\1/p' $1 | while read x; do [[ "BOND_STATUS_BONDED" == $x ]] && echo "Validator's status is bonded!" || exit 1 ; done
./wait.sh '[[ $(curl -s localhost:26657/block | sed -nr '"'"'s/.*signature": (.*?).*/\1/p'"'"' | wc -l) == 5 ]] && echo "There are 5 validators signatures in block!"'
./wait.sh '[[ $(curl -s localhost:26657/block | sed -nr '"'"'s/.*(signature": null).*/\1/p'"'"' | wc -l) == 0 ]] && echo "There are no null signatures in block!"'
