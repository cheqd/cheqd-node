# Frequently Asked Questions (FAQs) for Validators

## How do I **stake** more tokens after setting up a validator node?

When you set up your Validator node, it is recommended that you only stake a very small amount from the actual Validator node. This is to minimise the tokens that could be locked in an unbonding period, were your node to experience signficiant downtime.

You should delegate the rest of your tokens to your Validator node from a different key alias.

**How do I do this?**

You can add (as many as you want) additional keys you want using the function:

```bash
cheqd-noded keys add <alias>
```

When you create a new key, a mnemonic phrase and account address will be printed. Keep the mnemonic phrase safe as this is the only way to restore access to the account if they keyring cannot be recovered.

You can view all created keys using the function:

```bash
cheqd-noded keys list
```

You are able to transfer tokens between key accounts using the function.

```bash
cheqd-noded tx bank send <from> <to-address> <amount> --node <url> --chain-id <chain> --gas auto --gas-adjustment 1.2
```
  
You can then delegate to your Validator Node, using the function

```bash
cheqd-noded tx staking delegate <validator address> <amount to stake> --from <key alias> --gas auto --gas-adjustment 1.2 --gas-prices 25ncheq 
```

We use a second/different Virtual Machine to create these new accounts/wallets. In this instane, you only need to install cheqd-noded as a binary, you don't need to run it as a full node.

And then since this VM is not running a node, you can then append the --node parameter to any request and target the RPC port of the VM running the actual node.

That way:

1. The second node doesn't need to sync the full blockchain; and
2. You can separate out the keys/wallets, since the IP address of your actual node will be public by definition and people can attack it or try to break in

## What is Commission rate and is it important?

As a Validator Node, you should be familiar with the concept of commission. This is the percentage of tokens that you take as a fee for running the infrastructure on the network. Token holders are able to delegate tokens to you, with an understanding that they can earn staking rewards, but as consideration, you are also able to earn a flat percentage fee of the rewards on the delegated stake they supply.

There are three commission values you should be familiar with:

```text
max_commission_rate

max_commission_rate_change

commission_rate
```

The first is the maximum rate of commission that you will be able to move upwards to.

**Please note that this value cannot be changed once your Validator Node is set up**, so be careful and do your research.

The second parameter is the maximum amount of commission you will be able to increase by within a 24 hour period. For example if you set this as 0.01, you will be able to increase your commission by 1% a day.

The third value is your current commission rate.

Points to note: **lower commission rate** = **higher likelihood** of **more token holders delegating** tokens to you because they will **earn more rewards**. However, with a very low commission rate, in the future, you might find that the gas fees on the Network outweight the rewards made through commission.

**higher commission rate** = **you earn more tokens** from the existing stake + delegated tokens. But the tradeoff being that it **may appear less desirable for new delegators** when compared to other Validators.

You can have a look at other projects on Cosmos [here](https://www.mintscan.io/cosmos) to get an idea of the percentages that nodes set as commission.

## What is Gas and Gas Prices?

When setting up the Validator, the Gas parameter is the amount of tokens you are willing to spend on gas.

For simplicity, we suggest setting:

```text
--gas: auto
```

AND setting:

```text
--gas-adjustment: 1.2
```

These parameters, together, will make it highly likely that the transaction will go through and not fail. Having the gas set at auto, without the gas adjustment will endanger the transaction of failing, if the gas prices increase.

Gas prices also come into play here too, the lower your gas price, the more likely that your node will be considered in the active set for rewards.

We suggest the set:

```text
--gas-price
```

should fall within this recommended range:

Low: 25ncheq
Medium: 50ncheq
High: 100ncheq


## How do I change my public name and description

Your public name, is also known as your **moniker**.

You are able to change this, as well as the description of your node using the function:

```bash
cheqd-noded tx staking edit-validator --from validator1-eu --moniker "cheqd" --details "cheqd is building a private and secure decentralised digital identity network on the Cosmos ecosystem" --website "https://www.cheqd.io" --identity "F0669B9ACEE06ADC" --security-contact security@cheqd.io --gas auto --gas-adjustment 1.2 --gas-prices 25ncheq -chain-id cheqd-mainnet-1
```

## Should I set my firewall port 26656 open to the world?

Yes, this is how you should do it. Since it's a public permissionless network, there's no way of pre-determining what the set of IP addresses will be, as entities may leave and join the network. We suggest using a TCP/network load balancer and keeping your VM/node in a private subnet though for security reasons. The LB then becomes your network edge which if you're hosting on a cloud provider they manage/patch/run.
