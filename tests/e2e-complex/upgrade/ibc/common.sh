#!/bin/bash

set -euox pipefail


BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/../common.sh"

# shellcheck disable=SC1091
. "${BASE_DIR}/../../../tools/helpers.sh"


export COMPOSE_FILE="${BASE_DIR}/docker-compose.yml"

export CHEQD_SERVICE="validator-0"
export CHEQD_USER="operator-0"

function ibc_compose() {
    docker compose --file "${COMPOSE_FILE}" "$@"
}
