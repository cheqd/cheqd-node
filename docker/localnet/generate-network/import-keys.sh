#!/bin/bash

# Set default shell behaviour
set -euox pipefail

# Parameters
INPUT_FILE="${1:-accounts.csv}"
CHEQD_NODED_KEYRING_BACKEND="${2:-test}"

# Read accounts from CSV file
while IFS= read -r ACCOUNT MNEMONIC
do
  if cheqd-noded keys show "$ACCOUNT" --keyring-backend "$CHEQD_NODED_KEYRING_BACKEND"
  then
    echo "Key ${ACCOUNT} already exists"
  else
    echo "Importing account: ${ACCOUNT}"
    cheqd-noded keys add "$ACCOUNT" --recover --keyring-backend "$CHEQD_NODED_KEYRING_BACKEND" <<< "$MNEMONIC"
  fi
done < <(tail -n +2 "${INPUT_FILE}")

exit 0
