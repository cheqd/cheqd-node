#!/bin/bash

# Inits node configuration and runs the node.
# e -> exit immediately, u -> treat unset variables as errors and immediately, o -> sets the exit code to the rightmost command 
set -euo pipefail

## Parameters

# Within the container, $HOME=/home/cheqd
CHEQD_NODED_HOME="$HOME/.cheqdnode"
ENABLE_IMPORT="${ENABLE_IMPORT:-true}"

# Copy network configuration folder for specific moniker if it exists
if [ -d "${HOME}/network-config/${CHEQD_NODED_MONIKER}" ]
then
    echo "Network configuration directory found for moniker ${CHEQD_NODED_MONIKER}. Copying to node home."

    # Create .cheqdnode directory if it doesn't exist
    if [ ! -d "${CHEQD_NODED_HOME}" ]
    then
        echo "Node home directory not found. Creating."
        mkdir -p "${CHEQD_NODED_HOME}"
    else
        echo "Node home directory exists. Skipping creation."
    fi

    # Copy network configuration to .cheqdnode except for data directory
    cd "${HOME}/network-config/${CHEQD_NODED_MONIKER}"
    tar -cf -  --exclude "./data" . | tar -xC "${CHEQD_NODED_HOME}"
    cd "${HOME}"
else
    echo "Network configuration directory not found for moniker ${CHEQD_NODED_MONIKER}. Skipping copy."
fi

# Check if a priv_validator_state file has been passed in config
if [ -f "/private-validator-state" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "priv_validator_state.json file passed. Replacing current validator state."
    cp /private-validator-state "${CHEQD_NODED_HOME}/data/priv_validator_state.json"
else
    echo "No private validator state file passed. Skipping and retaining existing validator state."
fi

# Check if a upgrade_info file has been passed in config
if [ -f "/upgrade-info" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "upgrade_info.json file passed. Replacing current upgrade_info.json file."
    cp /upgrade-info "${CHEQD_NODED_HOME}/data/upgrade-info.json"
else
    echo "No upgrade_info.json file passed. Skipping and retaining existing upgrade_info.json file."
fi

# Check if an account key mnemonic file has been passed in config
# Proceed only if input file exists and import is enabled
if [ -f "/import-accounts" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "Account import file passed. Importing accounts."

    # Call account import script
    import-keys /import-accounts "$CHEQD_NODED_KEYRING_BACKEND"
else
    echo "Account import file not passed. Skipping account import."
fi

# Run node
cheqd-noded start
