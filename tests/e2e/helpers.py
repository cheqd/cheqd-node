import sys
import os
import pexpect
import re
import random
import string
import json

from vdrtools import wallet
from vdrtools import cheqd_keys, cheqd_pool, cheqd_ledger

IMPLICIT_TIMEOUT = 30
ENCODING = "utf-8"
READ_BUFFER = 6000

TEST_NET_NETWORK = "cheqd-testnet-2"
LOCAL_NET_NETWORK = "cheqd"
TEST_NET_NODE_TCP = "--node 'tcp://seed1.us.testnet.cheqd.network:26657'"
TEST_NET_NODE_HTTP = "--node http://node1.eu.testnet.cheqd.network:26657/"
LOCAL_NET_NODE_TCP = "--node 'tcp://localhost:26657'"
LOCAL_NET_NODE_HTTP = "--node http://localhost:26657/"
TEST_NET_DESTINATION = f"{TEST_NET_NODE_TCP} --chain-id 'cheqd-testnet-2'"
TEST_NET_DESTINATION_HTTP = f"{TEST_NET_NODE_HTTP} --chain-id 'cheqd-testnet-2'"
LOCAL_NET_DESTINATION = f"{LOCAL_NET_NODE_TCP} --chain-id 'cheqd'"
LOCAL_NET_DESTINATION_HTTP = f"{LOCAL_NET_NODE_HTTP} --chain-id 'cheqd'"
TEST_NET_FEES = "--fees 5000000ncheq"
TEST_NET_GAS_X_GAS_PRICES = "--gas 70000 --gas-prices 25ncheq"
YES_FLAG = "-y"
DENOM = "ncheq"
GAS_AMOUNT = 80000 # 70000 throws `out of gas` sometimes
GAS_PRICE = 25
TEST_NET_GAS_X_GAS_PRICES_INT = GAS_AMOUNT * GAS_PRICE

# SENDER_ADDRESS = "cheqd1ece09txhq6nm9fkft9jh3mce6e48ftescs5jsw"
# SENDER_MNEMONIC = "oil long siege student rent jar awkward park entry ripple enable company sort people little damp arrange wise slender push brief solve tattoo cycle"
# RECEIVER_ADDRESS= "cheqd16d72a6kusmzml5mjhzjv63c9j5xnpsyqs8f3sk"
# RECEIVER_MNEMONIC = "strike impact earth indoor man illness virus genuine rib control antenna loop neck rotate bargain original nasty size either try snap quiz stairs huge"

SENDER_ADDRESS = "cheqd1t07jem7yv740rq29ctwecxyvm9qegq6rvdwkcp"
SENDER_MNEMONIC = "exclude slam riot window wink peace lemon interest token accident pupil wall squirrel slight endless manage cereal celery local teach galaxy culture exact cliff"
RECEIVER_ADDRESS= "cheqd1xrhxp68j8wuy62uhrsje6g90pld55e4ncls4uz"
RECEIVER_MNEMONIC = "join coconut smooth number unfair future banner mad lawn deny virtual derive cradle brain business pyramid absorb crush couch cook cliff job poet differ"

LOCAL_SENDER_ADDRESS = os.environ["OP0_ADDRESS"]
LOCAL_RECEIVER_ADDRESS = os.environ["OP1_ADDRESS"]

CODE_0 = "\"code\":0"


def random_string(length):
    return ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(length))


def run(command_base, command, params, expected_output):
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)

    return cli


def run_interaction(cli, input_string, expected_output):
    cli.sendline(input_string)
    cli.expect(expected_output)


def get_balance(address, network_destination):
    cli = run("cheqd-noded query", "bank balances", f"{address} {network_destination}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:")
    balance = re.search(r"amount: \"(.+?)\"", cli.after).group(1).strip()

    return balance


async def get_balance_vdr(pool_alias, address):
    request = await cheqd_ledger.bank.build_query_balance(address, DENOM)
    res = await cheqd_pool.abci_query(pool_alias, request)
    res = await cheqd_ledger.bank.parse_query_balance_resp(res)
    sender_balance = json.loads(res)["balance"]["amount"]

    return sender_balance


def send_with_note(note):
    cli = run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} --note {note}", fr"{CODE_0}(.*?)\"value\":\"1000ncheq\"")
    tx_hash = re.search(r"\"txhash\":\"(.+?)\"", cli.before).group(1).strip()

    return tx_hash, note


async def send_tx_helper(pool_alias, wallet_handle, key_alias, public_key, sender_address, msg, memo):
    timeout_height = await get_timeout_height(pool_alias)
    account_number, sequence_number = await get_base_account_number_and_sequence(pool_alias, sender_address)
    tx = await cheqd_ledger.auth.build_tx(
        pool_alias, public_key, msg, account_number, sequence_number, GAS_AMOUNT, GAS_AMOUNT*GAS_PRICE, DENOM, sender_address, timeout_height, memo
    )
    tx_signed = await cheqd_ledger.auth.sign_tx(wallet_handle, key_alias, tx)
    res = json.loads(await cheqd_pool.broadcast_tx_commit(pool_alias, tx_signed))
    tx_hash = res["hash"]

    return res, tx_hash


async def create_did_helper(pool_alias, wallet_handle, key_alias, public_key, sender_address, fqdid, vk, memo):
    req = await cheqd_ledger.cheqd.build_msg_create_did(fqdid, vk)
    signed_req = await cheqd_ledger.cheqd.sign_msg_write_request(wallet_handle, fqdid, bytes(req))
    res, _ = await send_tx_helper(pool_alias, wallet_handle, key_alias, public_key, sender_address, bytes(signed_req), memo)

    return res


async def query_did_helper(pool_alias, fqdid):
    req = await cheqd_ledger.cheqd.build_query_get_did(fqdid)
    res = await cheqd_pool.abci_query(pool_alias, req)

    return res


async def update_did_helper(pool_alias, wallet_handle, key_alias, public_key, sender_address, fqdid, new_vk, version_id, memo):
    req = await cheqd_ledger.cheqd.build_msg_update_did(fqdid, new_vk, version_id)
    signed_req = await cheqd_ledger.cheqd.sign_msg_write_request(wallet_handle, fqdid, bytes(req))
    res, _ = await send_tx_helper(pool_alias, wallet_handle, key_alias, public_key, sender_address, bytes(signed_req), memo)

    return res


def set_up_operator():
    name = random_string(10)
    cli = run("cheqd-noded keys", "add", name, r"mnemonic: \"\"")
    address = re.search(r"address: (.+?)\n", cli.before).group(1).strip()
    print(address)
    pubkey = re.search(r"pubkey: (.+?)\n", cli.before).group(1).strip()
    print(pubkey)
    run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {address} 1100000000000000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"1100000000000000ncheq\"")

    return name, address, pubkey


async def wallet_helper(wallet_id=None, wallet_key="", wallet_key_derivation_method="ARGON2I_INT"):
    if not wallet_id:
        wallet_id = random_string(25)
    wallet_config = json.dumps({"id": wallet_id})
    wallet_credentials = json.dumps({"key": wallet_key, "key_derivation_method": wallet_key_derivation_method})
    await wallet.create_wallet(wallet_config, wallet_credentials)
    wallet_handle = await wallet.open_wallet(wallet_config, wallet_credentials)

    return wallet_handle, wallet_config, wallet_credentials


async def get_base_account_number_and_sequence(pool_alias, account_id):
    req = await cheqd_ledger.auth.build_query_account(account_id)
    resp = await cheqd_pool.abci_query(pool_alias, req)
    resp = await cheqd_ledger.auth.parse_query_account_resp(resp)
    account = json.loads(resp)["account"]
    base_account = account["value"]
    account_number = base_account["account_number"]
    account_sequence = base_account["sequence"]

    return account_number, account_sequence


async def get_timeout_height(pool_alias):
    TIMEOUT = 20
    info = await cheqd_pool.abci_info(pool_alias)
    info = json.loads(info)
    current_height = info["response"]["last_block_height"]

    return int(current_height) + TIMEOUT
