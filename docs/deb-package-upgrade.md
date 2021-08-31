# Overview

This document describes how to upgrade to a new release of `cheqd-node` using the Debian package \(.deb\) releases.

This document assumes that `cheqd-node` has already been installed previously using a .deb package. See [installation instructions](https://github.com/cheqd/cheqd-node/blob/main/docs/setting-up-a-new-node.md) if this your first time setting up a node.

## Get the latest package

The latest Debian package can be found in [releases](https://github.com/cheqd/cheqd-node/releases)

## Steps to upgrade via Debian package

1. Ensure that current `cheqd-noded` service is stopped. Make [backups of app data](https://github.com/cheqd/cheqd-node/blob/main/docs/deb-package-overview.md#directories-and-symlinks) and keys before package upgrading. To stop the node service:

   ```text
    systemctl stop cheqd-noded
   ```

   and

   ```text
    systemctl status cheqd-noded
   ```

   to confirming that service was stopped.

2. The new package version can be installed by calling:

   ```text
    dpkg -i <path/to/package>
   ```

   Depending on how your system is configured, you may need `sudo` or administrator permissions to carry out the step above.

3. Start `cheqd-noded` service by calling:

   ```text
    systemctl start cheqd-noded
   ```

   and confirm that the service is running:

   ```text
    systemctl status cheqd-noded
   ```

4. If the `cheqd-noded` service is running, the package upgrade has been successful.

## Steps to carry out after package upgrade

The package upgrade/installation process is idempotent and it should not change service files, configs or any other user data.

To check that your node is functioning correctly, it is recommended to attempt [querying the ledger](https://github.com/cheqd/cheqd-node/blob/main/docs/cosmos-cli.md) or [any of the other commands described in the cheqd Cosmos CLI guide](https://github.com/cheqd/cheqd-node/blob/main/docs/cosmos-cli.md).

If you are running a validator node, check that it's [connected to the network](setting-up-a-new-node.md) \(see check section\) and [signs blocks](setting-up-a-new-validator.md) \(see check section\).

## Further information

Additional information about [cheqd Debian package release can be found here in the overview](deb-package-overview.md).

