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


########## Creating DID ##########

ID="$(random_string)"
DID="did:cheqd:testnet:$ID"
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
RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


########## Creating Resource ##########

RESOURCE_ID=$(uuidgen)
RESOURCE_NAME="Resource 1"
RESOURCE_MIME_TYPE="application/json"
RESOURCE_RESOURCE_TYPE="CL-Schema"
RESOURCE_DATA='test data';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx resource create-resource ${ID} ${RESOURCE_ID} "${RESOURCE_NAME}" ${RESOURCE_RESOURCE_TYPE} ${RESOURCE_MIME_TYPE} <(echo "${RESOURCE_DATA}") "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"

########## Querying Resource ##########

# shellcheck disable=SC2086
RESULT=$(cheqd-noded query resource resource "${ID}" ${RESOURCE_ID}  ${QUERY_PARAMS})

assert_eq "$(echo "$RESULT" | jq -r ".resource.header.collection_id")" "${ID}"
assert_eq "$(echo "$RESULT" | jq -r ".resource.header.id")" "${RESOURCE_ID}"
assert_eq "$(echo "$RESULT" | jq -r ".resource.header.name")" "${RESOURCE_NAME}"
assert_eq "$(echo "$RESULT" | jq -r ".resource.header.resource_type")" "${RESOURCE_RESOURCE_TYPE}"
assert_eq "$(echo "$RESULT" | jq -r ".resource.header.mime_type")" "${RESOURCE_MIME_TYPE}"
assert_eq "$(echo "$RESULT" | jq -r ".resource.data")" "$(echo "${RESOURCE_DATA}" | base64 -w 0)"
