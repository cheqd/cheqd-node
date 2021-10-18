import sys
import pexpect


IMPLICIT_TIMEOUT = 30
ENCODING = "utf-8"
READ_BUFFER = 6000

TEST_NET_NETWORK = "cheqd-testnet-2"
TEST_NET_NODE_TCP = "--node 'tcp://seed1.us.testnet.cheqd.network:26657'"
TEST_NET_NODE_HTTP = "--node http://node1.eu.testnet.cheqd.network:26657/"
TEST_NET_DESTINATION = f"{TEST_NET_NODE_TCP} --chain-id 'cheqd-testnet-2'"
TEST_NET_DESTINATION_HTTP = f"{TEST_NET_NODE_HTTP} --chain-id 'cheqd-testnet-2'"
TEST_NET_FEES = "--fees 5000000ncheq"
TEST_NET_GAS_X_GAS_PRICES = "--gas 70000 --gas-prices 25ncheq"
YES_FLAG = "-y"

SENDER_ADDRESS = "cheqd1ece09txhq6nm9fkft9jh3mce6e48ftescs5jsw"
SENDER_MNEMONIC = "oil long siege student rent jar awkward park entry ripple enable company sort people little damp arrange wise slender push brief solve tattoo cycle"
RECEIVER_ADDRESS= "cheqd16d72a6kusmzml5mjhzjv63c9j5xnpsyqs8f3sk"
RECEIVER_MNEMONIC = "strike impact earth indoor man illness virus genuine rib control antenna loop neck rotate bargain original nasty size either try snap quiz stairs huge"

CODE_0 = "\"code\":0"


def run(command_base, command, params, expected_output):
    cli = pexpect.spawn(f"{command_base} {command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)
    print(f"BEFORE >>> {cli.before}") # FIXME DEBUG
    print(f"AFTER >>> {cli.after}") # FIXME DEBUG
    return cli


def run_interaction(cli, input_string, expected_output):
    cli.sendline(input_string)
    cli.expect(expected_output)
