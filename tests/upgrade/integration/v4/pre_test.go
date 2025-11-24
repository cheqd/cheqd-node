//go:build upgrade_integration

package integration

import (
	"fmt"
	"path/filepath"

	clihelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v3/cli"
	cliv4 "github.com/cheqd/cheqd-node/tests/upgrade/integration/v4/cli"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	var didFeeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams
	ProposalID := "1"

	BeforeEach(func() {
		// query fee params - case: did
		res, err := cli.QueryParams(cli.Validator0, didtypes.ModuleName, string(didtypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &didFeeParams)
		Expect(err).To(BeNil())

		// query fee params - case: resource
		res, err = cli.QueryParams(cli.Validator0, resourcetypes.ModuleName, string(resourcetypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &resourceFeeParams)
		Expect(err).To(BeNil())
	})

	Context("Before a software upgrade execution is initiated", func() {
		It("should wait for chain to bootstrap", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, cli.BootstrapHeight, cli.BootstrapPeriod+5)
			Expect(err).To(BeNil())
		})

		It("should load and run existing diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocCreatePayloads, err := RelGlob(GeneratedJSONDir, "pre", "create - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExistingDidDocCreatePayloads {
				var DidDocCreatePayload didcli.DIDDocument
				var DidDocCreateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				By("reading ")
				DidDocCreateSignInput, err = Loader(payload, &DidDocCreatePayload)
				Expect(err).To(BeNil())

				res, err := cliv4.CreateDid(DidDocCreatePayload, DidDocCreateSignInput, cli.Validator0, "", didFeeParams.CreateDid.String())
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))

				By("waiting for an additional set of blocks to be produced")
				height, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
				Expect(err).To(BeNil())
				err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, height+5, cli.VotingPeriod*2)
				Expect(err).To(BeNil())
			}
		})

		It("should load and run existing resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExistingResourceCreatePayloads, err := RelGlob(GeneratedJSONDir, "pre", "create - resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExistingResourceCreatePayloads {
				var ResourceCreatePayload resourcetypes.MsgCreateResourcePayload
				var ResourceCreateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				ResourceCreateSignInput, err = Loader(payload, &ResourceCreatePayload)
				Expect(err).To(BeNil())

				By("copying the existing resource file to the container")
				ResourceFile, err := CreateTestJSON(GinkgoT().TempDir(), ResourceCreatePayload.Data)
				Expect(err).To(BeNil())
				_, err = cli.LocalnetExecCopyAbsoluteWithPermissions(ResourceFile, cli.DockerHome, cli.Validator0)
				Expect(err).To(BeNil())

				res, err := cliv4.CreateResource(
					ResourceCreatePayload,
					filepath.Base(ResourceFile),
					ResourceCreateSignInput,
					cli.Validator0,
					resourceFeeParams.Json.String(),
				)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))

				By("waiting for an additional set of blocks to be produced")
				height, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
				Expect(err).To(BeNil())
				err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, height+5, cli.VotingPeriod*2)
				Expect(err).To(BeNil())
			}
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
			res, err := cliv4.SubmitUpgradeProposalLegacy(VotingEndHeight, cli.Validator0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*1)
			Expect(err).To(BeNil())
		})

		It("should deposit tokens for the software upgrade proposal", func() {
			By("sending a DepositGov transaction from `validator0` container")
			res, err := cli.DepositGov(cli.Validator0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for the proposal submission to be included in a block", func() {
			By("getting the current block height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for the proposal to be included in a block")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+3, cli.VotingPeriod*1)
			Expect(err).To(BeNil())
		})

		It("should vote for the software upgrade proposal from `validator0` container", func() {
			By("sending a VoteProposal transaction from `validator0` container")
			res, err := cli.VoteProposal(cli.Validator0, ProposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, cli.VotingPeriod*1)
			Expect(err).To(BeNil())
		})

		It("should vote for the software upgrade proposal from `validator1` container", func() {
			By("sending a VoteProposal transaction from `validator1` container")
			res, err := cli.VoteProposal(cli.Validator1, ProposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, cli.VotingPeriod*1)
			Expect(err).To(BeNil())
		})

		It("should vote for the software upgrade proposal from `validator2` container", func() {
			By("sending a VoteProposal transaction from `validator2` container")
			res, err := cli.VoteProposal(cli.Validator2, ProposalID, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, cli.VotingPeriod*1)
			Expect(err).To(BeNil())
		})

		It("should wait for the voting end height to be reached", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, VotingEndHeight, cli.VotingPeriod*10)
			Expect(err).To(BeNil())
		})

		It("should query the proposal status to ensure it has passed", func() {
			By("sending a QueryProposal Msg from `validator0` container")
			proposal, err := cli.QueryProposal(cli.Validator0, ProposalID)
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
		})

		It("should wait for the upgrade height to be reached", func() {
			By("pinging the node status until the upgrade height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, UpgradeHeight, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})
	})
})
