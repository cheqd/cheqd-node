#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -ux

if [ -z ${1+x} ];
then
  echo "Node moniker must be passed as the first parameter"
  exit 1
else
  export NODE_MONIKER=${1}
fi

# set environment variables
export GOPATH=~/go
export PATH=$PATH:~/go/bin
export RPC=https://rpc.cheqd.net:443
export RPCN=https://rpc.cheqd.net:443
export APPNAME=CHEQD_NODED

# Install Gaia
go install -tags rocksdb ./...

# MAKE HOME FOLDER AND GET GENESIS
cheqd-noded init $NODE_MONIKER --home /others/cheqd
wget -O /others/cheqd/config/genesis.json https://github.com/cheqd/cheqd-node/raw/main/persistent_chains/mainnet/genesis.json

INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s $RPC/block | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT-INTERVAL))
TRUST_HASH=$(curl -s "$RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"


# export state sync vars
export $(echo $APPNAME)_STATESYNC_ENABLE=true
export $(echo $APPNAME)_P2P_MAX_NUM_OUTBOUND_PEERS=500
export $(echo $APPNAME)_STATESYNC_RPC_SERVERS="$RPC,$RPCN"
export $(echo $APPNAME)_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export $(echo $APPNAME)_STATESYNC_TRUST_HASH=$TRUST_HASH
export $(echo $APPNAME)_P2P_SEEDS="258a9bfb822637bfca87daaab6181c10e7fd0910@seed1.eu.cheqd.net:26656,f565ff792b20977face9817df6acb268d41d4092@seed2.eu.cheqd.net:26656,388947cc7d901c5c06fedc4c26751634564d68e6@seed3.eu.cheqd.net:26656,9b30307a2a2819790d68c04bb62f5cf4028f447e@seed1.ap.cheqd.net:26656,debcb3fa7d40e681d98bcc7d22278fd58a34b73a@144.76.183.180:1234,abd4be300be882ae9a69ab0959260afe8871f7a6@165.232.167.156:26656,fd17fe46e8b69bfa006b3fba53cfc9df2b8f9512@161.35.177.151:26656"


cheqd-noded start --x-crisis-skip-assert-invariants --home /others/cheqd --grpc-web.address 127.0.0.1:5050


# THIS WILL FIX THE APP VERSION, contributed by callum and claimens
git clone https://github.com/tendermint/tendermint
cd tendermint
git checkout remotes/origin/callum/app-version
go install ./...
tendermint set-app-version 1 --home ~/others/cheqd


cheqd-noded start --x-crisis-skip-assert-invariants --home /others/cheqd --grpc-web.address 127.0.0.1:5050
