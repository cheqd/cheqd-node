#!/bin/bash

set -euox pipefail

for i in 1 2 3 4 5; do
    if eval "[[ $(curl -s localhost:26657/block | sed -nr 's/.*(signature": null).*/\1/p' | wc -l) == 0 ]] && echo 'There are no null signatures in block!'"; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned failure. Retrying..."
        sleep 60
    fi
done

exit 1