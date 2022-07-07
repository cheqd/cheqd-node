#!/bin/bash

set -euox pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function info() {
    printf "${GREEN}[info] %s${NC}\n" "${1}"
}

function err() {
    printf "${RED}[err] %s${NC}\n" "${1}"
}

function assert_tx_successful() {
  RES="$1"

  if [[ $(echo "${RES}" | jq --raw-output '.code') == 0 ]]
  then
    info "tx successful"
  else
    err "non zero tx return code"
    exit 1
  fi
}

function assert_network_running() {
  RES="$1"
  LATEST_HEIGHT=$(echo "${RES}" | jq --raw-output '.SyncInfo.latest_block_height')
  info "latest height: ${LATEST_HEIGHT}"

  if [[ $LATEST_HEIGHT -gt 1 ]]
  then
      info "network is running"
  else
      err "network is not running"
      exit 1
  fi
}

info "Tear down" # ---
docker-compose down --timeout 20
