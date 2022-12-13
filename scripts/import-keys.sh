#!/bin/bash

# Set default shell behaviour
set -euox pipefail

# Parameters
INPUT_FILE="${1:-accounts.csv}"
CHEQD_NODED_KEYRING_BACKEND="${2:-test}"

# Proceed only if input file exists
if [ -f "${INPUT_FILE}" ]
then
  # Count number of accounts in input file minus header
  EXPECTED_ACCOUNTS=$(tail -n +2 "${INPUT_FILE}" | wc -l)

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
  done < <(tail -n +2 "${INPUT_FILE}")

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
    exit 1
  fi

else
  echo "Input file ${INPUT_FILE} does not exist"
  exit 1
fi
