# Configure a cheqd node as validator

This document provides guidance on how configure and promote a cheqd node to validator status. Having a validator node is necessary to participate in staking rewards, block creation, and governance.

## Pre-requisites to promoting a node to validator on cheqd testnet

While the instructions listed here are specific to the cheqd testnet, a similar process is applicable to any network.

### Preparation steps

1. **Ensure you have a cheqd node installed as a service**

   You must already have a running `cheqd-node` instance installed using one of the supported methods.

   Please also ensure the node is fully caught up with the latest ledger updates.

   1. [Debian package install](debian/deb-package-install.md)
   2. [Docker install](docker-install.md)
   3. [Binary install](binary-install.md)

2. **Generate a new account key**

   Follow the guidance on [using cheqd CLI to manage keys](../cheqd-cli/cheqd-cli-key-management.md) to create a new account key.

   ```bash
   cheqd-noded keys add <alias>
   ```

   When you create a new key, a `mnemonic phrase` and `account address` will be printed. K**eep the mnemonic phrase safe** as this is the only way to restore access to the account if they keyring cannot be recovered.

3. **Get your node ID**

   Follow the guidance on [using cheqd CLI to manage nodes](../cheqd-cli/cheqd-cli-node-management.md) to fetch your node ID.

   ```bash
   cheqd-noded tendermint show-node-id
   ```

4. **Get your validator account address**

   The validator account address is generated in Step 1 above when a new key is added. To show the validator account address, follow the [cheqd CLI guide on managing accounts](../cheqd-cli/cheqd-cli-accounts.md).

   ```bash
   cheqd-noded keys list
   ```

   (The assumption above is that there is only one account / key that has been added on the node. In case you have multiple addresses, please jot down the preferred account address.)

### Requesting CHEQ tokens for cheqd mainnet
When you have a node successfully installed, please fill out our [**mainnet node operator onboarding form**](http://cheqd.link/mainnet-onboarding). You will need to have the following details on hand to fill out the form:
   1. Node ID for your node
   2. IP address / DNS record that points to the node \(if you're using an IP address, a static IP is recommended\)
   3. Peer-to-peer \(P2P\) connection port \(defaults to `26656`\)
   4. Validator account address (begins with `cheqd`)
   5. Moniker (Nickname/moniker that is set for your mainnet node)
3. Once you have received or purchased your tokens, [promote your node to a validator](docs/setup-and-configure/configure-new-validator.md).
4. If successfully configured, your node would become the latest validator on the cheqd mainnet!

### Requesting CHEQ tokens for cheqd testnet

Once you have successfully completed the steps above, please fill out our [**node operator onboarding form**](http://cheqd.link/join-testnet-form) so that you can acquire CHEQ testnet tokens required for staking on the network. The tokens will be send to your (validator) account address generated above.

You will need to have the following details on hand to fill out the form:

1. Node ID for your node
2. IP address / DNS record that points to the node \(if you're using an IP address, a static IP is recommended\)
3. Peer-to-peer \(P2P\) connection port \(default is `26656`\)
4. Validator account address (begins with `cheqd`)

If you need help or support, join our [**cheqd Community Slack**](http://cheqd.link/join-cheqd-slack) and [ask for help](https://cheqd-community.slack.com/archives/C02AQ9UK4HY).

## Promote a node to validator after acquiring CHEQ tokens for staking

1. **Ensure your account has a positive balance**

   Follow the guidance on [using cheqd CLI to manage accounts](../cheqd-cli/cheqd-cli-accounts.md) to check that your account is correctly showing the CHEQ testnet tokens provided to you.

   ```bash
   cheqd-noded query bank balances <address> --node <url>
   ```

2. **Get your node's validator public key**

   The node validator public key is required as a parameter for the next step. More details on validator public key is mentioned in the [cheqd CLI guide on managing nodes](../cheqd-cli/cheqd-cli-node-management.md).

   ```bash
   cheqd-noded tendermint show-validator
   ```

3. **Promote your node to validator status by staking your token balance**

   You can decide how many tokens you would like to stake from your account balance. For instance, you may want to leave a portion of the balance for paying transaction fees \(now and in the future\).

   To promote to validation, submit a `create-validator` transaction to the network:

   ```bash
   cheqd-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas <amount-gas> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
   ```

   Parameters required in the transaction above are:

   * **`amount`**: Amount of tokens to stake
   * **`from`**: Key alias of the node operator account that makes the initial stake
   * **`min-self-delegation`**: Minimum amount of tokens that the node operator promises to keep bonded
   * **`pubkey`**: Node's `bech32`-encoded validator public key from the previous step
   * **`commission-rate`**: Validator's commission rate
   * **`commission-max-rate`**: Validator's maximum commission rate, expressed as a number with up to two decimal points. The value for this cannot be changed later.
   * **`commission-max-change-rate`**: Maximum rate of change of a validator's commission rate per day, expressed as a number with up to two decimal points. The value for this cannot be changed later.
   * **`chain-id`**: Unique identifier for the chain. For cheqd's current testnet, this is `cheqd-testnet-2`.
   * **`gas`**: Maximum gas
   * **`gas-prices`**: Maximum gas price set by the validator

   _Example transaction:_

   ```bash
   cheqd-noded tx staking create-validator --amount 40000000000000000ncheq --from eu-node-operator --moniker node1-eu-testnet-cheqd --chain-id cheqd-testnet-2 --min-self-delegation="1" --gas="300000" --gas-prices="25ncheq" --pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"4anVUO8WhmRMqG1t4z6VxqmqZL3V7q6HqucjwZePiUw="}' --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01" --node http://node1.eu.testnet.cheqd.network:26657
   ```

4. **Check that your validator node is bonded**

   Checking that the validator is correctly bonded can be checked via any node:

   ```bash
   cheqd-noded query staking validators --node <any-rpc-url>
   ```

   Find your node by `moniker` and make sure that `status` is `BOND_STATUS_BONDED`.

5. **Check that your validator node is signing blocks and taking part in consensus**

   Find out your [validator node's hex-encoded address](../cheqd-cli/cheqd-cli-node-management.md) and look for `"ValidatorInfo":{"Address":"..."}`:

   ```bash
   cheqd-noded tendermint show-address
   ```

   Query the latest block. Open `<node-address:rpc-port/block` in a web browser. Make sure that there is a signature with your validator address in the signature list.

## Next steps

On completion of the steps above, you would have successfully bonded a node as validator to the cheqd testnet and participating in staking/consensus.

Learn more about what you can do with your new validator node in the [cheqd CLI guide](../cheqd-cli/README.md).
