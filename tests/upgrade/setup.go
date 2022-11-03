//go:build upgrade

package upgrade

import (
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Setup helpers for upgrade tests
func Setup() error {
	By("Copying keys to each container volume")
	_, err := cli.LocalnetExecCopyKeys()
	Expect(err).To(BeNil())

	return nil
}
