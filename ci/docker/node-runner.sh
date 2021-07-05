#!/bin/bash

# Inits node configuration and runs the node.

set -euox pipefail

NODE_HOME="$HOME/.verimnode"


# Init config directory
verim-noded init $NODE_MONIKER

# Update configs
echo "$GENESIS" | base64 --decode > $NODE_HOME/config/genesis.json
echo "$NODE_KEY" | base64 --decode > $NODE_HOME/config/node_key.json
echo "$PRIV_VALIDATOR_KEY" | base64 --decode > $NODE_HOME/config/priv_validator_key.json

# Run node
NODE_ARGS=${NODE_ARGS:-}  # Allo node args to be empty
verim-noded start $NODE_ARGS
