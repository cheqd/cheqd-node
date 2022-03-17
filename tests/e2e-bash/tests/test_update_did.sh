#!/bin/bash

set -euox pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/common.sh"


# Creating DID
ALICE_VER_KEY="$(cheqd-noded debug ed25519 random)"
ALICE_VER_PUB_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".pub_key_base_64")
ALICE_VER_PRIV_BASE_64=$(echo "${ALICE_VER_KEY}" | jq -r ".priv_key_base_64")
ALICE_VER_PUB_MULTIBASE_58=$(cheqd-noded debug encoding base64-multibase58 "${ALICE_VER_PUB_BASE_64}")

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

# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"

# Query DID to find out version id
# TODO: VersionId must be returned in MsgCreateDidResp
# shellcheck disable=SC2086
RESULT=$(cheqd-noded query cheqd did "${DID}" ${QUERY_PARAMS})
VERSION_ID=$(echo "${RESULT}" | jq -r ".metadata.version_id")


# Updating DID
MSG_UPDATE_DID='{
  "id": "'${DID}'",
  "version_id": "'$VERSION_ID'",
  "verification_method": [{
    "id": "'${KEY_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID}'",
    "public_key_multibase": "'${ALICE_VER_PUB_MULTIBASE_58}'"
  }],
  "authentication": [
    "'${KEY_ID}'"
  ],
  "capability_delegation": [
    "'${KEY_ID}'"
  ]
}';

# Post the message
# shellcheck disable=SC2086
RESULT=$(cheqd-noded tx cheqd update-did "${MSG_UPDATE_DID}" "${KEY_ID}" "${ALICE_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


# Query DID and assert resp
# shellcheck disable=SC2086
RESULT=$(cheqd-noded query cheqd did "${DID}" ${QUERY_PARAMS})

EXPECTED='{
   "context":[],
   "id":"'${DID}'",
   "controller":[],
   "verification_method":[
      {
         "id":"'${KEY_ID}'",
         "type":"Ed25519VerificationKey2020",
         "controller":"'${DID}'",
         "public_key_jwk":[],
         "public_key_multibase":"'${ALICE_VER_PUB_MULTIBASE_58}'"
      }
   ],
   "authentication":[
      "'${KEY_ID}'"
   ],
   "assertion_method":[],
   "capability_invocation":[],
   "capability_delegation":[
      "'${KEY_ID}'"
   ],
   "key_agreement":[],
   "service":[],
   "also_known_as":[]
}'

assert_json_eq "${EXPECTED}" "$(echo "$RESULT" | jq -r ".did")"
