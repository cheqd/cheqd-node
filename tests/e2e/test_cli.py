import sys
import os
import pexpect
import pytest


IMPLICIT_TIMEOUT = 30
ENCODING = "utf-8"
READ_BUFFER = 6000
TEST_NET_DESTINATION = "--node 'http://18.222.221.192:26657' --chain-id 'cheqd-testnet-2'"
TEST_NET_FEES = "--gas 70000 --gas-prices 25ncheq"
YES_FLAG = "-y"

sender = "cheqd1ece09txhq6nm9fkft9jh3mce6e48ftescs5jsw"
receiver = "cheqd16d72a6kusmzml5mjhzjv63c9j5xnpsyqs8f3sk"


def run(command_base, command, params, expected_output):
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)


@pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("help", "",r"cheqd App(.*?)Usage:(.*?)Available Commands:(.*?)Flags:"),
            # ("version", "",os.environ["RELEASE_NUMBER"]), # this works against deb package but not against starport build
            ("status", "",r"\"NodeInfo\"(.*?)\"network\":\"cheqd\"(.*?)\"moniker\":\"node0\""),
        ]
    )
def test_basic(command, params, expected_output):
    command_base = "cheqd-noded"
    run(command_base, command, params, expected_output)


@pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("add", "test1", r"- name: test1(.*?)type: local(.*?)address: (.*?)pubkey: (.*?)mnemonic: "),
            ("list", None, "- name: test1"),
            ("delete", "test1 -y", r"Key deleted forever \(uh oh!\)"),
            ("add", "test2", "- name: test2"),
            ("show", "test2", "- name: test2"),
            ("show", "test9", "Error: test9 is not a valid name or address"),
        ]
    )
def test_keys(command, params, expected_output):
    command_base = "cheqd-noded keys"
    run(command_base, command, params, expected_output)


@pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("staking validators", f"{TEST_NET_DESTINATION}", r"pagination:(.*?)validators:"),
            ("bank balances", f"{sender} {TEST_NET_DESTINATION}", r"balances:(.*?)amount:(.*?)denom: ncheq(.*?)pagination:"),
        ]
    )
def test_query(command, params, expected_output):
    command_base = "cheqd-noded query"
    run(command_base, command, params, expected_output)


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("bank send", "", r"Error: accepts 3 arg\(s\), received 0"),
            ("bank send", f"{sender} {receiver} 0ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"Error: : invalid coins"),
            ("bank send", f"{sender} {receiver} 1ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"1ncheq\""),
            ("bank send", f"{sender} {receiver} 99ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"99ncheq\""),
            ("bank send", f"{receiver} {sender} 1ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":0(.*?)\"value\":\"1ncheq\""),
            ("bank send", f"{receiver} {sender} 999999999ncheq {TEST_NET_DESTINATION} {TEST_NET_FEES} {YES_FLAG}", r"\"code\":5(.*?)insufficient funds"),
        ]
    )
def test_tx(command, params, expected_output):
    command_base = "cheqd-noded tx"
    run(command_base, command, params, expected_output)


@pytest.mark.skip
@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("show-validator", "", r"\"\@type\":(.*?)\"key\":"),
        ]
    )
def test_tendermint(command, params, expected_output):
    command_base = "cheqd-noded tendermint"
    run(command_base, command, params, expected_output)
