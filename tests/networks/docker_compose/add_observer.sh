#!/bin/bash

# Init node5 keys and configs
sudo su cheqd
cheqd-noded init node5
exit

sudo chmod -R 777 /etc/cheqd-node

NODE0_IP="127.0.0.1"
PEER0=$(cat ${NODE_CONFIGS_BASE}/node0/node_id.txt)@$NODE0_IP:26656

# Genesis
cp ${NODE_CONFIGS_BASE}/node0/.cheqdnode/config/genesis.json /etc/cheqd-node/

# Config
sed -ri "s|persistent_peers = \".*\"|persistent_peers = \"${PEER0}\"|" /etc/cheqd-node/config.toml
sed -ri "s|laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:26677\"|" /etc/cheqd-node/config.toml
sed -ri "s|laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://0.0.0.0:26676\"|" /etc/cheqd-node/config.toml

sudo systemctl start cheqd-noded
systemctl status cheqd-noded
