# Contributing

Thank you for considering contributing to cheqd!

You can contribute to the documentation in the following ways:

* Improving general documentation
* Contributing code to cheqd by fixing bugs or implementing features
* Improving the governance framework
* Finding and reporting bugs;
* Translating the cheqd documentation into various languages

You can contribute to the development of the Network in the following ways:

* Informal Proposal and discussion
* Formal 'on-ledger' Proposal
* Voting with your tokens
* Delegating your tokens to a Node Operator 

## Revisions and amendments

You are able to make revisions and amendments to the wiki and source code without making a formal or informal Proposal. 

cheqd is an Open Source project which means that anyone is free to propose a revision, addition or amendment. 

Changes can be made through two routes:

1. [GitBook](https://docs.cheqd.io/cheqd-node/)
2. [Github](https://github.com/cheqd)

### GitBook

GitBook is where cheqd's Wiki lives and where YOU can make suggested changes to the written documentation. 

To make a change, you need to use the link below to become an Open Source Contributor on cheqd's GitBook:

{% embed url="https://app.gitbook.com/invite/cheqd?invite=-MiVxCDUlLSB22RuQ6dl" %}

You should follow these instructions to make a change to any cheqd  GitBook content.

1. Make your own branch

To create a new branch, select 'New' followed by 'New variant'.

![](../../.gitbook/assets/image%20%281%29.png)

Next you need to name your new branch.

You should use the prefix:

* ADR/
* gov/

followed by a unique title indicating your draft number such as

* ADR/{username}draft1
* gov/{username}draft1

On this branch you should make and save all desired changes to the content.

Once you have finished your changes, save the branch but **do not merge the branch to main**. 

Your edits will be reviewed by a cheqd admin or moderator and then merged, amended or rejected.

### Github

You may also use Github to make changes and amendments to both the source code and the written content in this documentation.

You should similarly make a branch of the cheqd Github main. 

You should then make a pull request with your proposed changes, edits, revisions and additions. 

**Please use clean, concise titles for your pull requests.** Assume that the person reading the pull request title is not a programmer, but instead a cheqd Network user or lay-person, and **try to describe your change or fix from in plain English**. We use commit squashing, so the final commit in the main branch will carry the title of the pull request, and commits from the main branch are fed into the changelog. The changelog is separated into [keepachangelog.com categories](https://keepachangelog.com/en/1.0.0/), and while that spec does not prescribe how the entries ought to be named, for easier sorting, start your pull request titles using one of the verbs "Add", "Change", "Deprecate", "Remove", or "Fix" \(present tense\).

Example:

| Not ideal | Better |
| :--- | :--- |
| Fixed NoMethodError in RemovalWorker | Fix nil error when removing statuses caused by race condition |

It is not always possible to phrase every change in such a manner, but it is desired.

**The smaller the set of changes in the pull request is, the quicker it can be reviewed and merged.** Splitting tasks into multiple smaller pull requests is often preferable.

**Pull requests that do not pass automated checks may not be reviewed**. 

## Bug reports

Bug reports and feature suggestions must use descriptive and concise titles and be submitted to [GitHub Issues](https://github.com/cheqd/cheqd-node/issues). Please use the search function to make sure that you are not submitting duplicates, and that a similar report or request has not already been resolved or rejected.

## Translations

You can submit translations via Github or GitBook branches \(as above\). We would greatly appreciate this work!

## Proposals

One of the most important questions in this Governance Framework is explaining how any token holder can make a proposal or voice their opinion on the Network. 

There are two ways of doing this: 

1. **Informal ‘off-chain’ proposal**
2. **Formal ‘on-chain’ proposal’**

These will be discussed in turn.

### Informal off-chain proposal

Rather than making a proposal directly to the Network, proposals can be made off-chain. Off-chain governance is vital for building a healthy and active governance community.

Once consensus is reached in an off-chain forum, the parties to the discussion can have more confidence that a Proposal will reach minimum deposit and be approved on-chain. 

At present, Proposals SHOULD be made and discussed on:

1. Github Discussions

They can also be discussed by the community on:

1. [Telegram](https://t.me/cheqd), or
2. Discord

Engagement is likely to be critical to the success of a proposal.

The degree to which you engage with the cheqd community should be relative to the potential impact that your proposal may have on the stakeholders.

There are many different ways to engage. One strategy involves a few stages of engagement before and after submitting a proposal on chain. Why do it in stages? It's a more conservative approach to save resources. The idea is to check in with key stakeholders at each stage before investing more resources into developing your proposal.

#### Stage 1: Your Idea

In the first stage of this strategy, you should engage people \(ideally experts\) informally about your idea.

* Does it make sense?
* Are there critical flaws?
* Does it need to be reconsidered?

#### Not yet confident about your idea?

Don’t worry! Governance proposals potentially impact many stakeholders. Introduce your idea with known members of the community before investing resources into drafting a formal proposal. Don't let negative feedback dissuade you from exploring your idea if you think that it's still important.

If you know people who are very involved with cheqd, send them a private message with a concise overview of what you think will result from your idea or proposed changes. Wait for them to ask questions before providing details. Do the same in semi-private channels where people tend to be respectful \(and hopefully supportive\).  We recommend the [cheqd Telegram Community](https://t.me/cheqd).

#### Confident with your idea?

Great! However, we still recommend that you introduce your idea with members of the community before investing resources into drafting a proposal. At this point you should seek out and carefully consider critical feedback in order to protect yourself from [confirmation bias](https://en.wikipedia.org/wiki/Confirmation_bias). This is the ideal time to see a critical flaw, because submitting a flawed proposal will waste resources.

#### **Are you ready to draft a governance proposal?**

There will likely be differences of opinion about the value of what you're proposing to do and the strategy by which you're planning to do it. If you've considered feedback from broad perspectives and think that what you're doing is valuable and that your strategy should work, and you believe that others feel this way as well, it's likely worth drafting a proposal. 

A conservative approach is to have some confidence that you roughly have initial support from a good proportion of the voting power before proceeding to drafting your proposal. However, there are likely other approaches, and if your idea is important enough, you may want to pursue it regardless of whether or not you are confident that the voting power will support it.

#### Stage 2: Your Draft Proposal

Begin with a well-considered draft proposal. Please use our proposal template here.

The next major section outlines and describes some potential elements of drafting a proposal. Ensure that you have considered your proposal and anticipated questions that the community will likely ask.

The ideal format for a proposal is as a Markdown file \(ie. .md\) in a Github repo. Markdown is a simple and accessible format for writing plain text files that is easy to learn. See the [Github Markdown Guide](https://guides.github.com/features/mastering-markdown/) for details on writing markdown files.

If you don't have a [Github](http://github.com/) account already, register one. 

Then fork this repository, draft your proposal in the proposals directory, and make a pull-request back to this repository. For more details on using Github, see the [Github Forking Guide](https://guides.github.com/activities/forking/). If you need help using Github, don't be afraid to ask someone!

If you really don't want to deal with Github, you can always draft a proposal in Word or Google Docs, or directly in the forums, or otherwise. However, Markdown on Github is the ultimate standard for distributed collaboration on text files.

Engage the community with your draft proposal

1. Post a draft of your proposal as a topic in the 'governance' category of the cheqd forum. Ideally this should contain a link to this repository, either directly to your proposal if it has been merged, or else to a pull-request containing your proposal if it has not been merged yet.
2. Directly engage key members of the community for feedback. These could be large contributors, those likely to be most impacted by the proposal, and entities with high stake-backing \(eg. high-ranked Validators; large stakers\).
3. Target members of the community in a semi-public way before bringing the draft to a full public audience. 
4. Alert the entire community to the draft proposal via
   * Twitter, tagging accounts such as the [cheqd account](https://twitter.com/cheqd_io)
   * [Telegram](https://t.me/cheqd)

#### Submit your proposal to the testnet

We intend to expand this [guide to include testnet instructions](https://github.com/cosmos/governance/blob/master/submitting.md#submitting-your-proposal-to-the-testnet).

You may want to submit your proposal to the testnet chain before the mainnet for a number of reasons, such as wanting to see what the proposal description will look like, to share what the proposal will look like in advance with stakeholders, and to signal that your proposal is about to go live on the mainnet.

Perhaps most importantly, for parameter change proposals, you can test the parameter changes in advance \(if you have enough support from the voting power on the testnet\).

Submitting your proposal to the testnet increases the likelihood of engagement and the possibility that you will be alerted to a flaw before deploying your proposal to mainnet.

### Formal on-chain proposal

Once you have sensibly tested your proposal and bounced your ideas around the community, you are ready to submit a proposal on-chain.

#### Formatting the JSON file for the governance proposal

Prior to sending the transaction that submits your proposal on-chain, you must create a JSON file. This file will contain the information that will be stored on-chain as the governance proposal. Begin by creating a new text \(.txt\) file to enter this information. Use these best practices as a guide for the contents of your proposal. When you're done, save the file as a .json file. See the examples that follow to help format your proposal.

Each proposal type is unique in how the JSON should be formatted. See the relevant section for the type of proposal you are drafting:

1. **TextProposal**: All the proposals that do not involve a modification of the source code go under this type. For example, an opinion poll would use a proposal of type _**TextProposal**_.
2. **SoftwareUpgradeProposal**: If accepted, Validators are expected to update their software in accordance with the proposal. They must do so by following a 2-steps process described in the [Software Upgrade](https://docs.cosmos.network/v0.43/modules/gov/01_concepts.html#software-upgrade) section below. Software upgrade roadmap may be discussed and agreed on via _**TextProposals**_, but actual software upgrades must be performed via _**SoftwareUpgradeProposals**_.
3. **CommunityPoolSpendProposal**: details a proposal for use of community funds, together with how many coins are proposed to be spent, and to which recipient account.
4. **ParameterChangeProposal**: defines a proposal to change one or more parameters. If accepted, the requested parameter change is updated automatically by the proposal handler upon conclusion of the voting period.
5. **CancelSoftwareUpgradeProposal**: is a gov Content type for cancelling a software upgrade.

To create a new Proposal type, you can propose a _**ParameterChangeProposal**_ with a custom handler, to perform another type of state change. 

Once on-chain, most people will rely upon network explorers to interpret this information with a graphical user interface \(GUI\).

This is the command format for using cheqd’s CLI \(the command-line interface\) to submit your proposal on-chain:  


```text
VDR CLI tx gov submit-proposal \
  --title=<title> \
  --description=<description> \
  --type="Text" \
  --deposit="2170nanocheq" \
  --from=<name> \
  --chain-id=<chain_id>

```

#### Deposit

To prevent spam, proposals must be submitted with a deposit in the coins defined in the _**MinDeposit**_ param. The voting period will not start until the proposal's deposit equals _**MinDeposit**_.

When a proposal is submitted, it has to be accompanied by a deposit that must be strictly positive, but can be inferior to _**MinDeposit**_. The submitter doesn't need to pay for the entire deposit on their own. If a proposal's deposit is inferior to _**MinDeposit**_, other token holders can increase the proposal's deposit by sending a Deposit transaction. 

The deposit is kept in an escrow in the governance _**ModuleAccount**_ until the proposal is finalized \(passed or rejected\).

Once the proposal's deposit reaches _**MinDeposit**_, it enters the voting period. If a proposal's deposit does not reach _**MinDeposit**_ before _**MaxDepositPeriod**_, the proposal closes and nobody can deposit on it anymore.

In this scenario, the tokens spent on the Deposit which did not reach the _**MinDeposit**_ will be burnt, meaning that they will be removed from the active pool of tokens and put beyond use. 

The minimum deposit for cheqd will initially be 80,000 CHEQ tokens.   


#### Deposit refund and burn

When a proposal is finalized, the coins from the deposit are either refunded or burned, according to the final tally of the proposal:

* If the proposal is approved or if it's rejected but not vetoed, deposits will automatically be refunded to their respective depositor \(transferred from the governance _**ModuleAccount**_\).
* When the proposal is vetoed by 33.34%, deposits will be burned from the governance _**ModuleAccount**_.

#### Proposal type

If &lt;proposal type&gt; is left blank, the type will be a Text proposal. Otherwise, it can be set to _**param-change**_, _**community-pool-spend**_, _**software-upgdrade**_ or _**cancel-software-upgrade**_. Use _**--help**_ to get more info from the tool.

For instance, this is the complete command that I could use to submit a testnet parameter-change proposal right now: 

```text
VDR CLI tx gov submit-proposal \
--title=<Parameter change proposal> \
--description=<parameter change of min deposit> \
--type="param-change" \
--deposit="80000cheq" \
--from=<alex> \
--chain-id=<node 45.77.218.219:26657>
```

This is the complete command that I could use to submit a mainnet parameter-change proposal right now: 

```text
VDR CLI tx gov submit-proposal \
--title=<Parameter change proposal> \
--description=<parameter change of min deposit> \
--type="param-change" \
--deposit="80000cheq" \
--from=<alex> \
--chain-id=<cheqdnetwork--node cheqd-node-1.evernym.network:26657>
```

1. VDR CLI is the command-line interface client that is used to send transactions and query the Cosmos Hub
2. tx gov submit-proposal param-change indicates that the transaction is submitting a parameter-change proposal
3. --from alex is the account key that pays the transaction fee and deposit amount
4. --gas 500000 is the maximum amount of gas permitted to be used to process the transaction
   * the more content there is in the description of your proposal, the more gas your transaction will consume
   * if this number isn't high enough and there isn't enough gas to process your transaction, the transaction will fail
   * the transaction will only use the amount of gas needed to process the transaction
5. --fees is a flat-rate incentive for a Validator to process your transaction
   * the network still accepts zero fees, but many nodes will not transmit your transaction to the network without a minimum fee
   * many nodes \(including the Figment node\) use a minimum fee to disincentivize transaction spamming
   * 7500uCHEQ is equal to 0.0075 CHEQ
6. --chain-id cheqdnetwork is cheqd’s mainnet. For current and past chain-id's, please look at the cheqd/mainnetresource
   * the testnet chain ID is \[insert chain ID\]
7. --node cheqd-node-1.evernym.network:26657 is using Evernym Networks' node to send the transaction to the cheqd mainnet.

Note: be careful what you use for **--fees**. A mistake here could result in spending hundreds or thousands of CHEQs accidentally, which cannot be recovered.

####  

#### Verifying your transaction

After posting your transaction, your command line interface will provide you with the transaction's hash, which you can either query using _**gaiad**_ or by searching the hash using [Hubble](https://hubble.figment.network/cosmos/chains/cosmoshub-3/transactions/B8E2662DE82413F03919712B18F7B23AF00B50DAEB499DAD8C436514640EFC79). The hash should look something like this: **B8E2662DE82413F03919712B18F7B23AF00B50DAEB499DAD8C436514640EFC79**

####  

#### Troubleshooting a failed transaction

There are a number of reasons why a transaction may fail. Here are two examples:

1. Running out of gas - The more data there is in a transaction, the more gas it will need to be processed. If you don't specify enough gas, the transaction will fail.
2. Incorrect denomination - You may have specified an amount in 'microCHEQ' or 'CHEQ' instead of 'nanoCHEQ', causing the transaction to fail.

If you encounter a problem, try to troubleshoot it first, and then ask for help on the cheqd Governance forum. We can learn from failed attempts and use them to improve upon this document.

## How do I Vote?

#### Participants

Participants are users that have the right to vote on proposals. In the cheqd Network, participants are bonded CHEQ holders. Bonding means something different for Validators and for everyday Users:

1. Validators can ‘self-bond’ their staking tokens in order to vote on governance matters;
2. Everyday Users can ‘bond’ their tokens to a Validator, this is known as delegation. 

Unbonded CHEQ holders and other users do not get the right to participate in voting on proposals. However, they can submit and deposit on proposals.

Note that some participants can be forbidden to vote on a proposal under a certain Validator if:

* participant has bonded or unbonded CHEQ to a particular Validator after the proposal has entered its voting period.
* participant set up a node and became a Validator after the proposal entered its voting period.

This does not prevent the participant voting with CHEQ bonded to other Validators. For example, if a participant bonded some CHEQ to Validator A before a proposal entered voting period and other CHEQ to Validator B after proposal entered voting period, only the vote under Validator B will be forbidden.  


#### Inheritance

If a User does not vote, it will inherit the Validator’s vote which it is bonded to.

If the User votes before its Validator, it’s vote will take precedence; the Validator will not inherit the User’s vote. 

If the User votes after its Validator, it will override its Validator vote with its own. If the proposal is urgent, it is possible that the vote will close before User has a chance to react and override their Validator's vote.  


#### Voting period

Once a proposal reaches _**MinDeposit**_, it immediately enters Voting period. We define Voting period as the interval between the moment the vote opens and the moment the vote closes. Voting period should always be shorter than the Unbonding period to prevent double voting. 

The initial value of the cheqd Voting period is 2 weeks.  


#### Option set

The option set of a proposal refers to the set of choices a participant can choose from when casting its vote.

The initial option set includes the following options:

```text
Yes
No
NoWithVeto
Abstain
```

_**NoWithVeto**_ counts as _**No**_ but also adds a Veto vote. Abstain option allows voters to signal that they do not intend to vote in favor or against the proposal but accept the result of the vote.  


#### Weighted Votes

Users casting a vote on a proposal have the option to split their votes into several voting options. For example, a User could use 70% of its voting power to vote Yes and 30% of its voting power to vote No.

Often, the entity owning a particular governance address might not be a single individual. For example, a company might have different stakeholders who want to vote differently, and so it makes sense to allow them to split their voting power. Currently, it is not possible for them to do "passthrough voting" and give their users voting rights over their tokens. However, with this system, exchanges can poll their users for voting preferences off-chain, and then vote on-chain proportionally to the results of the poll.

For a weighted vote to be valid, the options field must not contain duplicate vote options, and the sum of weights of all options must be equal to 1.  


#### Quorum

Quorum is defined as the minimum percentage of voting power that needs to be cast on a proposal for the result to be valid. 

The quorum as a default setting is **33.34%**.

Going forward, more complex quorum mechanisms, such as Adaptive Quorum Biasing should be considered.   


#### Threshold

Threshold is defined as the minimum proportion of Yes votes \(excluding Abstain votes\) for the proposal to be accepted.

Initially, the threshold is set at **50%** with a possibility to veto if more than **33.34% of votes** \(excluding Abstain votes\) are _**NoWithVeto**_ votes. This means that proposals are accepted if the proportion of Yes votes \(excluding Abstain votes\) at the end of the voting period is superior to **50%** and if the proportion of _**NoWithVeto**_ votes is inferior to **33.34%** \(excluding Abstain votes\).  


#### Validator’s punishment for non-voting

At present, Validators are not punished for failing to vote.  


#### Governance address

For the MVP, the Governance address will be the main Validator address generated at account creation. This address corresponds to a different PrivKey than the Tendermint PrivKey which is responsible for signing consensus messages. Validators thus do not have to sign governance transactions with the sensitive Tendermint PrivKey.  


### Software Upgrade

If proposals are of type _**SoftwareUpgradeProposal**_, then nodes need to upgrade their software to the new version that was voted. This process is divided into two steps:

```text
Signal
Switch
```

#### Signal

After a _**SoftwareUpgradeProposal**_ is accepted, Validators are expected to download and install the new version of the software while continuing to run the previous version. Once a Validator has downloaded and installed the upgrade, it will start signaling to the network that it is ready to switch by including the proposal's proposalID in its precommits.\(Note: Confirmation that we want it in the precommit?\)

Note: There is only one signal slot per precommit. If several _**SoftwareUpgradeProposals**_ are accepted in a short timeframe, a pipeline will form and they will be implemented one after the other in the order that they were accepted.  


#### Switch

Once a block contains more than 2/3rd precommits where a common _**SoftwareUpgradeProposal**_ is signaled, all the nodes \(including Validator nodes, non-validating full nodes and light-nodes\) are expected to switch to the new version of the software.  


