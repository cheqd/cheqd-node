import re
import pytest
from helpers import run, run_interaction, \
    TEST_NET_DESTINATION, TEST_NET_FEES, TEST_NET_GAS_X_GAS_PRICES, YES_FLAG, \
    SENDER_ADDRESS, SENDER_MNEMONIC, RECEIVER_ADDRESS, RECEIVER_MNEMONIC


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("help", "",r"cheqd App(.*?)Usage:(.*?)Available Commands:(.*?)Flags:"),
            ("status", "--node 'tcp://seed1.us.testnet.cheqd.network:26657'",r"\"NodeInfo\"(.*?)\"network\":\"cheqd-testnet-2\"(.*?)\"moniker\":\"seed1-us-testnet-cheqd\""),
        ]
    )
def test_basic(command, params, expected_output):
    command_base = "cheqd-noded"
    run(command_base, command, params, expected_output)


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("add", "test1", r"- name: test1(.*?)type: local(.*?)address: (.*?)pubkey: (.*?)mnemonic: "),
            ("list", None, "- name: test1"),
            ("delete", f"test1 {YES_FLAG}", r"Key deleted forever \(uh oh!\)"),
            ("add", "test2", "- name: test2"),
            ("show", "test2", "- name: test2"),
            ("delete", f"test2 {YES_FLAG}", r"Key deleted forever \(uh oh!\)"),
            ("show", "test9", "Error: test9 is not a valid name or address"),
        ]
    )
def test_keys(command, params, expected_output):
    command_base = "cheqd-noded keys"
    run(command_base, command, params, expected_output)


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("staking validators", f"{TEST_NET_DESTINATION}", r"pagination:(.*?)validators:"),
            ("bank balances", f"{SENDER_ADDRESS} {TEST_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),
        ]
    )
def test_query(command, params, expected_output):
    command_base = "cheqd-noded query"
    run(command_base, command, params, expected_output)


@pytest.mark.usefixtures('restore_test_keys')
# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("bank send", "", r"Error: accepts 3 arg\(s\), received 0"), # no args
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 0ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"Error: : invalid coins"), # 0
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 1ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"1ncheq\""), # 1 + fees
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 2ncheq {TEST_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"2ncheq\""), # 2 + gas x price
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 99ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"99ncheq\""),
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 99ncheq {TEST_NET_DESTINATION} {YES_FLAG}", r"\"code\":13(.*?)insufficient fees"),
            ("bank send", f"{RECEIVER_ADDRESS} {SENDER_ADDRESS} 2ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"2ncheq\""), # transfer back 2 + fees
            ("bank send", f"{RECEIVER_ADDRESS} {SENDER_ADDRESS} 1ncheq {TEST_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"1ncheq\""), # transfer back 1 + gas x price
            ("bank send", f"{RECEIVER_ADDRESS} {SENDER_ADDRESS} 999999999ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":5(.*?)insufficient funds"),
            ("bank send", f"{SENDER_ADDRESS} {RECEIVER_ADDRESS} 1000ncheq {TEST_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} --note 'test123!=$'", r"\"code\":0(.*?)\"value\":\"1000ncheq\""),
        ]
    )
def test_tx(command, params, expected_output):
    command_base = "cheqd-noded tx"
    run(command_base, command, params, expected_output)


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("show-address", "", r"cheqd(.*?)"),
            ("show-node-id", "", r"^[a-z\d]{40}"),
            ("show-validator", "", r"\"\@type\":(.*?)\"key\":"),
        ]
    )
def test_tendermint(command, params, expected_output):
    command_base = "cheqd-noded tendermint"
    run(command_base, command, params, expected_output)


# @pytest.mark.skip
def test_production(send_with_note):
    tx_hash, tx_memo = send_with_note
    run("cheqd-noded query", "tx", f"{tx_hash} {TEST_NET_DESTINATION}", fr"code: 0(.*?)memo: {tx_memo}(.*?)txhash: {tx_hash}")
