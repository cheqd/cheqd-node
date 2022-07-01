#!/bin/bash

set -euox pipefail

# TODO: Assert that transactions are successful

CHEQD_IMAGE_FROM="ghcr.io/cheqd/cheqd-cli:0.5.0"
# shellcheck disable=SC2034
CHEQD_IMAGE_TO="cheqd-cli:latest"
# shellcheck disable=SC2034
CHEQD_VERSION_TO=$(git describe --always --tag --match "v*" | sed 's/^v//')
# shellcheck disable=SC2034
UPGRADE_NAME="v0.6"
# shellcheck disable=SC2034
UPGRADE_INFO="{
  \"binaries\": {
    \"linux/amd64\":\"http://10.5.0.10:8000/cheqd-noded.tar.gz\"
  }
}"
VOTING_PERIOD=15
EXPECTED_BLOCK_SECOND=5
EXTRA_BLOCKS=10
# shellcheck disable=SC2034
UPGRADE_HEIGHT=$((VOTING_PERIOD / EXPECTED_BLOCK_SECOND + EXTRA_BLOCKS))
# shellcheck disable=SC2034
DEPOSIT_AMOUNT=10000000
CHAIN_ID="cheqd"
CHEQD_USER="cheqd"
FNAME_TXHASHES="txs.hashes"
AMOUNT_BEFORE="19000000000000000"
CHEQ_AMOUNT="1ncheq"
CHEQ_AMOUNT_NUMBER="1"
# shellcheck disable=SC2034
DID_1="did:cheqd:testnet:1111111111111111"
# shellcheck disable=SC2034
DID_2="did:cheqd:testnet:2222222222222222"
RESOURCE_1="82aadc50-58e4-4e00-bf35-36062c2784be"
CHEQD_HOME="/home/cheqd"


GAS="auto"
GAS_ADJUSTMENT="1.3"
GAS_PRICES="25ncheq"
TX_PARAMS="--gas ${GAS} --gas-adjustment ${GAS_ADJUSTMENT} --gas-prices ${GAS_PRICES} --chain-id ${CHAIN_ID} -y"

# cheqd_noded docker wrapper
cheqd_noded_docker() {
    docker run --rm \
        -v "$(pwd):${CHEQD_HOME}" \
        --network host \
        -u root \
        -e HOME=${CHEQD_HOME} \
        --entrypoint "cheqd-noded" \
        ${CHEQD_IMAGE_FROM} "$@"
}

# Parameters
# $1 - Name of container to run command inside
# $2 - The full command to run
function docker_exec () {
    NODE_CONTAINER="$1"

    docker exec -u $CHEQD_USER "$NODE_CONTAINER" "${@:2}"
}

# Parameters
# $1 - Version of base image
# $2 - Root path for making directories for volumes
function docker_compose_up () {
    pushd "node_configs/node0"
    NODE_0_ID=$(cheqd_noded_docker tendermint show-node-id | sed 's/\r//g')
    export NODE_0_ID="$NODE_0_ID"
    popd
    
    export CHEQD_IMAGE_NAME="$1"
    export MOUNT_POINT="$2"

    docker compose --env-file .env up -d
}

# Stop docker-compose
function docker_compose_down () {
    docker compose --env-file .env down 
}

# Clean environment
function clean_env () {
    sudo rm -rf node_configs
    sudo rm -f $FNAME_TXHASHES
}

# Run command using local generated keys from node_configs/client
function local_client_tx () {
    cheqd_noded_docker "$@" --home node_configs/client/.cheqdnode/ --keyring-backend test
}

function make_775 () {
    sudo chmod -R 777 node_configs
}


# Transaction related funcs

function random_string() {
  echo $RANDOM | base64 | head -c 20
  return 0
}

function get_addresses () {
    all_keys=$(local_client_tx keys list)
    mapfile -t addresses < <(echo "$all_keys" | grep -o 'cheqd1.*')
    # addresses=( $(echo "$all_keys" | grep -o 'cheqd1.*') )
    echo "${addresses[@]}"
}

# Send tokens from the first address in the list to another one
# Input: address to send to
function send_tokens() {
    get_addresses 

    OP_ADDRESS_TO="$1"
    OP0_ADDRESS=${addresses[0]}

    # shellcheck disable=SC2086
    send_res=$(local_client_tx tx \
                    bank \
                    send "$OP0_ADDRESS" "$OP_ADDRESS_TO" $CHEQ_AMOUNT \
                    ${TX_PARAMS})
    txhash="$(echo "$send_res" | jq ".txhash" | tr -d '"')"
    echo "$txhash" >> "$FNAME_TXHASHES"
}

# Send DID
# input: DID to write
function send_did_new () {
    did_to_write=$1

    # Generate Alice identity key
    ALICE_VER_KEY="$(cheqd_noded_docker debug ed25519 random)"
    ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
    ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
    ALICE_VER_PUB_MULTIBASE_58=$(cheqd_noded_docker debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

    # Build CreateDid message
    KEY_ID="${did_to_write}#key1"

    # shellcheck disable=SC2089
    MSG_CREATE_DID='{"id":"'${did_to_write}'","verification_method":[{"id":"'"${KEY_ID}"'","type":"Ed25519VerificationKey2020","controller":"'${did_to_write}'","public_key_multibase":"'${ALICE_VER_PUB_MULTIBASE_58}'"}],"authentication":["'${KEY_ID}'"]}'

    # Post the message
    did=$(local_client_tx tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
        --from operator0 \
        --gas-prices "25ncheq" \
        --chain-id $CHAIN_ID \
        --output json \
        -y)


    txhash=$(echo "$did" | jq ".txhash" | tr -d '"')
    echo "$txhash" >> $FNAME_TXHASHES
}

# Send resource
# input: resource to write
function send_resource_new () {
    collection_id_to_write=$1
    resource_to_write=$2

    # Generate Alice identity key
    ALICE_VER_KEY="$(cheqd_noded_docker debug ed25519 random)"
    ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
    ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
    ALICE_VER_PUB_MULTIBASE_58=$(cheqd_noded_docker debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

    # Build CreateDid message
    KEY_ID="${collection_id_to_write}#key1"

    RESOURCE_NAME="Resource 1"
    RESOURCE_RESOURCE_TYPE="CL-Schema"
    RESOURCE_DATA='{ "content": "test data" }';

    # Post the message
    # shellcheck disable=SC2086
    resource=$(cheqd-noded tx resource create-resource \
    --collection-id ${collection_id_to_write} \
    --resource-id ${resource_to_write} \
    --resource-name "${RESOURCE_NAME}" \
    --resource-type ${RESOURCE_RESOURCE_TYPE} \
    --resource-file <(echo "${RESOURCE_DATA}") \
    "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
    --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

    txhash=$(echo "$resource" | jq ".txhash" | tr -d '"')
    echo "$txhash" >> $FNAME_TXHASHES
}

# Send DID
# input: DID to write
function send_did () {
    did_to_write=$1

    # Generate Alice identity key
    ALICE_VER_KEY="$(cheqd_noded_docker debug ed25519 random)"
    ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
    ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
    ALICE_VER_PUB_MULTIBASE_58=$(cheqd_noded_docker debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

    # Build CreateDid message
    KEY_ID="${did_to_write}#key1"

    # shellcheck disable=SC2089
    MSG_CREATE_DID='{"id":"'${did_to_write}'","verification_method":[{"id":"'"${KEY_ID}"'","type":"Ed25519VerificationKey2020","controller":"'${did_to_write}'","public_key_multibase":"'${ALICE_VER_PUB_MULTIBASE_58}'"}],"authentication":["'${KEY_ID}'"]}'

    # Post the message
    did=$(local_client_tx tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" \
        --ver-key "${ALICE_VER_PRIV_BASE_64}" \
        --from operator0 \
        --gas-prices "25ncheq" \
        --chain-id $CHAIN_ID \
        --output json \
        -y)


    txhash=$(echo "$did" | jq ".txhash" | tr -d '"')
    echo "$txhash" >> $FNAME_TXHASHES
}

# Check transaction hashes
function check_tx_hashes () {
    while IFS= read -r txhash
    do
        txhash=$(echo "${txhash}" | tr -d '"')
        result=$(cheqd_noded_docker query tx "${txhash}" --output json)
        tx_exist=$(echo "$result" | jq ".code")
        if [ "$tx_exist" != "0" ] ; then
            echo "Error was in checking tx with hash: $txhash"
            exit 1
        fi
    done < $FNAME_TXHASHES
}

function get_balance () {
    address=$1
    cheqd_noded_docker query bank balances "$address" | grep amount | sed 's/[^0-9]//g'
}

function get_did () {
    requested_did=$1
    cheqd_noded_docker query cheqd did "$requested_did" --output json
}

function get_resource () {
    collection_id=$1
    resource_id=$2
    cheqd_noded_docker query resource resource "$collection_id" "$resource_id" --output json
}

# Check that balance of operator3 increased to CHEQ_AMOUNT
# Input: Address to check
function check_balance () {
    address_to_check=$1
    new_balance=$(get_balance "$address_to_check")
    if [ $((new_balance - AMOUNT_BEFORE)) != $CHEQ_AMOUNT_NUMBER ];
    then
        echo "Balance after token send is not expected"
        exit 1
    fi
}

# Check that $DID exists
function check_did () {
    did_to_check=$1
    did_from=$(get_did "$did_to_check" | jq ".did.id" | tr -d '"')
    if [ "$did_from" != "$did_to_check" ];
    then
        echo "There is no any $did_to_check on server"
        exit 1
    fi
}

# Check that $DID exists
function check_resource () {
    collection_id_to_check=$1
    resource_to_check=$2
    resource_from=$(get_resource "$collection_id_to_check" "resource_to_check" | jq ".resource.id" | tr -d '"')
    if [ "$resource_from" != "$resource_to_check" ];
    then
        echo "There is no any $did_to_check on server"
        exit 1
    fi
}
