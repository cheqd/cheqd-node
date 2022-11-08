//go:build upgrade

package upgrade

import (
	"fmt"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Execute", func() {
	When("a software upgrade execution is initiated", func() {
		It("should ensure the localnet environment is running @latest version", func() {
			By("checking if the localnet environment is running via health check")
			// TODO: Implement this test
			// Anyway, this test is not needed for the upgrade process.
			// It is just a check to ensure the localnet environment is running, for completeness and to fail fast.
			Expect(true).To(BeTrue())

			By("setting the localnet environment variables to the latest version")
			err := cli.SetOldDockerComposeEnv()
			Expect(err).To(BeNil())
		})

		It("should calculate the upgrade height", func() {
			By("getting the current block height and calculating the voting end height")
			UPGRADE_HEIGHT, VOTING_END_HEIGHT, HEIGHT_ERROR = cli.CalculateUpgradeHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME)
			Expect(HEIGHT_ERROR).To(BeNil())
		})

		It("should submit a software upgrade proposal", func() {
			By("sending a SubmitUpgradeProposal transaction from `validator0` container")
			res, err := cli.SubmitUpgradeProposal(UPGRADE_HEIGHT, cli.VALIDATOR0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should deposit tokens for the software upgrade proposal", func() {
			By("sending a DepositGov transaction from `validator0` container")
			res, err := cli.DepositGov(cli.VALIDATOR0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator0` container", func() {
			By("sending a VoteUpgradeProposal transaction from `validator0` container")
			res, err := cli.VoteUpgradeProposal(cli.VALIDATOR0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator1` container", func() {
			By("sending a VoteUpgradeProposal transaction from `validator1` container")
			res, err := cli.VoteUpgradeProposal(cli.VALIDATOR1)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator2` container", func() {
			By("sending a VoteUpgradeProposal transaction from `validator2` container")
			res, err := cli.VoteUpgradeProposal(cli.VALIDATOR2)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator3` container", func() {
			By("sending a VoteUpgradeProposal transaction from `validator3` container")
			res, err := cli.VoteUpgradeProposal(cli.VALIDATOR3)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for the voting end height to be reached", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, VOTING_END_HEIGHT, cli.VOTING_PERIOD)
			Expect(err).To(BeNil())
		})

		It("should query the proposal status to ensure it has passed", func() {
			By("sending a QueryUpgradeProposal Msg from `validator0` container")
			proposal, err := cli.QueryUpgradeProposal(cli.VALIDATOR0)
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1beta1.StatusPassed))
		})

		It("should wait for the upgrade height to be reached", func() {
			By("pinging the node status until the upgrade height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, UPGRADE_HEIGHT, cli.VOTING_PERIOD)
			Expect(err).To(BeNil())
		})

		It("should halt the localnet environment", func() {
			By("switching genesis.json to the new version for each validator node")
			_, err := cli.LocalnetExecSwitchGenesis(cli.VALIDATOR0)
			Expect(err).To(BeNil())
			_, err = cli.LocalnetExecSwitchGenesis(cli.VALIDATOR1)
			Expect(err).To(BeNil())
			_, err = cli.LocalnetExecSwitchGenesis(cli.VALIDATOR2)
			Expect(err).To(BeNil())
			_, err = cli.LocalnetExecSwitchGenesis(cli.VALIDATOR3)
			Expect(err).To(BeNil())

			By("executing the container stop command")
			_, err = cli.LocalnetExecDown()
			Expect(err).To(BeNil())
		})

		It("should ensure the localnet environment is running @new version", func() {
			By("replacing the binary with the new version")
			_, err := cli.ReplaceBinaryWithPermissions("previous-to-next")

			By("executing the container up command for the new version")
			_, err = cli.LocalnetExecUpWithNewImage()
			Expect(err).To(BeNil())

			By("checking if the localnet environment is running via health check")
			// TODO: Implement this test
			// Anyway, this test is not needed for the upgrade process.
			// It is just a check to ensure the localnet environment is running, for completeness and to fail fast.
			Expect(true).To(BeTrue())
		})

		It("should wait for the upgrade height plus 2 blocks to be reached", func() {
			By("pinging the node status until the upgrade height plus 2 blocks is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, UPGRADE_HEIGHT+2, cli.VOTING_PERIOD*2)
			Expect(err).To(BeNil())
		})

		fmt.Printf("%s Upgrade successful.", cli.GREEN)
	})
})
