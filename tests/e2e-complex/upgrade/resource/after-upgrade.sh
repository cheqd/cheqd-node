#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. "common.sh"


# Use new version
cheqd_noded_docker() {
    # echo "new docker"
    docker run --rm \
        -v "$(pwd):/home/cheqd" \
        --network host \
        -u root \
        -e HOME=/home/cheqd \
        --entrypoint "cheqd-noded" \
        "${CHEQD_IMAGE_TO}" "$@"
}

get_addresses
# "To" address was used for sending tokens before upgrade
# shellcheck disable=SC2154
OP_ADDRESS_BEFORE=${addresses[2]}
# "To" address was used for sending tokens after upgrade
OP_ADDRESS_AFTER=${addresses[3]}

# Check balances after tokens sending
check_balance "$OP_ADDRESS_BEFORE"

# Check that did written before upgrade stil exist
check_did "$DID_1"

# Send tokens for checking functionality after upgrade
send_tokens "$OP_ADDRESS_AFTER"

# Send DID after upgrade
send_did_new "$DID_2"

# Check balance after token sending
check_balance "$OP_ADDRESS_AFTER"

# Check that did written before upgrade stil exist
check_did "$DID_2"

# Check that token transaction exists after upgrade too
check_tx_hashes

# Send new resource
send_resource_new "$DID_2_IDENTIFIER" "$RESOURCE_1"

# Check new resource
check_resource "$DID_2_IDENTIFIER" "$RESOURCE_1"
