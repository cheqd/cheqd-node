//go:build upgrade

package upgrade

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Fee parameter change proposal", func() {
	It("should wait for node catching up", func() {
		By("pinging the node status until catching up is flagged as false")
		err := cli.WaitForCaughtUp(cli.VALIDATOR0, cli.CLI_BINARY_NAME, cli.VOTING_PERIOD*6)
		Expect(err).To(BeNil())
	})

	It("should submit a parameter change proposal for did module", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join(GENERATED_JSON_DIR, "proposal", "existing", "param_change_did.json"), cli.DOCKER_HOME, cli.VALIDATOR0)
		Expect(err).To(BeNil())

		By("sending a SubmitParamChangeProposal transaction from `validator0` container")
		res, err := cli.SubmitParamChangeProposal(cli.VALIDATOR0, "param_change_did.json")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposal(cli.VALIDATOR0, "2", "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposal(cli.VALIDATOR1, "2", "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the parameter change proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposal(cli.VALIDATOR2, "2", "yes")
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, currentHeight+20, cli.VOTING_PERIOD*3)
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposalLegacy(cli.VALIDATOR0, "2")
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1beta1.StatusPassed))
	})

	It("should validate the param change result with the expected outcome", func() {
		By("sending a QueryParams query from `validator0` container")
		feeParams, err := cli.QueryDidFeeParams(cli.VALIDATOR0, didtypes.ModuleName, string(didtypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())

		By("checking against the expected fee params")
		var expectedFeeParams didtypes.FeeParams
		err = Loader(filepath.Join(GENERATED_JSON_DIR, "expected", "param_change_did.json"), &expectedFeeParams)
		Expect(err).To(BeNil())
		Expect(feeParams).To(Equal(expectedFeeParams))
	})
})
