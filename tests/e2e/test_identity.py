import pytest
import json
import logging

from vdrtools import wallet, did
from vdrtools import cheqd_keys, cheqd_pool, cheqd_ledger
from vdrtools.error import CommonInvalidStructure

from helpers import random_string, wallet_helper, get_base_account_number_and_sequence, get_timeout_height, get_balance_vdr, transfer_tokens_vdr, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, TEST_NET_NETWORK, TEST_NET_GAS_X_GAS_PRICES_INT, GAS_PRICE, DENOM

# logger = logging.getLogger(__name__)
# logging.basicConfig(level=logging.DEBUG)

key_alias = "qaatests"
default_amount = 1000
default_memo = "test_memo"
MAX_GAS_MAGIC_NUMBER_NEGATIVE = 1.2
MAX_GAS_MAGIC_NUMBER = 1.3
TEST_NET_HTTP = "http://seed1.us.testnet.cheqd.network:26657/"
LOCAL_POOL_HTTP = "http://localhost:26657/"


@pytest.mark.parametrize(
    "magic_number_negative, magic_number_positive",
    [
        (1.2, 1.3), # base case, corner values
        (1.0, 1.5),
        (0.5, 2),
    ]
)
@pytest.mark.parametrize("transfer_amount", ["1", "2002", "3000003"])
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
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, 0, 0, DENOM, timeout_height, default_memo
    )
    request = await cheqd_ledger.tx.build_query_simulate(test_tx)
    response = await cheqd_pool.abci_query(pool_alias, request)
    response = await cheqd_ledger.tx.parse_query_simulate_resp(response)
    gas_estimation = json.loads(response)["gas_info"]["gas_used"]
    print(gas_estimation)

    # negative case
    prod_tx_negative = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, int(gas_estimation*magic_number_negative), int(gas_estimation*magic_number_negative*GAS_PRICE), DENOM, timeout_height, default_memo
    )
    prod_tx_negative_signed = await cheqd_keys.sign(wallet_handle, key_alias, prod_tx_negative)
    with pytest.raises(CommonInvalidStructure):
        await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_negative_signed)

    # positive case
    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, SENDER_ADDRESS) # get this one more time to avoid `incorrect account sequence` error
    prod_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, json.loads(res3)["pub_key"], msg, account_number, sequence_number, int(gas_estimation*magic_number_positive), int(gas_estimation*magic_number_positive*GAS_PRICE), DENOM, timeout_height, default_memo
    )
    prod_tx_signed = await cheqd_keys.sign(wallet_handle, key_alias, prod_tx)
    positive_res = await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_signed)
    assert json.loads(positive_res)["check_tx"]["code"] == 0


@pytest.mark.parametrize("transfer_amount", [1, 999, 987654321])
@pytest.mark.asyncio
async def test_token_transfer(transfer_amount):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, TEST_NET_HTTP, TEST_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, key_alias, SENDER_MNEMONIC, "")
    )["pub_key"]

    sender_balance = await get_balance_vdr(pool_alias, SENDER_ADDRESS)
    receiver_balance = await get_balance_vdr(pool_alias, RECEIVER_ADDRESS)

    await transfer_tokens_vdr(pool_alias, wallet_handle, key_alias, public_key, SENDER_ADDRESS, RECEIVER_ADDRESS, str(transfer_amount), default_memo)

    new_sender_balance = await get_balance_vdr(pool_alias, SENDER_ADDRESS)
    new_receiver_balance = await get_balance_vdr(pool_alias, RECEIVER_ADDRESS)

    assert int(new_sender_balance) == (int(sender_balance) - transfer_amount - TEST_NET_GAS_X_GAS_PRICES_INT)
    assert int(new_receiver_balance) == (int(receiver_balance) + transfer_amount)


@pytest.mark.parametrize("memo", ["a", "test_memo_test", "123qwe$%^&"])
@pytest.mark.asyncio
async def test_memo(memo):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, TEST_NET_HTTP, TEST_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, key_alias, SENDER_MNEMONIC, "")
    )["pub_key"]

    tx_hash = await transfer_tokens_vdr(pool_alias, wallet_handle, key_alias, public_key, SENDER_ADDRESS, RECEIVER_ADDRESS, str(default_amount), memo) 

    request = await cheqd_ledger.tx.build_query_get_tx_by_hash(tx_hash)
    res = await cheqd_pool.abci_query(pool_alias, request)
    res = json.loads(await cheqd_ledger.tx.parse_query_get_tx_by_hash_resp(res))

    assert memo == res["tx"]["body"]["memo"]


@pytest.mark.asyncio
async def test_did():
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, TEST_NET_HTTP, TEST_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    _did, vk = await did.create_and_store_my_did(wallet_handle, '{}')
    fqdid = "did:cheqd:cheqd:" + _did
    print(fqdid, vk)

    # create FQDID, build -> sign -> build tx -> sign -> broadcast

    req = await cheqd_ledger.cheqd.build_msg_create_did(fqdid, vk)
    print(req)
    signed_req = await cheqd_ledger.cheqd.sign_msg_write_request(wallet_handle, fqdid, req)
    print(signed_req)
