#!/bin/bash

# TODO: Utilize home env

set -euo pipefail

rm -rf ~/.verimcosmos

rm -rf localnet
mkdir localnet localnet/client localnet/node0 localnet/node1 localnet/node2 localnet/node3

# client

# verim-cosmosd config chain-id verim-cosmoschain
# verim-cosmosd config output json
# verim-cosmosd config indent true
# verim-cosmosd config trust-node false

echo 'test1234' | verim-cosmosd keys add jack
echo 'test1234' | verim-cosmosd keys add alice
echo 'test1234' | verim-cosmosd keys add bob
echo 'test1234' | verim-cosmosd keys add anna

cp -r ~/.verimcosmos/* localnet/client

# node 0

verim-cosmosd init node0 --chain-id verim-cosmoschain

jack_address=$(verim-cosmosd keys show jack -a)
jack_pubkey=$(verim-cosmosd keys show jack -p)

alice_address=$(verim-cosmosd keys show alice -a)
alice_pubkey=$(verim-cosmosd keys show alice -p)

bob_address=$(verim-cosmosd keys show bob -a)
bob_pubkey=$(verim-cosmosd keys show bob -p)

anna_address=$(verim-cosmosd keys show anna -a)
anna_pubkey=$(verim-cosmosd keys show anna -p)

verim-cosmosd add-genesis-account --address=$jack_address --pubkey=$jack_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$alice_address --pubkey=$alice_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$bob_address --pubkey=$bob_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$anna_address --pubkey=$anna_pubkey --roles="NodeAdmin"

echo 'test1234' | verim-cosmosd gentx --from jack

mv ~/.verimcosmos/* localnet/node0

# node 1

verim-cosmosd init node1 --chain-id verim-cosmoschain

verim-cosmosd add-genesis-account --address=$jack_address --pubkey=$jack_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$alice_address --pubkey=$alice_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$bob_address --pubkey=$bob_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$anna_address --pubkey=$anna_pubkey --roles="NodeAdmin"

echo 'test1234' | verim-cosmosd gentx --from alice

mv ~/.verimcosmos/* localnet/node1

# node 2

verim-cosmosd init node2 --chain-id verim-cosmoschain

verim-cosmosd add-genesis-account --address=$jack_address --pubkey=$jack_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$alice_address --pubkey=$alice_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$bob_address --pubkey=$bob_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$anna_address --pubkey=$anna_pubkey --roles="NodeAdmin"

echo 'test1234' | verim-cosmosd gentx --from bob

mv ~/.verimcosmos/* localnet/node2

# node 3

verim-cosmosd init node3 --chain-id verim-cosmoschain

verim-cosmosd add-genesis-account --address=$jack_address --pubkey=$jack_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$alice_address --pubkey=$alice_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$bob_address --pubkey=$bob_pubkey --roles="Trustee,NodeAdmin"
verim-cosmosd add-genesis-account --address=$anna_address --pubkey=$anna_pubkey --roles="NodeAdmin"

echo 'test1234' | verim-cosmosd gentx --from anna

cp -r ~/.verimcosmos/* localnet/node3

# Collect all validator creation transactions

cp localnet/node0/config/gentx/* ~/.verim-cosmosd/config/gentx
cp localnet/node1/config/gentx/* ~/.verim-cosmosd/config/gentx
cp localnet/node2/config/gentx/* ~/.verim-cosmosd/config/gentx
cp localnet/node3/config/gentx/* ~/.verim-cosmosd/config/gentx

# Embed them into genesis

verim-cosmosd collect-gentxs
verim-cosmosd validate-genesis

# Update genesis for all nodes

cp ~/.verimcosmos/config/genesis.json localnet/node0/config/
cp ~/.verimcosmos/config/genesis.json localnet/node1/config/
cp ~/.verimcosmos/config/genesis.json localnet/node2/config/
cp ~/.verimcosmos/config/genesis.json localnet/node3/config/

# Find out node ids

id0=$(ls localnet/node0/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id1=$(ls localnet/node1/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id2=$(ls localnet/node2/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id3=$(ls localnet/node3/config/gentx | sed 's/gentx-\(.*\).json/\1/')

# Update address book of the first node
peers="$id0@192.167.10.2:26656,$id1@192.167.10.3:26656,$id2@192.167.10.4:26656,$id3@192.167.10.5:26656"
sed -i "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" localnet/node0/config/config.toml

# Make RPC enpoint available externally

sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node0/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node1/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node2/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node3/config/config.toml
