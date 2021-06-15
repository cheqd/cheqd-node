#!/bin/bash

# Documentation: https://github.com/cosmos/gaia/blob/0ecb6ed8a244d835807f1ced49217d54a9ca2070/docs/resources/gaiad.md#multisig-transactions

set -euox pipefail

NODE_CONFIGS="../local_net/node_configs"
CHAIN_ID="verim"

JACK_HOME="$NODE_CONFIGS/node0"
ALICE_HOME="$NODE_CONFIGS/node1"
BOB_HOME="$NODE_CONFIGS/node2"
ANNA_HOME="$NODE_CONFIGS/node3"

JACK_PUBKEY=$(verim-noded keys show jack -p --home $JACK_HOME)
ALICE_PUBKEY=$(verim-noded keys show alice -p --home $ALICE_HOME)
BOB_PUBKEY=$(verim-noded keys show bob -p --home $BOB_HOME)
ANNA_PUBKEY=$(verim-noded keys show anna -p --home $ANNA_HOME)

JACK_ADDRESS=$(verim-noded keys show jack -a --home $JACK_HOME)

echo "jack: $JACK_PUBKEY"
echo "alice: $ALICE_PUBKEY"
echo "bob: $BOB_PUBKEY"
echo "anna: $ANNA_PUBKEY"


echo "############ Jack imports other's public keys and creates multisig key: ffour"

verim-noded keys add alice_pub --pubkey $ALICE_PUBKEY --home $JACK_HOME
verim-noded keys add bob_pub --pubkey $BOB_PUBKEY --home $JACK_HOME
verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $JACK_HOME

verim-noded keys add ffour --multisig=jack,alice_pub,bob_pub,anna_pub --multisig-threshold=3 --home $JACK_HOME
verim-noded keys show ffour --home $JACK_HOME

echo "############ Someone (Jack) transfers money to the ffour account"

FFOUR_ADDRESS=$(verim-noded keys show ffour -a --home $JACK_HOME)
FFOUR_PUBKEY=$(verim-noded keys show ffour -p --home $JACK_HOME)

verim-noded tx bank send jack $FFOUR_ADDRESS 1000000token \
  --fees 200000token \
  --chain-id=$CHAIN_ID \
  --home $JACK_HOME \
  --yes

echo "############ Alice imports other's public keys and creates multisig key: ffour"

verim-noded keys add jack_pub --pubkey $JACK_PUBKEY --home $ALICE_HOME
verim-noded keys add bob_pub --pubkey $BOB_PUBKEY --home $ALICE_HOME
verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $ALICE_HOME

# Key order and trashhold must be the same. Multisig key import is not supported.
verim-noded keys add ffour --multisig=jack_pub,alice,bob_pub,anna_pub --multisig-threshold=3 --home $ALICE_HOME
verim-noded keys show ffour --home $ALICE_HOME

echo "############ Bob imports other's public keys and creates multisig key: ffour"

verim-noded keys add jack_pub --pubkey $JACK_PUBKEY --home $BOB_HOME
verim-noded keys add alice_pub --pubkey $ALICE_PUBKEY --home $BOB_HOME
verim-noded keys add anna_pub --pubkey $ANNA_PUBKEY --home $BOB_HOME

# Key order and trashhold must be the same. Multisig key import is not supported.
verim-noded keys add ffour --multisig=jack_pub,alice_pub,bob,anna_pub --multisig-threshold=3 --home $BOB_HOME
verim-noded keys show ffour --home $BOB_HOME

echo "############ Jack generates a transaction"

verim-noded tx bank send $FFOUR_ADDRESS $JACK_ADDRESS 1000token \
  --fees 200000token \
  --generate-only \
  --home $JACK_HOME \
  > $NODE_CONFIGS/unsignedTx.json

echo "############ Jack signs the transaction"

verim-noded tx sign $NODE_CONFIGS/unsignedTx.json \
  --multisig=$FFOUR_ADDRESS \
  --from=jack \
  --output-document=$NODE_CONFIGS/jack_signature.json \
  --home $JACK_HOME \
  --chain-id=$CHAIN_ID

echo "############ Alice signs the transaction"

verim-noded tx sign $NODE_CONFIGS/unsignedTx.json \
  --multisig=$FFOUR_ADDRESS \
  --from=alice \
  --output-document=$NODE_CONFIGS/alice_signature.json \
  --home $ALICE_HOME \
  --chain-id=$CHAIN_ID

echo "############ Bob signs the transaction. 3 out of 4 signatures is enough."

verim-noded tx sign $NODE_CONFIGS/unsignedTx.json \
  --multisig=$FFOUR_ADDRESS \
  --from=bob \
  --output-document=$NODE_CONFIGS/bob_signature.json \
  --home $BOB_HOME \
  --chain-id=$CHAIN_ID

echo "############ Jack (can be anyone) composes signatures"

# chain-id is important
verim-noded tx multisign \
  $NODE_CONFIGS/unsignedTx.json \
  ffour \
  $NODE_CONFIGS/jack_signature.json $NODE_CONFIGS/alice_signature.json $NODE_CONFIGS/bob_signature.json \
  --chain-id=$CHAIN_ID \
  --home $JACK_HOME \
  > $NODE_CONFIGS/signedTx.json

echo "############ Jack (can be anyone) broadcasts the transaction"

verim-noded tx broadcast $NODE_CONFIGS/signedTx.json --home $JACK_HOME --yes
