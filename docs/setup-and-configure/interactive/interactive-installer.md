# Overview

This document provides information about how-to use new interactive installer and describes the main options.

## Installation process

By default, our target platform already has `python3` under the hood and no additional packages are needed and other preparation steps.

For running installer the next command can be used:

```bash
wget -q https://raw.githubusercontent.com/cheqd/cheqd-node/3866783bd3282dcb7fb908cc6b72840cf137a41f/installer/installer.py && sudo python3 installer.py
```

## Questions

All the questions at the end have the default value in [] brackets, like `[v0.5.0]`. If a default value exists you can just press `Enter` without needing to type the whole answer.

During the installation process you will see the following questions:

* `Which version do you want to install? Or type 'list' to get the list of releases: [v0.5.0]`. Possible answers are the exact versions that currently exist, for example type `0.4.0` or `0.5.0` or `list`(in this case the last 5 releases will be printed and you can choose what the version is needed)

```text
Which version below do you want to install? Please enter the list item number below [1]: 
list
1) v0.5.0
2) v0.4.1
3) v0.3.5
4) v0.3.4
5) v0.3.3
Please insert the number by picking up the version: 1
```

* `Please, type the path to home directory for user cheqd. To keep the default value, just type 'Enter': [/home/cheqd]`. Here you need to specify the path to home directory for the new user `cheqd`. By default `/home/cheqd` will be used.
* `Do you want to run a node or just use the CLI? [No]`. If you are installing the node from the beginning, you can use this ability to setup your node. Possible variants for answering `y, n, yes, no`. If the answer was `Yes`, the next questions will be about the config settings.
* `Do you want to use Cosmovisor? Please type any kind of variants (yes/no) [yes]:`. With current installer we are proposing the ability to setup cosmovisor. It will help you with upgrades, it allows to do it in the full automatic mode. Possible variants for answering (yes/no) [yes].
* `Which chain do you want to use? Possible variants are: testnet, mainnet [testnet]`. For now, we have 2 networks, `testnet` and `mainnet`. Please, type here which chain you want to use or just keep the default by clicking `Enter`.
* `Do you want to download the latest ledger snapshot to speed up node setup process? Please type any kind of variants: yes, no, y, n. [No]`. This can help you speed up the catchup to cheqd network. Possible variants for answering `y, n, yes, no`.
* `Do you want to setup systemctl services after installation? Please type any kind of variants: yes, no, y, n. [Yes]`. If you are installing the node from the beginning, you can use this ability to setup your node. Possible variants for answering `y, n, yes, no`. If the answer was `Yes`, the next questions will be about the config settings.
* `Please, type the moniker for your node:`. Here you need to specify a nickname for your node. This will be used in various user facing places such as on block explorers where a user is deciding where to delegate their tokens.
* `What is the external address of your node? Please type in format: <ip-address-or-dns-name>:<port>`. Here you need to specify the external address of your machine and P2P port also. For example, `8.8.8.8:26656`.

P.S. Due to the size of the snapshots it will take a long time to download. During this period the script will print some message about the process each minute.

P.P.S It's possible to run the installer again in case of failure or a typo however it **will not override already created files.**

### Example of installing

```text
Which version do you want to install? Or type 'list' to get the list of releases: [v0.5.0]

Please type the path to the home directory for user cheqd. To keep the default value, just type 'Enter': [/home/cheqd]
/root
Do you want to use Cosmovisor? Please type any kind of variants: yes, no, y, n. [yes]

Which chain do you want to use? Possible variants are: testnet, mainnet [testnet]

Do you want to deploy the latest snapshot from https://snapshots.cheqd.net? Please type any kind of variants: yes, no, y, n. [No]
y
Which snapshot do you want to use? Please type the full URL to archive or click return to use the latest [https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/testnet/latest/cheqd-testnet-4_2022-06-10.tar.gz]

Do you want to setup a node after installation? Please type any kind of variants: yes, no, y, n. [No]
y
Please, type the moniker for your node:
test
What is external IP address for your node? Please type in format: <ip_address>:<port>

*********  Download the binary
*********  Executing command: wget -qO - https://github.com/cheqd/cheqd-node/releases/download/v0.5.0/cheqd-node_0.5.0.tar.gz  | tar xz
*********  Create a user cheqd cause it's not created yet
*********  Create group, cheqd by default
*********  Executing command: addgroup cheqd --quiet
*********  Create user, cheqd by default
*********  Executing command: adduser --system cheqd --home /root --shell /bin/bash --ingroup cheqd --quiet
*********  Make root directory for cheqd-node
*********  Chown to default cheqd user: cheqd
*********  Executing command: chown -R cheqd:cheqd /root/.cheqdnode
*********  Setup log directory
*********  Executing command: chown -R syslog:syslog /root/.cheqdnode/log
*********  Configure rsyslog
*********  Executing command: systemctl restart rsyslog
*********  Add config for logrotation
*********  Executing command: systemctl restart rsyslog
*********  Restart logrotate services
*********  Executing command: systemctl restart logrotate.service
*********  Executing command: systemctl restart logrotate.timer
*********  Setup systemctl service config
*********  Enable systemctl service
*********  Executing command: systemctl enable cheqd-noded
*********  Setup the cosmovisor
*********  Executing command: wget -qO - https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz  | tar xz
*********  Moving binary from /root/cheqd-noded to /root/.cheqdnode/cosmovisor/genesis/bin/cheqd-noded
*********  Executing command: sudo mv /root/cheqd-noded /root/.cheqdnode/cosmovisor/genesis/bin/cheqd-noded
*********  Making symlink to /root/.cheqdnode/cosmovisor/genesis/bin/cheqd-noded
*********  Changing owner to cheqd user
*********  Executing command: chown -R cheqd:cheqd /root/.cheqdnode/cosmovisor
*********  Executing command: sudo -u cheqd cheqd-noded init test
*********  Executing command: curl -s https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/genesis.json > /root/.cheqdnode/config/genesis.json
*********  Executing command: curl -s https://raw.githubusercontent.com/cheqd/cheqd-node/main/persistent_chains/testnet/seeds.txt
*********  Executing command: sudo -u cheqd cheqd-noded configure p2p seeds 658453f9578d82f0897f13205ca2e7ad37279f95@seed1.eu.testnet.cheqd.network:26656,eec97b12f7271116deb888a8d62e0739b4350fbd@seed1.us.testnet.cheqd.network:26656,32d626260f74f3c824dfa15a624c078f27fc31a2@seed1.ap.testnet.cheqd.network:26656
*********  Going to download the archive and untar it on a fly. It can take a really LONG TIME
*********  Directory /root/.cheqdnode/data already exists
*********  Executing command: wget -qO - https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/testnet/latest/cheqd-testnet-4_2022-06-10.tar.gz  | sudo -u cheqd tar xzf - -C /root/.cheqdnode/data
*********  Downloading is alive, it already took: 0:01:00
*********  Downloading is alive, it already took: 0:02:00
*********  Downloading is alive, it already took: 0:03:00
*********  Downloading is alive, it already took: 0:04:00
*********  Downloading is alive, it already took: 0:05:00
*********  Downloading is alive, it already took: 0:06:00
*********  Downloading is alive, it already took: 0:07:00
*********  Downloading is alive, it already took: 0:08:00
*********  Making symlink current -> genesis
*********  Copying upgrade-info.json file to cosmovisor/current/
*********  Changing owner to cheqd user
*********  Executing command: chown -R cheqd:cheqd /root/.cheqdnode/cosmovisor
*********  Executing command: chown -R cheqd:cheqd /root/.cheqdnode/data
```

After the installation process ends you can start the `systemctl` service:

```bash
sudo systemctl start cheqd-noded
```
