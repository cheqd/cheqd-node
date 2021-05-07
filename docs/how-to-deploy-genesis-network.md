# Deploy genesis network

This document describes in details how to configure a genesis (first) validator node.

Here we assume (for simplicity) that the genesis block consists of a single node only. Please note that nothing prevents you from adding more nodes to the genesis file by adapting the instructions accordingly.

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

- Add key for new client `verim-cosmosd keys add <name-client> --home localnet/client`
- Optionally, add other clients using the same command.
2. Initialize node
    - Initialize new configuration `verim-cosmosd init <name-node> --chain-id <chain-id> --home localnet/<name-node>`
    - Copy genesis file `cp -r localnet/client/* localnet/<name-node>`
    - Add genesis account with the generated key `verim-cosmosd add-genesis-account <name-client> 1000token,100000000stake --home localnet/<name-node>`
    - Optionally, add other genesis accounts using the same command.
    - Create genesis transaction: `verim-cosmosd gentx <name-client> 1000000stake --chain-id <chain-id> --home localnet/<name-node>`
3. Collect all validator creation transactions:
    - Create `gentx` directory `mkdir localnet/client/config/gentx`
    - Copy `gentx` for every node `cp localnet/<name-node>/config/gentx/* localnet/client/config/gentx`
5. Embed them into genesis:
    - `verim-cosmosd init dummy-node --chain-id <chain-id> --home localnet/client`
    - `verim-cosmosd add-genesis-account <name-client> 1000token,100000000stake --home localnet/client`
    - Repeat last substep for every client.
    - `verim-cosmosd collect-gentxs --home localnet/client`
    - `verim-cosmosd validate-genesis --home localnet/client`
6. Update config for every node:
    - `cp localnet/client/config/genesis.json localnet/<name-node>/config/`
    - Find out ip address and id of all nodes, then open file `localnet/<name-node>/config/config.toml`, where `<name-node>` is name of first node, and set value for param `persistent_peers="<id-node1>@<ip-node1>,<id-node2>@<ip-node2>"`.
    - Set `laddr="tcp:\/\/127.0.0.1:26657"` for RPC server to listen on, and set `laddr = "tcp:\/\/0.0.0.0:26657"` for incoming connections, for every node config `localnet/<name-node>/config/config.toml`.
7. Run your nodes:
    - `verim-cosmosd start --home localnet/<name-node>`.
8. Congrats! You deploy own network!
