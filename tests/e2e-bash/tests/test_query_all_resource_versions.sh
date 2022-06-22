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


########## Creating DID 1 ##########

ID1="$(random_string)"
DID1="did:cheqd:testnet:$ID1"
KEY1_ID="${DID1}#key1"

MSG_CREATE_DID_1='{
  "id": "'${DID1}'",
  "verification_method": [{
    "id": "'${KEY1_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID1}'",
    "public_key_multibase": "'${ALICE_VER_PUB_MULTIBASE_58}'"
  }],
  "authentication": [
    "'${KEY1_ID}'"
  ]
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID_1}" "${KEY1_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


########## Creating Resource 1 ##########

RESOURCE1_V1_ID=$(uuidgen)
RESOURCE1_NAME="Resource 1"
RESOURCE1_MIME_TYPE="application/json"
RESOURCE1_RESOURCE_TYPE="CL-Schema"
RESOURCE1_V1_DATA='dGVzdCBiYXNlNTYgZW5jb2RlZCBkYXRh';

MSG_CREATE_RESOURCE1='{
  "collection_id": "'${ID1}'",
  "id": "'${RESOURCE1_V1_ID}'",
  "name": "'${RESOURCE1_NAME}'",
  "mime_type": "'${RESOURCE1_MIME_TYPE}'",
  "resource_type": "'${RESOURCE1_RESOURCE_TYPE}'",
  "data": "'${RESOURCE1_V1_DATA}'"
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx resource create-resource-raw "${MSG_CREATE_RESOURCE1}" "${KEY1_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"

########## Querying Resource 1 ##########

# shellcheck disable=SC2086
RESULT=$(cheqd-noded query resource resource "${ID1}" ${RESOURCE1_V1_ID}  ${QUERY_PARAMS})

EXPECTED_RES1_V1='{
  "collection_id": "'${ID1}'",
  "id": "'${RESOURCE1_V1_ID}'",
  "name": "'${RESOURCE1_NAME}'",
  "mime_type": "'${RESOURCE1_MIME_TYPE}'",
  "resource_type": "'${RESOURCE1_RESOURCE_TYPE}'",
  "data": "'${RESOURCE1_V1_DATA}'"
}'

DEL_FILTER='del(.checksum, .created, .next_version_id, .previous_version_id)'
assert_json_eq "$(echo "$RESULT" | jq -r ".resource | ${DEL_FILTER}")" "${EXPECTED_RES1_V1}"


########## Creating Resource 1 v2 ##########

RESOURCE1_V2_ID=$(uuidgen)
RESOURCE1_V2_DATA='dGVzdCBiYXNlNTYgZW5jb2RlZCBkYXRhLg==';

MSG_CREATE_RESOURCE1_V2='{
  "collection_id": "'${ID1}'",
  "id": "'${RESOURCE1_V2_ID}'",
  "name": "'${RESOURCE1_NAME}'",
  "mime_type": "'${RESOURCE1_MIME_TYPE}'",
  "resource_type": "'${RESOURCE1_RESOURCE_TYPE}'",
  "data": "'${RESOURCE1_V2_DATA}'"
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx resource create-resource-raw "${MSG_CREATE_RESOURCE1_V2}" "${KEY1_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


########## Creating Resource 2 ##########

RESOURCE2_ID=$(uuidgen)
RESOURCE2_DATA='dGVzdCBiYXNlNTYgZW5jb2RlZCBkYXRhdGVzdCBiYXNlNTYgZW5jb2RlZCBkYXRh';
RESOURCE2_NAME="Resource 2"
RESOURCE2_MIME_TYPE="application/json"
RESOURCE2_RESOURCE_TYPE="CL-Schema"

MSG_CREATE_RESOURCE2='{
  "collection_id": "'${ID1}'",
  "id": "'${RESOURCE2_ID}'",
  "name": "'${RESOURCE2_NAME}'",
  "mime_type": "'${RESOURCE2_MIME_TYPE}'",
  "resource_type": "'${RESOURCE2_RESOURCE_TYPE}'",
  "data": "'${RESOURCE2_DATA}'"
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx resource create-resource-raw "${MSG_CREATE_RESOURCE2}" "${KEY1_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


########## Querying All Resource 1 versions ##########

EXPECTED_RES1_V2='{
  "collection_id": "'${ID1}'",
  "id": "'${RESOURCE1_V2_ID}'",
  "name": "'${RESOURCE1_NAME}'",
  "mime_type": "'${RESOURCE1_MIME_TYPE}'",
  "resource_type": "'${RESOURCE1_RESOURCE_TYPE}'",
  "data": "'${RESOURCE1_V2_DATA}'"
}'

# shellcheck disable=SC2086
RESULT=$(cheqd-noded query resource all-resource-versions "${ID1}" "${RESOURCE1_NAME}" ${RESOURCE1_RESOURCE_TYPE} ${RESOURCE1_MIME_TYPE} ${QUERY_PARAMS})

assert_eq "$(echo "$RESULT" | jq -r ".resources | length")" "2"
assert_json_eq "$(echo "$RESULT" | jq -r '.resources[] | select(.id == "'"${RESOURCE1_V1_ID}"'") | '"${DEL_FILTER}"'')" "${EXPECTED_RES1_V1}"
assert_json_eq "$(echo "$RESULT" | jq -r '.resources[] | select(.id == "'"${RESOURCE1_V2_ID}"'") | '"${DEL_FILTER}"'')" "${EXPECTED_RES1_V2}"
