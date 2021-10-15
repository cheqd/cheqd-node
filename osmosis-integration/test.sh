#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi


docker-compose exec osmosis osmosisd keys list
docker-compose exec cheqd cheqd-noded keys list