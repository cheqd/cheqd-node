# Voting on cheqd

## Learning the basics

Before you dive into this section, we suggest that you familiarise yourself with the basic concepts surrounding Governance which can be found here.

## How do I Vote?

### Participants

Participants are users that have the right to vote on proposals. In the cheqd Network, participants are bonded CHEQ holders. Bonding means something different for Validators and for everyday Users:

1. Validators can ‘self-bond’ their staking tokens in order to vote on governance matters;
2. Everyday Users can ‘bond’ their tokens to a Validator, this is known as delegation. 

Unbonded CHEQ holders and other users do not get the right to participate in voting on proposals. However, they can submit and deposit on proposals.

Note that some participants can be forbidden to vote on a proposal under a certain Validator if:

* participant has bonded or unbonded CHEQ to a particular Validator after the proposal has entered its voting period.
* participant set up a node and became a Validator after the proposal entered its voting period.

This does not prevent the participant voting with CHEQ bonded to other Validators. For example, if a participant bonded some CHEQ to Validator A before a proposal entered voting period and other CHEQ to Validator B after proposal entered voting period, only the vote under Validator B will be forbidden.  


### Inheritance

If a User does not vote, it will inherit the Validator’s vote which it is bonded to.

If the User votes before its Validator, it’s vote will take precedence; the Validator will not inherit the User’s vote. 

If the User votes after its Validator, it will override its Validator vote with its own. If the proposal is urgent, it is possible that the vote will close before User has a chance to react and override their Validator's vote.  


### Voting period

Once a proposal reaches _**MinDeposit**_, it immediately enters Voting period. We define Voting period as the interval between the moment the vote opens and the moment the vote closes. Voting period should always be shorter than the Unbonding period to prevent double voting. 

The initial value of the cheqd Voting period is 2 weeks.  


### Option set

The option set of a proposal refers to the set of choices a participant can choose from when casting its vote.

The initial option set includes the following options:

```text
Yes
No
NoWithVeto
Abstain
```

_**NoWithVeto**_ counts as _**No**_ but also adds a Veto vote. Abstain option allows voters to signal that they do not intend to vote in favor or against the proposal but accept the result of the vote.  


### Weighted Votes

Users casting a vote on a proposal have the option to split their votes into several voting options. For example, a User could use 70% of its voting power to vote Yes and 30% of its voting power to vote No.

Often, the entity owning a particular governance address might not be a single individual. For example, a company might have different stakeholders who want to vote differently, and so it makes sense to allow them to split their voting power. Currently, it is not possible for them to do "passthrough voting" and give their users voting rights over their tokens. However, with this system, exchanges can poll their users for voting preferences off-chain, and then vote on-chain proportionally to the results of the poll.

For a weighted vote to be valid, the options field must not contain duplicate vote options, and the sum of weights of all options must be equal to 1.  


### Quorum

Quorum is defined as the minimum percentage of voting power that needs to be cast on a proposal for the result to be valid. 

The quorum as a default setting is **33.34%**.

Going forward, more complex quorum mechanisms, such as Adaptive Quorum Biasing should be considered.   


### Threshold

Threshold is defined as the minimum proportion of Yes votes \(excluding Abstain votes\) for the proposal to be accepted.

Initially, the threshold is set at **50%** with a possibility to veto if more than **33.34% of votes** \(excluding Abstain votes\) are _**NoWithVeto**_ votes. This means that proposals are accepted if the proportion of Yes votes \(excluding Abstain votes\) at the end of the voting period is superior to **50%** and if the proportion of _**NoWithVeto**_ votes is inferior to **33.34%** \(excluding Abstain votes\).  


### Validator’s punishment for non-voting

At present, Validators are not punished for failing to vote.  


### Governance address

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

