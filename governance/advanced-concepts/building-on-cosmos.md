# Building on Cosmos

## Introduction to Cosmos SDK

cheqd’s Network is built using the Cosmos SDK, which means that it can inherit a lot of the core functionality of Cosmos. 

Notably, for this Governance Framework, we will often refer to specifications laid out in Cosmos’ SDK documentation, [here](https://docs.cosmos.network/v0.42/modules/gov/).

Importantly, cheqd is an [Application Specific Blockchain](https://docs.cosmos.network/master/intro/why-app-specific.html). This means that Users of the cheqd Network have full control over the entire chain and are not reliant on Cosmos at all. This ensures the community will not be stuck if, for example, a bug in the Cosmos hub or Inter-Blockchain Communication is discovered, and that it has the entire freedom to choose how it is going to evolve.

However, it is worth noting that once the Network has departed from Cosmos’ initial foundation, it will no longer be able to seamlessly update itself to Cosmos’ latest functionality. 

For this reason, it is suggested that the Network mirrors Cosmos’ architecture and updates, especially during the initial stages of its lifecycle.

## Looking beyond Indy

For those close to the SSI space, a large proportion of projects around the world are built upon [Hyperledger Indy](https://www.hyperledger.org/use/hyperledger-indy). This was a natural place for us to start, especially since there was an [existing SSI & token implementation](https://sovrin.atlassian.net/jira/software/c/projects/ST/issues/) from [Sovrin](https://sovrin.org/).

However, a few things stopped us from going down this route:

* **Limited transactions per second \(TPS\):** Indy benchmarked in our tests at ~4 TPS for token-related transactions. This is manageable for identity-only implementations since only minimal data is written to the ledger in Indy, with a bulk of the data such as Verifiable Credentials being off-ledger. This is important in the context of payments that are related to identity interactions, as low transaction speeds can result in delays for interactions that are real-time. A useful benchmark here is [Visa is capable of ~24,000 TPS](https://howmuch.net/articles/crypto-transaction-speeds-compared).
* **Limited decentralisation:** Linked to the TPS, Indy networks are limited to an upper limit of ~25 nodes beyond which the consensus mechanism breaks down. Equally, the fault tolerance of the network degrades below 7 live nodes. As a result, there is a small ‘goldilocks’ window for a viable network configuration.

As a result, Indy is fundamentally limited both in terms of throughput but also scalability.

Beyond these two limitations, the decision brought to mind a conversation between the two of us [two years ago when Hyperledger Aries and Ursa](https://www.evernym.com/blog/hyperledger-aries/) were broken out into separate projects. With standardised wallet technologies and cryptography, wouldn’t it make sense to put the identity capabilities onto another ledger? A more performant, resilient, feature-full one? One that would help us realise our roadmap quicker? And so the hunt began…

