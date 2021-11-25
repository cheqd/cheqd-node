# Node Operator Setup Help

### Staking help

**Stake** is the amount of tokens a **Node Operator** puts aside and dedicates to a network’s active pool, in order to contribute to governance and earn rewards. **Staking** is the verb used to describe this contribution. As cheqd is a Proof of Stake (PoS) Network, rewards can be earnt in direct correlation with the amount of stake a Node Operator contributes. Tokens which are in the active pool are known as ‘bonded’ tokens.

The goal for bonded tokens on the Network is **60%**. Below **60%** bonded tokens, the rate of inflation will increase, tending to **4%** **(inflation max)**. Above **60%** bonded tokens, the rate of inflation will decrease, tending towards **1% (inflation min)**.

When you setup your Validator node, you also need to set a minimum self-delegation which is the minimum amount of tokens you promise to keep bonded. This should be approximately 0.001 CHEQ (1000000 ncheq).

Instead of keeping all our tokens in a key directly associated with the node, we suggest allocating your validator 0,001 CHEQ, and then using a separate wallet on a different machine to delegate tokens to your own validator.

Since the delegator key/wallet is on a completely different server, and can even be offline, even if the validator server gets breached or loses its keys, you don't lose control of a large balance. Validators don't have the private keys or the mnemonics of those who account that delegate to them, which allows you to maintain an air gap between what's stored on the validator vs your keyring on other machines.

You can redelegate tokens easily, rather than having to carry out a full unbonding. This may be hugely beneficial going forward as redelegating is much faster than unbonding. 

Through your initial stake and delegated stake, you can earn rewards proportionate to your stake in the Network. These include:

#### 1. Block rewards

cheqd naturally has inflation, which means that it creates more blocks over time. Each time a new block is created, there is a chance that your Validator node will 'propose' that block (because it is in the validator pool). The higher the stake you hold, the higher this chance is. Over a period of time, by the law of averages, you will earn a reasonably consistent proportion of block rewards.

#### 2. Transaction fees

The more that DIDs are written to the Network (and in the future) verifiable credentials are verified by third parties, there will be transaction fees distributed to the network, which will be split between the Node Operators and those delegated to the Node Operators.

Gas prices also come into play here too, the lower your gas price, the more likely that your node will be considered in the active set for rewards.

You can read more about gas fees and prices [here](https://https://docs.cosmos.network/master/basics/gas-fees.html).

For context, we suggest the set gas-price should fall within this range: Low: 25ncheq Medium: 50ncheq High: 100ncheq

### How do I **stake** more tokens after setting up a validator node?

You can add (as many as you want) additional keys you want using the function **_cheqd-noded keys_** add and then transfer tokens over to it. Then use the new account to delegate to the validator's operator account

We use a second/different VM sometimes to create these new accounts/wallets. You only need to install cheqd-noded as a binary, you don't need to run it as a full node.
And then since this VM is not running a node, you append the --node parameter to any request and target the RPC port of the VM running the actual node.

That way:
1.  The second node doesn't need to sync the full blockchain; and
2.  You can separate out the keys/wallets, since the IP address of your actual node will be public by definition and people can attack it or try to break in

### What is **delegation**?

Token holders, ‘**users’, **can** delegate** their tokens to Node Operators in order to earn rewards and participate in governance. Once a user has delegated tokens to a Node Operator and has tokens added to the active pool, they are known as **Participants.**

Users can delegate to multiple Node Operators at the same time to essentially diversify their bonded token portfolio.

### What does **bonded** mean?

Bonded tokens are those present in the active pool.

**Bonded** tokens = **staked** tokens by Node Operator + **delegated** tokens by user.


### What is Commission rate and is it important?

As a Validator Node, you must set a rate of commission. This is the percentage of tokens that you take as a fee for running the infrastructure on the network. Token holders are able to delegate tokens to you, with an understanding that they can earn staking rewards, but as consideration, you are also able to earn a flat percentage fee of the rewards on the delegated stake they supply.

A **lower commission rate** = **higher likelihood** of **more token holders delegating** tokens to you because they will **earn more rewards**.

A **higher commission rate** = **you earn more tokens** from the existing stake + delegated tokens. But the tradeoff being that it **may appear less desirable for new delegators** when compared to other Validators.

You can have a look at other projects on Cosmos [here](https://www.mintscan.io/cosmos) to get an idea of the percentages that nodes set as commission.

Please note, that once you set a maximum commission rate, **you are not able to change this going forward**. You will only be able to have a commission at that max or lower.

You also need to set a parameter called Commission Rate Max Change. This parameter sets the limit of which a commission rate can be changed within a single day. This helps protect those users who have delegated to Node Operators from delegating to a Node Operator with a low commission and it changing overnight. This parameter is also public, so a lower figure here is likely better if your goal as a Validator is to attract more users to delegate to you. We hope that you do your own research and make an informed decision on these parameters.

### What is Gas and Gas Prices?

When setting up the Validator, the Gas parameter is the amount of tokens you are willing to spend on gas.

For simplicity, we suggest setting --gas: auto

AND setting --gas-adjustment: 1.2

These parameters, together, will make it highly likely that the transaction will go through and not fail. Having the gas set at auto, without the gas adjustment will endanger the transaction of failing, if the gas prices increase. 

Gas prices also come into play here too, the lower your gas price, the more likely that your node will be considered in the active set for rewards.

We suggest the set gas-price should fall within this range:

Low: 25ncheq(recommended) 
Medium: 50ncheq 
High: 100ncheq 

### Should I set my firewall port 26656 open to the world?

Yes, this is how you should do it. Since it's a public permissionless network, there's no way of pre-determining what the set of IP addresses will be, as entities may leave and join the network. We suggest using a TCP/network load balancer and keeping your VM/node in a private subnet though for security reasons. The LB then becomes your network edge which if you're hosting on a cloud provider they manage/patch/run.
