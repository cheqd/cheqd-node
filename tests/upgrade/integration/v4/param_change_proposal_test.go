//go:build upgrade_integration

package integration

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v4/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Fee parameter change proposal", func() {
	var Proposal_id string
	It("should wait for node catching up", func() {
		By("pinging the node status until catching up is flagged as false")
		err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
		Expect(err).To(BeNil())
	})

	It("should submit a parameter change proposal for did module (optimistic)", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join(GeneratedJSONDir, ProposalJSONDir, "existing", "param_change_did.json"), cli.DockerHome, cli.Validator0)
		Expect(err).To(BeNil())

		By("sending a SubmitParamChangeProposal transaction from `validator0` container")
		res, err := cli.SubmitProposal(cli.Validator0, "param_change_did.json")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
		res, err = cli.QueryTxn(cli.Validator0, res.TxHash)
		Expect(err).To(BeNil())
		proposal_id, err := cli.GetProposalID(res.Events)
		Expect(err).To(BeNil())
		Proposal_id = proposal_id
	})

	It("should wait for the proposal submission to be included in a block", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to be included in a block")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+2, cli.VotingPeriod*1)
		Expect(err).To(BeNil())
	})

	It("should vote for the parameter change proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposal(cli.Validator0, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposal(cli.Validator1, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposal(cli.Validator2, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator3` container", func() {
		By("sending a VoteProposal transaction from `validator3` container")
		res, err := cli.VoteProposal(cli.Validator3, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*1)
		Expect(err).To(BeNil())
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})

	It("should validate the param change result with the expected outcome", func() {
		By("sending a QueryParams query from `validator0` container")
		feeParams, err := cli.QueryDidFeeParams(cli.Validator0)
		Expect(err).To(BeNil())

		By("checking against the expected fee params")
		var expectedFeeParams didtypes.FeeParams
		_, err = Loader(filepath.Join(GeneratedJSONDir, ProposalJSONDir, "expected", "param_change_did.json"), &expectedFeeParams)
		Expect(err).To(BeNil())
		Expect(feeParams.Params).To(Equal(expectedFeeParams))
	})

	It("should submit a parameter change proposal for resource module (optimistic)", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join(GeneratedJSONDir, ProposalJSONDir, "existing", "param_change_resource.json"), cli.DockerHome, cli.Validator0)
		Expect(err).To(BeNil())

		By("sending a SubmitParamChangeProposal transaction from `validator0` container")
		res, err := cli.SubmitProposal(cli.Validator0, "param_change_resource.json")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		res, err = cli.QueryTxn(cli.Validator0, res.TxHash)
		Expect(err).To(BeNil())

		proposal_id, err := cli.GetProposalID(res.Events)
		Expect(err).To(BeNil())

		Proposal_id = proposal_id
	})

	It("should wait for the proposal submission to be included in a block", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to be included in a block")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+2, cli.VotingPeriod*1)
		Expect(err).To(BeNil())
	})

	It("should vote for the parameter change proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposal(cli.Validator0, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposal(cli.Validator1, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposal(cli.Validator2, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator3` container", func() {
		By("sending a VoteProposal transaction from `validator3` container")
		res, err := cli.VoteProposal(cli.Validator3, Proposal_id, "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, cli.VotingPeriod*2)
		Expect(err).To(BeNil())
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})

	It("should validate the param change result with the expected outcome", func() {
		By("sending a QueryParams query from `validator0` container")
		feeParams, err := cli.QueryResourceFeeParams(cli.Validator0)
		Expect(err).To(BeNil())

		By("checking against the expected fee params")
		var expectedFeeParams resourcetypes.FeeParams
		_, err = Loader(filepath.Join(GeneratedJSONDir, ProposalJSONDir, "expected", "param_change_resource.json"), &expectedFeeParams)
		Expect(err).To(BeNil())
		Expect(feeParams.Params).To(Equal(expectedFeeParams))
	})
})
