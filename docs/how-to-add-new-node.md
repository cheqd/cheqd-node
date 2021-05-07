# Running a Validator Node

This document describes in details how to configure a validator node, and add it to the existing network.

The existing network can be either a custom one, or one of the persistent networks (such as a Test Net).

If a new network needs to be initialized, please follow the Running Genesis Node instructions first. After this more validator nodes can be added by following the instructions from this doc.

### Hardware requirements
Minimal:
- 1GB RAM
- 25GB of disk space
- 1.4 GHz CPU

Recommended (for highload applications):
- 2GB RAM
- 100GB SSD
- x64 2.0 GHz 2v CPU

### Operating System
Current delivery is compiled and tested under `Ubuntu 20.04 LTS` so we recommend using this distribution for now. In future, it will be possible to compile the application for a wide range of operating systems thanks to Go language.

## Deployment steps

1. In the beginning you need to prepare your generate keys:
    - Create keys `verim-cosmosd keys add <key-name>`;
    - You need to get amount currency by another member:
        <!-- - Find out your address `verim-cosmosd keys list`; -->
        - Share your address for another member (command for show addresses `verim-cosmosd keys list`);
        - Wait transaction for your account.
2. Init your node and connect it to the network as observer:
    - Find out `<chain-id>` via command `verim-cosmosd status` (it will be showed as `network`);
    - Initialize node's configuration files `verim-cosmosd init <node-name> --chain-id <chain-id>`;
    - Copy `genesis.json` to node's config directory`$HOME/.verimcosmos/config`;
    - Run node `verim-cosmosd start`;
    - Find out ip address and id of another node, then open file `<key-name>/config/config.toml` and set value for param `persistent_peers="<id-node>@<ip-node>"`.
3. Create validator:
    - Run command with your params:
    ```verim-cosmosd tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation="1" --gas <amount-gas> --gas-prices <price-gas> --pubkey <pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>```
    where commission-max-change-rate, commission-max-rate and commission-rate may take value fraction number as `0.01`



