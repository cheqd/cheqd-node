#!/bin/bash
sudo chown -R runner:docker ${NODE_CONFIGS_BASE}/client
export HOME=${NODE_CONFIGS_BASE}/client
sudo -u cheqd cheqd-noded keys list --keyring-backend "test"
OP0_ADDRESS=$(sudo -u cheqd cheqd-noded keys list --keyring-backend "test" --home ${NODE_CONFIGS_BASE}/node0 | sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')
sudo -u cheqd cheqd-noded keys add node5-operator --keyring-backend "test"
OP5_ADDRESS=$(sudo -u cheqd cheqd-noded keys list --keyring-backend "test"| sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')
export HOME=/home/runner
NODE5_PUBKEY=$(sudo -u cheqd cheqd-noded tendermint show-validator | sed 's/\r//g')
HOME=${NODE_CONFIGS_BASE}/client sudo -u cheqd cheqd-noded tx bank send ${OP0_ADDRESS} ${OP5_ADDRESS} 1100000000000000ncheq --chain-id cheqd --fees 5000000ncheq --node "http://localhost:26657" -y --keyring-backend "test"
HOME=${NODE_CONFIGS_BASE}/client sudo -u cheqd cheqd-noded tx staking create-validator --amount 1000000000000000ncheq --from node5-operator --chain-id cheqd --min-self-delegation="1" --gas-prices="25ncheq" --pubkey ${NODE5_PUBKEY} --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01" --gas 500000 --node "http://localhost:26657" -y --keyring-backend "test"
