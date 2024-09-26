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

CHANNEL_ID="channel-0"

export OWNER=$(osmosisd keys show osmosis-user -a --keyring-backend test )
echo $OWNER=$(osmosisd keys show osmosis-user -a --keyring-backend test )


info "store crosschain_registry wasm"
#store the SetupCrosschainRegistry
RES=$(osmosisd tx wasm store /osmosis/bytecode/crosschain_registry.wasm --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --gas 10000000 --fees 25000uosmo --yes)
assert_tx_successful $RES

sleep 5

# instantiate
INIT_SWAPREGISTRY='{"owner":"'$OWNER'"}'
info "initiate crosschain_registry wasm"
RES=$(osmosisd tx wasm instantiate 1 "$INIT_SWAPREGISTRY" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --label "test" --no-admin --yes --fees 5000uosmo)
assert_tx_successful $RES
SWAPREGISTRY_ADDRESS=osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9
sleep 5

# execute
info "execute crosschain_registry modify_chain_channel_links"
EXE_MSG='{"modify_chain_channel_links": {"operations": [{"operation": "set","source_chain": "cheqd","destination_chain": "osmosis","channel_id": "channel-0"},{"operation": "set","source_chain": "osmosis","destination_chain": "cheqd","channel_id": "channel-0"}]}}'
RES=$(osmosisd tx wasm execute "$SWAPREGISTRY_ADDRESS" "$EXE_MSG" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

info "execute crosschain_registry modify_bech32_prefixes"
PREFIX='{"modify_bech32_prefixes": {"operations": [{"operation": "set", "chain_name": "cheqd", "prefix": "cheqd"},{"operation": "set", "chain_name": "osmosis", "prefix": "osmo"}]}}'
RES=$(osmosisd tx wasm execute "$SWAPREGISTRY_ADDRESS" "$PREFIX" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

info "execute crosschain_registry propose_pfm"
FEEABS_PFM='{"propose_pfm":{"chain": "cheqd"}}'
RES=$(osmosisd tx wasm execute "$SWAPREGISTRY_ADDRESS" "$FEEABS_PFM" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

# GAIA_PFM='{"propose_pfm":{"chain": "gaiad-t1"}}'
# RES=$(osmosisd tx wasm execute "$SWAPREGISTRY_ADDRESS" "$GAIA_PFM" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

info "store swaprouter wasm"
# Store the swaprouter contract
RES=$(osmosisd tx wasm store /osmosis/bytecode/swaprouter.wasm --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --gas 10000000 --fees 25000uosmo --yes)
assert_tx_successful $RES
sleep 5

info "instantiate swaprouter"
# Instantiate the swaprouter contract
INIT_SWAPROUTER='{"owner":"'$OWNER'"}'
RES=$(osmosisd tx wasm instantiate 2 "$INIT_SWAPROUTER" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --label "test" --no-admin --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

SWAPROUTER_ADDRESS=osmo1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqvlx82r
echo $SWAPROUTER_ADDRESS

info "execute swaprouter set_route"
# Configure the swaprouter
CONFIG_SWAPROUTER='{"set_route":{"input_denom":"uosmo","output_denom":"ibc/19D515BB82FEAFCEC357D04C1B75DBF123DB5479041A9BAE38BFCF295477404D","pool_route":[{"pool_id":"1","token_out_denom":"ibc/19D515BB82FEAFCEC357D04C1B75DBF123DB5479041A9BAE38BFCF295477404D"}]}}'
RES=$(osmosisd tx wasm execute $SWAPROUTER_ADDRESS "$CONFIG_SWAPROUTER" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis -y --fees 5000uosmo)
assert_tx_successful $RES
sleep 5

info "store crosschain_swaps wasm"
# Store the crosschainswap contract
RES=$(osmosisd tx wasm store /osmosis/bytecode/crosschain_swaps.wasm --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --gas 10000000 --fees 25000uosmo --yes)
assert_tx_successful $RES
sleep 10

info "instantiate swaprouter"
# Instantiate the crosschainswap contract
INIT_CROSSCHAIN_SWAPS='{"swap_contract":"'$SWAPROUTER_ADDRESS'","governor": "'$OWNER'"}'
RES=$(osmosisd tx wasm instantiate 3 "$INIT_CROSSCHAIN_SWAPS" --keyring-backend=test --home=$HOME/.osmosisd --from osmosis-user --chain-id osmosis --label "test" --no-admin --yes --fees 5000uosmo)
assert_tx_successful $RES
sleep 5
CROSSCHAIN_SWAPS_ADDRESS=osmo17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgs5yczr8
# 1