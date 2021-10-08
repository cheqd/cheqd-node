import sys
import os
import pexpect
import pytest


IMPLICIT_TIMEOUT = 30
ENCODING = "utf-8"
READ_BUFFER = 6000


# @pytest.mark.skip
@pytest.mark.parametrize(
        "command, expected_output",
        [
            ("help", r"cheqd App(.*?)Usage:(.*?)Available Commands:(.*?)Flags:"),
            # ("version", os.environ["RELEASE_NUMBER"]), # this works against deb package but not against starport build
            ("status", r"\"NodeInfo\"(.*?)\"network\":\"cheqd\"(.*?)\"moniker\":\"node0\""),
        ]
    )
def test_basic(command, expected_output):
    command_base = "cheqd-noded"
    cli = pexpect.spawn(f"{command_base} {command}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)


# @pytest.mark.skip
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
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)


@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("staking", "validators", r"pagination:(.*?)validators:"),
        ]
    )
def test_query(command, params, expected_output):
    command_base = "cheqd-noded query"
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)


# def test_tendermint():
#     command_base = "cheqd-noded tendermint"
#     cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT)
#     cli.logfile = sys.stdout
#     cli.expect(expected_output)


# def test_tx():
#     command_base = "cheqd-noded tx"
#     cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT)
#     cli.logfile = sys.stdout
#     cli.expect(expected_output)
