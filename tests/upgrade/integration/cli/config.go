package cli

import (
	cheqdapp "github.com/cheqd/cheqd-node/app"
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	integrationnetwork "github.com/cheqd/cheqd-node/tests/integration/network"
)

const (
	CLIBinaryName = integrationcli.CliBinaryName
	Green         = integrationcli.Green
	Purple        = integrationcli.Purple
)

const (
	CLIBinaryNamePrevious = CLIBinaryName + "-previous"
	CLIBinaryNameNext     = CLIBinaryName + "-next"
)

const (
	KeyringBackend = integrationcli.KeyringBackend
	OutputFormat   = integrationcli.OutputFormat
	Gas            = integrationcli.Gas
	GasAdjustment  = integrationcli.GasAdjustment
	GasPrices      = integrationcli.GasPrices

	BootstrapPeriod            = 20
	BootstrapHeight            = 1
	VotingPeriod         int64 = 10
	ExpectedBlockSeconds int64 = 1
	ExtraBlocks          int64 = 5
	UpgradeName                = cheqdapp.UpgradeName
	DepositAmount              = "10000000ncheq"
	NetworkConfigDir           = "network-config"
	KeyringDir                 = "keyring-test"
)

var (
	TXParams = []string{
		"--gas", Gas,
		"--gas-adjustment", GasAdjustment,
		"--gas-prices", GasPrices,
		"--keyring-backend", KeyringBackend,
		"--chain-id", integrationnetwork.ChainID,
		"-y",
	}
	QueryParamsConst = []string{
		"--chain-id", integrationnetwork.ChainID,
		"--output", OutputFormat,
	}
)
