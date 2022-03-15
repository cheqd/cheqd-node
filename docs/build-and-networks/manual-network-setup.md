# Setting up a new network manually

## Creating a new network from genesis

This document describes in details how to configure a genesis network with any amount of participants.

### Requirements

#### Hardware requirements

Minimal:

* 1GB RAM
* 25GB of disk space
* 1.4 GHz CPU

Recommended (for highload applications):

* 2GB RAM
* 100GB SSD
* x64 2.0 GHz 2v CPU

#### Operating System

Current delivery is compiled and tested for `Ubuntu 20.04 LTS` so we recommend using this distribution for now. In the future, it will be possible to compile the application for a wide range of operating systems thanks to the Go language.

### Deployment steps

#### Setting up nodes

All participants should [setup their nodes](../setup-and-configure/README.md) using default genesis and empty list of peers.

#### Generating genesis file

1. Participants must choose `<chain_id>` for the network.
2. Each participant (one by one) should:
   * **Generate local keys for the future account:**

     Command `cheqd-noded keys add <key_name>`

     Example `cheqd-noded keys add alice`

   * **(Each participant except the first one) Get genesis config from the another participant:**

     Location on the previous participant's machine: `$HOME/.cheqdnode/config/genesis.json`

     Destination folder on the current participant's machine: `$HOME/.cheqdnode/config/`

   * **(Each participant except the first one) Get genesis node transactions from the previous participant:**

     Location on the previous participant's machine: `$HOME/.cheqdnode/config/gentx/`

     Destination folder on the current participant's machine: `$HOME/.cheqdnode/config/gentx/`

   * **Add a genesis account with a public key:**

     Command: `cheqd-noded add-genesis-account <key_name> 10000000cheq,100000000stake`

     Example: `cheqd-noded add-genesis-account alice 10000000cheq,100000000stake`

   * **Generate genesis node transaction:**

     Command: `cheqd-noded gentx <key_name> 1000000stake --chain-id <chain_id>`

     Example: `cheqd-noded gentx alice 1000000stake --chain-id cheqd-node`

     **TODO: Node owner should specify gas prices here. This work is in progress.**
3. The last participant:
   * **Add genesis node transactions into genesis:**

     Command: `cheqd-noded collect-gentxs`

   * **Verify genesis:**

     Command: `cheqd-noded validate-genesis`

   * **Share genesis with other nodes:**

     Location on the last participant's machine: `$HOME/.cheqdnode/config/genesis.json`

#### Sharing peer list

All participants should share their peer info with each other. See [node setup instruction](../setup-and-configure/README.md) for more information.

#### Updating genesis and persistent peers

* Each participant should:
  * **Stop the node:** `systemctl stop cheqd-noded`
  * **Make sure the node is stopped** `systemctl status cheqd-noded`
  * **Update the genesis file:**

    File location:

    * Deb destribution: `/etc/cheqd-node/genesis.json`
    * Binary destribution: `$HOME/.cheqdnode/config/genesis.json`

  * **Update peer list:**

    Open node's config file:

    * Deb destribution: `/etc/cheqd-node/config.toml`
    * Binary destribution: `$HOME/.cheqdnode/config/config.toml`

      Search for `persistent_peers` parameter and set it's value to a comma separated list of peers.

      Format: `<node-0-id>@<node-0-ip>, <node-1-id>@<node-1-ip>, <node-2-id>@<node-2-id>, <node-3-id>@<node-3-id>`.

      Domain names can be used instead of IP adresses.

      Example:

      ```text
      persistent_peers = "d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656, 9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656, abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656, cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656"
      ```

  * **Start node:** `systemctl start cheqd-noded`
  * **Make sure the node process is running:** `systemctl status cheqd-noded`

Congratulations, you should have node(s) deployed and running on a network if the above steps succeed.

## Support

Please log issues and any questions via GitHub Issues

