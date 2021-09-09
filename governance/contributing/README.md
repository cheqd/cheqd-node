# Contributing

## How do I make changes on the Network?

Every User is encouraged to contribute to the cheqd Network. To make this easy, it is important to explain HOW a User can contribute, and what the best practices are for making changes. 

For this reason, it is important to distinguish when a change can be made entirely off-ledger, and when a change is needed to be voted on, and made on-ledger. To do this we must differentiate between **minor changes** which are able to take place **entirely off-ledger** and **major network changes** which need to be **formalised on-ledger.**

### **Minor Network changes**

These are changes that are insignificant, and do not affect the way the Network runs overall. Minor Network changes include, but are not limited to:

* Textual edits to the Governance Framework or written general documentation;
* Sensible additions to general documentation to improve clarity;
* Minor code changes;
* Finding, reporting and suggesting fixes to bugs;
* Translating the cheqd documentation into various languages.

### Major Network changes 

These are changes that have a materially significant effect on the Network. Such changes SHOULD be made via a Proposal, following the steps in the decision tree diagram below. 

Major Network changes include, but are not limited to:

* Materially significant Architecture Decisions \(**ADs**\), such as:
  * An additional feature to cheqd;
  * Removal of a feature of cheqd;
* Parameter changes for the Network;
* Community pool decisions;
* Materially significant changes to a cheqd Principle;
* Partnerships or connections to other infrastructure.

## Decision Tree

To help YOU understand how to make changes on the cheqd Network, the decision tree below visualises how changes should be carried out.

\(Placeholder\)





## Learning the basics

For general holders of coins or tokens across the industry, governance is often seen as something inaccessible and complex. This is because education about governance is often lackluster. At cheqd, we want to make participating in governance easy.

#### Frequently asked questions:

1. What is a **Validator** or **Node Operator**?  The terms Validator and Node Operator are synonymous. In blockchain ecosystems, the **Node Operator** runs what is called a **node**. A node can be thought of like a power pylon in the physical world, which helps to distribute electricity around a wide network of users. Without these pylons, electricity would be largely centralised in one location; the pylons help to distribute power to entire wide-scale populations. And if one pylon fails, the grid is set up to circumvent this pylon and re-route the electricity a different route. Similarly, in blockchain infrastructure, each node runs an instance of the consensus protocol and helps to create a broad, robust network, with no single points of failure. A node failing will have no impact on the Network as a whole.  
2. What does **staking** mean?  **Stake** is the amount of tokens a **Node Operator** puts aside and dedicates to cheqd, in order to contribute to governance and earn rewards. cheqd is a Proof of Stake Network. This means that rewards can be earnt in direct correlation with the amount of stake a Node Operator holds.  
3. What does it mean to have **bonded** tokens?  In order to participate in governance, Users need to **bond** their tokens to a Node Operator. Users with bonded tokens are known as **Participants**. This is a beneficial arrangement for both the Participant and the Node Operator. The Participant may not be interested in voting on the Network, or may be busy, and in this instance the Node Operator can cast a vote on behalf of the Participant. Bonded tokens from Users add to the **stake** of Node Operators. If the vote is accepted, the Node Operator may receive rewards in proportion to how many tokens it has staked. A percentage of these rewards can be distributed to the Users as a commission.  ****
4. Can **Participants** earn money?  


   In short, yes. Participants may be eligible for a percentage of the rewards that Node Operators earn during the course of running a node. Node Operators may offer a certain commission for bonding to them. These rewards may come in two forms:  


   1. **Transaction fees**  ****Read and writes to the cheqd Network incur what is known as a **transaction fee,** which is a calculated based on **gas.** Gas may be higher when there are high transaction volumes on the Network, and vice versa. Node Operators may also set their own gas prices, above which they are considered in the pool of who creates that transaction block. However, we will not get into the nuances of gas here.  
   2. **Block rewards**

  
      Block rewards depend on **inflation**. Inflation is the gradual increase in the number of total tokens on the Network. A Node Operator may earn block rewards during a period of inflation, which can be disseminated to the Users with bonded tokens.   
  
      For this reason, it is suggested that token holders **bond** and **delegate** their tokens, to create a healthy Network and earn passive income.   

5. What if I want to **vote unilaterally**, without a Node Operator/Validator?  


   If you are particularly interested or passionate about a specific governance proposal, or do not agree with your bonded Node Operator, it is absolutely possible to vote unilaterally. To do this, follow the instructions below in the section How do I Vote?  

6. Is my **choice of Node Operator** important?  Yes, your choice of bonded Node Operator is very important. You should be mindful about what reputation the Node Operator has. This is because the Node Operator may use your votes against the best interest of yourself and the community. Furthermore, the Node Operator might exhibit bad behaviour on the Network and have some of its stake slashed. Slashing means taking stake away from the Node Operator and adding it to the Community Pool. This may involve those who have bonded and delegated tokens to have their tokens slashed as well.  
7. Is it possible to **unbond** from a Node Operator?  Yes, is possible to instantly redelegate to a new Node Operator or **unbond** your tokens. You cannot 'hop' between Node Operators quickly however. You must complete your redelegation to a Node Operator before moving again. If you want to completely withdraw and **unbond** your tokens, you need to be mindful that this process takes **2 weeks** for the tokens to become available for use again.  

## 

## Community

You do not have to own any CHEQ to participate in our community discussion. You are free to learn about cheqd and participate in our community discussions across multiple platforms and forums.

We want learning about cheqd and participating in the community to be easy and accessible. For this reason, we have decided to make sure our [wiki](https://docs.cheqd.io/cheqd-node/) and our [source code](https://github.com/cheqd/cheqd-node) are natively interconnected. 

**What does this mean?**

If you are a technical person and want to make changes directly to our source code, on a text-based or code-based change, you can do this directly on our [Github](https://github.com/cheqd/cheqd-node).

If you are non-technical and want to make edits in a more visual way, you can do this on our [GitBook](https://docs.cheqd.io/cheqd-node/).

**Let's summarise this easily:**

1. [Github](https://github.com/cheqd): Here we will host:
   1. Our open sourced code;
   2. Documentation to help YOU onboard to the Network as a Node Operator;
   3. Our Wiki page, with all formal documentation YOU need to know about cheqd, including this Governance Framework;
   4. A forum for discussion relevant to specific topics, issues and proposals. This is a space for you to make suggestions and proposals on how you want to improve cheqd, as well as participate in discussion with the community.   You can edit our Github through branches and pull requests, explained further in our [Contributing ](https://docs.cheqd.io/cheqd-node/v/gov%2Fdraft1/governance/contributing)document.  
2. [GitBook](https://docs.cheqd.io/cheqd-node/): Here we will host:
   1. Documentation to help YOU onboard to the Network as a Node Operator;
   2. Our Wiki page, with all formal documentation YOU need to know about cheqd, including this Governance Framework;  GitBook will be a more accessible and easy-to-read layout.  You can edit our GitBook through branching our main repository and making edit requests, explained further in our [Contributing](https://docs.cheqd.io/cheqd-node/v/gov%2Fdraft1/governance/contributing) document.

If you want to stay updated with cheqd news, we recommend that you join us on:

1. [Twitter](https://twitter.com/cheqd_io): Follow us here to keep up to date with the latest cheqd news and to be the first to hear about special announcements;
2. [Telegram](https://t.me/cheqd): Join our community here to participate in general conversation about cheqd with our core community following and be the first to hear about special announcements.
3. [Discord:](https://discord.gg/SQA8NpVe2v) Discord is a fantastic way to discuss topics in a more structured way, across specific channels related to your interests. 

We also kindly ask you to familiarise yourself with our Code of Conduct which sets our clearly defined expectations and behaviour that we want to uphold in the community.  

1. [cheqd Code of Conduct](https://docs.google.com/document/d/1Rbw-0TMg8PZO85R0SuDKffXnwrDTzcPrSrURPI3el7c/edit)

This ensures that the cheqd community remains a safe and welcoming place, for any person regardless of who they are. 

