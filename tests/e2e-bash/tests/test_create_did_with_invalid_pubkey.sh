#!/bin/bash

set -euox pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/common.sh"


# Generate Alice identity key
ALICE_VER_KEY="$(cheqd-noded debug ed25519 random)"
ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")

# Build CreateDid message
DID="did:cheqd:testnet:$(random_string)"
KEY_ID="${DID}#key1"

MSG_CREATE_DID='{
  "id": "'${DID}'",
  "verification_method": [{
    "id": "'${KEY_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID}'",
    "public_key_multibase": "z6qdb7tjBeqaaWHU9W8nLHoxtZzro4ka"
  }],
  "authentication": [
    "'${KEY_ID}'"
  ]
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS} 2>&1  || true)  # Allow command to fail, redirect stderr to stdout

assert_str_contains "$RESULT" "public_key_multibase: ed25519: bad public key length: 23"
