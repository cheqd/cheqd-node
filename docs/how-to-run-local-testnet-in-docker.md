# Running a Local Testnet in Docker

This document describes steps how to run local testnet in docker.

## Prerequisites

You need to have Go and Node.js toolchains installed.

## Deploy steps

1. Build node executable:
    ```
    starport build
    ```
2. Use `genlocalnetconfig.sh` script to generate node configurations:
    ```
    ./genlocalnetconfig.sh
    ```
3. Use docker-compose to run nodes:
    ```
    docker-compose up --build
    ```

This will setup 4 nodes listening on the following ports:

|     | Node0 | Node1 | Node2 | Node3 |
|-----|-------|-------|-------|-------|
| P2P | 26656 | 26666 | 26676 | 26686 |
| RPC | 26657 | 26667 | 26677 | 26687 |


You can tests connection to a node using browser: `http://localhost:<rpc_port>`. Example for the fitst node: `http://localhost:26657`.

Keys of all accounts are located in `localnet/client`. When connecting using CLI, point path to home directory: `--home localnet/client`. This directory contains keys from genesis acounts.

## Command examples:

Show balances:

```
verim-cosmosd query bank balances $(verim-cosmosd keys show anna -a --home localnet/client) --home localnet/client
```

Create NYM:

```
verim-cosmosd tx verimcosmos create-nym "alias" "verkey" "did" "role" --from anna --gas-prices 1token --chain-id verimcosmos --home localnet/client
```

List NYMs:

```
verim-cosmosd query verimcosmos list-nym --home localnet/client
```
