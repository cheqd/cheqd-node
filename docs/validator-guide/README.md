# Configure a cheqd node as validator

This document provides guidance on how configure and promote a cheqd node to validator status. Having a validator node is necessary to participate in staking rewards, block creation, and governance.

## Preparation steps

### Step 1: Ensure you have a cheqd node installed as a service

You must already have a running `cheqd-node` instance installed using one of the supported methods.

Please also ensure the node is fully caught up with the latest ledger updates.

1. [Debian package install](../setup-and-configure/debian/deb-package-install.md)
2. [Docker install](../setup-and-configure/docker-install.md)
3. [Binary install](../setup-and-configure/binary-install.md)

### Step 2: Generate a new account key

Follow the guidance on [using cheqd CLI to manage keys](../cheqd-cli/cheqd-cli-accounts.md) to create a new account key.

```bash
cheqd-noded keys add <alias>
```

When you create a new key, a new **account address** and **mnemonic backup phrase** will be printed. Keep the mnemonic phrase safe as this is the only way to restore access to the account if they keyring cannot be recovered.

P.S. in case of using Ledger Nano device it would be helpful to use [this instructions](#using-ledger-nano-device)

1. **Get your node ID**

   Follow the guidance on [using cheqd CLI to manage nodes](../cheqd-cli/cheqd-cli-node-management.md) to fetch your node ID.

   ```bash
   cheqd-noded tendermint show-node-id
   ```

2. **Get your validator account address**

   The validator account address is generated in Step 1 above when a new key is added. To show the validator account address, follow the [cheqd CLI guide on managing accounts](../cheqd-cli/cheqd-cli-accounts.md).

   ```bash
   cheqd-noded keys list
   ```

   (The assumption above is that there is only one account / key that has been added on the node. In case you have multiple addresses, please jot down the preferred account address.)

## Promote a node to validator after acquiring CHEQ tokens for staking

1. **Ensure your account has a positive balance**

   Follow the guidance on [using cheqd CLI to manage accounts](../cheqd-cli/cheqd-cli-accounts.md) to check that your account is correctly showing the CHEQ testnet tokens provided to you.

   ```bash
   cheqd-noded query bank balances <address>
   ```

2. **Get your node's validator public key**

   The node validator public key is required as a parameter for the next step. More details on validator public key is mentioned in the [cheqd CLI guide on managing nodes](../cheqd-cli/cheqd-cli-node-management.md).

   ```bash
   cheqd-noded tendermint show-validator
   ```

3. **Promote your node to validator status by staking your token balance**

   You can decide how many tokens you would like to stake from your account balance. For instance, you may want to leave a portion of the balance for paying transaction fees (now and in the future).

   To promote to validation, submit a `create-validator` transaction to the network:

   ```bash
   cheqd-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas auto --gas-adjustment <multiplication-factor> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
   ```

   Parameters required in the transaction above are:

   * **`amount`**: Amount of tokens to stake. You should stake at least 1 CHEQ (= 1,000,000,000ncheq) to successfully complete a staking transaction.
   * **`from`**: Key alias of the node operator account that makes the initial stake
   * **`min-self-delegation`**: Minimum amount of tokens that the node operator promises to keep bonded
   * **`pubkey`**: Node's `bech32`-encoded validator public key from the previous step
   * **`commission-rate`**: Validator's commission rate
   * **`commission-max-rate`**: Validator's maximum commission rate, expressed as a number with up to two decimal points. The value for this cannot be changed later.
   * **`commission-max-change-rate`**: Maximum rate of change of a validator's commission rate per day, expressed as a number with up to two decimal points. The value for this cannot be changed later.
   * **`chain-id`**: Unique identifier for the chain.
     * For cheqd's current mainnet, this is `cheqd-mainnet-1`
     * For cheqd's current testnet, this is `cheqd-testnet-4`
   * **`gas`**: Maximum gas to use for *this specific* transaction. Using `auto` uses Cosmos's auto-calculation mechanism, but can also be specified manually as an integer value.
   * **gas-adjustment** (optional): If you're using `auto` gas calculation, this parameter multiplies the auto-calculated amount by the specified factor, e.g., `1.2`. This is recommended so that it leaves enough margin of error to add a bit more gas to the transaction and ensure it successfully goes through.
   * **`gas-prices`**: Maximum gas price set by the validator

Please note the parameters below are just an “**example**”.

When setting parameters such as the commission rate, a good benchmark is to consider the [commission rates set by validators on existing networks such as Cosmos ATOM chain](https://www.mintscan.io/cosmos/validators).

You will see the commission they set, the max rate they set, and the rate of change. Please use this as a guide when thinking of your own commission configurations. This is important to get right, because the `commission-max-rate` and `commission-max-change-rate` cannot be changed after they are initially set.

   ```bash
   cheqd-noded tx staking create-validator --amount 1000000000ncheq --from key-alias-name --moniker mainnet-validator-name --chain-id cheqd-mainnet-1 --min-self-delegation="1" --gas auto --gas-adjustment 1.2 --gas-prices="25ncheq" --pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"4anVUO8WhmRMqG1t4z6VxqmqZL3V7q6HqucjwZePiUw="}' --commission-max-change-rate 0.01 --commission-max-rate 0.2 --commission-rate 0.01 --node https://rpc.cheqd.net:443
   ```

1. **Check that your validator node is bonded**

   Checking that the validator is correctly bonded can be checked via any node:

   ```bash
   cheqd-noded query staking validators --node <any-rpc-url>
   ```

   Find your node by `moniker` and make sure that `status` is `BOND_STATUS_BONDED`.

2. **Check that your validator node is signing blocks and taking part in consensus**

   Find out your [validator node's hex-encoded address](../cheqd-cli/cheqd-cli-node-management.md) and look for `"ValidatorInfo":{"Address":"..."}`:

   ```bash
   cheqd-noded tendermint show-address
   ```

   Query the latest block. Open `<node-address:rpc-port/block` in a web browser. Make sure that there is a signature with your validator address in the signature list.

## Using Ledger Nano device

To use your Ledger Nano you will need to complete the following steps:

* Set-up your wallet by creating a PIN and passphrase, which must be stored securely to enable recovery if the device is lost or damaged.
* Connect your device to your PC and update the firmware to the latest version using the Ledger Live application.
* Install the Cosmos application using the software manager (Manager > Cosmos > Install).
* Adding a new key
In order to use the hardware wallet address with the cli, the user must first add it via `cheqd-noded`. This process only records the public information about the key.

To import the key first plug in the device and enter the device pin. Once you have unlocked the device navigate to the Cosmos app on the device and open it.

To add the key use the following command:

```bash
cheqd-noded keys add <name for the key> --ledger
```

Note

The `--ledger` flag tells the command line tool to talk to the ledger device and the `--index` flag selects which HD index should be used.

When running this command, the Ledger device will prompt you to verify the genereated address. Once you have done this you will get an output in the following form:

```bash
$ cheqd-noded keys add test --ledger
- name: test
 type: ledger
 address: cheqd1zx9a7rsrmy5a2hakas0vnfwpadqwp3m327f2yt
 pubkey: ‘{“@type”:“/cosmos.crypto.secp256k1.PubKey”,“key”:“Akm0MdDZpTVltoCpRmmWd/wxiosA9edjPlbNcirs4YO1"}’
 mnemonic: “”
```

## Next steps

On completion of the steps above, you would have successfully bonded a node as validator to the cheqd testnet and participating in staking/consensus.

Learn more about what you can do with your new validator node in the [cheqd CLI guide](../cheqd-cli/README.md).
