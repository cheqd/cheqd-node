# Setup a new cheqd node

This document describes in detail how to configure infrastructure and deploy a new node as an observer or validator.

After creating the nodes, if a new network needs to be initialized, please follow the instructions for [creating a new network from genesis](../setting-up-a-new-network.md).

If a new validator needs to be added to the existing network, please refer to [joining existing network](configure-new-validator.md) instruction.

## Setting up infrastructure

### Hardware requirements

Minimal:

* 1GB RAM
* 25GB of disk space
* 1.4 GHz CPU

Recommended \(for highload applications\):

* 2GB RAM
* 100GB SSD
* x64 2.0 GHz 2v CPU

Extended information on [recommended hardware requirements is available in Tendermint documentation](https://docs.tendermint.com/master/nodes/running-in-production.html#hardware).

### Operating System

Our [packaged releases](https://github.com/cheqd/cheqd-node/releases) are currently compiled and tested for `Ubuntu 20.04 LTS`, which is the recommended operating system in case the installation is carried out using Debian package or binaries.

For other operating systems, we recommend using [pre-built Docker image releases for `cheqd-node`](https://github.com/orgs/cheqd/packages?repo_name=cheqd-node).

We plan on supporting other operating systems in the future based on demand for specific platforms by the community.

### Ports

To function properly, `cheqd-node` requires two types of ports to be configured.

#### P2P port

* This port is used for peer-to-peer node communication
* Incoming and outcoming TCP connections must be allowed from any IPv4 address
* `26656` by default
* Can be configured in `/etc/cheqd-node/config.toml`

#### RPC port

* This port is used by client applications. Open it only if you want clients to be able to connect to your node.
* Incoming tcp connections should be allowed.
* SSL can also be configured separately
* `26657` by default
* Can be configured in `/etc/cheqd-node/config.toml`

### Volumes

We recommend using a separate storage volume for the `data` directory where the node's copy of the ledger is stored.

The default directory location depends on the installation method used:

* For binary distribution, it is `$HOME/.cheqdnode/data`
* For installations done using our Debian packages, it is `/var/lib/cheqd/.cheqdnode/data`.

### Sentry nodes \(optional\)

Tendermint allows more complex setups in production, where the ingress/egress to a validator node is [proxied behind a "sentry" node](https://docs.tendermint.com/master/nodes/validators.html#setting-up-a-validator).

While this setup is not compulsory, node operators with higher stakes or a need to have more robust network security may consider setting up a sentry-validator node architecture.

## Installing and configuring software

### Installing using .deb package

The recommended way to install `cheqd-node` on a standalone \(virtual\) machine is to use our Debian package installer on Ubuntu 20.04 LTS. Detailed information about changes made by the package can be found [here](https://github.com/cheqd/cheqd-node/tree/8ea79cda7c3fb5f41ab90faba7ed9e42d2096dce/docs/deb-package-overview.md)

1. Get `deb` package for Ubuntu 20.04 in [releases](https://github.com/cheqd/cheqd-node/releases):

   ```text
   wget https://github.com/cheqd/cheqd-node/releases/download/v0.2.3/cheqd-node_0.2.3_amd64.deb
   ```

2. Install the package:

   ```text
   sudo dpkg -i cheqd-node_0.2.3_amd64.deb
   ```

3. Switch to `cheqd` system user:

   You should always switch to `cheqd` user before managing node. That's because node stores configuration files in home directory which is different for each user.

   ```text
   sudo su cheqd
   ```

4. Initialize node config files:

   ```text
   cheqd-noded init <your-node-name>
   ```

5. Set genesis:

   Genesis files for persistent chains are published in [this directory](https://github.com/cheqd/cheqd-node/tree/main/persistent_chains). Download corresponding `genesis.json` and put it to the `/etc/cheqd-node/`.

   For `testnet`:

   ```bash
   $ wget -O /etc/cheqd-node/genesis.json https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/genesis.json
   ```

6. Set seeds:

   Open node's config file: `/etc/cheqd-node/config.toml`. Search for `seeds` parameter and set it's value to a comma separated list of seed node addresses from `seeds.txt` from [this directory](https://github.com/cheqd/cheqd-node/tree/main/persistent_chains).

   For `testnet`:

   ```bash
   $ SEEDS=$(wget -qO- https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/seeds.txt)

   $ echo $SEEDS
   # Comma separated list should be printed

   $ sed -i.bak 's/seeds = ""/seeds = "'$SEEDS'"/g' /etc/cheqd-node/config.toml
   ```

7. Set gas prices:

   Open app's config file: `/etc/cheqd-node/app.toml`. Search for `minimum-gas-prices` parameter and set it to a non-empty value. Recommended one is `25ncheq`.

   For `testnet`:

   ```bash
   sed -i.bak 's/minimum-gas-prices = ""/minimum-gas-prices = "25ncheq"/g' /etc/cheqd-node/app.toml
   ```

8. \(optional\) Make RPC endpoint available externally:

   This step is necessary if you want to allow incoming client application connections to your node. Otherwise, the node will be accessible only locally.

   Open the node configuration file `/etc/cheqd-node/config.toml` in any text editor. Search for `laddr` parameter in `RPC Server Configuration Options` section and replace it's value to `0.0.0.0:26657`

   For `testnet`:

   ```bash
   sed -i.bak 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' /etc/cheqd-node/config.toml
   ```

9. Enable `cheqd-noded` service and start it:

   ```text
   systemctl enable cheqd-noded

   systemctl start cheqd-noded
   ```

   Check that the service is running:

   ```text
   systemctl status cheqd-noded
   ```

   Status should be: `Active: active (running)`

10. Check that the node is connected and catching up:

    ```text
    cheqd-noded status
    ```

    Make sure that `latest_block_height` is increasing over time.

    Wait for `catching_up` to become `false`.

    * You can query status remotely: `cheqd-noded status --node <rpc-address>`;
    * Another way is to open status page in browser: `<rpc-address>/status`.

### Installing using binary

1. Get binary

   You can get the binary in several ways:

   * Compile from source code - [instruction](../../);
   * Get `tar` archive with the binary compiled for Ubuntu 20.04 in [releases](https://github.com/cheqd/cheqd-node/releases);

2. Set up `cheqd-noded` binary as a service

   It is highly recommended to run the `cheqd-node` as a system service using a supervisor such as `systemd`.

   Our Debian package uses [postinst](https://github.com/cheqd/cheqd-node/blob/main/build_tools/postinst) script for setting up our binary as a service. The same tool can be used to set up the binary as a service.

   There is only one input parameter for `postinst` script, it's a path to where binary is.

   To set up the binary using `postint`, execute the following with sudo privileges:

   ```text
   # bash postinst <path/to/cheqd-noded/binary>
   ```

   This will add a service file and prepare all needed directories for `configs/keys` and `data`. The script also creates a new service user called `cheqd`, to ensure that all processes and directorioes related to `cheqd-noded` are isolated under that service user.

3. Configure node and run service

   Follow `Installing using .deb package` section starting form step 2.

### Other ways

* Get docker image form [packages](https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-node).

## Additional information

You can read other advices about running node in production [here](https://docs.tendermint.com/master/nodes/running-in-production.html).

[Сosmovisor](https://docs.cosmos.network/master/run-node/cosmovisor.html) can be used for automatic updates.

