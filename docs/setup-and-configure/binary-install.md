# Installing a cheqd node from binary package releases

## Context

This document provides guidance on how to install and configure a node for the cheqd testnet from [our binary package releases](https://github.com/cheqd/cheqd-node/releases/latest).

If you are planning to install a node on a host machine running Ubuntu 20.04 LTS, our recommended method is to use the [Debian package installation guide](debian/deb-package-install.md) instead.

The [node setup guide provides pre-requisites](README.md) needed before the steps below should be attempted.

## Installation steps for `cheqd-node` .tar.gz binary

1. Get binary

   You can get the binary in several ways:

   * [Compile from source](../build-and-networks/README.md)
   * Get `tar` archive with the binary compiled for Ubuntu 20.04 in [releases](https://github.com/cheqd/cheqd-node/releases);

2. Define configuration for running the `cheqd-noded` binary as a service

   It is highly recommended to run the `cheqd-node` as a system service using a supervisor such as `systemd`.

   Our Debian package uses [postinst](../../build-tools/postinst) script for setting up our binary as a service. The same tool can be used to set up the binary as a service.

   There is only one input parameter for `postinst` script, it's a path to where binary is.

   To set up the binary using `postint`, execute the following with sudo privileges:

   ```bash
   bash postinst <path/to/cheqd-noded/binary>
   ```

   This will add a service file and prepare all needed directories for `configs/keys` and `data`. The script also creates a new service user called `cheqd`, to ensure that all processes and directories related to `cheqd-noded` are isolated under that service user.

3. Configure node and run service

   Continue following the instructions from [Step 4 in the Debian package installation guide](debian/deb-package-install.md) to complete the installation process.

## Further information

After successful installation, learn [how to configure your node as a validator node](../validator-guide/README.md) to participate in staking rewards, block creation, and governance.
