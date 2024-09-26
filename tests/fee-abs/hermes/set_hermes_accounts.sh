#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

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

info "Create relayer user on cheqd"
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"
CHEQD_RELAYER_ACCOUNT=$(cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json 2>&1)
CHEQD_RELAYER_ADDRESS=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "${CHEQD_RELAYER_MNEMONIC}" > cheqd_relayer_mnemonic.txt

info "Send some tokens to it"
RES=$(cheqd-noded tx bank send cheqd-user "${CHEQD_RELAYER_ADDRESS}" 500000000000000000ncheq --gas-prices 50ncheq --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

info "Create relayer user on osmosis"
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test 2>&1)
OSMOSIS_RELAYER_ADDRESS=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "${OSMOSIS_RELAYER_MNEMONIC}" > osmo_relayer_mnemonic.txt

info "Send some tokens to it"
RES=$(osmosisd tx bank send osmosis-user "${OSMOSIS_RELAYER_ADDRESS}" 1000000000uosmo --fees 500uosmo --chain-id osmosis -y --keyring-backend test --output json --node http://localhost:26667)
assert_tx_successful "${RES}"

sleep 10 

# Create dirs for keys
mkdir -p ~/.hermes/keys/cheqd/keyring-test
mkdir -p ~/.hermes/keys/osmosis/keyring-test

# Copy keys
cp cheqd_relayer_mnemonic.txt ~/.hermes/keys/cheqd/keyring-test/
cp osmo_relayer_mnemonic.txt ~/.hermes/keys/osmosis/keyring-test/

# Import keys
hermes keys add --chain cheqd --mnemonic-file ~/.hermes/keys/cheqd/keyring-test/cheqd_relayer_mnemonic.txt --key-name cheqd-key --overwrite
hermes keys add --chain osmosis --mnemonic-file ~/.hermes/keys/osmosis/keyring-test/osmo_relayer_mnemonic.txt --key-name osmosis-key --overwrite

info "Open channel"

hermes create channel --a-chain cheqd --b-chain osmosis --a-port transfer --b-port transfer --new-client-connection --yes
# hermes create connection --a-chain cheqd --b-chain osmosis --new-client-connection --yes
hermes create channel --a-chain osmosis --b-chain cheqd --a-port icqhost --b-port feeabs --new-client-connection --yes
# hermes create channel --a-chain cheqd --b-chain osmosis --a-port icqhost --b-port feeabs --new-client-connection --yes

hermes start