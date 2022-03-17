import json

import pytest

from helpers import run, LOCAL_SENDER_ADDRESS, LOCAL_RECEIVER_ADDRESS, LOCAL_NET_DESTINATION, GAS_PRICE, YES_FLAG, \
    KEYRING_BACKEND_TEST, get_gas_extimation, CODE_0, TEST_NET_GAS_X_GAS_PRICES, generate_ed25519_key, random_string, \
    build_create_did_msg, json_loads, build_update_did_msg, CODE_1100, CODE_1203, get_balance, GAS_AMOUNT, CODE_5, \
    CODE_11


@pytest.mark.parametrize("magic_number_positive", [1.3, 2, 3, 10])
@pytest.mark.parametrize("transfer_amount", ["1", "20002", "3000003", "9000000009"])
def test_gas_estimation_positive(magic_number_positive, transfer_amount):
    # Get the gas_wanted value and use it as gas_estimation 
    cli = run(
        "cheqd-noded tx", 
        "bank send", 
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {transfer_amount}ncheq {LOCAL_NET_DESTINATION} --gas auto --gas-prices {GAS_PRICE}ncheq --gas-adjustment {magic_number_positive} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        r"")
    # Previous command returns the string, like:
    # gas estimate: 123456
    # at the beginning and it will be caught by this function
    gas_estimation = get_gas_extimation(str(cli.read()))

    gas = int(gas_estimation * magic_number_positive)
    # Send the same request for checking that gas was calculated in a right way
    # The main difference between previous run is that `--gas` parameter is set with particular value
    run(
        "cheqd-noded tx",
        "bank send",
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {transfer_amount}ncheq {LOCAL_NET_DESTINATION} --gas {gas} --gas-prices {GAS_PRICE}ncheq --gas-adjustment {magic_number_positive} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_0}(.*?)\"value\":\"{transfer_amount}ncheq\"")


@pytest.mark.parametrize("magic_number_negative", [1, 0.5, 0.1])
@pytest.mark.parametrize("transfer_amount", ["1", "20002", "3000003", "9000000009"])
def test_gas_estimation_negative(magic_number_negative, transfer_amount):
    # Get the gas_wanted value and use it as gas_estimation
    cli = run(
        "cheqd-noded tx",
        "bank send",
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {transfer_amount}ncheq {LOCAL_NET_DESTINATION} --gas auto --gas-prices {GAS_PRICE}ncheq --gas-adjustment {magic_number_negative} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        r"")

    # Previous command returns the string, like:
    # gas estimate: 123456
    # at the beginning and it will be caught by this function
    gas_estimation = get_gas_extimation(str(cli.read()))

    gas = int(gas_estimation * magic_number_negative)
    # Send the same request for checking that gas was calculated in a right way
    # The main difference between previous run is that `--gas` parameter is set with particular value
    run(
        "cheqd-noded tx",
        "bank send",
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {transfer_amount}ncheq {LOCAL_NET_DESTINATION} --gas {gas} --gas-prices {GAS_PRICE}ncheq --gas-adjustment {magic_number_negative} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_11}(.*?)\"raw_log\":\"out of gas in location")


@pytest.mark.parametrize("transfer_amount", [1111111111111111111111, 55555555555555555555555, 999999999999999999999999, 10000000000000000000000000000]) # TODO: hypothesis
def test_token_transfer_negative(transfer_amount):
    # Get balances before sending tx
    sender_balance = int(get_balance(LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION))
    receiver_balance = int(get_balance(LOCAL_RECEIVER_ADDRESS, LOCAL_NET_DESTINATION))

    # Send token transfer
    cli = run(
        "cheqd-noded tx",
        "bank send",
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} {transfer_amount}ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_5}(.*?)\"raw_log\":\"(.*?) is smaller than {transfer_amount}ncheq: insufficient funds")

    # Get balances after
    sender_balance_after = int(get_balance(LOCAL_SENDER_ADDRESS, LOCAL_NET_DESTINATION))
    receiver_balance_after = int(get_balance(LOCAL_RECEIVER_ADDRESS, LOCAL_NET_DESTINATION))

    # Compare them ad make sure that only fee amount from sender was burned
    assert sender_balance - sender_balance_after == GAS_AMOUNT * GAS_PRICE
    assert receiver_balance - receiver_balance_after == 0


@pytest.mark.parametrize("note", ["a", "1", "test_memo_test", "123qwe$%^", "______________________________"]) # TODO: hypothesis
def test_transfer_memo(note):
    # Send token transaction with memo set up
    cli = run(
        "cheqd-noded tx",
        "bank send",
        f"{LOCAL_SENDER_ADDRESS} {LOCAL_RECEIVER_ADDRESS} 1ncheq {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST} --note {note}",
        fr"")

    # Get transaction hash
    result_str = cli.read()
    tx_hash = json_loads(result_str)["txhash"]

    # Get transaction from the pool by hash
    cli = run(
        "cheqd-noded query",
        "tx",
        fr"{tx_hash} --output json",
        "")

    # Get memo from
    memo_json = json_loads(cli.read())

    # Compare it with what was sent
    assert note == memo_json["tx"]["body"]["memo"]


def test_did_query_non_existent():
    did = fr"did:cheqd:testnet:AbCdEfGh"

    # Try to get did
    run(
        "cheqd-noded query",
        "cheqd did",
        fr"{did}",
        fr"Error: rpc error: code = InvalidArgument desc = not found: invalid request")


def test_did_wrong_version_update():
    did = fr"did:cheqd:testnet:{random_string(5)}"
    key_id = fr"{did}#key1"

    # Generate ed25519 key
    ed25519_key = generate_ed25519_key()
    pub_key_base_64 = ed25519_key["pub_key_base_64"]
    priv_key_base_64 = ed25519_key["priv_key_base_64"]

    # Get multibase58 represantation
    cli = run(
        "cheqd-noded debug",
        "encoding base64-multibase58",
        fr"{pub_key_base_64}",
        "")

    ver_pub_multibase_58 = cli.read().strip()

    # Send request to create a DID
    msg_create_did = build_create_did_msg(did,
                                          key_id,
                                          ver_pub_multibase_58)

    run(
        "cheqd-noded tx",
        "cheqd create-did",
        f" '{msg_create_did}' {key_id} {priv_key_base_64} --from {LOCAL_SENDER_ADDRESS} {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_0}")

    # Get the created DID for getting version_id
    cli = run(
        "cheqd-noded query",
        "cheqd did",
        fr"{did} --output json",
        "")

    did_json = json_loads(cli.read())
    version_id_orig = did_json["metadata"]["version_id"]

    # Change version_id to the wrong one
    wrong_version_id = version_id_orig + "abc"

    # Prepare and send update did message for getting an error
    msg_update_did = build_update_did_msg(did,
                                          key_id,
                                          ver_pub_multibase_58,
                                          wrong_version_id)
    msg_update_did["capability_delegation"] = [key_id]

    # here we are expecting an 1203 error about wrong version_id
    run(
        "cheqd-noded tx",
        "cheqd update-did",
        f" '{json.dumps(msg_update_did)}' {key_id} {priv_key_base_64} --from {LOCAL_SENDER_ADDRESS} {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_1203}(.*?)\"raw_log\":\"(.*?)Expected(.*?)unexpected DID version")


def test_did_wrong_verkey_update():
    did = fr"did:cheqd:testnet:{random_string(5)}"
    key_id = fr"{did}#key1"

    # Generate ed25519 key
    ed25519_key = generate_ed25519_key()
    pub_key_base_64 = ed25519_key["pub_key_base_64"]
    priv_key_base_64 = ed25519_key["priv_key_base_64"]

    # Get multibase58 represantation
    cli = run(
        "cheqd-noded debug",
        "encoding base64-multibase58",
        fr"{pub_key_base_64}",
        "")

    ver_pub_multibase_58 = cli.read().strip()

    # Send request to create a DID
    msg_create_did = build_create_did_msg(did,
                                          key_id,
                                          ver_pub_multibase_58)

    run(
        "cheqd-noded tx",
        "cheqd create-did",
        f" '{msg_create_did}' {key_id} {priv_key_base_64} --from {LOCAL_SENDER_ADDRESS} {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_0}")

    # Get the created DID for getting version_id
    cli = run(
        "cheqd-noded query",
        "cheqd did",
        fr"{did} --output json",
        "")

    did_json = json_loads(cli.read())
    version_id = did_json["metadata"]["version_id"]

    # Prepare and send update did message for getting an error
    msg_update_did = build_update_did_msg(did,
                                          key_id,
                                          ver_pub_multibase_58 + "abc",
                                          version_id)

    msg_update_did["capability_delegation"] = [key_id]

    # Create another ed25519 key for using the new one for signing
    new_priv_key_base_64 = generate_ed25519_key()["priv_key_base_64"]

    # here we are expecting an 1203 error about wrong version_id
    run(
        "cheqd-noded tx",
        "cheqd update-did",
        f" '{json.dumps(msg_update_did)}' {key_id} {new_priv_key_base_64} --from {LOCAL_SENDER_ADDRESS} {LOCAL_NET_DESTINATION} {TEST_NET_GAS_X_GAS_PRICES} {YES_FLAG} {KEYRING_BACKEND_TEST}",
        fr"{CODE_1100}(.*?)\"raw_log\":\"(.*?)invalid signature detected")
