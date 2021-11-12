#!/bin/bash

# Inits node configuration and runs the node.

set -euox pipefail

NODE_HOME="$HOME/.cheqdnode"


# Init node config directory
if [ ! -d "${NODE_HOME}/config" ]
then
    echo "Node home not found. Initializing."
    cheqd-noded init $NODE_MONIKER
else
    echo "Node home exists. Skipping initialization."
fi


# Update configs
set - # Disable ditailed command logging

echo "Updating genesis"
echo "$GENESIS" | base64 --decode > $NODE_HOME/config/genesis.json

echo "Updating node key"
echo "$NODE_KEY" | base64 --decode > $NODE_HOME/config/node_key.json

echo "Updating validator key"
echo "$PRIV_VALIDATOR_KEY" | base64 --decode > $NODE_HOME/config/priv_validator_key.json

set -x # Re-enable ditailed command logging

# Run node
NODE_ARGS=${NODE_ARGS:-}  # Allo node args to be empty
cheqd-noded start $NODE_ARGS
