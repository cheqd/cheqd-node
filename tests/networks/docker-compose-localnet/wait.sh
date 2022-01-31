#!/bin/bash

# TODO: Improve or get rid

set -euox pipefail

cmd="$1"

for _ in 1 2 3; do
    sleep 60
    if eval "$cmd"; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned failure. Retrying..."
    fi
done

exit 1
