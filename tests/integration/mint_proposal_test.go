package integration

import (
	"time"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration - Mint coins to given address", func() {
	It("should submit a mint  proposal ", func() {
		By("sending a SubmitParamChangeProposal transaction from `validator0` container")
		res, err := cli.SubmitProposal(cli.Validator0, cli.GasParams, "proposal.json")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposal(cli.Validator1, "1", "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposal(cli.Validator2, "1", "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator3` container", func() {
		By("sending a VoteProposal transaction from `validator3` container")
		res, err := cli.VoteProposal(cli.Validator3, "1", "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
	It("should vote for the mint proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposal(cli.Validator0, "1", "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
	It("wait until the voting period is done", func() {
		time.Sleep(time.Second * 40)
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, "1")
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})
})
