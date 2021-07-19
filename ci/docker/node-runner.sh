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
echo "$GENESIS" | base64 --decode > $NODE_HOME/config/genesis.json
echo "$NODE_KEY" | base64 --decode > $NODE_HOME/config/node_key.json
echo "$PRIV_VALIDATOR_KEY" | base64 --decode > $NODE_HOME/config/priv_validator_key.json

# Run node
NODE_ARGS=${NODE_ARGS:-}  # Allo node args to be empty
cheqd-noded start $NODE_ARGS
