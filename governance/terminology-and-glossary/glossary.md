# Glossary

<table>
  <thead>
    <tr>
      <th style="text-align:left"><b>Term</b>
      </th>
      <th style="text-align:left"><b>Definition</b>
      </th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td style="text-align:left"><b>Agency</b>
      </td>
      <td style="text-align:left">A service provider that hosts Cloud Agents and may provision Edge Agents
        on behalf of Entities. Agencies may be Unaccredited, Self-Certified, or
        Accredited.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Agent</b>
      </td>
      <td style="text-align:left">A software program or process used by or acting on behalf of an Entity
        to interact with other Agents or with the cheqd Network or other distributed
        ledgers. Agents are of two types: Edge Agents run at the edge of the network
        on a local device; Cloud Agents run remotely on a server or cloud hosting
        service. Agents require access to a Wallet in order to perform cryptographic
        operations on behalf of the Entity they represent.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Agent-to-Agent (A2A) Protocol</b>
      </td>
      <td style="text-align:left">A protocol for communicating between Agents to form Connections, exchange
        Credentials, and have other secure private Interactions. A less technical
        synonym is DID Communication (DIDComm)<b>.</b>
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Agent-to-Agent (A2A) Protocol Layer</b>
      </td>
      <td style="text-align:left">The technical stack &apos;Layer&apos; for peer-to-peer Connections and
        Interactions over the Agent-to-Agent Protocol.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Anonym</b>
      </td>
      <td style="text-align:left">A DID used exactly once, so it cannot be contextualized or correlated
        beyond that single usage. See also Pseudonym and Verinym.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Anywise</b>
      </td>
      <td style="text-align:left">A non-reciprocal relationship rooted in the Identity of one party, where
        the other party is the public (a faceless &#x201C;other&#x201D; that can
        be instantiated without bound). For an Organization to issue publicly Verifiable
        Credentials, its Issuer DID must be on a public ledger such as cheqd. It
        is thus an Anywise DID&#x2014;a DID to which any other Entity may refer
        without coordination. The term &#x201C;public DID&#x201D; is sometimes
        used as a casual synonym for &#x201C;Anywise DID&#x201D;. However, &#x201C;public
        DID&#x201D; is deprecated because it is ambiguous, i.e., it may refer to
        a DID that is world-visible but usable only in pairwise mode, or to a DID
        that is not published in a central location but nonetheless used in many
        contexts, or to a DID that is both publicly visible and used in Anywise
        mode. Compare N-wise and Pairwise.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>ATOM</b>
      </td>
      <td style="text-align:left">The exchange token that runs and transacts natively on the Cosmos ledger.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Attribute</b>
      </td>
      <td style="text-align:left">An Identity trait, property, or quality of an Entity.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Authentic data</b>
      </td>
      <td style="text-align:left">A relatively new term, intended to replace &apos;Self-sovereign identity&apos;.
        It relates to signed, verifiable and cryptographically resolvable data
        using the <a href="https://www.w3.org/TR/vc-data-model/">W3C Verifiable Credential data model</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Bitcoin</b>
      </td>
      <td style="text-align:left">A type of cryptocurrency, famously created by pseudonymous Satoshi Nakamoto,
        using a Proof-of-Work (PoW) consensus model to mine blocks.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Blocks</b>
      </td>
      <td style="text-align:left">A set of transaction data, forming part of a Blockchain.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Blockchain</b>
      </td>
      <td style="text-align:left">A system in which a record of transactions are maintained across several
        computers that are linked in a distributed, peer-to-peer network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Bonding</b>
      </td>
      <td style="text-align:left">Attaching tokens to a specific Node Operator, to participate in cheqd
        Network Governance and delegate unused votes to that specific Node Operator.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>CHEQ</b>
      </td>
      <td style="text-align:left">The native medium of exchange, governance and transaction fee on the cheqd
        Network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>cheqd Network</b>
      </td>
      <td style="text-align:left">The blockchain, built on the Cosmos SDK, that cheqd uses for transactions,
        governance and identity interactions.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Claim</b>
      </td>
      <td style="text-align:left">An assertion about an Attribute of a Subject. Examples of a Claim include
        date of birth, height, government ID number, or postal address&#x2014;all
        of which are possible Attributes of an Individual. A Credential is comprised
        of a set of Claims.<em> (Note: Early in the development of Self-Sovereign Identity technology, this term was used the same way it was used in the early W3C Verifiable Claims Working Group specifications&#x2014;as a synonym for what is now a Credential. That usage is now deprecated.)</em>
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Cloud Agent</b>
      </td>
      <td style="text-align:left">An Agent that is hosted in the cloud. It typically operates on a computing
        device over which the Identity Owner does not have direct physical control
        or access. Mutually exclusive with Edge Agent. A Cloud Agent requires a
        Wallet and typically has a Service Endpoint. Cloud agents may be hosted
        by an Agency.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Coin</b>
      </td>
      <td style="text-align:left">A coin operates on its own independent blockchain and acts like a native
        currency within a specific financial system.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Connection</b>
      </td>
      <td style="text-align:left">A cryptographically verifiable communications channel established using
        an Agent-to-Agent Protocol between two DIDs representing two Entities and
        their associated Agents. Connections may be Edge-to-Edge Connections or
        Cloud-to-Cloud Connections. Connections may be used to exchange Verifiable
        Credentials or for any other communications purpose. Connections may be
        encrypted and decrypted using the Public Keys and Private Keys for each
        DID. A Connection may be temporary or it may last as long as the two Entities
        desire to keep it. Two Entities may have multiple Connections between them,
        however each Connection must be between a unique pair of DIDs. A relationship
        between more than two Entities may be modeled either as Pairwise connections
        between all of the Entities (Peering) or each Entity can form a Connection
        with an Entity representing a Group.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Connection Invitation</b>
      </td>
      <td style="text-align:left">An Agent-to-Agent Protocol message type sent from one Entity to a second
        Entity to invite the second Entity to send a Connection Request.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Controlled Document</b>
      </td>
      <td style="text-align:left">A subdocument of a Governance Framework as a normative component of the
        framework. These are often referred to in the Trust over IP metamodel,
        which cheqd intends to comply with.
        <br />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Controller</b>
      </td>
      <td style="text-align:left">An Entity that has the Private Keys and responsibility to take actions
        on behalf of another Entity.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Cosmos</b>
      </td>
      <td style="text-align:left">The distributed ledger, with the coin ATOM, which cheqd is building its
        infrastructure on top of.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Cosmos SDK</b>
      </td>
      <td style="text-align:left">The development kit that cheqd is using to build its infrastructure.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Core Principles</b>
      </td>
      <td style="text-align:left">The Principles published in this Governance Framework that seek to govern
        the behaviour of participants in the cheqd Network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential</b>
      </td>
      <td style="text-align:left">A digital assertion containing a set of Claims made by an Entity about
        itself or another Entity. Credentials are a subset of Identity Data. A
        Credential is based on a Credential Definition. The Entity described by
        the Claims is called the Subject of the Credential. The Entity creating
        the Credential is called the Issuer. The Entity holding the issued Credential
        is called the Holder. The Entity to whom a Credential is presented is generally
        called the Relying Party, and specifically called the Verifier if the Credential
        is a Verifiable Credential. Once issued, a Credential is typically stored
        by an Agent. (In cheqd&apos;s infrastructure, Credentials are never stored
        on the Ledger.)
        <br />
        <br />Credentials are very broad in their potential use: Examples of Credentials
        include college transcripts, driver licenses, health insurance cards, and
        building permits. See also Verifiable Credential.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Definition (CredDef)</b>
      </td>
      <td style="text-align:left">A machine-readable definition of the semantic structure of a Credential
        based on one or more Schemas. Credential Definitions are stored on the
        cheqd Network. Credential Definitions must include an Issuer Public Key.
        Credential Definitions facilitate interoperability of Credentials and Proofs
        across multiple Issuers, Holders, and Verifiers.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Exchange</b>
      </td>
      <td style="text-align:left">A set of Interaction Patterns within an Agent-to-Agent Protocol for exchange
        of Credentials between Entities acting in Credential Exchange Roles.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Exchange Layer</b>
      </td>
      <td style="text-align:left">The technical infrastructure Layer for Credential Exchange.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Offer</b>
      </td>
      <td style="text-align:left">An Agent-to-Agent Protocol message type sent from an Issuer to a Holder
        to invite the Holder to send a Credential Request to the Issuer.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Registry</b>
      </td>
      <td style="text-align:left">An Entity that serves as a Holder of Credentials issued by Trust Community
        Members in order to provide a cryptographically verifiable directory service
        to the Trust Community or to the public. The term also refers to the actual
        repository of Credentials maintained by this Entity. An informal Credential
        Registry may accept Credentials from participants whose purpose is to cross-certify
        each other&#x2019;s roles in the Trust Community. A formal Credential Registry
        may be authorized directly by a Governance Authority or Accredited by an
        authorized Auditor for the relevant Governance Framework.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Credential Request</b>
      </td>
      <td style="text-align:left">An Agent-to-Agent Protocol message type sent from a Holder to an Issuer
        to request the issuance of a Credential to that Holder.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Cryptocurrency</b>
      </td>
      <td style="text-align:left">A digital currency in which transactions are verified and records maintained
        by a decentralized system using cryptography, rather than by a centralized
        authority.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Cryptographic Trust</b>
      </td>
      <td style="text-align:left">Trust bestowed in a set of machines (Man-Made Things) that are operating
        a set of cryptographic algorithms will behave as expected. This form of
        trust is based in mathematics and computer hardware/software engineering.
        Compare with Human Trust.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Data Controller</b>
      </td>
      <td style="text-align:left">As defined by the <a href="https://en.wikipedia.org/wiki/General_Data_Protection_Regulation">EU General Data Protection Regulation</a> (GDPR),
        the natural or legal person, public authority, agency, or other body which,
        alone or jointly with others, determines the purposes and means of the
        processing of Personal Data.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Data Processor</b>
      </td>
      <td style="text-align:left">As defined by the <a href="https://en.wikipedia.org/wiki/General_Data_Protection_Regulation">EU General Data Protection Regulation</a> (GDPR),
        a natural or legal person, public authority, agency, or other body which
        processes Personal Data on behalf of a Data Controller.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Data Protection by Design</b>
      </td>
      <td style="text-align:left">A widely recognized <a href="https://ico.org.uk/for-organisations/guide-to-the-general-data-protection-regulation-gdpr/principles/">set of principles</a> for
        protecting Personal Data. Specific cheqd Data Protection by Design principles
        are a subset of the General Principles in the cheqd Governance Framework.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Data Subject</b>
      </td>
      <td style="text-align:left">As defined by the <a href="https://en.wikipedia.org/wiki/General_Data_Protection_Regulation">EU General Data Protection Regulation</a> (GDPR),
        any person whose Personal Data is being collected, held, or processed.
        In the cheqd Governance Framework, a Data Subject is referred to as an
        Individual.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Decentralised identity</b>
      </td>
      <td style="text-align:left">Synonymous with Self-Sovereign Identity, decentralised identity refers
        to the control and management of identity Credentials, Claims and Attributes
        by an Entity which the data contained in the Credentials, Claims and Attributes
        is about.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Decentralised Identifier (DID)</b>
      </td>
      <td style="text-align:left">A globally unique identifier developed specifically for decentralized
        systems as defined by the <a href="https://w3c-ccg.github.io/did-spec/">W3C DID specification</a>.
        DIDs enable interoperable decentralized Self-Sovereign Identity management.
        A DID is associated with exactly one DID Document.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Delegate</b>
      </td>
      <td style="text-align:left">
        <p>This term has two meanings in different contexts.</p>
        <p></p>
        <p>Firstly, it can mean an Identity Controller that acts on behalf of another
          Identity Controller to assist or manage Credentials, Claims or Attributes
          on behalf of that secondary Identity Controller.</p>
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Delegate</b>
      </td>
      <td style="text-align:left">Secondly, it can mean delegating tokens for the purpose of participating
        in cheqd on-chain Governance. Delegating tokens means bonding tokens to
        a specific Node Operator.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID</b>
      </td>
      <td style="text-align:left">Acronym for Decentralized Identifier.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID Communication</b>
      </td>
      <td style="text-align:left">Synonym for Agent-to-Agent Protocol.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID Document</b>
      </td>
      <td style="text-align:left">The machine-readable document to which a DID points as defined by the
        <a
        href="https://w3c-ccg.github.io/did-spec/">W3C DID specification</a>. A DID document describes the Public Keys, Service
          Endpoints, and other metadata associated with a DID. A DID Document is
          associated with exactly one DID.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID Method</b>
      </td>
      <td style="text-align:left">A specification that defines a particular type of DID conforming to the
        <a
        href="https://w3c-ccg.github.io/did-spec/">W3C DID specification</a>. A DID Method specifies both the format of the
          particular type of DID as well as the set of operations for creating, reading,
          updating, and deleting (revoking) it.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID Resolver</b>
      </td>
      <td style="text-align:left">A software module that takes a DID as input and returns a DID document
        by invoking the DID Method used by that particular DID. Analogous to the
        function of a <a href="https://en.wikipedia.org/wiki/Domain_Name_System#DNS_resolvers">DNS resolver</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DID Subject</b>
      </td>
      <td style="text-align:left">The Entity identified by a DID.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DKMS</b>
      </td>
      <td style="text-align:left"><a href="http://bit.ly/dkmsv3">Decentralized Key Management System</a>,
        an emerging standard for interoperable cryptographic key management based
        on DIDs.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>DKMS Protocol</b>
      </td>
      <td style="text-align:left">A subset of the Agent-to-Agent Protocol that enables Agents to perform
        DKMS functions for interoperable digital Wallet management, e.g., key exchange,
        automated backup, offline recovery, social recovery, etc.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Edge Agent</b>
      </td>
      <td style="text-align:left">An Agent that operates at the edge of the network on a local device, such
        as a smartphone, tablet, laptop, automotive computer, etc. The device owner
        usually has local access to the device and can exert control over its use
        and authorization. Mutually exclusive with Cloud Agent. An Edge Agent may
        be an app used directly by an Identity Owner, or it may be an operating
        system module or background process called by other apps. Edge Agents typically
        do not have a publicly exposed Service Endpoint in a DID Document, but
        do have access to a Wallet. Note that the local device may itself be an
        Active Thing with its own Agent, and for which the Identity Owner is the
        Thing Controller.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Edge-to-Edge Connection</b>
      </td>
      <td style="text-align:left">A Connection that forms and/or communicates directly between two Edge
        Agents.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Entity</b>
      </td>
      <td style="text-align:left">As used in <a href="https://tools.ietf.org/html/rfc3986">IETF RFC 3986, Uniform Resource Identifier (URI)</a>,
        a resource of any kind that can be uniquely and independently identified.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Entropy</b>
      </td>
      <td style="text-align:left">The determinate and irreversible transition of the cheqd Network&#x2019;s
        application of Governance from centralisation to decentralisation.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Fee</b>
      </td>
      <td style="text-align:left">A proportion of Network transaction costs that is taken used to remunerate
        Network participants or the Community Pool.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Genesis Parameters</b>
      </td>
      <td style="text-align:left">The initial Network parameters governing how cheqd works at an architectural
        level.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Governance Authority (GA)</b>
      </td>
      <td style="text-align:left">The Entity (typically an Organization) governing and making decisions
        related to a particular Governance Framework. cheqd does not have a Governance
        Authority in its traditional sense, its governance is conducted by the
        distributed consensus of the Network itself.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Guardian</b>
      </td>
      <td style="text-align:left">An Identity Controller who administers Identity Data, Wallets, and/or
        Agents on behalf of a Dependent. A Guardian is different than a Delegate&#x2014;in
        Delegation, the Identity Controller still retains control of one or more
        Wallets. With Guardianship, an Identity Controller is wholly dependent
        on the Guardian to manage the Identity Controllers&apos; Wallet.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Holder</b>
      </td>
      <td style="text-align:left">
        <p>A role played by an Entity when it is issued a Credential by an Issuer.
          The Holder may or may not be the Subject of the Credential. (There are
          many use cases in which the Holder is not the Subject, e.g., a birth certificate
          where the Subject is a baby and both the mother and father may be Holders.)</p>
        <p></p>
        <p>Holders are also those who own and hold CHEQ.</p>
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Hyperledger</b>
      </td>
      <td style="text-align:left">An initiative of the <a href="https://www.linuxfoundation.org/">Linux Foundation</a> to
        develop open source distributed ledger and blockchain technology. The Hyperledger
        home page is <a href="https://wiki.hyperledger.org/">https://wiki.hyperledger.org/</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Hyperledger Aries</b>
      </td>
      <td style="text-align:left">Provides a shared, reusable, interoperable tool kit designed for initiatives
        and solutions focused on creating, transmitting and storing verifiable
        digital credentials. It is infrastructure for blockchain-rooted, peer-to-peer
        interactions.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Hyperledger Indy</b>
      </td>
      <td style="text-align:left">
        <p>An open source project under the Hyperledger umbrella for decentralized
          Self-Sovereign Identity. The source code for Hyperledger Indy was originally
          contributed to the Linux Foundation by the Sovrin Foundation.</p>
        <p></p>
        <p>cheqd does not use Hyperledger Indy for its Network, instead it uses Cosmos.</p>
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Hyperledger Ursa</b>
      </td>
      <td style="text-align:left">A shared cryptographic library that would enable people (and projects)
        to avoid duplicating other cryptographic work and hopefully increase security
        in the process.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Identifier</b>
      </td>
      <td style="text-align:left">A text string or other atomic data structure used to provide a base level
        of Identity for an Entity in a specific context. In Self-Sovereign Identity
        systems, Decentralized Identifiers (DIDs) are the standard Identifier.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Identity Data</b>
      </td>
      <td style="text-align:left">The set of data associated with an Identity that permits identification
        of the underlying Entity.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Identity Controller</b>
      </td>
      <td style="text-align:left">The person, organisation, group or thing that retains control over the
        Private Key(s) relating to specific identity data.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Interaction</b>
      </td>
      <td style="text-align:left">A set of messages exchanged over a Connection using an Agent-to-Agent
        Protocol.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Issuer</b>
      </td>
      <td style="text-align:left">The Entity that issues a Credential to a Holder. Based on the definition
        provided by the <a href="https://www.w3.org/2017/vc/">W3C Verifiable Claims Working Group</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>JSON</b>
      </td>
      <td style="text-align:left">Open standard data format used for some Verifiable Credentials</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>JSON-LD</b>
      </td>
      <td style="text-align:left">Open standard data format used for some Verifiable Credentials, specifically
        for linking data to other datasets.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Key Recover</b>y</td>
      <td style="text-align:left">The process of recovering access to and control of a set of Private Keys&#x2014;or
        an entire Wallet&#x2014;after loss or compromise. Key Recovery is a major
        focus of the emerging DKMS standard for cryptographic key management.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Layer 1</b>
      </td>
      <td style="text-align:left">The core, foundational infrastructure that an SSI ecosystem is built upon.
        In cheqd&apos;s case, the cheqd Network is Layer 1.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Layer 2</b>
      </td>
      <td style="text-align:left">Within blockchains a Layer 2 is a separate ledger running adjacent to
        the Layer 1. Layer 2 Networks are used for efficiency and scaling. The
        outcome of Layer 2 transactions and interactions are recorded periodically
        back into Layer 1.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Level of Assurance (LOA)</b>
      </td>
      <td style="text-align:left">A measure, usually numeric, of the Trust Assurance that one Entity has
        in another Entity based on a defined set of criteria that establish the
        amount of reliance the first Entity may accept from the second Entity in
        the performance of the criteria. LOAs are often defined in or referenced
        by Governance Frameworks.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Liquidity pool</b>
      </td>
      <td style="text-align:left">Pools of tokens locked in smart contracts that provide liquidity in decentralized
        exchanges in an attempt to attenuate the problems caused by the illiquidity
        typical of such systems.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Main net</b>
      </td>
      <td style="text-align:left">cheqd&apos;s Network has both a test net and a main net. The main net
        is the Network where live and public transactions will take place after
        launch.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Minumum deposit</b>
      </td>
      <td style="text-align:left">The minimum amount of tokens needed for a governance proposal to reach
        the stage where it is voted upon.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Mint</b>
      </td>
      <td style="text-align:left">Minting a token simply means creating a token on-ledger</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Node</b>
      </td>
      <td style="text-align:left">A computer network server running an instance of the code necessary to
        operate a distributed ledger or blockchain. In cheqd Infrastructure, a
        Node is operated by a Node Operator running an instance of the cheqd Open
        Source Code.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Node Operator</b>
      </td>
      <td style="text-align:left">The Entity responsible for running a node.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Open Governance</b>
      </td>
      <td style="text-align:left">A governance model in which the Governance Authority is open to public
        participation, operates with full transparency, and does not favor any
        particular contributor or constituency.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Open Source License</b>
      </td>
      <td style="text-align:left">Any form of intellectual property license approved and published by the
        <a
        href="https://opensource.org/">Open Source Initiative</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Open Standards</b>
      </td>
      <td style="text-align:left">Technical standards that are developed under an Open Governance process;
        are publicly available for anyone to use; and which do not lock in users
        of the standard to a specific vendor or implementation. Open Standards
        facilitate interoperability and data exchange among different products
        or services and are intended for widespread adoption. Many Open Standards
        have implementations that are available under an Open Source License.
        <br
        />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Organization</b>
      </td>
      <td style="text-align:left">A legal Entity that is not a natural person (i.e., not an Individual).
        Examples of Organizations include a Group, sole proprietorship, partnership,
        corporation, LLC, association, NGO, cooperative, government, etc. Mutually
        exclusive with Individual.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Overlay</b>
      </td>
      <td style="text-align:left">A data structure that provides an extra layer of contextual and/or conditional
        information to a Schema. This extra context can be used by an Agent to
        transform how information is displayed to a viewer or to guide the Agent
        in how to apply a custom process to Schema data.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Pairwise</b>
      </td>
      <td style="text-align:left">A direct relationship between exactly two Entities. Most relationships
        in the cheqd ecosystem will be likely Pairwise, even when one or both Entities
        are not Individuals. For example, business-to-business relationships are
        pairwise by default. A DID or a Public Key or a Service Endpoint is Pairwise
        if it is used exclusively in a Pairwise relationship. Pairwise relationships
        can exist entirely off-ledger.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Payment</b>
      </td>
      <td style="text-align:left">A transfer of CHEQ or other cryptographically verifiable units of value
        from one Entity to another Entity.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Payment Address</b>
      </td>
      <td style="text-align:left">The address of a Payment Transaction on the cheqd Network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Permissionless</b>
      </td>
      <td style="text-align:left">Permissionless blockchains are blockchains that require no permission
        to join and interact with.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Personal Data</b>
      </td>
      <td style="text-align:left">As defined by the <a href="https://en.wikipedia.org/wiki/General_Data_Protection_Regulation">EU General Data Protection Regulation</a> (GDPR),
        any information relating to an identified or identifiable natural person.
        In the GDPR, this natural person is called the Data Subject. Personal data
        SHOULD never be written to the cheqd Network.
        <br />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Privacy by Design</b>
      </td>
      <td style="text-align:left">A set of <a href="https://en.wikipedia.org/wiki/Privacy_by_design">seven foundational principles</a> for
        taking privacy into account throughout the entire design and engineering
        of a system, product, or service. Originally defined by the <a href="https://www.ipc.on.ca/">Information and Privacy Commissioner of Ontario, Canada</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Private Key</b>
      </td>
      <td style="text-align:left">The half of a cryptographic key pair designed to be kept as the Private
        Data of an Entity. In elliptic curve cryptography, a Private Key is called
        a signing key.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Proof</b>
      </td>
      <td style="text-align:left">Cryptographic verification of a Claim or a Credential. A <a href="https://en.wikipedia.org/wiki/Digital_signature">digital signature</a> is
        a simple form of Proof. A <a href="https://en.wikipedia.org/wiki/Cryptographic_hash_function">cryptographic hash</a> is
        also a form of Proof. Zero Knowledge Proofs enable <a href="https://en.wikipedia.org/wiki/Selective_disclosure">selective disclosure</a> of
        the information in a Credential.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Proof request</b>
      </td>
      <td style="text-align:left">The data structure sent by a Verifier to a Holder that describes the Proof
        required by the Verifier.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Pseudonym</b>
      </td>
      <td style="text-align:left">A DID used to prevent correlation outside of a specific context. A Pseudonym
        may be Pairwise, N-wise, or Anywise.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Public Key</b>
      </td>
      <td style="text-align:left">The half of a cryptographic key pair designed to be shared with other
        parties in order to decrypt or verify encrypted communications from an
        Entity. In digital signature schemes, a Public Key is also called a verification
        key. A Public Key may be either Public Data or Private Data depending on
        the policies of the Entity.
        <br />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Quorum</b>
      </td>
      <td style="text-align:left">The minimum number of participants in the Network who need to vote on
        a governance proposal for the vote to be valid.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Recovery Key</b>
      </td>
      <td style="text-align:left">A special Private Key used for purposes of recovering a Wallet after loss
        or compromise. In the DKMS key management protocol, a Recovery Key may
        be cryptographically sharded for <a href="https://en.wikipedia.org/wiki/Secret_sharing">secret sharing</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Recovery Key Trustee</b>
      </td>
      <td style="text-align:left">A Trustee trusted by another Identity Controller to authorise sharing
        back a Recovery Key for purposes of restoring a Wallet after loss or compromise.
        <br
        />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Relying Party</b>
      </td>
      <td style="text-align:left">An Entity that consumes Identity Data and accepts some Level of Assurance
        from another Entity for some purpose. Verifiers are one type of Relying
        Party.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Resolver</b>
      </td>
      <td style="text-align:left">A software module that accepts an Identifier as input, looks up the Identifier
        in a database or ledger, and returns metadata describing the identified
        Entity. The Domain Name System (DNS) uses a <a href="https://en.wikipedia.org/wiki/Domain_Name_System#DNS_resolvers">DNS resolver</a>.
        Self-Sovereign Identity uses a DID Resolver.
        <br />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Revocation</b>
      </td>
      <td style="text-align:left">The act of an Issuer revoking the validity of a Claim or a Credential.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Revocation Registry</b>
      </td>
      <td style="text-align:left">An online repository of data needed for Revocation. In cheqd&apos;s Network
        Infrastructure, a Revocation Registry is a privacy-respecting cryptographic
        data structure maintained on the Ledger by an Issuer in order to support
        Revocation of a Credential.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Schema</b>
      </td>
      <td style="text-align:left">A machine-readable definition of the semantics of a data structure. Schemas
        are used to define the Attributes used in one or more Credential Definitions.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Security by Design</b>
      </td>
      <td style="text-align:left">A widely recognized <a href="https://en.wikipedia.org/wiki/Secure_by_design">set of principles</a> for
        building security into systems, products, and services from the very start.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Selective Disclosure</b>
      </td>
      <td style="text-align:left">A Privacy by Design principle of revealing only the subset of the data
        described in a Claim, Credential, or other set of Private Data that is
        required by a Verifier.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Self-Sovereign Identity (SSI)</b>
      </td>
      <td style="text-align:left">An identity system architecture based on the core principle that individual
        Identity Controllers have the right to permanently control one or more
        Identifiers together with the usage of the associated Identity Data.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Service Endpoint</b>
      </td>
      <td style="text-align:left">An addressable network location offering a service operated on behalf
        of an Entity. As defined in the <a href="https://w3c-ccg.github.io/did-spec/">DID specification</a>,
        a Service Endpoint is expressed as a <a href="https://tools.ietf.org/html/rfc3986">URI (Uniform Resource Identifier)</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Slashing</b>
      </td>
      <td style="text-align:left">The potential deduction of tokens for bad behaviour on the Network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>SSI Stack</b>
      </td>
      <td style="text-align:left">A general representation of different technological functions and protocols
        that provide different functions on top of each other. The most common
        SSI stack is the <a href="https://trustoverip.org/toip-model/">Trust over IP stack</a>.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Staking</b>
      </td>
      <td style="text-align:left">In order to participate in cheqd Network governance, participants must
        allocate a proportion of their tokens to float on the Network. Staking
        is required to validate transactions and earn Staking rewards.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Sufficient decentralisation</b>
      </td>
      <td style="text-align:left">The point at which accountability over the cheqd Network is spread amongst
        enough nodes to render each node legally unaccountable for the decisions
        made on the Network. This is reached through the increasing of Entropy.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Token</b>
      </td>
      <td style="text-align:left">A medium of interaction, either for pure utility or for payment transactions.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Tokenomics</b>
      </td>
      <td style="text-align:left">The fundamental code and parameters that determines how the Network architecture
        runs.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Tombstone</b>
      </td>
      <td style="text-align:left">A mark associated with a Transaction to suggest that the Transaction should
        no longer be returned in response to requests for read access.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Transaction</b>
      </td>
      <td style="text-align:left">A record of any type written to the cheqd Network.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Verifiable Credential</b>
      </td>
      <td style="text-align:left">A Credential that includes a Proof from the Issuer. Typically this proof
        is in the form of a digital signature.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Verifier</b>
      </td>
      <td style="text-align:left">An Entity who requests a Credential or Proof from a Holder and verifies
        it in order to make a trust decision about an Entity.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Verinym</b>
      </td>
      <td style="text-align:left">A DID that it is directly or indirectly associated with the Legal Identity
        of the Identity Controller.</td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Wallet</b>
      </td>
      <td style="text-align:left">A software module, and optionally an associated hardware module, for securely
        storing and accessing Private Keys other sensitive cryptographic key material,
        and other Private Data used by an Entity. A Wallet is accessed by an Agent.
        <br
        />
      </td>
    </tr>
    <tr>
      <td style="text-align:left"><b>Zero Knowledge Proof</b>
      </td>
      <td style="text-align:left">A Proof that uses cryptography to support Selective Disclosure of information
        about a set of Claims from a set of Credentials. A Zero Knowledge Proof
        provides cryptographic proof about some or all of the data in a set of
        Credentials without revealing the actual data or any additional information,
        including the Identity of the Holder.
        <br />
      </td>
    </tr>
  </tbody>
</table>

