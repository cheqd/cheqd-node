import pytest
from helpers import run, KEYRING_BACKEND_TEST


@pytest.fixture(scope='session')
def create_export_keys():
    command_base = "cheqd-noded keys"
    run(command_base, "add", f"export_key {KEYRING_BACKEND_TEST}", "name: export_key")
