//go:build upgrade_integration

package integration

import (
	"fmt"
	"os"
	// "path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	// didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	// didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	// resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	Context("Before a softare upgrade execution is initiated", func() {

		BeforeEach(func() {
			cli.RUN_INSIDE_DOCKER = false
		})

		It("should wait for chain to bootstrap", func() {
			By("pinging the node status until the dvoting end height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, cli.BOOTSTRAP_HEIGHT, cli.BOOTSTRAP_PERIOD)
			Expect(err).To(BeNil())
		})

		var UPGRADE_HEIGHT int64
		var VOTING_END_HEIGHT int64
		var proposalID = os.Getenv("PROPOSAL_ID")

		if proposalID == "" {
			proposalID = "1"
		}

		fmt.Println("Proposal ID: ", proposalID)

		It("should calculate the upgrade height", func() {
			By("getting the current block height and calculating the voting end height")
			var err error
			UPGRADE_HEIGHT, VOTING_END_HEIGHT, err = cli.CalculateUpgradeHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME)
			Expect(err).To(BeNil())
			fmt.Printf("Upgrade height: %d\n", UPGRADE_HEIGHT)
			fmt.Printf("Voting end height: %d\n", VOTING_END_HEIGHT)
		})

		It("should submit a software upgrade proposal", func() {
			By("sending a SubmitUpgradeProposal transaction from `validator0` container")
			res, err := cli.SubmitUpgradeProposal(
				UPGRADE_HEIGHT, 
				cli.VALIDATOR0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should deposit tokens for the software upgrade proposal", func() {
			By("sending a DepositGov transaction from `validator0` container")
			res, err := cli.DepositGov(cli.VALIDATOR0, proposalID)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator0` container", func() {
			By("sending a VoteProposal transaction from `validator0` container")
			res, err := cli.VoteProposal(cli.VALIDATOR0, proposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator1` container", func() {
			By("sending a VoteProposal transaction from `validator1` container")
			res, err := cli.VoteProposal(cli.VALIDATOR1, proposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator2` container", func() {
			By("sending a VoteProposal transaction from `validator2` container")
			res, err := cli.VoteProposal(cli.VALIDATOR2, proposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator3` container", func() {
			By("sending a VoteProposal transaction from `validator3` container")
			res, err := cli.VoteProposal(cli.VALIDATOR3, proposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for the voting end height to be reached", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, VOTING_END_HEIGHT, cli.VOTING_PERIOD)
			Expect(err).To(BeNil())
		})

		It("should query the proposal status to ensure it has passed", func() {
			By("sending a QueryProposal Msg from `validator0` container")
			proposal, err := cli.QueryProposalLegacy(cli.VALIDATOR0, proposalID)
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1beta1.StatusPassed))
		})

		It("should wait for the upgrade height to be reached", func() {
			By("pinging the node status until the upgrade height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, UPGRADE_HEIGHT, cli.VOTING_PERIOD)
			Expect(err).To(BeNil())
		})
	})
})
