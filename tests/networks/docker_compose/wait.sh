#!/bin/bash

cmd="$1"

for i in 1 2 3; do
    if eval $cmd; then
        echo "success"
        exit 0
    else
        echo "fail"
        sleep 15
    fi
done

exit 1