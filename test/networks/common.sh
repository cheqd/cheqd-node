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
  docker run --rm -it \
    -v "$(pwd)/.cheqdnode":/home/cheqd/.cheqdnode \
    cheqd-node "$@"
}
