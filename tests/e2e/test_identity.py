import pytest
import json
import logging
import time
from hypothesis import settings, given, strategies, Phase, Verbosity
from string import digits, ascii_letters

from vdrtools import wallet, did
from vdrtools import cheqd_keys, cheqd_pool, cheqd_ledger
from vdrtools.error import CommonInvalidStructure, PoolLedgerConfigAlreadyExistsError

from helpers import create_did_helper, query_did_helper, random_string, update_did_helper, wallet_helper, get_base_account_number_and_sequence, get_timeout_height, get_balance_vdr, send_tx_helper, get_tx_helper, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, LOCAL_NET_NETWORK, TEST_NET_GAS_X_GAS_PRICES_INT, GAS_AMOUNT, GAS_PRICE, DENOM, IMPLICIT_TIMEOUT, CODE_0_DIGIT

# logger = logging.getLogger(__name__)
# logging.basicConfig(level=logging.DEBUG)

KEY_ALIAS = "operator0"
DEFAULT_AMOUNT = 1000
DEFAULT_MEMO = "test_memo"
MAX_GAS_MAGIC_NUMBER_NEGATIVE = 1.2
MAX_GAS_MAGIC_NUMBER = 1.3
LOCAL_POOL_HTTP = "http://localhost:26657"
FQ_PREFIX = "did:cheqd:testnet:"
MAX_GAS = 0
MAX_COIN_AMOUNT = 0


@pytest.mark.parametrize(
    "magic_number_negative, magic_number_positive",
    [
        (1.2, 1.3), # base case, corner values
        (1.0, 1.5),
        (0.5, 2),
    ]
)
@pytest.mark.parametrize("transfer_amount", ["1", "20002", "3000003", "9000000009"])
@pytest.mark.asyncio
async def test_gas_estimation(magic_number_negative, magic_number_positive, transfer_amount):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    msg = await cheqd_ledger.bank.build_msg_send(
        SENDER_ADDRESS, RECEIVER_ADDRESS, transfer_amount, DENOM
    )
    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, SENDER_ADDRESS)
    timeout_height = await get_timeout_height(pool_alias)
    test_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, public_key, msg, account_number, sequence_number, MAX_GAS, MAX_COIN_AMOUNT, DENOM, SENDER_ADDRESS, timeout_height, DEFAULT_MEMO
    )
    request = await cheqd_ledger.tx.build_query_simulate(test_tx)
    response = await cheqd_pool.abci_query(pool_alias, request)
    response = await cheqd_ledger.tx.parse_query_simulate_resp(response)
    gas_estimation = json.loads(response)["gas_info"]["gas_used"]

    # negative case
    prod_tx_negative = await cheqd_ledger.auth.build_tx(
        pool_alias, public_key, msg, account_number, sequence_number, int(gas_estimation*magic_number_negative), int(gas_estimation*magic_number_negative*GAS_PRICE), DENOM, SENDER_ADDRESS, timeout_height, DEFAULT_MEMO
    )
    prod_tx_negative_signed = await cheqd_ledger.auth.sign_tx(wallet_handle, KEY_ALIAS, prod_tx_negative)
    with pytest.raises(CommonInvalidStructure):
        await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_negative_signed)

    # positive case
    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, SENDER_ADDRESS) # get this one more time to avoid `incorrect account sequence` error
    prod_tx = await cheqd_ledger.auth.build_tx(
        pool_alias, public_key, msg, account_number, sequence_number, int(gas_estimation*magic_number_positive), int(gas_estimation*magic_number_positive*GAS_PRICE), DENOM, SENDER_ADDRESS, timeout_height, DEFAULT_MEMO
    )
    prod_tx_signed = await cheqd_ledger.auth.sign_tx(wallet_handle, KEY_ALIAS, prod_tx)
    positive_res = await cheqd_pool.broadcast_tx_commit(pool_alias, prod_tx_signed)
    assert json.loads(positive_res)["check_tx"]["code"] == CODE_0_DIGIT


@pytest.mark.parametrize("transfer_amount", [1, 999, 1001, 987654321]) # TODO: hypothesis
@pytest.mark.asyncio
async def test_token_transfer(transfer_amount):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    sender_balance = await get_balance_vdr(pool_alias, SENDER_ADDRESS)
    receiver_balance = await get_balance_vdr(pool_alias, RECEIVER_ADDRESS)

    msg = await cheqd_ledger.bank.build_msg_send(
        SENDER_ADDRESS, RECEIVER_ADDRESS, str(transfer_amount), DENOM
    )
    await send_tx_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, msg, DEFAULT_MEMO)

    new_sender_balance = await get_balance_vdr(pool_alias, SENDER_ADDRESS)
    new_receiver_balance = await get_balance_vdr(pool_alias, RECEIVER_ADDRESS)

    assert int(new_sender_balance) == (int(sender_balance) - transfer_amount - TEST_NET_GAS_X_GAS_PRICES_INT)
    assert int(new_receiver_balance) == (int(receiver_balance) + transfer_amount)


@pytest.mark.parametrize("note", ["a", "1", "test_memo_test", "123qwe$%^&", "______________________________"]) # TODO: hypothesis
@pytest.mark.asyncio
async def test_memo(note): # intermittent failures here due to `Internal error: timed out waiting for tx to be included in a block`
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    msg = await cheqd_ledger.bank.build_msg_send(
        SENDER_ADDRESS, RECEIVER_ADDRESS, str(DEFAULT_AMOUNT), DENOM
    )
    res, tx_hash = await send_tx_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, msg, note)
    assert res["check_tx"]["code"] == CODE_0_DIGIT

    try:
        res = await get_tx_helper(pool_alias, tx_hash)
        assert res["tx"]["body"]["memo"] == note
    except TypeError:
        time.sleep(IMPLICIT_TIMEOUT)
        res = await get_tx_helper(pool_alias, tx_hash)
        assert res["tx"]["body"]["memo"] == note



@pytest.mark.asyncio
async def test_did_positive():
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    # create
    _did, vk = await did.create_and_store_my_did(wallet_handle, '{}')
    fqdid = FQ_PREFIX + _did

    res = await create_did_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, fqdid, vk, DEFAULT_MEMO)
    assert res["check_tx"]["code"] == CODE_0_DIGIT
    assert res["deliver_tx"]["code"] == CODE_0_DIGIT
    parsed_res = json.loads(await cheqd_ledger.cheqd.parse_msg_create_did(json.dumps(res)))
    assert parsed_res["id"] == fqdid

    # query
    res = await query_did_helper(pool_alias, fqdid)
    parsed_res = json.loads(await cheqd_ledger.cheqd.parse_query_get_did_resp(res))
    version_id = parsed_res["metadata"]["version_id"]
    assert parsed_res["did"]["id"] == fqdid
    assert parsed_res["did"]["verification_method"][0]["public_key_multibase"] == f"z{vk}"

    # update
    new_vk = await did.replace_keys_start(wallet_handle, _did, '{}')
    await did.replace_keys_apply(wallet_handle, _did)

    res = await update_did_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, fqdid, new_vk, version_id, DEFAULT_MEMO)
    assert res["check_tx"]["code"] == CODE_0_DIGIT
    assert res["deliver_tx"]["code"] == CODE_0_DIGIT
    parsed_res = json.loads(await cheqd_ledger.cheqd.parse_msg_update_did(json.dumps(res)))
    assert parsed_res["id"] == fqdid

    # query
    res = await query_did_helper(pool_alias, fqdid)
    parsed_res = json.loads(await cheqd_ledger.cheqd.parse_query_get_did_resp(res))
    new_version_id = parsed_res["metadata"]["version_id"]
    assert version_id != new_version_id
    assert parsed_res["did"]["id"] == fqdid
    assert parsed_res["did"]["verification_method"][0]["public_key_multibase"] == f"z{new_vk}" # new vk


@pytest.mark.parametrize("version_id", ["test",  "12345", "+M/qoYGJqKE1mwRFSINHIW9cKNGcskTGqPww2kk9aes="])
@pytest.mark.asyncio
async def test_did_update_wrong_version(version_id):
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    # create
    _did, vk = await did.create_and_store_my_did(wallet_handle, '{}')
    fqdid = FQ_PREFIX + _did

    await create_did_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, fqdid, vk, DEFAULT_MEMO)

    # update
    new_vk = await did.replace_keys_start(wallet_handle, _did, '{}')
    await did.replace_keys_apply(wallet_handle, _did)

    req = await cheqd_ledger.cheqd.build_msg_update_did(fqdid, new_vk, version_id)
    signed_req = await cheqd_ledger.cheqd.sign_msg_write_request(wallet_handle, fqdid, bytes(req))
    with pytest.raises(CommonInvalidStructure):
        await send_tx_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, bytes(signed_req), DEFAULT_MEMO)


@pytest.mark.asyncio
async def test_did_update_wrong_vk():
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()
    public_key = json.loads(
        await cheqd_keys.add_from_mnemonic(wallet_handle, KEY_ALIAS, SENDER_MNEMONIC, "")
    )["pub_key"]

    # create
    _did, vk = await did.create_and_store_my_did(wallet_handle, '{}')
    fqdid = FQ_PREFIX + _did

    await create_did_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, fqdid, vk, DEFAULT_MEMO)

    # query
    res = await query_did_helper(pool_alias, fqdid)
    parsed_res = json.loads(await cheqd_ledger.cheqd.parse_query_get_did_resp(res))
    version_id = parsed_res["metadata"]["version_id"]

    # update
    _, new_vk = await did.create_and_store_my_did(wallet_handle, '{}')

    req = await cheqd_ledger.cheqd.build_msg_update_did(fqdid, new_vk, version_id) # new vk
    signed_req = await cheqd_ledger.cheqd.sign_msg_write_request(wallet_handle, fqdid, bytes(req))
    with pytest.raises(CommonInvalidStructure):
        await send_tx_helper(pool_alias, wallet_handle, KEY_ALIAS, public_key, SENDER_ADDRESS, bytes(signed_req), DEFAULT_MEMO)


@pytest.mark.asyncio
async def test_did_query_non_existent():
    pool_alias = random_string(5)
    await cheqd_pool.add(pool_alias, LOCAL_POOL_HTTP, LOCAL_NET_NETWORK)
    wallet_handle, _, _ = await wallet_helper()

    _did, _ = await did.create_and_store_my_did(wallet_handle, '{}')
    fqdid = FQ_PREFIX + _did

    # query
    res = await query_did_helper(pool_alias, fqdid)
    with pytest.raises(CommonInvalidStructure):
        await cheqd_ledger.cheqd.parse_query_get_did_resp(res)
