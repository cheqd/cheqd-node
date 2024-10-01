package cli

import (
	upgradeV3 "github.com/cheqd/cheqd-node/app/upgrades/v3"
	integrationnetwork "github.com/cheqd/cheqd-node/tests/integration/network"
)

const CliBinaryName = "cheqd-noded"

const (
	KeyringBackend = "test"
	OutputFormat   = "json"
	Gas            = "auto"
	GasAdjustment  = "2.5"
	GasPrices      = "60ncheq"
)

const (
	Green  = "\033[32m"
	Purple = "\033[35m"
)

const (
	BootstrapPeriod            = 20
	BootstrapHeight            = 1
	VotingPeriod         int64 = 10
	ExpectedBlockSeconds int64 = 1
	ExtraBlocks          int64 = 10
	UpgradeName                = upgradeV3.UpgradeName
	DepositAmount              = "10000000ncheq"
	NetworkConfigDir           = "network-config"
	KeyringDir                 = "keyring-test"
)

var (
	TXParams = []string{
		"--keyring-backend", KeyringBackend,
		"--chain-id", integrationnetwork.ChainID,
		"-y",
	}
	GasParams = []string{
		"--gas", Gas,
		"--gas-adjustment", GasAdjustment,
		"--gas-prices", GasPrices,
	}
	QueryParamsConst = []string{
		"--chain-id", integrationnetwork.ChainID,
		"--output", OutputFormat,
	}
	KeysParams = []string{
		"--keyring-backend", KeyringBackend,
	}
)
