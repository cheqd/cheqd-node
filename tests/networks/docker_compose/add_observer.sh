#!/bin/bash
sudo chmod -R 777 /etc/cheqd-node
sudo chmod -R 777 /var/lib/cheqd/data
cd
cheqd-noded init node5
NODE0_IP="127.0.0.1"
PEER0=$(cat ${NODE_CONFIGS_BASE}/node0/node_id.txt)@$NODE0_IP:26656
sed -ri "s|peers = \".*\"|peers = \"${PEER0}\"|" ~/.cheqdnode/config/config.toml
sed -ri "s|laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:26677\"|" ~/.cheqdnode/config/config.toml
sed -ri "s|laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://0.0.0.0:26676\"|" ~/.cheqdnode/config/config.toml
# cp ~/.cheqdnode/config/* /etc/cheqd-node/               # /var/lib/cheqd/.cheqdnode/config    ->    /etc/cheqd-node/
# cp ~/.cheqdnode/data/* /var/lib/cheqd/data/             # /var/lib/cheqd/.cheqdnode/data    ->    /var/lib/cheqd/data
cp ${NODE_CONFIGS_BASE}/node0/.cheqdnode/config/genesis.json ~/.cheqdnode/config/
sudo chmod -R 777 /etc/cheqd-node
sudo chmod -R 777 /var/lib/cheqd/data
sudo systemctl start cheqd-noded
systemctl status cheqd-noded