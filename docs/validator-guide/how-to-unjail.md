# Unjailing a Validator Node

Jailing of a Validator Node occurs when it experiences downtime, and is unable to sign a requisite amount of block precommits in a given period.

The specific parameters for jailing are:

1. Validator signs less than **50%** of the blocks in the signed block window
2. The signed block window is 12,100 blocks
3. The number of blocks per second is roughly 1 block every 6 seconds
4. The Validator would therefore need to be down for 50% of a 50 hour timeframe

If your Validator Node is jailed, do not worry - this guide is intended to help you get it back up and running as quickly as possible.

## Step 1: Check your Node is up to date

During the downtime of a Validator Node, it is common for the Node to miss important software upgrades, since they are no longer in the active set of nodes on the main ledger. 

Therefore, the first step is checking that your node is up to date. You can execute the command

~~~
cheqd-noded version
~~~

The expected response will be the latest cheqd-noded software release. At the time of writing, the expected response would be 

~~~
0.5.0
~~~

## Step 2: Upgrading to latest software

If your node is not up to date, please [follow the instructions here](https://github.com/cheqd/cheqd-node/blob/main/docs/setup-and-configure/debian/deb-package-upgrade.md)

## Step 3: Conforming the Node is up to date

Once again, check if your node is up to date, following Step 1.

Expected response: In the output, look for the text ```latest_block_height``` and note the value. Execute the status command above a few times and make sure the value of ```latest_block_height``` has increased each time.

The node is fully caught up when the parameter ```catching_up``` returns the output false.

Additionally,, you can check this has worked:
~~~
http://<your node ip or domain name>:26657/abci_info
~~~
It shows you a page and field "version": "0.5.0"
____

## Step 4: Unjailing command

If everything is up to date, and the node has fully caught, you can now unjail your node using this command in the cheqd CLI:

~~~
cheqd-noded tx slashing unjail --from <address_alias> --gas auto --gas-adjustment 1.2 --gas-prices 25ncheq --chain-id cheqd-mainnet-1
~~~
