#!/bin/bash

set -euox pipefail

DOCKER_COMPOSE_DIR="../networks/docker-compose-localnet"
CHEQD_IMAGE_FROM="ghcr.io/cheqd/cheqd-node:0.4.0"
CHEQD_IMAGE_TO="cheqd-node"
CHEQD_VERSION_TO=`echo $(git describe --always --tag --match "v*") | sed 's/^v//'`
UPGRADE_NAME="v0.4"
VOTING_PERIOD=30
EXPECTED_BLOCK_SECOND=5
EXTRA_BLOCKS=5
UPGRADE_HEIGHT=$(echo $VOTING_PERIOD/$EXPECTED_BLOCK_SECOND+$EXTRA_BLOCKS | bc)
DEPOSIT_AMOUNT=10000000
CHAIN_ID="cheqd"
CHEQD_USER="cheqd"
FNAME_TXHASHES="txs.hashes"

# cheqd_noded docker wrapper

cheqd_noded_docker() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    --network host \
    -u root \
    -e HOME=/cheqd \
    --entrypoint "cheqd-noded" \
    ${CHEQD_IMAGE_FROM} $@
}

# Parameters
# $1 - Name of container to run command inside
# $2 - The full command to run
function docker_exec () {
    NODE_CONTAINER="$1"

    docker exec -u $CHEQD_USER $NODE_CONTAINER "${@:2}"
}

# Parameters
# $1 - Version of base image
# $2 - Root path for making directories for volumes
function docker_compose_up () {
    CURR_DIR=$(pwd)
    MOUNT_POINT="$2"
    pushd "node_configs/node0"
    export NODE_0_ID=$(cheqd_noded_docker tendermint show-node-id | sed 's/\r//g')
    export CHEQD_IMAGE_NAME="$1"
    export MOUNT_POINT=$MOUNT_POINT
    docker-compose -f ../../$DOCKER_COMPOSE_DIR/docker-compose.yml --env-file ../../$DOCKER_COMPOSE_DIR/.env up -d
    pushd $CURR_DIR
}

# Stop docker-compose
function docker_compose_down () {
    docker-compose -f $DOCKER_COMPOSE_DIR/docker-compose.yml --env-file $DOCKER_COMPOSE_DIR/.env down 
}

# Clean environment
function clean_env () {
    rm -rf node_configs
    rm -f $FNAME_TXHASHES
}

# Run command using local generated keys from node_configs/client
function local_client_tx () {
    result="$(cheqd_noded_docker $@ --home node_configs/client/.cheqdnode/ --keyring-backend test)"
}

function make_777 () {
    sudo chmod -R 777 node_configs
}



# Transaction related funcs

function random_string() {
  echo $RANDOM | base64 | head -c 20
  return 0
}

function get_addresses () {
    local_client_tx keys list
    addresses=( $(echo "$result" | grep -o 'cheqd1.*') )
    echo "${addresses[@]}"
}

CHEQ_AMOUNT="1ncheq"

# Send tokens from the first address in the list to another one
function send_tokens() {
    get_addresses 

    OP0_ADDRESS=${addresses[0]}
    OP1_ADDRESS=${addresses[1]}
    OP2_ADDRESS=${addresses[2]}
    OP3_ADDRESS=${addresses[3]}

    local_client_tx tx \
                    bank \
                    send $OP0_ADDRESS $OP1_ADDRESS $CHEQ_AMOUNT \
                    --gas auto \
                    --gas-adjustment 1.2 \
                    --gas-prices "25ncheq" \
                    --chain-id $CHAIN_ID \
                    -y
    txhash=$(echo $result | jq ".txhash" | tr -d '"')
    echo $txhash >> $FNAME_TXHASHES
}


function send_did () {

    # Generate Alice identity key
    ALICE_VER_KEY="$(cheqd_noded_docker debug ed25519 random)"
    ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
    ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
    ALICE_VER_PUB_MULTIBASE_58=$(cheqd_noded_docker debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

    # Build CreateDid message
    DID="did:cheqd:testnet:$(random_string)"
    KEY_ID="${DID}#key1"

    MSG_CREATE_DID='{"id":"'${DID}'","verification_method":[{"id":"'${KEY_ID}'","type":"Ed25519VerificationKey2020","controller":"'${DID}'","public_key_multibase":"'${ALICE_VER_PUB_MULTIBASE_58}'"}],"authentication":["'${KEY_ID}'"]}'

    # Post the message
    local_client_tx tx cheqd create-did ${MSG_CREATE_DID} ${KEY_ID} \
        --ver-key "${ALICE_VER_PRIV_BASE_64}" \
        --from operator0 \
        --gas-prices "25ncheq" \
        --chain-id $CHAIN_ID \
        --output json \
        -y


    txhash=$(echo $result | jq ".txhash" | tr -d '"')
    echo $txhash >> $FNAME_TXHASHES
}


function check_tx_hashes () {
    for txhash in $(cat $FNAME_TXHASHES); 
    do
        txhash=$(echo ${txhash} | tr -d '"')
        result=$(cheqd_noded_docker query tx ${txhash} --output json)
        tx_exist=$(echo $result | jq ".code")
        if [ $tx_exist != "0" ] ; then
            echo "Error was in checking tx with hash: $txhash"
            exit 1
        fi
    done    
}
