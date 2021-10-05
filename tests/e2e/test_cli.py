import sys
import pexpect
from pytest import *


def test_version():
    pexpect.spawn("cheqd-noded version", encoding="utf-8").expect("0.2.2")
