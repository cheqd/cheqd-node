#!/bin/bash

set -euo pipefail

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

# cheqd_noded docker wrapper

cheqd_noded_docker() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    cheqd-node "$@"
}

cheqd_noded_docker_alt() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    "$@"
}
