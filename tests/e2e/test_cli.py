import sys
import os
import pexpect
import pytest


@pytest.mark.skip
@pytest.mark.parametrize(
        "command, expected_output",
        [
            ("help", "cheqd App"),
            ("version", os.environ["RELEASE_NUMBER"]), # this works against deb package but not against starport build
        ]
    )
def test_basic(command, expected_output):
    cli = pexpect.spawn(f"cheqd-noded {command}", encoding="utf-8")
    cli.logfile = sys.stdout
    cli.expect(expected_output)

@pytest.mark.parametrize(
        "command, params, expected_output",
        [
            ("list", None, "- name: node5-operator"),
            ("add", "test4", "- name: test4"),
            ("list", None, "- name: test4"),
            ("delete", "test4 -y", "Key deleted forever"),
            ("show", "test", "- name: test"),
            ("show", "test4", "Error: test4 is not a valid name or address"),
        ]
    )
def test_keys(command, params, expected_output):
    cli = pexpect.spawn(f"cheqd-noded keys {command} {params}", encoding="utf-8")
    cli.logfile = sys.stdout
    cli.expect(expected_output)
