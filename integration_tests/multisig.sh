#!/bin/bash

set -euox pipefail

JACK_HOME="localnet/node0"
ALICE_HOME="localnet/node1"
BOB_HOME="localnet/node2"
ANNA_HOME="localnet/node3"

JACK_PUBKEY=$(verim-cosmosd keys show jack -p --home $JACK_HOME)
ALICE_PUBKEY=$(verim-cosmosd keys show alice -p --home $ALICE_HOME)
BOB_PUBKEY=$(verim-cosmosd keys show bob -p --home $BOB_HOME)
ANNA_PUBKEY=$(verim-cosmosd keys show anna -p --home $ANNA_HOME)

JACK_ADDRESS=$(verim-cosmosd keys show jack -a --home $JACK_HOME)

echo "jack: $JACK_PUBKEY"
echo "alice: $ALICE_PUBKEY"
echo "bob: $BOB_PUBKEY"
echo "anna: $ANNA_PUBKEY"


############ Jack creates multisig key: ffour (fantastic four)

verim-cosmosd keys add alice_pub --pubkey $ALICE_PUBKEY --home $JACK_HOME
verim-cosmosd keys add bob_pub --pubkey $BOB_PUBKEY --home $JACK_HOME
verim-cosmosd keys add anna_pub --pubkey $ANNA_PUBKEY --home $JACK_HOME

verim-cosmosd keys add ffour --multisig=jack,alice_pub,bob_pub,anna_pub --multisig-threshold=3 --home $JACK_HOME

verim-cosmosd keys show ffour --home $JACK_HOME
FFOUR_ADDRESS=$(verim-cosmosd keys show ffour -a --home $JACK_HOME)

############ Someone (Jack) transfers money to the ffour account

verim-cosmosd tx bank send jack $FFOUR_ADDRESS 1token --home $JACK_HOME


############ Alice creates multisig key: ffour (fantastic four)

verim-cosmosd keys add jack_pub --pubkey $JACK_PUBKEY --home $ALICE_HOME
verim-cosmosd keys add bob_pub --pubkey $BOB_PUBKEY --home $ALICE_HOME
verim-cosmosd keys add anna_pub --pubkey $ANNA_PUBKEY --home $ALICE_HOME

# is 3 imoprtant?
verim-cosmosd keys add ffour --multisig=jack_pub,alice,bob_pub,anna_pub --multisig-threshold=3 --home $ALICE_HOME

verim-cosmosd keys show ffour --home $ALICE_HOME
FFOUR_ADDRESS_2=$(verim-cosmosd keys show ffour -a --home $ALICE_HOME)


############ Bob creates multisig key: ffour (fantastic four)

verim-cosmosd keys add jack_pub --pubkey $JACK_PUBKEY --home $BOB_HOME
verim-cosmosd keys add alice_pub --pubkey $ALICE_PUBKEY --home $BOB_HOME
verim-cosmosd keys add anna_pub --pubkey $ANNA_PUBKEY --home $BOB_HOME

# is 3 imoprtant?
verim-cosmosd keys add ffour --multisig=jack_pub,alice_pub,bob,anna_pub --multisig-threshold=3 --home $BOB_HOME

verim-cosmosd keys show ffour --home $BOB_HOME
FFOUR_ADDRESS_3=$(verim-cosmosd keys show ffour -a --home $BOB_HOME)


############ Ffour

verim-cosmosd tx bank send $FFOUR_ADDRESS $JACK_ADDRESS 1token --generate-only > localnet/unsignedTx.json --home $JACK_HOME

verim-cosmosd tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=jack --output-document=localnet/jack_signature.json --home $JACK_HOME --chain-id=verimcosmos

verim-cosmosd tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=alice --output-document=localnet/alice_signature.json --home $ALICE_HOME --chain-id=verimcosmos

verim-cosmosd tx sign localnet/unsignedTx.json --multisig=$FFOUR_ADDRESS --from=bob --output-document=localnet/bob_signature.json --home $BOB_HOME --chain-id=verimcosmos # --account-number=11

verim-cosmosd tx multisign \
  localnet/unsignedTx.json \
  ffour \
  localnet/alice_signature.json localnet/jack_signature.json localnet/bob_signature.json > localnet/signedTx.json --home $JACK_HOME

verim-cosmosd tx broadcast localnet/signedTx.json --home $JACK_HOME
