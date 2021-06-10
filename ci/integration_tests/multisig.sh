#!/bin/bash

set -euox pipefail

JACK_HOME="localnet/node0"
ALICE_HOME="localnet/node1"
BOB_HOME="localnet/node2"
ANNA_HOME="localnet/node3"

JACK_PUBKEY=$(verim-noded keys show jack -p --home $JACK_HOME)
ALICE_PUBKEY=$(verim-noded keys show alice -p --home $ALICE_HOME)
BOB_PUBKEY=$(verim-noded keys show bob -p --home $BOB_HOME)
ANNA_PUBKEY=$(verim-noded keys show anna -p --home $ANNA_HOME)

JACK_ADDRESS=$(verim-noded keys show jack -a --home $JACK_HOME)

echo "jack: $JACK_PUBKEY"
echo "alice: $ALICE_PUBKEY"
echo "bob: $BOB_PUBKEY"
echo "anna: $ANNA_PUBKEY"


############ Jack imports other's public keys and creates multisig key: ffour

verim-noded keys add alice_pub --pubkey $ALICE_PUBKEY --home $JACK_HOME
verim-noded keys add bob_pub --pubkey $BOB_PUBKEY --home $JACK_HOME
verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $JACK_HOME

verim-noded keys add ffour --multisig=jack,alice_pub,bob_pub,anna_pub --multisig-threshold=3 --home $JACK_HOME

verim-noded keys show ffour --home $JACK_HOME
FFOUR_ADDRESS=$(verim-noded keys show ffour -a --home $JACK_HOME)

############ Someone (Jack) transfers money to the ffour account

verim-noded tx bank send jack $FFOUR_ADDRESS 1token --home $JACK_HOME


############ Alice imports other's public keys and creates multisig key: ffour (she doesn't trust shared ffour key)

verim-noded keys add jack_pub --pubkey $JACK_PUBKEY --home $ALICE_HOME
verim-noded keys add bob_pub --pubkey $BOB_PUBKEY --home $ALICE_HOME
verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $ALICE_HOME

# is 3 imoprtant?
verim-noded keys add ffour --multisig=jack_pub,alice,bob_pub,anna_pub --multisig-threshold=3 --home $ALICE_HOME

verim-noded keys show ffour --home $ALICE_HOME
FFOUR_ADDRESS_2=$(verim-noded keys show ffour -a --home $ALICE_HOME)
FFOUR_PUBKEY=$(verim-noded keys show ffour -p --home $ALICE_HOME)


############ Bob imports ffour key directly (he trusts shared multisig key)

# verim-noded keys add jack_pub --pubkey $JACK_PUBKEY --home $BOB_HOME
# verim-noded keys add alice_pub --pubkey $ALICE_PUBKEY --home $BOB_HOME
# verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $BOB_HOME

# # is 3 imoprtant?
# verim-noded keys add ffour --multisig=jack_pub,alice_pub,bob,anna_pub --multisig-threshold=3 --home $BOB_HOME

verim-noded keys add ffour_pub --pubkey $FFOUR_PUBKEY --home $BOB_HOME

verim-noded keys show ffour_pub --home $BOB_HOME
FFOUR_ADDRESS_3=$(verim-noded keys show ffour_pub -a --home $BOB_HOME)


############ Ffour

verim-noded tx bank send $FFOUR_ADDRESS $JACK_ADDRESS 1token --generate-only > localnet/unsignedTx.json --home $JACK_HOME

verim-noded tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=jack --output-document=localnet/jack_signature.json --home $JACK_HOME --chain-id=verim

verim-noded tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=alice --output-document=localnet/alice_signature.json --home $ALICE_HOME --chain-id=verim

verim-noded tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=bob --output-document=localnet/bob_signature.json --home $BOB_HOME --chain-id=verim # --account-number=11

verim-noded tx multisign \
  localnet/unsignedTx.json \
  ffour \
  localnet/alice_signature.json localnet/jack_signature.json localnet/bob_signature.json > localnet/signedTx.json --home $JACK_HOME

verim-noded tx broadcast localnet/signedTx.json --home $JACK_HOME
