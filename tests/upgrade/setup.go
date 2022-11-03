//go:build upgrade

package upgrade

import (
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Setup helpers for upgrade tests
func Setup() error {
	By("waiting for chain to start")
	err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, 1, 60)
	Expect(err).To(BeNil())

	By("Copying keys to each container volume")
	_, err = cli.LocalnetExecCopyKeys()
	Expect(err).To(BeNil())

	return nil
}
