#!/bin/bash

set -euox pipefail

for i in 1 2 3 4 5; do
    if eval "[[ $(cheqd-noded status -n 'tcp://localhost:26677' 2>&1 | wc -l) == 1 ]] && echo 'New node returns status!'"; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned failure. Retrying..."
        sleep 60
    fi
done

exit 1