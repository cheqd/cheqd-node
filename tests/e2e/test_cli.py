import sys
import os
import pexpect
import pytest


@pytest.mark.parametrize(
        "command, expected_output",
        [
            ("help", "cheqd App"),
            # ("version", os.environ["RELEASE_NUMBER"]), # this works against deb package but not against starport build
        ]
    )
def test_basic(command, expected_output):
    cli = pexpect.spawn(f"cheqd-noded {command}", encoding="utf-8")
    cli.logfile = sys.stdout
    cli.expect(expected_output)

@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("add", "test1", "- name: test1"),
            ("list", None, "- name: test1"),
            ("delete", "test1 -y", "Key deleted forever"),
            ("add", "test2", "- name: test2"),
            ("show", "test2", "- name: test2"),
            ("show", "test9", "Error: test9 is not a valid name or address"),
        ]
    )
def test_keys(command, params, expected_output):
    cli = pexpect.spawn(f"cheqd-noded keys {command} {params}", encoding="utf-8")
    cli.logfile = sys.stdout
    cli.expect(expected_output)

def test_query():
    pass

def test_tendermint():
    pass

def test_tx():
    pass
