#!/bin/bash

set -euox pipefail

. common.sh

# Wait for upgrade height
bash ../networks/wait_for_chain.sh $UPGRADE_HEIGHT $(echo "2*$VOTING_PERIOD" | bc)

# Stop docker-compose service
docker_compose_down

# Make all the data accessable
make_777

# Start docker-compose with new base image on new version
docker_compose_up "$CHEQD_IMAGE_TO" $(pwd)

# Check that upgrade was successful

# Wait for upgrade height
bash ../networks/wait_for_chain.sh $(echo $UPGRADE_HEIGHT+2 | bc)

CURRENT_VERSION=$(curl -s http://localhost:26657/abci_info | jq '.result.response.version' | grep -Eo '[0-9]*\.[0-9]*\.[0-9]*')

if [ $CURRENT_VERSION != $CHEQD_VERSION_TO ] ; then
     echo "Upgrade to version $CHEQD_VERSION_TO was not successful"
     exit 1
fi

# Stop docker compose
docker_compose_down

# Clean environment after test
clean_env