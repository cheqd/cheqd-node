# Running a new Validator Node

This document describes in detail how to join existing network as validator.

## Steps

1. Setup a node:

    Follow [this instruction](how-to-setup-a-new-node.md).

    Use corresponding `genesis.json` and `persistent peers list` form `persistent_chains` folder in the root of the repository.

2. Create an account:

    - **Generate local keys for the future account:**

        Command: `cheqd-noded keys add <key_name>`

        Example: `cheqd-noded keys add alice`

    - **Ask another member to transfer some tokens:**

        Tokens are used to post transactions. It also used to create a stake for new nodes.

        Another mmber will ask for the address of the new participant. Cosmos account address is a function of the public key.

        Use this command to find out your adress and other key information: `cheqd-noded keys show <key_name>`

3. Promote the node to the validator:

    - **Post `create-validator` transaction to the network:**
    
        ```
        cheqd-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas <amount-gas> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
        ```

        `commission-max-change-rate`, `commission-max-rate` and `commission-rate` may take fraction number as `0.01`

        Use this command to find out the `<validator-pubkey>`: `cheqd-noded tendermint show-validator`. This command **MUST** be run on the node's machine.
        
        Example:
        
        ```
        cheqd-noded tx staking create-validator --amount 50000000stake --from steward1 --moniker steward1 --chain-id cheqdnode --min-self-delegation="1" --gas="auto" --gas-prices="1token" --pubkey cosmosvalconspub1zcjduepqpmyzmytdzjhf2fjwttjsrv49t62gdexm2yttpmgzh38p0rncqg8ssrxm2l --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01"
        ```
