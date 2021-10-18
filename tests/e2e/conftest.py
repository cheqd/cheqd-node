import pytest
import re
from helpers import run, run_interaction, \
    TEST_NET_DESTINATION, TEST_NET_FEES, TEST_NET_GAS_X_GAS_PRICES, YES_FLAG, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, RECEIVER_MNEMONIC


# Recover sender and receiver keys for TESTNET
@pytest.fixture(scope="session")
def restore_test_keys():
    cli1 = run("cheqd-noded keys", "add", "qaatests --recover", r"Enter your bip39 mnemonic")
    run_interaction(cli1, SENDER_MNEMONIC, r"- name: qaatests")

    cli2 = run("cheqd-noded keys", "add", "qaatests2 --recover", r"Enter your bip39 mnemonic")
    run_interaction(cli2, RECEIVER_MNEMONIC, r"- name: qaatests2")


# Send txn with memo to check `query tx` in test
@pytest.fixture(scope="session")
def send_with_note():
    tx_memo = "test_memo_value"
    cli = run("cheqd-noded tx", "bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 1000ncheq {TEST_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} --note {tx_memo}", r"\"code\":0(.*?)\"value\":\"1000ncheq\"")
    tx_hash = re.search(r"\"txhash\":\"(.+?)\"", cli.before).group(1).strip()
    yield tx_hash, tx_memo
