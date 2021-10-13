#!/bin/bash

cmd="$1"

for i in 1 2 3 4; do
    if eval $cmd; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned fail. Retrying..."
        sleep 30
    fi
done

exit 1