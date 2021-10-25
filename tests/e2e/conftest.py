import pytest
import re

from helpers import run, run_interaction, random_string, \
    TEST_NET_DESTINATION, TEST_NET_FEES, TEST_NET_GAS_X_GAS_PRICES, YES_FLAG, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, RECEIVER_MNEMONIC, LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION, CODE_0


# Recover sender and receiver keys for TESTNET
@pytest.fixture(scope="session")
def restore_test_keys():
    cli1 = run("cheqd-noded keys", "add", "qaatests --recover", r"Enter your bip39 mnemonic")
    run_interaction(cli1, SENDER_MNEMONIC, r"- name: qaatests")

    cli2 = run("cheqd-noded keys", "add", "qaatests2 --recover", r"Enter your bip39 mnemonic")
    run_interaction(cli2, RECEIVER_MNEMONIC, r"- name: qaatests2")


@pytest.fixture(scope='session')
def create_export_keys():
    command_base = "cheqd-noded keys"
    run(command_base, "add", "export_key", "name: export_key")
