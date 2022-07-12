#!/bin/bash

set -euox pipefail
sudo chown -R cheqd:cheqd "/home/runner/cheqd/"

sudo -u cheqd -H cheqd-noded init node5

if [ -z ${GENESIS_PATH+x} ]; then
  GENESIS_PATH=${NODE_CONFIGS_BASE}/node0/.cheqdnode/config/genesis.json
fi

VALIDATOR_0_ID=`cheqd-noded tendermint show-node-id --home ${NODE_CONFIGS_BASE}/node0/.cheqdnode`

PERSISTENT_PEERS="${VALIDATOR_0_ID}@127.0.0.1:26656"
sudo -u cheqd -H cheqd-noded configure p2p persistent-peers "${PERSISTENT_PEERS}"

sudo cp "${GENESIS_PATH}" "/home/runner/cheqd/.cheqdnode/config"

sudo chmod -R 755 "/home/runner/cheqd/.cheqdnode"

# Configure ports because they conflict with localnet
sudo -u cheqd -H cheqd-noded configure p2p laddr "tcp://0.0.0.0:26676"
sudo -u cheqd -H cheqd-noded configure rpc-laddr "tcp://0.0.0.0:26677"

# TODO: Use environment variables
sudo sed -i.bak 's|pprof_laddr = "localhost:6060"|pprof_laddr = "localhost:6070"|g' /home/runner/cheqd/.cheqdnode/config/config.toml
sudo sed -i.bak 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9100"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo sed -i.bak 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9101"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo sed -i.bak 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:1327"|g' /home/runner/cheqd/.cheqdnode/config/app.toml
sudo sed -i.bak 's|address = ":8080"|address = ":8090"|g' /home/runner/cheqd/.cheqdnode/config/app.toml

sudo chown -R cheqd:cheqd "/home/runner/cheqd/"

sudo systemctl start cheqd-cosmovisor
systemctl status cheqd-cosmovisor
sleep 10
journalctl --since "2 days ago" | grep cosmovisor

bash wait.sh "[[ $(cheqd-noded status -n 'tcp://localhost:26677' 2>&1 | wc -l) == 1 ]] && echo \"Observer node is up\""

NODE_CONFIGS_BASE="${NODE_CONFIGS_BASE}" bash promote-validator.sh

bash check-promotion.sh
# shellcheck disable=SC2016
bash wait.sh '[[ $(curl -s localhost:26657/block | sed -nr '"'"'s/.*signature": (.*?).*/\1/p'"'"' | wc -l) == 5 ]] && echo "There are 5 validators signatures in block!"'
# shellcheck disable=SC2016
bash wait.sh '[[ $(curl -s localhost:26657/block | sed -nr '"'"'s/.*(signature": null).*/\1/p'"'"' | wc -l) == 0 ]] && echo "There are no null signatures in block!"'

