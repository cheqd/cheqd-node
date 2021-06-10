# Localnet

## Build docker image

See [instruction](sdfsdf)

Commands to setup localnet:

```
starport build
./genlocalnetconfig.sh
docker-compose up --build
```

This will setup 4 nodes listening on the following ports:

- Node0:
    - p2p: 26656
    - rpc: 26657
- Node1:
    - p2p: 26666
    - rpc: 26667
- Node2:
    - p2p: 26676
    - rpc: 26677
- Node3:
    - p2p: 26686
    - rpc: 26687

You can tests connection to a node using browser: `http://localhost:<rpc_port>`. Example for the fitst node: `http://localhost:26657`.

When connecting using CLI, point path to home directory: `--home localnet/client`. This directory contains keys from genesis acounts.

Demo commands:

```
# Show balances
vc query bank balances (vc keys show anna -a --home localnet/client) --home localnet/client

# Create NYM
vc tx verimcosmos create-nym "alias" "verkey" "did" "role" --from anna --gas-prices 1token --chain-id verim-cosmos-chain --home localnet/client

# List nym
vc query verimcosmos list-nym --home localnet/client
```
