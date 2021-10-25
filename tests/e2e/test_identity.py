import pytest
import json
import logging

from vdrtools import wallet
from vdrtools import cheqd_keys, cheqd_pool, cheqd_ledger
from vdrtools.error import CommonInvalidStructure

from helpers import random_string, wallet_helper, get_base_account_number_and_sequence, get_timeout_height, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, TEST_NET_NETWORK

# logger = logging.getLogger(__name__)
# logging.basicConfig(level=logging.DEBUG)

key_alias = "qaatests"
memo = "test_memo"
MAX_GAS_MAGIC_NUMBER_NEGATIVE = 1.2
MAX_GAS_MAGIC_NUMBER = 1.3
GAS_PRICE = 25
DENOM = "ncheq"
TEST_NET_HTTP = "http://seed1.us.testnet.cheqd.network:26657/"


@pytest.mark.parametrize(
    "magic_number_negative, magic_number_positive",
    [
        (1.2, 1.3), # base case
        (1.0, 1.5),
        (0.5, 2),
    ]
)
@pytest.mark.parametrize("transfer_amount", ["1", "1000", "1000000"])
@pytest.mark.asyncio
async def test_gas_estimation(magic_number_negative, magic_number_positive, transfer_amount):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, TEST_NET_HTTP, TEST_NET_NETWORK)

    wallet_handle, _, _ = await wallet_helper()

    res3 = await cheqd_keys.add_from_mnemonic(
        wallet_handle, key_alias, SENDER_MNEMONIC, ""
    )

    msg = await cheqd_ledger.bank.build_msg_send(
        SENDER_ADDRESS, RECEIVER_ADDRESS, transfer_amount, DENOM
    )

    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, SENDER_ADDRESS)

    timeout_height = await get_timeout_height(pool_alias)

    test_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, 0, 0, DENOM, timeout_height, memo
    )

    request = await cheqd_ledger.tx.build_query_simulate(test_tx)

    response = await cheqd_pool.abci_query(pool_alias, request)

    response = await cheqd_ledger.tx.parse_query_simulate_resp(response)

    gas_estimation = json.loads(response)["gas_info"]["gas_used"]
    print(gas_estimation)

    # negative case
    prod_tx_negative = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, int(gas_estimation*magic_number_negative), int(gas_estimation*magic_number_negative*GAS_PRICE), DENOM, timeout_height, memo
    )

    prod_tx_negative_signed = await cheqd_keys.sign(wallet_handle, key_alias, prod_tx_negative)

    with pytest.raises(CommonInvalidStructure):
        await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_negative_signed)

    # positive case
    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, SENDER_ADDRESS) # get this one more time to avoid `incorrect account sequence` error

    prod_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, int(gas_estimation*magic_number_positive), int(gas_estimation*magic_number_positive*GAS_PRICE), DENOM, timeout_height, memo
    )

    prod_tx_signed = await cheqd_keys.sign(wallet_handle, key_alias, prod_tx)

    positive_res = await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_signed)

    assert json.loads(positive_res)["check_tx"]["code"] == 0


# @pytest.mark.asyncio
# async def test_token_transfer():
#     pass


# @pytest.mark.asyncio
# async def test_memo():
#     pass
