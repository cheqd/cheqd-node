# Upgrading a cheqd node using Debian package releases

## Context

This document provides guidance on how to upgrade to an [existing installation of `cheqd-node` that was done using the Debian package](deb-package-install.md) release to a new release version.

It is assumed that the [pre-requisites mentioned in the node setup guide](../README.md) are satisfied, as a node has already been installed.

Before carrying out an upgrade, please read our [guide to Debian packages for `cheqd-node`](README.md) to understand an overview of what configuration actions are carried out by the installer.

## Upgrade steps for `cheqd-node` .deb

| :warning: WARNING |
| :--- |
| Please make sure you have backed up `node_key.json`, `priv_validator_key.json` and any account keys are backed up or exported before attempting any upgrade |

The package upgrade process is idempotent and it should not affect service files, configurations or any other user data.

However, as best practice we recommend backing up the [app data directories for `cheqd-node`](README.md) and Cosmos account keys before attempting the upgrade process.

1. **Download** [**the latest release of `cheqd-node` .deb**](https://github.com/cheqd/cheqd-node/releases/latest) **package**

   ```bash
    wget https://github.com/cheqd/cheqd-node/releases/download/v0.5.0/cheqd-node_0.5.0_amd64.deb
   ```

2. **Stop the existing `cheqd-noded` service**

   To stop the `cheqd-noded` service (with `sudo` privileges or as `root` user, if necessary):

   ```bash
    systemctl stop cheqd-noded
   ```

   Confirm the `cheqd-noded` service has been successfully stopped:

   ```bash
    systemctl status cheqd-noded
   ```

3. **Install the new .deb package version**

   | :warning: WARNING |
   | :--- |
   | If you are [upgrading from v0.2.x to any higher release version](#upgrade-from-02x), the default home directory folder has changed and may need to be manually configured |

   | :warning: WARNING |
   | :--- |
   | If you are [upgrading from `0.4.0` to any higher release version](#upgrade-from-040), we recommend to remove previous one and install the new package instead of "simple installing" (over `0.4.x`). |

   Install the `cheqd-node` package downloaded (with `sudo` privileges or as `root` user, if necessary):

   ```bash
   dpkg -i <path/to/package>
   ```

   To specify [a custom home directory location](deb-package-install.md), use the following command instead:

   ```bash
   sudo CHEQD_HOME_DIR=/path/to/home/directory dpkg -i cheqd-node_0.5.0_amd64.deb
   ```

4. **Re-start the `cheqd-noded` service and confirm it is running**

   To start the `cheqd-noded` service (with `sudo` privileges or as `root` user, if necessary):

   ```bash
   systemctl start cheqd-noded
   ```

   Check that the `cheqd-noded` service is running. If successfully started, the status output should return `Active: active (running)`

   ```bash
   systemctl status cheqd-noded
   ```

## Post-upgrade steps

The package upgrade process is successful once the service re-started.

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

## Upgrade from `0.2.x`

One of the changes made from v0.2.x to the higher software release versions was to allow the home directory for cheqd data to be customised.

The default home directory in v0.2.x used to be `/var/lib/cheqd`, but can be modified as desc:

```bash
sudo CHEQD_HOME_DIR=/path/to/home/directory dpkg -i cheqd-node_0.3.3_amd64.deb
```

In general, it's not required and up to system administrators how to ensure safe recovering after crashes.

If you have `0.2.3` version installed and you want to follow the new `$HOME` directory approach the next steps can help with it:

* Please define the mount point for `cheqd` root directory where all the configs and data will be placed. For example, let it be `/cheqd`.
* Stop `cheqd-noded` service by running:

   ```bash
   sudo systemctl stop cheqd-noded
   ```

* Install `.deb` package for `0.3.1` version:

   ```bash
   sudo CHEQD_HOME_DIR=/cheqd dpkg -i cheqd-node_0.3.1_amd64.deb
   ```

* After that the next directory tree is expected:

   ```bash
   /cheqd/.cheqdnode/data
   /cheqd/.cheqdnode/config
   /cheqd/.cheqdnode/log
   ```

* After that you should move all the configs from previous location into the new one `/cheqd/.cheqdnode/config`, data into `/cheqd/.cheqdnode/data`. It's assumed that root directory `/cheqd` will be stored and mounted as external resource and will not be removed after potential instance crashing.
* For logs symlink can be created by using command:

   ```bash
   ln -s /cheqd/.cheqdnode/log /var/log/cheqd-node
   ```

* Start `cheqd-noded` service by running:

   ```bash
   sudo systemctl start cheqd-noded
   ```

and check the service status or just check RPC endpoint.

## Upgrade from `0.4.0`

Due to changes in the `postremove` script it's highly recommended not to install packages `0.4.1` and higher version just over `0.4.0` one cause it requires double installing for it.
We recommend using the next schema for upgrading from `0.4.0` version:

* Remove `0.4.0` package by calling:

    ```bash
    sudo dpkg -r cheqd-node
    ```

* Install new package in a general way:

    ```bash
    sudo dpkg -i <cheqd-node-package>
    ```

## Next steps

For further confirmation on whether your node is working correctly, we recommend attempting to [run commands from the cheqd CLI guide](../../cheqd-cli/README.md); e.g., query the ledger for transactions, account balances etc.
