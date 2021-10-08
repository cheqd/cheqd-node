# Upgrading a cheqd node using Debian package releases

## Context

This document provides guidance on how to upgrade to an [existing installation of `cheqd-node` that was done using the Debian package](deb-package-install.md) release to a new release version.

It is assumed that the [pre-requisites mentioned in the node setup guide](../readme.md) are satisfied, as a node has already been installed.

Before carrying out an upgrade, please read our [guide to Debian packages for `cheqd-node`](readme.md) to understand an overview of what configuration actions are carried out by the installer.

## Upgrade steps for `cheqd-node` .deb

| :warning: WARNING          |
|:---------------------------|
| Please make sure any accounts keys are backed up or exported before attempting uninstallation |

The package upgrade process is idempotent and it should not affect service files, configurations or any other user data.

However, as best practice we recommend backing up the [app data directories for `cheqd-node`](readme.md)  and Cosmos account keys before attempting the upgrade process.

1. **Download [the latest release of `cheqd-node` .deb](https://github.com/cheqd/cheqd-node/releases/latest) package**

    ```bash
    wget https://github.com/cheqd/cheqd-node/releases/download/v0.2.3/cheqd-node_0.2.3_amd64.deb
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

   Install the `cheqd-node` package downloaded in step 1 (with `sudo` privileges or as `root` user, if necessary):

   ```bash
   dpkg -i <path/to/package>
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

## Next steps

For further confirmation on whether your node is working correctly, we recommend attempting to [run commands from the cheqd CLI guide](../../cheqd-cli.md); e.g., query the ledger for transactions, account balances etc.
