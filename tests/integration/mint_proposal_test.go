//go:build integration

package integration

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration - Mint coins to given address", func() {
	It("should wait for node catching up", func() {
		By("pinging the node status until catching up is flagged as false ")
		err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
		Expect(err).To(BeNil())
	})
	It("should submit a mint  proposal ", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join("proposal.json"), cli.DockerHome, cli.Validator0)
		Expect(err).To(BeNil())

		By("sending a SubmitParamChangeProposal transaction from `validator0` container")
		res, err := cli.SubmitProposal(cli.Validator0, cli.GasParams, "proposal.json")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal submission to be included in a block", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to be included in a block")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*3)
		Expect(err).To(BeNil())
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
	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, cli.VotingPeriod*3)
		Expect(err).To(BeNil())
	})
	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, "1")
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})
})