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

## Finding a new home <a id="f5a9"></a>

### Defining our criteria <a id="0893"></a>

Firstly, we needed to define our criteria and priorities. What would we value when making this decision? Like any other project, this broke down into functionals, non-functionals and execution related requirements.

![cheqd&#x2019;s criteria&#x200A;&#x2014;&#x200A;how to choose the right technology](https://miro.medium.com/max/1120/1*Ssqpz_tAsIoL1nHxAwYLdA.png)

The majority of the requirements above are easy to understand without additional context. However, it is worth further explaining both “no requirement to purchase Layer 1 tokens” and “ability to have sovereign network” requirements. The first was to avoid spending capital simply buying tokens rather than developing functionality. The second was to avoid a scenario where other activities on the network impacted any identity use-case. For example, someone waiting to board a plane should not have to wait minutes because Elon Musk has tweeted about Doge and caused a stampede that congests the network.

### Narrowing down <a id="5ad7"></a>

After these clear exclusions, the choice largely boiled down to [Polkadot](https://polkadot.network/) versus [Cosmos](https://cosmos.network/) with assessments across the various requirements above. This then primarily came down to available functionality and available implementations \(identity and token\) and infrastructure \(e.g. exchanges, custodians, etc.\). As can be seen in the comparison below, both Polkadot and Cosmos have very healthy developer communities.

![cheqd&#x200A;&#x2014;&#x200A;Developer&#x2019;s community in Hyperledger Indy, Polkadot, Stellar &amp; Cosmos](https://miro.medium.com/max/1120/0*aJZ0HwW6y7vpeZlD)

However, when looking at identity and token implementations and infrastructure, we felt that Cosmos was the stronger of the two due to its longer history and adoption with many of the companies below either directly or tangentially involved in SSI. These were signposts for us that there is a valuable community to learn from and contribute back to, regardless of whether this is for a token, SSI or [decentralised autonomous organisation \(DAO\)](https://www.investopedia.com/tech/what-dao/) implementations.

### Timeline impacts <a id="f720"></a>

Like any other Network, our speed to market is key. Not only did we need the right protocol for the long-term, but it also couldn’t impact our execution timelines. We are strong believers in releasing early, releasing often and incorporating the feedback of the SSI and crypto communities as well as end-users. As a result, each option was compared against each other not just based on functional and non-functional requirements but also on delivery risk. The figure below outlines this between Cosmos and Indy. 

![Risk over time depending on the platform](https://miro.medium.com/max/1120/0*lVe_-7FySUoNOxas)

As can be seen above, there is a much higher risk with Indy in the initial stages, all due to a lack of a supporting blockchain product ecosystem to support an Indy token. Since there are no existing Indy token projects on exchanges or with custodians, they would need to build out entirely new support with no expectation of reusing existing integrations they have with other blockchains. Given how in-demand exchanges and custodians are, we felt this would impact the reach and availability of functionality we plan to build for enterprise and individual users. There is a slightly higher risk on Cosmos in the medium term as we need to convince existing projects built on Hyperledger Indy to migrate but in the long term allow us to build richer functionality over time.

### Validating our decision <a id="db76"></a>

Rather than jumping in with both feet, we took a more tactical approach by executing multiple proof-of-concepts to validate our decision. We have seen enough cases where project documentation does not reflect reality to know that these things need testing. Our development team built quick PoCs on everything from basic token implementations to decentralised identifiers \(DIDs\) on the ledger to check the protocol capabilities and the quality of the open source projects in the ecosystem. This process of experimentation and evaluation gave us the conviction to explore using the Cosmos blockchain framework seriously.

### Making migration easy <a id="df93"></a>

![World map showing some of the world&#x2019;s self-sovereign identity \(SSI\) projects](https://miro.medium.com/max/1120/0*Jnysby2gg9CAVa8Y)

To come full circle, we return to the currently deployed SSI projects. Although there are many existing networks in various stages of adoption, we firmly believe they will be turbocharged by adopting a new network that facilitates value transfer / payment for trusted data. To make this migration from Indy \(and others\) easy, we are partnering with open source projects that have modularised client-side SDKs with support for multiple DID methods / SSI networks. Our goal is to have client SDKs available \(initially for iOS and Android\) that are drop-in replacements for existing Hyperledger Indy client SDKs, as we expect Indycreds to be in circulation for a while. In addition, the same SDK would also support the new Cosmos-based DID method and Verifiable Credentials we plan on building as open source.Ultimately we don’t just want migration to be easy, but outright adoption of the technology by any app developer to be as simple as possible while acknowledging the current leading position Indycreds occupy in the field of SSI.

![cheqd&#x200A;&#x2014;&#x200A;Making migration as easy as possible through pluggable SDKs](https://miro.medium.com/max/1120/1*Zb8WDVIBEEuk1vu6cefZOQ.png)

## Looking back and stargazing <a id="018b"></a>

Our team have been through enough technology or platform selections to have some scars and to be a bit jaded, so this experience was a breath of fresh air. Historically, our experience has been that you start off with rose-tinted spectacles and slowly uncover more and more issues until you regret ever deciding to start the process. Thankfully, this has been the exact opposite.

Throughout the entire process, we have discovered a wonderful, friendly community, implementations galore, high-quality infrastructure available to us and crucially open source or native functionality that helps accelerate our roadmap for a public-permissionless decentralised identity network with significantly differentiated mechanisms to make trusted data ecosystems economically viable. Perfect examples of this are the ability to reward node operators and implement governance via voting on-chain, which are immediately available to us from Cosmos.

This brings us to now and the future. Now that [we have a testnet live with external partners](https://blog.cheqd.io/announcing-cheqds-testnet-for-a-new-incentivised-decentralised-identity-4f625ea77076), we are firmly focused on launching the mainnet and leveraging everything Cosmos and its community has available to do so successfully. We will then turn our focus to the product roadmap, looking at how much functionality we can bring to users quickly, whilst contributing back meaningfully, such as being used across [Cosmos networks via IBC](https://docs.cosmos.network/master/ibc/overview.html).

The time is ripe for making SSI truly successful, and with Cosmos, we have the right tools to do exactly that. But…even with the right tools, we can’t do this by ourselves. We invite everyone in the SSI community to be part of this journey. It would be awesome to get feedback, discuss what we got wrong, and ways to improve our network.

