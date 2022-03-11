# Installing a cheqd node using Debian package releases

## Context

This document provides guidance on how to install and configure a node for the cheqd testnet.

* Our [guide to Debian packages for `cheqd-node`](README.md) provides an overview of what configuration actions are carried out by the installer.
* The [node setup guide provides pre-requisites](../README.md) needed before the steps below should be attempted.
* Separate instructions apply if you're looking to [upgrade an existing cheqd node](deb-package-upgrade.md).

## Installation steps for `cheqd-node` .deb

1. **Download** [**the latest release of `cheqd-node` .deb**](https://github.com/cheqd/cheqd-node/releases/latest) **package**

   For example:

   ```bash
   wget https://github.com/cheqd/cheqd-node/releases/download/v0.4.1/cheqd-node_0.4.1_amd64.deb
   ```

2. **Install the package**

   For example:

   ```bash
   sudo dpkg -i cheqd-node_0.4.1_amd64.deb
   ```

   As a part of installation `cheqd` user will be created. By default, `HOME` directory for the user is `/home/cheqd`, but it can be changed by setting `CHEQD_HOME_DIR` environment variable before running `dpkg` command. Additionnally, a custom logging directory can also be defined by passing the environment variable `CHEQD_LOG_DIR` (defaults to `/home/cheqd/.cheqdnode/log`):

   Example custom directories:

   ```bash
   sudo CHEQD_HOME_DIR=/path/to/desired/home/directory dpkg -i cheqd-node_0.4.1_amd64.deb
   ```

3. **Switch to the `cheqd` system user**

   Always switch to `cheqd` user before managing node. By default, the node stores configuration files in the home directory for the user that initialises the node. By switching to the special `cheqd` system user, you can ensure that the configuration files are stored in the [system app data directories](README.md) configured by the Debian package.

   ```bash
   sudo su cheqd
   ```

4. **Initialise the node configuration files**

   The "moniker" for your node is a "friendly" name that will be available to peers on the network in their address book, and allows easily searching a peer's address book.

   ```bash
   cheqd-noded init <node-moniker>
   ```

5. **Download the genesis file for a persistent chain, such as the cheqd testnet**

   Download the `genesis.json` file for the relevant [persistent chain](https://github.com/cheqd/cheqd-node/tree/main/persistent_chains/) and put it in the `$HOME/.cheqdnode/config` directory.

   For cheqd mainnet:

   ```bash
   wget -O $HOME/.cheqdnode/config/genesis.json https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/mainnet/genesis.json
   ```

   For cheqd testnet:

   ```bash
   wget -O $HOME/.cheqdnode/config/genesis.json https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/genesis.json
   ```

6. **Define the seed configuration for populating the list of peers known by a node**

   Update `seeds` with a comma separated list of seed node addresses specified in `seeds.txt` for the relevant [persistent chain](https://github.com/cheqd/cheqd-node/tree/main/persistent_chains/).

   For cheqd mainnet, set the `SEEDS` environment variable:

   ```bash
   SEEDS=$(wget -qO- https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/mainnet/seeds.txt)
   ```

   For cheqd testnet, set the `SEEDS` environment variable:

   ```bash
   SEEDS=$(wget -qO- https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/seeds.txt)
   ```

   After the `SEEDS` variable is defined, pass the values to the `cheqd-noded configure` tool to set it in the configuration file.

   ```bash
   $ echo $SEEDS
   # Comma separated list should be printed
   
   $ cheqd-noded configure p2p seeds "$SEEDS"
   ```

7. **Set gas prices accepted by the node**

   Update `minimum-gas-prices` parameter if you want to use custom value. The default is `25ncheq`.

   ```bash
   cheqd-noded configure min-gas-prices "25ncheq"
   ```

8. **Turn off empty block creation**

   By default, the underlying Tendermint consensus creates blocks even when there are no transactions ("empty blocks").

   Turning off empty block creation is an optimisation strategy to limit growth in chain size, although this only works if a majority of the nodes opt-in to this setting.

   ```bash
   cheqd-noded configure create-empty-blocks false
   ```

9. **Define the external peer-to-peer address**

   Unless you are running a node in a sentry/validator two-tier architecture, your node should be reachable on its peer-to-peer (P2P) port by other other nodes. This can be defined by setting the `external-address` property which defines the externally reachable address. This can be defined using either IP address or DNS name followed by the P2P port (Default: 26656).

   ```bash
   cheqd-noded configure p2p external-address <ip-address-or-dns-name:p2p-port>
   # Example
   # cheqd-noded configure p2p external-address 8.8.8.8:26656
   ```

   This is especially important if the node has no public IP address, e.g., if it's in a private subnet with traffic routed via a load balancer or proxy. Without the `external-address` property, the node will report a private IP address from its own host network interface as its `remote_ip`, which will be unreachable from the outside world. The node still works in this configuration, but only with limited unidirectional connectivity.

10. **Make the RPC endpoint available externally** (optional)

      This step is necessary only if you want to allow incoming client application connections to your node. Otherwise, the node will be accessible only locally. Further details about the RPC endpoints is available in the [cheqd node setup guide](../README.md).

      ```bash
      cheqd-noded configure rpc-laddr "tcp://0.0.0.0:26657"
      ```

11. **Enable and start the `cheqd-noded` system service**

      If you are prompted for a password for the `cheqd` user, type `exit` to logout and then attempt to execute this as a privileged user (with `sudo` privileges or as root user, if necessary).

      ```bash
      $ systemctl enable cheqd-noded
      Created symlink /etc/systemd/system/multi-user.target.wants/cheqd-noded.service → /lib/systemd/system/cheqd-noded.service.

      $ systemctl start cheqd-noded
      ```

      Check that the `cheqd-noded` service is running. If successfully started, the status output should return `Active: active (running)`

      ```bash
      systemctl status cheqd-noded
      ```

## Post-installation checks

Once the `cheqd-noded` daemon is active and running, check that the node is connected to the cheqd testnet and catching up with the latest updates on the ledger.

### Checking node status via terminal

```bash
cheqd-noded status
```

In the output, look for the text `latest_block_height` and note the value. Execute the status command above a few times and make sure the value of `latest_block_height` has increased each time.

The node is fully caught up when the parameter `catching_up` returns the output `false`.

### Checking node status via the RPC endpoint

An alternative method to check a node's status is via the RPC interface, if it has been configured.

* Remotely via the RPC interface: `cheqd-noded status --node <rpc-address>`
* By opening the JSONRPC over HTTP status page through a web browser: `<node-address:rpc-port>/status`

## Next steps

At this stage, your node would be connected to the cheqd testnet as an observer node. Learn [how to configure your node as a validator node](../../validator-guide/README.md) to participate in staking rewards, block creation, and governance.
