#!/bin/bash

set -euox pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/common.sh"


# Generate Alice identity key
ALICE_VER_KEY="$(cheqd-noded debug ed25519 random)"
ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
ALICE_VER_PUB_MULTIBASE_58=$(cheqd-noded debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

# Build CreateDid message
DID="did:cheqd:testnet:$(random_string)"
KEY_ID="${DID}#key1"

MSG_CREATE_DID='{
  "id": "'${DID}'",
  "verification_method": [{
    "id": "'${KEY_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID}'",
    "public_key_multibase": "'${ALICE_VER_PUB_MULTIBASE_58}'"
  }],
  "authentication": [
    "'${KEY_ID}'"
  ]
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" --ver-key "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"
