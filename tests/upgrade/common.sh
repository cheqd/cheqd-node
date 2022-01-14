set -euox pipefail

DOCKER_COMPOSE_DIR="../networks/docker_compose"
CHEQD_VERSION_FROM="v0.3.1"
CHEQD_VERSION_TO="0.4.0"
UPGRADE_NAME="v0.4"
VOTING_PERIOD=30
EXPECTED_BLOCK_SECOND=5
EXTRA_BLOCKS=5
UPGRADE_HEIGHT=$(echo $VOTING_PERIOD/$EXPECTED_BLOCK_SECOND+$EXTRA_BLOCKS | bc)
DEPOSIT_AMOUNT=10000000
CHAIN_ID="cheqd"
CHEQD_USER="cheqd"

# cheqd_noded docker wrapper

cheqd_noded_docker() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    --network host \
    ghcr.io/cheqd/cheqd-node:${CHEQD_VERSION_FROM} "$@"
}

# cheqd-noded docker wrapper for init purposes. It's needed because of "permission denied" issue on github runners

cheqd_noded_docker_init() {
  docker run --rm \
    -v "$(pwd)":"/cheqd" \
    --network host \
    -u root \
    ghcr.io/cheqd/cheqd-node:${CHEQD_VERSION_FROM} init "$@" --home /cheqd/.cheqdnode
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
    export CHEQD_VERSION="$1"
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
}

# Run command using local generated keys from node_configs/client
function local_client_exec () {
    cheqd_noded_docker "$@" --home node_configs/client/.cheqdnode/ --keyring-backend test
}