package cli

import (
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	cheqdapp "github.com/cheqd/cheqd-node/app"
)

func SubmitUpgradeProposal(from string) (string, error) {
	return Exec(
		integrationcli.CLI_BINARY_NAME,
		"tx",
		"gov",
		"submit-proposal",
		"upgrade",
		cheqdapp.UpgradeName,
		"--from",
		from,
		"--deposit",
		"1000000ncheq",
		"--upgrade-height",
		"1",
		"--upgrade-info",
		"cosmovisor_test",
		"--chain-id",
		"cheqd",
		"--keyring-backend",
		"test",
		"--home",
		"tests/integration/cli/home/alice",
	)
}