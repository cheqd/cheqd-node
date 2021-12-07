#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

sed_extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "$SCRIPT_DIR/common.sh"


# Creating DID
OLD_VER_KEY="$(cheqd-noded debug ed25519 random)"
OLD_VER_PUB_BASE_64=$(echo "${OLD_VER_KEY}" | jq -r ".pub_key_base_64")
OLD_VER_PRIV_BASE_64=$(echo "${OLD_VER_KEY}" | jq -r ".priv_key_base_64")
OLD_VER_PUB_MULTIBASE_58=$(cheqd-noded debug encoding base64-multibase58 "${OLD_VER_PUB_BASE_64}")

DID="did:cheqd:testnet:$(random_string)"
KEY_ID="${DID}#key1"

MSG_CREATE_DID='{
  "id": "'${DID}'",
  "verification_method": [{
    "id": "'${KEY_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID}'",
    "public_key_multibase": "'${OLD_VER_PUB_MULTIBASE_58}'"
  }],
  "authentication": [
    "'${KEY_ID}'"
  ]
}';

RESULT=$(cheqd-noded tx cheqd create-did "${MSG_CREATE_DID}" "${KEY_ID}" --ver-key "${OLD_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"

# Query DID to find out version id
# TODO: VersionId must be returned in MsgCreateDidResp
RESULT=$(cheqd-noded query cheqd did "${DID}" ${QUERY_PARAMS})
VERSION_ID=$(echo "${RESULT}" | jq -r ".metadata.version_id")


# Updating DID
NEW_VER_KEY="$(cheqd-noded debug ed25519 random)"
NEW_VER_PUB_BASE_64=$(echo "${NEW_VER_KEY}" | jq -r ".pub_key_base_64")
NEW_VER_PRIV_BASE_64=$(echo "${NEW_VER_KEY}" | jq -r ".priv_key_base_64")
NEW_VER_PUB_MULTIBASE_58=$(cheqd-noded debug encoding base64-multibase58 "${NEW_VER_PUB_BASE_64}")

MSG_UPDATE_DID='{
  "id": "'${DID}'",
  "version_id": "'$VERSION_ID'",
  "verification_method": [{
    "id": "'${KEY_ID}'",
    "type": "Ed25519VerificationKey2020",
    "controller": "'${DID}'",
    "public_key_multibase": "'${NEW_VER_PUB_MULTIBASE_58}'"
  }],
  "authentication": [
    "'${KEY_ID}'"
  ],
  "capability_delegation": [
    "'${KEY_ID}'"
  ]
}';

# Post the message
RESULT=$(cheqd-noded tx cheqd update-did "${MSG_UPDATE_DID}" "${KEY_ID}" --ver-key "${OLD_VER_PRIV_BASE_64}" \
  --from "${BASE_ACCOUNT_1}" ${TX_PARAMS})

assert_tx_successful "$RESULT"


# Query DID and assert resp
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
         "public_key_multibase":"'${NEW_VER_PUB_MULTIBASE_58}'"
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
