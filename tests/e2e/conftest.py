import pytest
import re
from helpers import run, run_interaction, random_string, \
    TEST_NET_DESTINATION, TEST_NET_FEES, TEST_NET_GAS_X_GAS_PRICES, YES_FLAG, \
    LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION, CODE_0


@pytest.fixture(scope='session')
def create_export_keys():
    command_base = "cheqd-noded keys"
    run(command_base, "add", "export_key", "name: export_key")
