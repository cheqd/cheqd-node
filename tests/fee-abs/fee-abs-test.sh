#!/bin/bash
# shellcheck disable=SC2086

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

  if [[ $(echo "${RES}" | jq --raw-output '.code') == 0 ]]; then
    info "tx successful"
  else
    err "non zero tx return code"
    exit 1
  fi
}

function exit_if_tx_successful() {
  RES="$1"

  if [[ $(echo "${RES}" | jq --raw-output '.code') == 0 ]]; then
    info "tx successful"
    exit 0
  else
    err "tx failed, retrying"
  fi
}

function assert_network_running() {
  RES="$1"
  LATEST_HEIGHT=$(echo "${RES}" | jq --raw-output '.SyncInfo.latest_block_height')
  info "latest height: ${LATEST_HEIGHT}"

  if [[ $LATEST_HEIGHT -gt 1 ]]; then
    info "network is running"
  else
    err "network is not running"
    exit 1
  fi
}

function assert_network_running_comet_v38_or_above() {
  RES="$1"
  LATEST_HEIGHT=$(echo "${RES}" | jq --raw-output '.sync_info.latest_block_height')
  info "latest height: ${LATEST_HEIGHT}"

  if [[ $LATEST_HEIGHT -gt 1 ]]; then
    info "network is running"
  else
    err "network is not running"
    exit 1
  fi
}

info "Cleanup"
docker compose down --volumes --remove-orphans

info "Running cheqd network"
docker compose up -d cheqd
docker compose cp ./cheqd cheqd:/
docker compose exec cheqd bash /cheqd/cheqd-init.sh
docker compose exec -d cheqd cheqd-noded start

info "Running osmosis network"
docker compose up -d osmosis
docker compose cp ./osmosis osmosis:/
docker compose exec osmosis apk add bash jq
docker compose exec osmosis bash /osmosis/osmosis-init.sh
docker compose exec -d osmosis osmosisd start

info "Waiting for chains"
# TODO: Get rid of this
sleep 10

info "Checking statuses"
CHEQD_STATUS=$(docker compose exec cheqd cheqd-noded status 2>&1)
assert_network_running_comet_v38_or_above "${CHEQD_STATUS}"

OSMOSIS_STATUS=$(docker compose exec osmosis osmosisd status 2>&1)
assert_network_running_comet_v38_or_above "${OSMOSIS_STATUS}"

info "Create relayer user on cheqd" # ---
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"
CHEQD_RELAYER_ACCOUNT=$(docker compose exec cheqd cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json 2>&1)
CHEQD_RELAYER_ADDRESS=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "${CHEQD_RELAYER_MNEMONIC}" >cheqd_relayer_mnemonic.txt

info "Send some cheqd tokens to it" # ---
RES=$(docker compose exec cheqd cheqd-noded tx bank send cheqd-user "${CHEQD_RELAYER_ADDRESS}" 500000000000000000ncheq --fees 2000000000ncheq --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

info "Create relayer user on osmosis" # ---
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(docker compose exec osmosis osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test 2>&1)
OSMOSIS_RELAYER_ADDRESS=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "${OSMOSIS_RELAYER_MNEMONIC}" >osmo_relayer_mnemonic.txt

info "Send some osmo tokens to it" # ---
RES=$(docker compose exec osmosis osmosisd tx bank send osmosis-user "${OSMOSIS_RELAYER_ADDRESS}" 1000000000uosmo --fees 500uosmo --chain-id osmosis -y --keyring-backend test --output json)
assert_tx_successful "${RES}"
sleep 10 # Wait for state

info "Import accounts in hermes" # ---
docker compose up -d hermes

# Create dirs for keys
docker compose exec --user root hermes mkdir -p /home/hermes/.hermes/keys/cheqd/keyring-test
docker compose exec --user root hermes mkdir -p /home/hermes/.hermes/keys/osmosis/keyring-test

# Hand over ownership to hermes user
docker compose exec --user root hermes chown -R hermes:hermes /home/hermes/.hermes/keys

# Copy keys
docker compose cp cheqd_relayer_mnemonic.txt hermes:/home/hermes
docker compose cp osmo_relayer_mnemonic.txt hermes:/home/hermes

# Import keys
docker compose exec hermes hermes keys add --chain cheqd --mnemonic-file cheqd_relayer_mnemonic.txt --key-name cheqd-key
docker compose exec hermes hermes keys add --chain osmosis --mnemonic-file osmo_relayer_mnemonic.txt --key-name osmosis-key

info "Open channel" # ---
docker compose exec hermes hermes create channel --a-chain cheqd --b-chain osmosis --a-port transfer --b-port transfer --new-client-connection --yes
docker compose exec hermes hermes create channel --a-chain osmosis --b-chain cheqd --a-port icqhost --b-port feeabs --new-client-connection --yes
info "Start hermes" # ---
docker compose exec -d hermes hermes start

info "Deploy the smart contracts in osmosis"
docker compose cp osmosis/deploy_osmosis_contract.sh osmosis:/osmosis/deploy_osmosis_contract.sh
docker compose exec osmosis bash /osmosis/deploy_osmosis_contract.sh

CHEQD_USER_ADDRESS=$(docker compose exec cheqd cheqd-noded keys show --address cheqd-user --keyring-backend test | tr -d '\r')
OSMOSIS_USER_ADDRESS=$(docker compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | tr -d '\r')

CHEQD_RELAYER_ADDRESS=$(docker compose exec cheqd cheqd-noded keys show --address cheqd-relayer --keyring-backend test | tr -d '\r')
OSMOSIS_RELAYER_ADDRESS=$(docker compose exec osmosis osmosisd keys show --address osmosis-relayer --keyring-backend test | tr -d '\r')

info "Transfer cheqd -> osmosis" # ---
PORT="transfer"
CHANNEL="channel-0"
docker compose exec cheqd cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL "$OSMOSIS_USER_ADDRESS" 10000000000ncheq --from cheqd-user --chain-id cheqd --gas-prices 10000ncheq --keyring-backend test -y
sleep 30 # Wait for relayer

info "Get balances" # ---
CHEQD_BALANCE_2=$(docker compose exec cheqd cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json 2>&1)
BALANCES=$(docker compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS" --output json 2>&1)

info "Denom trace" # ---
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_TRACE=$DENOM
docker compose exec osmosis osmosisd query ibc-transfer denom-trace "$DENOM_TRACE"

info "Send 500000OSMO to cheqd"
docker compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL "$CHEQD_USER_ADDRESS" 500000000000uosmo --from osmosis-user --chain-id osmosis --fees 500uosmo --keyring-backend test -y
sleep 30

CHEQD_BALANCE_2=$(docker compose exec cheqd cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)

info "balances before"
echo $CHEQD_BALANCE_2

info "create pool"
# create pool
TX_HASH=$(docker compose exec osmosis osmosisd tx gamm create-pool --pool-file /osmosis/pool.json --from $OSMOSIS_USER_ADDRESS --keyring-backend test --fees 100000000uosmo --gas auto --gas-adjustment 2 -y --chain-id osmosis --output json | jq -r '.txhash')
echo "tx hash: $TX_HASH"
sleep 5

POOL_ID=$(docker compose exec osmosis osmosisd q tx $TX_HASH --output json | jq -r '.events[] | select(.type == "pool_created") | .attributes[] | select(.key == "pool_id") | .value')
echo "pool id: $POOL_ID"

info "enable fee abs"
RES=$(docker compose exec cheqd cheqd-noded tx gov submit-proposal /cheqd/proposal.json --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 10000ncheq --gas 350000)
assert_tx_successful "${RES}"
sleep 5

info "vote on param change proposal"
RES=$(docker compose exec cheqd cheqd-noded tx gov vote 1 yes --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 10000ncheq --gas 350000)
assert_tx_successful "${RES}"
sleep 5

info "add host zone config"
RES=$(docker compose exec cheqd cheqd-noded tx gov submit-proposal /cheqd/host_zone.json --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 10000ncheq --gas 350000)
assert_tx_successful "${RES}"
sleep 5
RES=$(docker compose exec cheqd cheqd-noded tx gov vote 2 yes --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 10000ncheq --gas 350000)
assert_tx_successful "${RES}"
sleep 5

info "fund fee-abs module account with IBC osmo"
RES=$(docker compose exec cheqd cheqd-noded tx feeabs fund 1000000000ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518 --from $CHEQD_USER_ADDRESS --gas-prices 10000ncheq --gas 350000 --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"
sleep 5

info "fund fee-abs module account with native asset, cheq"
RES=$(docker compose exec cheqd cheqd-noded tx feeabs fund 200000000000000000ncheq --from $CHEQD_USER_ADDRESS --gas-prices 10000ncheq --gas 350000 --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

info "wait for exchange rate to be updated"
sleep 150

info "pay fees using osmo in cheqd (recursively)"
# shellcheck disable=SC2034
for i in {1..20}; do
  RES=$(docker compose exec cheqd cheqd-noded tx bank send cheqd-user "$CHEQD_RELAYER_ADDRESS" 50000000ncheq --fees 200000000000ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518 --chain-id cheqd -y --keyring-backend test)
  exit_if_tx_successful "${RES}"
  sleep 6
done
