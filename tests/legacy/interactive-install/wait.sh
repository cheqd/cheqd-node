#!/bin/bash

# TODO: Improve or get rid

set -euox pipefail

cmd="$1"

for _ in 1 2 3 4 5 6 7 8 9; do
    sleep 20
    if eval "$cmd"; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned failure. Retrying..."
    fi
done

exit 1
