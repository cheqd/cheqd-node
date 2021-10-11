#!/bin/bash

cd
cheqd-noded init node5
cp ${{ env.NODE_CONFIGS_BASE }}/node0/.cheqdnode/config/genesis.json ~/.cheqdnode/config/
NODE0_IP=$(docker inspect -f {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} docker_compose_node0_1)
echo $NODE0_IP
PEER0=$(cat ${{ env.NODE_CONFIGS_BASE }}/node0/node_id.txt)@$NODE0_IP:26656
echo $PEER0
sed -ri "s|persistent_peers = \".*\"|persistent_peers = \"${PEER0}\"|" ~/.cheqdnode/config/config.toml
sed -ri "s|laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:26677\"|" ~/.cheqdnode/config/config.toml
sed -ri "s|laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://0.0.0.0:26676\"|" ~/.cheqdnode/config/config.toml
cheqd-noded start