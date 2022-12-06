package cli

import (
	cheqdapp "github.com/cheqd/cheqd-node/app"
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	integrationnetwork "github.com/cheqd/cheqd-node/tests/integration/network"
)

const (
	CLI_BINARY_NAME = integrationcli.CLI_BINARY_NAME
	GREEN           = integrationcli.GREEN
	PURPLE          = integrationcli.PURPLE
)

const (
	CLI_BINARY_NAME_PREVIOUS = CLI_BINARY_NAME + "-previous"
	CLI_BINARY_NAME_NEXT     = CLI_BINARY_NAME + "-next"
)

const (
	KEYRING_BACKEND = integrationcli.KEYRING_BACKEND
	OUTPUT_FORMAT   = integrationcli.OUTPUT_FORMAT
	GAS             = integrationcli.GAS
	GAS_ADJUSTMENT  = integrationcli.GAS_ADJUSTMENT
	GAS_PRICES      = integrationcli.GAS_PRICES

	CHEQD_IMAGE_FROM             = "cheqd/cheqd-node:latest"
	CHEQD_TAG_FROM               = "v0.6.9"
	CHEQD_IMAGE_TO               = "cheqd/cheqd-node:production-latest"
	CHEQD_TAG_TO                 = "v1.0.0"
	BOOTSTRAP_PERIOD             = 20
	BOOTSTRAP_HEIGHT             = 1
	VOTING_PERIOD          int64 = 10
	EXPECTED_BLOCK_SECONDS int64 = 1
	EXTRA_BLOCKS           int64 = 5
	UPGRADE_NAME                 = cheqdapp.UpgradeName
	DEPOSIT_AMOUNT               = "10000000ncheq"
	NETWORK_CONFIG_DIR           = "network-config"
	KEYRING_DIR                  = "keyring-test"
)

var (
	TX_PARAMS = []string{
		"--gas", GAS,
		"--gas-adjustment", GAS_ADJUSTMENT,
		"--gas-prices", GAS_PRICES,
		"--keyring-backend", KEYRING_BACKEND,
		"--chain-id", integrationnetwork.CHAIN_ID,
		"-y",
	}
	QUERY_PARAMS = []string{
		"--chain-id", integrationnetwork.CHAIN_ID,
		"--output", OUTPUT_FORMAT,
	}
)
