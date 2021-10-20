import json
import pytest
from hypothesis import settings, given, strategies, Phase, Verbosity
from string import digits, ascii_letters
from helpers import run, run_interaction, get_balance, send_with_note, set_up_operator, \
    TEST_NET_NETWORK, TEST_NET_NODE_TCP, TEST_NET_NODE_HTTP, TEST_NET_DESTINATION, TEST_NET_DESTINATION_HTTP, \
    LOCAL_NET_NETWORK, LOCAL_NET_NODE_TCP, LOCAL_NET_NODE_HTTP, LOCAL_NET_DESTINATION, LOCAL_NET_DESTINATION_HTTP, \
    TEST_NET_FEES, TEST_NET_GAS_X_GAS_PRICES, YES_FLAG, \
    SENDER_ADDRESS, RECEIVER_ADDRESS, LOCAL_SENDER_ADDRESS, LOCAL_RECEIVER_ADDRESS,CODE_0, TEST_NET_GAS_X_GAS_PRICES_INT


@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("help", "",r"cheqd App(.*?)Usage:(.*?)Available Commands:(.*?)Flags:"),
            ("status", TEST_NET_NODE_TCP, fr"\"NodeInfo\"(.*?)\"network\":\"{TEST_NET_NETWORK}\"(.*?)\"moniker\":\"seed1-us-testnet-cheqd\""), # tcp + us node
            ("status", TEST_NET_NODE_HTTP, fr"\"NodeInfo\"(.*?)\"network\":\"{TEST_NET_NETWORK}\"(.*?)\"moniker\":\"node1-eu-testnet-cheqd\""), # http + eu node
            ("status", LOCAL_NET_NODE_TCP, fr"\"NodeInfo\"(.*?)\"network\":\"{LOCAL_NET_NETWORK}\"(.*?)\"moniker\":\"node0\""), # tcp + local
            ("status", LOCAL_NET_NODE_HTTP, fr"\"NodeInfo\"(.*?)\"network\":\"{LOCAL_NET_NETWORK}\"(.*?)\"moniker\":\"node0\""), # http + local
        ]
    )
def test_basic(command, params, expected_output):
    command_base = "cheqd-noded"
    run(command_base, command, params, expected_output)


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


@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("staking validators", f"{TEST_NET_DESTINATION}", r"pagination:(.*?)validators:"), # test net
            ("bank balances", f"{SENDER_ADDRESS} {TEST_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),
            ("bank balances", f"{RECEIVER_ADDRESS} {TEST_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),

            ("staking validators", f"{LOCAL_NET_DESTINATION}", r"pagination:(.*?)validators:"), # local net
            ("bank balances", f"{LOCAL_SENDER_ADDRESS} {LOCAL_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),
            ("bank balances", f"{LOCAL_RECEIVER_ADDRESS} {LOCAL_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),
        ]
    )
def test_query(command, params, expected_output):
    command_base = "cheqd-noded query"
    run(command_base, command, params, expected_output)


@pytest.mark.usefixtures('restore_test_keys') # for test net
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("bank send", "", r"Error: accepts 3 arg\(s\), received 0"), # no args
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} -1ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"Error: unknown shorthand flag: '1' in -1ncheq"), # -1
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 0ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"Error: : invalid coins"), # 0
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"1ncheq\""), # 1 + fees
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 2ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"2ncheq\""), # 2 + gas x price
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 99ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"99ncheq\""),
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1ncheq {LOCAL_NET_DESTINATION} {YES_FLAG}", r"\"code\":13(.*?)insufficient fees"), # no fees
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 2ncheq {LOCAL_NET_DESTINATION} --fees 4000000ncheq {YES_FLAG}", r"\"code\":13(.*?)insufficient fees"), # bad fees
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 3ncheq {LOCAL_NET_DESTINATION} --gas 70000 --gas-prices 1ncheq {YES_FLAG}", r"\"code\":13(.*?)insufficient fees"), # bad gas price
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 4ncheq {LOCAL_NET_DESTINATION} --gas 1 --gas-prices 25ncheq {YES_FLAG}", r"\"code\":11(.*?)out of gas"), # bad gas amount
            ("bank send", f"{LOCAL_RECEIVER_ADDRESS} {LOCAL_SENDER_ADDRESS} 2ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"2ncheq\""), # transfer back: 2 + fees
            ("bank send", f"{LOCAL_RECEIVER_ADDRESS} {LOCAL_SENDER_ADDRESS} 1ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"1ncheq\""), # transfer back: 1 + gas x price
            ("bank send", f"{LOCAL_RECEIVER_ADDRESS} {LOCAL_SENDER_ADDRESS} 99999999999999999ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":5(.*?)insufficient funds"),
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1000ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} --note 'test123!=$'", fr"{CODE_0}(.*?)\"value\":\"1000ncheq\""), # note
            ("bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 9999ncheq {LOCAL_NET_DESTINATION_HTTP} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"9999ncheq\""), # http + gas x price
            ("bank send", f"{LOCAL_RECEIVER_ADDRESS} {LOCAL_SENDER_ADDRESS} 9999ncheq {LOCAL_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"9999ncheq\""), # transfer back: tcp + fees
        ]
    )
def test_tx_bank_send(command, params, expected_output):
    command_base = "cheqd-noded tx"
    run(command_base, command, params, expected_output)


# TODO get observers' pubkeys to promote them
NODE_PUBKEY = json.dumps({"@type":"/cosmos.crypto.ed25519.PubKey","key":"+Gt8W3guq0TE0HuVuJBI3maNhj2uCW02CZE9pAbkiA8="})
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("staking create-validator", "", r"Error: required flag\(s\) \"amount\", \"from\", \"pubkey\" not set"), # no args
            ("staking create-validator", fr"--amount 1000000000000000ncheq --from {set_up_operator()[0]} --pubkey {set_up_operator()[2]} --min-self-delegation='1' --commission-max-change-rate='0.02' --commission-max-rate='0.02' --commission-rate='0.01' {TEST_NET_FEES} {LOCAL_NET_DESTINATION} {YES_FLAG}", r"\"code\":6(.*?)validator pubkey type is not supported"), # wrong pubkey type
            ("staking create-validator", fr"--amount 1000000000000000ncheq --from {set_up_operator()[0]} --pubkey '{NODE_PUBKEY}' --min-self-delegation='1' --commission-max-change-rate='0.02' --commission-max-rate='0.02' --commission-rate='0.01' {TEST_NET_FEES} {LOCAL_NET_DESTINATION} {YES_FLAG}", fr"{CODE_0}"),
            ("staking create-validator", fr"--amount 1000000000000000ncheq --from {set_up_operator()[0]} --pubkey '{NODE_PUBKEY}' --min-self-delegation='1' --commission-max-change-rate='0.02' --commission-max-rate='0.02' --commission-rate='0.01' {TEST_NET_FEES} {LOCAL_NET_DESTINATION} {YES_FLAG}", fr"\"code\":5(.*?)validator already exist for this pubkey"), # promote twice the same pubkey
        ]
)
def test_tx_staking_create(command, params, expected_output):
    command_base = "cheqd-noded tx"
    run(command_base, command, params, expected_output)


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


@settings(deadline=None, max_examples=25)
@given(note=strategies.text(ascii_letters, min_size=1, max_size=1024))
def test_memo(note):
    tx_hash, tx_memo = send_with_note(note)
    run("cheqd-noded query", "tx", f"{tx_hash} {LOCAL_NET_DESTINATION}", fr"code: 0(.*?)memo: {tx_memo}(.*?)txhash: {tx_hash}") # check that txn has correct memo value


@settings(deadline=None, max_examples=25)
@given(value=strategies.integers(min_value=1, max_value=999999999))
def test_token_transfer(value):
    sender_balance = get_balance(LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION)
    receiver_balance = get_balance(LOCAL_RECEIVER_ADDRESS, LOCAL_NET_DESTINATION)

    run("cheqd-noded tx", "bank send", f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {value}ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG}", fr"{CODE_0}(.*?)\"value\":\"{value}ncheq\"")

    new_sender_balance = get_balance(LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION)
    new_receiver_balance = get_balance(LOCAL_RECEIVER_ADDRESS, LOCAL_NET_DESTINATION)

    assert int(new_sender_balance) == (int(sender_balance) - value - TEST_NET_GAS_X_GAS_PRICES_INT)
    assert int(new_receiver_balance) == (int(receiver_balance) + value)
