#!/bin/bash

# Inits node configuration and runs the node.
# e -> exit immediately, u -> treat unset variables as errors and immediately, o -> sets the exit code to the rightmost command 
set -euo pipefail

## Parameters

# Within the container, $HOME=/home/cheqd
CHEQD_ROOT_DIR="$HOME/.cheqdnode"
ENABLE_IMPORT="${ENABLE_IMPORT:-true}"

# Initialise node config directory
if [ ! -d "${CHEQD_ROOT_DIR}/config" ]
then
    echo "Node configuration directory not found. Initializing."
    cheqd-noded init "${CHEQD_NODED_MONIKER}"
else
    echo "Node config exists. Skipping initialization."
fi

# Check if a genesis file has been passed in config
if [ -f "/genesis" ]
then
    echo "Genesis file passed. Replacing current genesis file."
    cp /genesis "${CHEQD_ROOT_DIR}/config/genesis.json"
else
    echo "No genesis file passed. Skipping and retaining existing genesis."
fi

# Check if a seeds file has been passed in config
if [ -f "/seeds" ]
then
    echo "Seeds file passed. Replacing current seeds."
    cp /seeds "${CHEQD_ROOT_DIR}/config/seeds.txt"
    CHEQD_NODED_P2P_SEEDS="$(cat "${CHEQD_ROOT_DIR}/config/seeds.txt")"
    export CHEQD_NODED_P2P_SEEDS
else
    echo "No seeds file passed. Skipping and retaining existing seeds."
fi

# Check if a node_key file has been passed in config
if [ -f "/node-key" ]
then
    echo "node_key.json file passed. Replacing existing node_key.json."
    cp /node-key "${CHEQD_ROOT_DIR}/config/node_key.json"
else
    echo "No node key file passed. Skipping and retaining existing node key."
fi

# Check if a priv_validator_key file has been passed in config
if [ -f "/private-validator-key" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "priv_validator_key.json file passed. Replacing current validator key."
    cp /private-validator-key "${CHEQD_ROOT_DIR}/config/priv_validator_key.json"
else
    echo "No private validator key file passed. Skipping and retaining existing key."
fi

# Check if a priv_validator_state file has been passed in config
if [ -f "/private-validator-state" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "priv_validator_state.json file passed. Replacing current validator state."
    cp /private-validator-state "${CHEQD_ROOT_DIR}/data/priv_validator_state.json"
else
    echo "No private validator state file passed. Skipping and retaining existing validator state."
fi

# Check if a validator account has been passed in config
if [ -f "/validator-account" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "Validator account key file passed. Replacing current validator account key file."
    # TODO
else
    echo "No validator account key file passed. Skipping and retaining existing validator account."
fi

# Check if a upgrade_info file has been passed in config
if [ -f "/upgrade-info" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "upgrade_info.json file passed. Replacing current upgrade_info.json file."
    cp /upgrade-info "${CHEQD_ROOT_DIR}/data/upgrade-info.json"
else
    echo "No upgrade_info.json file passed. Skipping and retaining existing upgrade_info.json file."
fi

# Check if an account key mnemonic file has been passed in config
# Proceed only if input file exists and import is enabled
if [ -f "/import-accounts" ] && [ "$ENABLE_IMPORT" = true ]
then
    echo "Account import file passed. Importing accounts."

    # Count number of accounts in input file minus header
    EXPECTED_ACCOUNTS=$(tail -n +2 "/import-accounts" | wc -l)

    # Read accounts from CSV file
    while IFS="," read -r ACCOUNT MNEMONIC
    do
        if cheqd-noded keys show "$ACCOUNT" --keyring-backend "$CHEQD_NODED_KEYRING_BACKEND"
        then
            echo "Key ${ACCOUNT} already exists"
        else
            echo "Importing account: ${ACCOUNT}"
            cheqd-noded keys add "$ACCOUNT" --recover --keyring-backend "$CHEQD_NODED_KEYRING_BACKEND" <<< "$MNEMONIC"
        fi
    done < <(tail -n +2 "/import-accounts")

    # Count number of imported accounts
    IMPORTED_ACCOUNTS=$(cheqd-noded keys list --keyring-backend "$CHEQD_NODED_KEYRING_BACKEND" | grep -c "cheqd1")  

    if [ "$IMPORTED_ACCOUNTS" -eq "$EXPECTED_ACCOUNTS" ]
    then
        echo "All accounts imported successfully"
        echo "Imported accounts: ${IMPORTED_ACCOUNTS}"
    else
        echo "Mismatch in number of imported accounts"
        echo "Imported accounts: ${IMPORTED_ACCOUNTS}"
        echo "Expected accounts: ${EXPECTED_ACCOUNTS}"
    fi

else
    echo "Account import file not passed. Skipping account import."
fi

# Run node
cheqd-noded start
