import copy
import sys
import os
import pexpect
import re
import random
import string
import json
import time

IMPLICIT_TIMEOUT = 30
ENCODING = "utf-8"
READ_BUFFER = 60000

TEST_NET_NETWORK = "cheqd-testnet-4"
LOCAL_NET_NETWORK = "cheqd"
TEST_NET_NODE_TCP = "--node tcp://rpc.cheqd.network:443"
TEST_NET_NODE_HTTP = "--node https://rpc.cheqd.network/"
LOCAL_NET_NODE_TCP = "--node tcp://localhost:26657"
LOCAL_NET_NODE_HTTP = "--node http://localhost:26657/"
TEST_NET_DESTINATION = f"{TEST_NET_NODE_TCP} --chain-id 'cheqd-testnet-4'"
TEST_NET_DESTINATION_HTTP = f"{TEST_NET_NODE_HTTP} --chain-id 'cheqd-testnet-4'"
LOCAL_NET_DESTINATION = f"{LOCAL_NET_NODE_TCP} --chain-id 'cheqd'"
LOCAL_NET_DESTINATION_HTTP = f"{LOCAL_NET_NODE_HTTP} --chain-id 'cheqd'"
GAS_AMOUNT = 90000 # 70000 throws `out of gas` sometimes
GAS_PRICE = 25
TEST_NET_FEES = "--fees 5000000ncheq"
TEST_NET_GAS_X_GAS_PRICES = "--gas 90000 --gas-prices 25ncheq"
YES_FLAG = "-y"
KEYRING_BACKEND_TEST = "--keyring-backend test"
DENOM = "ncheq"

TEST_NET_GAS_X_GAS_PRICES_INT = GAS_AMOUNT * GAS_PRICE
MAX_GAS_MAGIC_NUMBER = 1.3

# addresses and mnemonics for test net docker image
SENDER_ADDRESS = "cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96"
SENDER_MNEMONIC = "sketch mountain erode window enact net enrich smoke claim kangaroo another visual write meat latin bacon pulp similar forum guilt father state erase bright"
RECEIVER_ADDRESS= "cheqd1l9sq0se0jd3vklyrrtjchx4ua47awug5vsyeeh"
RECEIVER_MNEMONIC = "ugly dirt sorry girl prepare argue door man that manual glow scout bomb pigeon matter library transfer flower clown cat miss pluck drama dizzy"

LOCAL_SENDER_ADDRESS = os.environ["OP0_ADDRESS"]
LOCAL_RECEIVER_ADDRESS = os.environ["OP1_ADDRESS"]

CODE_0 = "\"code\":0"
CODE_5 = "\"code\":5"
CODE_11 = "\"code\":11"
CODE_1203 = "\"code\":1203"
CODE_1100 = "\"code\":1100"
CODE_1101 = "\"code\":1101"
CODE_0_DIGIT = 0


def random_string(length):
    return ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(length))


def run(command_base, command, params, expected_output):
    # ToDo: Make it more clear.
    # Quick hack for getting passing timouted transactions
    timeout_str = "Error(.*?)timed out waiting for tx to be included in a block"
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    try:
        cli.expect(expected_output)
    except pexpect.exceptions.EOF as err:
        if re.search(timeout_str, cli.before):
            get_balance(LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION)
            return cli
        raise err

    return cli


def run_interaction(cli, input_string, expected_output):
    cli.sendline(input_string)
    cli.expect(expected_output)


def get_balance(address, network_destination):
    cli = run("cheqd-noded query", "bank balances", f"{address} {network_destination}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:")
    balance = re.search(r"amount: \"(.+?)\"", cli.after).group(1).strip()

    return balance


def json_loads(s_to_load: str) -> dict:
    s = copy.copy(s_to_load)
    s = s.replace("\\", "")
    s = s.replace("\"[", "[")
    s = s.replace("]\"", "]")
    return json.loads(s)


def get_gas_extimation(s_to_search: str) -> int:
    return int(re.findall(r"\d+", s_to_search)[0])


def send_with_note(note):
    try:
        cli = run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST} --note {note}", fr"{CODE_0}(.*?)\"value\":\"1000ncheq\"")
    except pexpect.exceptions.EOF:
        time.sleep(IMPLICIT_TIMEOUT)
        cli = run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST} --note {note}", fr"{CODE_0}(.*?)\"value\":\"1000ncheq\"")

    tx_hash = re.search(r"\"txhash\":\"(.+?)\"", cli.before).group(1).strip()

    return tx_hash, note


def set_up_operator():
    name = random_string(10)
    cli = run("cheqd-noded keys", "add", f"{name} {KEYRING_BACKEND_TEST}", r"mnemonic: \"\"")
    address = re.search(r"address: (.+?)\n", cli.before).group(1).strip()
    print(address)
    pubkey = re.search(r"pubkey: (.+?)\n", cli.before).group(1).strip()
    print(pubkey)
    run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {address} 1100000000000000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}", fr"{CODE_0}(.*?)\"value\":\"1100000000000000ncheq\"")

    return name, address, pubkey


def build_create_did_msg(did: str,
                         key_id: str,
                         ver_pub_multibase_58: str) -> str:
    return f'{{ "id": "{did}", \
    "verification_method": [{{ \
       "id": "{key_id}", \
       "type": "Ed25519VerificationKey2020", \
       "controller": "{did}", \
       "public_key_multibase": "{ver_pub_multibase_58}" \
     }}], \
     "authentication": [ \
       "{key_id}" \
     ] \
    }} \ '


def build_update_did_msg(did: str,
                         key_id: str,
                         ver_pub_multibase_58: str,
                         version_id: str) -> dict:
    return json.loads(f'{{ "id": "{did}", \
     "version_id": "{version_id}", \
     "verification_method": [{{ \
       "id": "{key_id}", \
       "type": "Ed25519VerificationKey2020", \
       "controller": "{did}", \
       "public_key_multibase": "{ver_pub_multibase_58}" \
     }}], \
     "authentication": [ \
       "{key_id}" \
     ] \
    }}')


def generate_ed25519_key() -> dict:
    cli = run(
        "cheqd-noded debug",
        "ed25519 random",
        "",
        "")
    return json_loads(cli.read())

def generate_public_multibase() -> str:
    ed25519_key = generate_ed25519_key()
    pub_key_base_64 = ed25519_key["pub_key_base_64"]

    # Get multibase58 represantation
    cli = run(
        "cheqd-noded debug",
        "encoding base64-multibase58",
        fr"{pub_key_base_64}",
        "")

    return cli.read().strip()

def generate_did() -> str:
    letters = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
    return ''.join(random.choice(letters) for i in range(16))
