#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. "common.sh"


###
# Test data
###


# Get address of operator which will be used for sending tokens before upgrade
get_addresses
# shellcheck disable=SC2154
OP2_ADDRESS=${addresses[2]}

# Send tokens operator-0 -> operator-1
send_tokens "$OP2_ADDRESS"

# Send DID transaction
send_did_new "$DID_1"

# Check that token transaction exists
check_tx_hashes

# Check that $DID was written
check_did "$DID_1"

# Check balance after token sending
check_balance "$OP2_ADDRESS"
