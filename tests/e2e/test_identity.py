import pytest
import json
import logging

from vdrtools import wallet
from vdrtools import cheqd_keys, cheqd_pool, cheqd_ledger

from helpers import random_string

# logger = logging.getLogger(__name__)
# logging.basicConfig(level=logging.DEBUG)

account_id = "cheqd1ece09txhq6nm9fkft9jh3mce6e48ftescs5jsw"
pool_alias = random_string(5)
MAX_GAS_MAGIC_NUMBER = 1.3
GAS_PRICE = 25


async def wallet_helper(wallet_id=None, wallet_key="", wallet_key_derivation_method="ARGON2I_INT"):
    if not wallet_id:
        wallet_id = random_string(25)
    wallet_config = json.dumps({"id": wallet_id})
    wallet_credentials = json.dumps({"key": wallet_key, "key_derivation_method": wallet_key_derivation_method})
    await wallet.create_wallet(wallet_config, wallet_credentials)
    wallet_handle = await wallet.open_wallet(wallet_config, wallet_credentials)

    return wallet_handle, wallet_config, wallet_credentials


async def get_base_account_number_and_sequence(account_id):
    req = await cheqd_ledger.auth.build_query_account(account_id)
    resp = await cheqd_pool.abci_query(pool_alias, req)
    resp = await cheqd_ledger.auth.parse_query_account_resp(resp)
    account = json.loads(resp)["account"]
    base_account = account["value"]
    account_number = base_account["account_number"]
    account_sequence = base_account["sequence"]

    return account_number, account_sequence


async def get_timeout_height():
    TIMEOUT = 20
    info = await cheqd_pool.abci_info(pool_alias)
    info = json.loads(info)
    current_height = info["response"]["last_block_height"]

    return int(current_height) + TIMEOUT


@pytest.mark.asyncio
async def test_basic():
    await cheqd_pool.add(pool_alias, "http://seed1.us.testnet.cheqd.network:26657/", "cheqd-testnet-2")

    wallet_handle, _, _ = await wallet_helper()

    res3 = await cheqd_keys.add_from_mnemonic(
        wallet_handle, "qaatests", "oil long siege student rent jar awkward park entry ripple enable company sort people little damp arrange wise slender push brief solve tattoo cycle", ""
    )

    msg = await cheqd_ledger.bank.build_msg_send(
        account_id, "cheqd16d72a6kusmzml5mjhzjv63c9j5xnpsyqs8f3sk", "1000", "ncheq"
    )

    account_number, sequence_number = await get_base_account_number_and_sequence(account_id)

    timeout_height = await get_timeout_height()

    test_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, 1000000, 0, "ncheq", timeout_height, "test_memo"
    )

    request = await cheqd_ledger.tx.build_query_simulate(test_tx)

    response = await cheqd_pool.abci_query(pool_alias, request)

    response = await cheqd_ledger.tx.parse_query_simulate_resp(response)

    gas_estimation = json.loads(response)["gas_info"]["gas_used"]
    print(gas_estimation)

    prod_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, int(gas_estimation*MAX_GAS_MAGIC_NUMBER), int(gas_estimation*MAX_GAS_MAGIC_NUMBER*GAS_PRICE), "ncheq", timeout_height, "test_memo"
    )

    prod_tx_signed = await cheqd_keys.sign(wallet_handle, "qaatests", prod_tx)

    final_res = await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_signed)
    print(final_res)
