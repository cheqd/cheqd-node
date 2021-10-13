#!/bin/bash

set -euox pipefail

for i in 1 2 3 4 5; do
    if eval "[[ $(curl -s -N localhost:26657/block | jq -cr '.result.block.last_commit.height') -gt 1 ]] && echo 'Height is more than 1'"; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned failure. Retrying..."
        sleep 60
    fi
done

exit 1