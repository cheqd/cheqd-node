# Running a new Validator Node

This document describes in detail how to join existing network as validator.

## Steps

1. Setup a node:

   Follow [this instruction](setting-up-a-new-node.md).

   Use corresponding `genesis.json` and `persistent peers list` form `persistent_chains` folder in the root of the repository.

2. Generate a user key:

   See [the instruction](cosmos-cli.md#managing-keys).

3. Get some tokens:

   If you would like to participate `testnet`, join [cheqd Slack community](http://cheqd.link/join-cheqd-slack) and ask for tokens here.

   Make sure that your balance is positive using [this reference](cosmos-cli.md#managing-account-balances).

4. Promote the node to the validator:

   Post `create-validator` transaction to the network:

   ```text
    cheqd-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas <amount-gas> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
   ```

   Parameters:

   * `amount` - amount of tokens to stake;
   * `from` - key alias of account that will become node operator and will make initial stake;
   * `min-self-delegation` - minimal amount of tokens that the node operator promises to keep bonded;
   * `pubkey` - validator's public key. See [this reference](cosmos-cli.md#managing-node) on how to get it;
   * `commission-rate` - validator's commission;
   * `commission-max-rate` - maximum validator's commission. Can't be changed;
   * `commission-max-change-rate` - maximum validator's commission change per day. Can't be changed;
   * `chain-id` - chain id;
   * `gas` - max gas;
   * `gas-prices` - gas price;

     `commission-max-change-rate`, `commission-max-rate` and `commission-rate` may take fraction number as `0.01`

     Example:

     ```text
     cheqd-noded tx staking create-validator --amount 50000000cheq --from steward1 --moniker steward1 --chain-id cheqd-testet --min-self-delegation="1" --gas="auto" --gas-prices="1cheq" --pubkey cosmosvalconspub1zcjduepqpmyzmytdzjhf2fjwttjsrv49t62gdexm2yttpmgzh38p0rncqg8ssrxm2l --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01"
     ```

5. Check that the validator is bonded and taking part in consensus:

   Perform the following query:

   ```text
    cheqd-noded query staking validators --node <any-rpc-url>
   ```

   Find your node by `moniker` and make sure that `status` is `BOND_STATUS_BONDED`.

   Make sure that your node is signing blocks:

   * Find out [hex encoded validator address](cosmos-cli.md#managing-node) of your node;
   * Query the latest block. Open `<rpc-url>/block` in browser;
   * Make sure that there is a signature with your validator address in the signature list.

