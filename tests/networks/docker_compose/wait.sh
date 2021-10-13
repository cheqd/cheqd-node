#!/bin/bash

cmd="$1"

for i in 1 2 3 4 5; do
    if eval $cmd; then
        echo "Waiter returned success!"
        exit 0
    else
        echo "Waiter returned fail. Retrying..."
        sleep 15
    fi
done

exit 1