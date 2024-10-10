//go:build upgrade_integration

package integration

import (
	"fmt"
	"path/filepath"

	clihelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v3/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Post", func() {
	var feeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams

	BeforeEach(func() {
		// query fee params - case: did
		res, err := cli.QueryParams(cli.Validator0, didtypes.ModuleName, string(didtypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())

		// query fee params - case: resource
		res, err = cli.QueryParams(cli.Validator0, resourcetypes.ModuleName, string(resourcetypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &resourceFeeParams)
		Expect(err).To(BeNil())
	})

	Context("After a software upgrade execution has concluded", func() {
		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should match the expected module version map", func() {
			By("loading the expected module version map")
			var expected upgradetypes.QueryModuleVersionsResponse
			_, err := Loader(filepath.Join(GeneratedJSONDir, "post", "query - module-version-map", "v3.json"), &expected)
			Expect(err).To(BeNil())

			By("matching the expected module version map")
			actual, err := cli.QueryModuleVersionMap(cli.Validator0)
			Expect(err).To(BeNil())

			Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		})

		It("should load and run expected diddoc payloads", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocExistingRecords, err := RelGlob(GeneratedJSONDir, "post", "query - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocExistingRecords {
				var DidDocExistingRecord didtypes.DidDoc

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				fmt.Println("Running: " + testCase)

				_, err = Loader(payload, &DidDocExistingRecord)
				Expect(err).To(BeNil())

				res, err := cli.QueryDid(DidDocExistingRecord.Id, cli.Validator0)
				Expect(err).To(BeNil())

				if DidDocExistingRecord.Context == nil {
					DidDocExistingRecord.Context = []string{}
				}
				if DidDocExistingRecord.Authentication == nil {
					DidDocExistingRecord.Authentication = []string{}
				}
				if DidDocExistingRecord.AssertionMethod == nil {
					DidDocExistingRecord.AssertionMethod = []string{}
				}
				if DidDocExistingRecord.CapabilityInvocation == nil {
					DidDocExistingRecord.CapabilityInvocation = []string{}
				}
				if DidDocExistingRecord.CapabilityDelegation == nil {
					DidDocExistingRecord.CapabilityDelegation = []string{}
				}
				if DidDocExistingRecord.KeyAgreement == nil {
					DidDocExistingRecord.KeyAgreement = []string{}
				}
				if DidDocExistingRecord.Service == nil {
					DidDocExistingRecord.Service = []*didtypes.Service{}
				}
				if DidDocExistingRecord.AlsoKnownAs == nil {
					DidDocExistingRecord.AlsoKnownAs = []string{}
				}

				Expect(*res.Value.DidDoc).To(Equal(DidDocExistingRecord))
			}
		})

		It("should load and run expected resource payloads", func() {
			By("matching the glob pattern for existing resource payloads")
			ExpectedResourceCreateRecords, err := RelGlob(GeneratedJSONDir, "post", "query - resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedResourceCreateRecords {
				var ResourceCreateRecord resourcetypes.ResourceWithMetadata

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				fmt.Println("Running: " + testCase)

				_, err = Loader(payload, &ResourceCreateRecord)
				Expect(err).To(BeNil())

				res, err := cli.QueryResource(ResourceCreateRecord.Metadata.CollectionId, ResourceCreateRecord.Metadata.Id, cli.Validator0)

				Expect(err).To(BeNil())
				Expect(res.Resource.Metadata.Id).To(Equal(ResourceCreateRecord.Metadata.Id))
				Expect(res.Resource.Metadata.CollectionId).To(Equal(ResourceCreateRecord.Metadata.CollectionId))
				Expect(res.Resource.Metadata.Name).To(Equal(ResourceCreateRecord.Metadata.Name))
				Expect(res.Resource.Metadata.Version).To(Equal(ResourceCreateRecord.Metadata.Version))
				Expect(res.Resource.Metadata.ResourceType).To(Equal(ResourceCreateRecord.Metadata.ResourceType))
				Expect(res.Resource.Metadata.AlsoKnownAs).To(Equal(ResourceCreateRecord.Metadata.AlsoKnownAs))
				Expect(res.Resource.Metadata.MediaType).To(Equal(ResourceCreateRecord.Metadata.MediaType))
				Expect(res.Resource.Metadata.Checksum).To(Equal(ResourceCreateRecord.Metadata.Checksum))
				Expect(res.Resource.Metadata.PreviousVersionId).To(Equal(ResourceCreateRecord.Metadata.PreviousVersionId))
				Expect(res.Resource.Metadata.NextVersionId).To(Equal(ResourceCreateRecord.Metadata.NextVersionId))
			}
		})
	})
	var Proposal_id string
	Context("After a software upgrade execution has concluded", func() {
		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})
		It("should submit a parameter change proposal for did module (optimistic)", func() {
			By("passing the proposal file to the container")
			_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join("proposal.json"), cli.DockerHome, cli.Validator0)
			Expect(err).To(BeNil())
			By("sending a SubmitParamChangeProposal transaction from `validator0` container")
			res, err := cli.SubmitProposal(cli.Validator0, "proposal.json")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
			res, err = cli.QueryTxn(cli.Validator0, res.TxHash)
			Expect(err).To(BeNil())

			proposal_id, err := cli.GetProposalID(res.RawLog)
			Proposal_id = proposal_id
			Expect(err).To(BeNil())
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
		It("should vote for the parameter change proposal from `validator3` container", func() {
			By("sending a VoteProposal transaction from `validator3` container")
			res, err := cli.VoteProposal(cli.Validator0, Proposal_id, "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})
		It("should wait for the proposal to pass", func() {
			By("getting the current block height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for the proposal to pass")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, cli.VotingPeriod*100)
			Expect(err).To(BeNil())
		})

		It("should check the proposal status to ensure it has passed", func() {
			By("sending a QueryProposal query from `validator0` container")
			proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
		})

		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})
		It("should burn the coins from the given address (here container/validator)", func() {
			coins := sdk.NewCoins(sdk.Coin{Denom: "ncheq", Amount: sdk.NewInt(1000)})
			res, err := cli.BurnMsg(cli.Validator0, coins.String())
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})
	})
})
