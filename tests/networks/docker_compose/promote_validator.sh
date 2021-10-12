sudo chown -R runner:docker ${NODE_CONFIGS_BASE}/client
export HOME=${NODE_CONFIGS_BASE}/client
cheqd-noded keys list
OP0_ADDRESS=$(cheqd-noded keys list | sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')
cheqd-noded keys add node5-operator
OP5_ADDRESS=$(cheqd-noded keys list | sed -nr 's/.*address: (.*?).*/\1/p' | sed -n 1p | sed 's/\r//g')
export HOME=/home/runner
NODE5_PUBKEY=$(cheqd-noded tendermint show-validator | sed 's/\r//g')
HOME=${NODE_CONFIGS_BASE}/client cheqd-noded tx bank send ${OP0_ADDRESS} ${OP5_ADDRESS} 1100000000000000ncheq --chain-id cheqd --fees 5000000ncheq --node "http://localhost:26657" -y
HOME=${NODE_CONFIGS_BASE}/client cheqd-noded tx staking create-validator --amount 1000000000000000ncheq --from node5-operator --chain-id cheqd --min-self-delegation="1" --gas="auto" --gas-prices="25ncheq" --pubkey $NODE5_PUBKEY --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01" --gas 239933 --node "http://localhost:26657" -y