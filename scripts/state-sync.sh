#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -euox pipefail

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
export RPC=https://eu-rpc.cheqd.net:443
export RPCN=https://ap-rpc.cheqd.net:443
export APPNAME=CHEQD_NODED

# Install Gaia
go install -tags goleveldb ./...

# MAKE HOME FOLDER AND GET GENESIS
cheqd-noded init "$NODE_MONIKER" --home /home/cheqd
wget -O /home/cheqd/.cheqdnode/config/genesis.json https://github.com/canow-co/cheqd-node/raw/main/networks/mainnet/genesis.json
wget -O seeds.txt https://github.com/canow-co/cheqd-node/raw/main/networks/mainnet/seeds.txt

INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s $RPC/block | jq -r .result.block.header.height);
BLOCK_HEIGHT=$((LATEST_HEIGHT-INTERVAL))
TRUST_HASH=$(curl -s "$RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"


# export state sync vars
export ${APPNAME}_STATESYNC_ENABLE=true
export ${APPNAME}_P2P_MAX_NUM_OUTBOUND_PEERS=500
export ${APPNAME}_STATESYNC_RPC_SERVERS="$RPC,$RPCN"
export ${APPNAME}_STATESYNC_TRUST_HEIGHT="$BLOCK_HEIGHT"
export ${APPNAME}_STATESYNC_TRUST_HASH="$TRUST_HASH"
export ${APPNAME}_P2P_SEEDS="$(cat seeds.txt)"


cheqd-noded start --x-crisis-skip-assert-invariants --home /home/cheqd/.cheqdnode --grpc-web.address 127.0.0.1:9091


# THIS WILL FIX THE APP VERSION, contributed by callum and claimens
git clone https://github.com/tendermint/tendermint
cd tendermint
git checkout remotes/origin/callum/app-version
go install ./...
tendermint set-app-version 1 --home ~/home/cheqd/.cheqdnode


cheqd-noded start --x-crisis-skip-assert-invariants --home /home/cheqd/.cheqdnode --grpc-web.address 127.0.0.1:9091
