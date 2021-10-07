import sys
import os
import pexpect
import pytest


@pytest.mark.parametrize(
        "command, expected_output",
        [
            ("help", "cheqd App"),
            ("version", os.environ["RELEASE_NUMBER"]),
        ]
    )
def test_basic(command, expected_output):
    pexpect.spawn(f"cheqd-noded {command}", encoding="utf-8").expect(expected_output)
