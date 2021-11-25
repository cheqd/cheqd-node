# ❓ Key terms and FAQs

### What is a **Validator** or **Node Operator**

In blockchain ecosystems, the **Node Operator** runs what is called a **node**. A node can be thought of like a power pylon in the physical world, which helps to distribute electricity around a wide network of users.

Without these pylons, electricity would be largely centralised in one location; the pylons help to distribute power to entire wide-scale populations. And if one pylon fails, the grid is set up to circumvent this pylon and re-route the electricity a different route.

Similarly, in blockchain infrastructure, each node runs an instance of the consensus protocol and helps to create a broad, robust network, with no single points of failure. A node failing will have little impact on the Network as a whole; however, if multiple nodes fail, or disagree with information entered into the transaction, then the block may not be signed, and there are fail-safe measures to notify the rest of the Node Operators of this.&#x20;

The terms Validator and Node Operator are somewhat synonymous. Validator is the term used more commonly in the[ Cosmos documentation](https://docs.cosmos.network) when referring to a Node Operator that is validating transactions on a blockchain. The only point worth mentioning is you can have a Node Operator that is NOT a Validator. These are known as Observer nodes which play a more passive role on the network, as they don’t stake on the network or validate transactions, but can observe them.&#x20;

&#x20;

### What does a **Validator** actually do?

The [Cosmos Hub](https://hub.cosmos.network/main/gaia-tutorials/what-is-gaia.html) is based on [Tendermint](https://tendermint.com/docs/introduction/what-is-tendermint.html), which relies on a set of validators to secure the network. By ‘secure the network’ this refers to the way in which validators “_participate in consensus by broadcasting votes which contain _[_cryptographic signatures signed by their private key_](https://hub.cosmos.network/main/validators/validator-faq.html)”.

A cryptographic signature is a mathematical scheme for verifying the authenticity of digital messages or documents. A private key is like a password — a string of letters and numbers — that allows you to access and manage your crypto funds (your mnemonic is a version of this). So, the above is saying validators can broadcast that they agree with transactions in a block, using their password to sign their agreement in a mathematical way which ensures security and privacy.&#x20;

### What does **staking** mean?

**Stake** is the amount of tokens a **Node Operator** puts aside and dedicates to a network’s active pool, in order to contribute to governance and earn rewards. **Staking** is the verb used to describe this contribution. As cheqd is a Proof of Stake (PoS) Network, rewards can be earnt in direct correlation with the amount of stake a Node Operator contributes. Tokens which are in the active pool are known as ‘bonded’ tokens.

The goal for bonded tokens on the Network is **60%**. Below **60%** bonded tokens, the rate of inflation will increase, tending to **4%** **(inflation max)**. Above **60%** bonded tokens, the rate of inflation will decrease, tending towards **1% (inflation min)**.

When you setup your Validator node, you also need to set a minimum self-delegation which is the minimum amount of tokens you promise to keep bonded. This can be as low as 1 ncheq.

You can earn rewards proportionate to your stake in the Network. These include:

#### 1. Block rewards

cheqd naturally has inflation, which means that it creates more blocks over time. Each time a new block is created, there is a chance that your Validator node will 'propose' that block (because it is in the validator pool). The higher the stake you hold, the higher this chance is. Over a period of time, by the law of averages, you will earn a reasonably consistent proportion of block rewards.

#### 2. Transaction fees

The more that DIDs are written to the Network (and in the future) verifiable credentials are verified by third parties, there will be transaction fees distributed to the network, which will be split between the Node Operators and those delegated to the Node Operators.

Gas prices also come into play here too, the lower your gas price, the more likely that your node will be considered in the active set for rewards.

You can read more about gas fees and prices here.

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

Bonded tokens are those present in the active pool.&#x20;

**Bonded **tokens = **staked** tokens by Node Operator + **delegated** tokens by user.


### As a Validator, is my commission rate important?

As a Validator Node, you must set a rate of commission. This is the percentage of tokens that you take as a fee for running the infrastructure on the network. Token holders are able to delegate tokens to you, with an understanding that they can earn staking rewards, but as consideration, you are also able to earn a flat percentage fee of the rewards on the delegated stake they supply.

A **lower commission rate** = **higher likelihood** of **more token holders delegating** tokens to you because they will **earn more rewards**.

A **higher commission rate** = **you earn more tokens** from the existing stake + delegated tokens. But the tradeoff being that it **may appear less desirable for new delegators** when compared to other Validators.

You can have a look at other projects on Cosmos here to get an idea of the percentages that nodes set as commission.

Please note, that once you set a maximum commission rate, **you are not able to change this going forward**. You will only be able to have a commission at that max or lower.

You also need to set a parameter called Commission Rate Max Change. This parameter sets the limit of which a commission rate can be changed within a single day. This helps protect those users who have delegated to Node Operators from delegating to a Node Operator with a low commission and it changing overnight. This parameter is also public, so a lower figure here is likely better if your goal as a Validator is to attract more users to delegate to you. We hope that you do your own research and make an informed decision on these parameters.


### **Governance** and **voting**

Users with **bonded **tokens, **Participants**, are able to vote on **Governance Proposals**. The weight of a vote is directly tied to the amount of bonded tokens delegated to Node Operators.&#x20;

The specifics of how a participant can vote on Proposals, or create Proposals, is detailed further in the rest of our [Governance Framework](../../contributing/major-network-changes/).&#x20;

If the User does not want to vote on a Governance Proposal or misses it for any particular reason, the Node Operator will **inherit **the **delegated **tokens and can use them to vote with.&#x20;

**Node Operator voting power** = **initial stake** + **inherited delegated tokens** (if participants do not vote)

**Participant voting power** = **delegated tokens** (if participant chooses to vote)\
\


### What if I want to **vote unilaterally**, i.e. without a Node Operator/Validator?

If you are particularly interested or passionate about a specific governance proposal, or do not agree with your bonded Node Operator, it is absolutely possible to vote unilaterally. However, you must still have delegated tokens, bonded with a Node Operator to do so. To do this, follow the instructions in the section[ Voting on cheqd](https://docs.cheqd.io/governance/contributing/voting-on-cheqd).

### Can **Participants** earn tokens?

In short, yes. **Participants **may be eligible for a percentage of the rewards that Node Operators earn during the course of running a node. Node Operators may offer a certain **commission **for delegating to them. Participants earn rewards on the delegated tokens which are bonded to the Validator. \


### How do I choose which **Node Operator** to **delegate** to?

Choosing your Node Operator or multiple Node Operators is an important decision. There a few things you can take into consideration:

1.  **Commission rate**



    The incentive for delegating tokens to a Node Operator, is that you, as a participant, can earn rewards based on the stake of the Node Operator. Each Node Operator will have a **commission rate. **This is the percentage of tokens that the Node Operator will take as **commission **for running the infrastructure - the remaining percentage is distributed to the **participants**.

2.  **Reputation**



    You should be mindful about what reputation the Node Operator has. This is because the Node Operator may use your votes against the best interest of yourself and the community. As cheqd evolves, it is likely that there will become a political spectrum of Node Operators, who will cast their vote in different directions. Some may want to create chaos on the network and vote to disrupt the established paradigms, for example. A chaotic actor may lure users to delegate to them with a favourable commission rate, and use the accumulated bonded tokens against the best interests of the network. For this reason the choice of Node Operator you delegate to is very important.&#x20;


3.  **Slashing and Validator Jail**



    As the name would suggest, staking is not risk free. As the word stake literally means “having something to gain or lose by having a form of ownership of something”, individuals should be wary of the risk, as we’ll come on to.&#x20;

    Think of it like this, if someone says to you “_what’s at stake_?” they are essentially asking: “_what am I risking in return for the potential rewards?_”

    Node Operators might exhibit bad behaviour on the Network and, as a result, have their stake slashed. Slashing means taking **stake** away from the Node Operator and adding it to the Community Pool.\
    \
    Bad behaviour in this context usually means that the Node Operator has not signed a sufficient number of blocks as ‘pre commits’ over a certain period of time. This could be due to inactivity or potential malicious intent.&#x20;

    For example in June 2019, CosmosPool a former Cosmos validator, experienced a server outage on their main node; downtime that resulted in its validator being temporarily jailed and its stake being slashed by 0.01%, including that of its delegators. This was what’s call a downtime slashing incident (soft slashing) whereby the validator and delegators were punished for downtime proportionally to their stake on the network ([JohnnieCosmos](https://johnniecosmos.medium.com/what-you-need-to-know-when-staking-on-the-cosmos-ecosystem-e6fc13a1b0e3)). On top of this, further slashing later occurred as evidence was found of double block signing. Both CosmosPool AND the delegators’ stakes were slashed an additional 5% and the validator was peremandiet removed (‘tombstoned’).&#x20;

    Slashing can therefore certainly affect your delegated and bonded tokens, so it is important to consider your choice.&#x20;

### What if I change my mind about a Node Operator? Is it possible to **redelegate** or **unbond** from a Node Operator?&#x20;

Yes, it is possible to instantly **redelegate** to a new Node Operator; however, you cannot 'hop' between Node Operators quickly. You must complete your redelegation to a Node Operator before moving again. If you want to completely withdraw and **unbond** your tokens, you need to be mindful that this process takes **2 weeks** for the tokens to become available for use again. This **unbonding period** is an example of a parameter which can be adjusted by the Governance Framework process.
