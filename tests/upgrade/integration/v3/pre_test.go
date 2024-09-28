package integration

import (
	"fmt"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v3/cli"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {

	Context("Before a software upgrade execution is initiated", func() {
		It("should wait for chain to bootstrap", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, cli.BootstrapHeight, cli.BootstrapPeriod)
			Expect(err).To(BeNil())
		})

		var UpgradeHeight int64
		var VotingEndHeight int64

		It("should calculate the upgrade height", func() {
			By("getting the current block height and calculating the voting end height")
			var err error
			UpgradeHeight, VotingEndHeight, err = cli.CalculateUpgradeHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())
			fmt.Printf("Upgrade height: %d\n", UpgradeHeight)
			fmt.Printf("Voting end height: %d\n", VotingEndHeight)
		})

		It("should submit a software upgrade proposal", func() {
			By("sending a SubmitUpgradeProposal transaction from `validator0` container")
			res, err := cli.SubmitUpgradeProposalLegacy(VotingEndHeight, cli.Validator0)
			Expect(err).To(BeNil())
			fmt.Println("response is>>>>>>.", res)
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should deposit tokens for the software upgrade proposal", func() {
			By("sending a DepositGov transaction from `validator0` container")
			res, err := cli.DepositGov(cli.Validator0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
			fmt.Println("response>>>>>>>>>>>>", res)
		})

		It("should wait for the proposal submission to be included in a block", func() {
			By("getting the current block height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for the proposal to be included in a block")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*3)
			Expect(err).To(BeNil())
		})

		It("should vote for the software upgrade proposal from `validator0` container", func() {
			By("sending a VoteProposal transaction from `validator0` container")
			res, err := cli.VoteProposal(cli.Validator0, "1", "yes")
			Expect(err).To(BeNil())
			fmt.Println("response>>>>>>>>>>>>>>", res)

			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})
		It("should vote for the software upgrade proposal from `validator1` container", func() {
			By("sending a VoteProposal transaction from `validator1` container")
			res, err := cli.VoteProposal(cli.Validator1, "1", "yes")
			Expect(err).To(BeNil())
			fmt.Println("response>>>>>>>>>>>>>>", res)

			Expect(res.Code).To(BeEquivalentTo(0))
		})
		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should vote for the software upgrade proposal from `validator2` container", func() {
			By("sending a VoteProposal transaction from `validator2` container")
			res, err := cli.VoteProposal(cli.Validator2, "1", "yes")
			Expect(err).To(BeNil())
			fmt.Println("response>>>>>>>>>>>>>>", res)
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})
		It("should vote for the software upgrade proposal from `validator3` container", func() {
			By("sending a VoteProposal transaction from `validator3` container")
			res, err := cli.VoteProposal(cli.Validator3, "1", "yes")
			Expect(err).To(BeNil())
			fmt.Println("response>>>>>>>>>>>>>>", res)

			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for the voting end height to be reached", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, VotingEndHeight, cli.VotingPeriod*10)
			Expect(err).To(BeNil())
		})

		It("should query the proposal status to ensure it has passed", func() {
			By("sending a QueryProposal Msg from `validator0` container")
			proposal, err := cli.QueryProposal(cli.Validator0, "1")
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
		})

		It("should wait for the upgrade height to be reached", func() {
			By("pinging the node status until the upgrade height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, UpgradeHeight, cli.VotingPeriod)
			Expect(err).To(BeNil())
		})
	})
})
