//go:build upgrade_integration

package integration

import (
	"fmt"
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types" <-- TODO: uncomment when whole sequence of upgrade tests is ready
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	Context("Before a softare upgrade execution is initiated", func() {
		It("should wait for chain to bootstrap", func() {
			By("pinging the node status until the dvoting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CLIBinaryName, cli.BootstrapHeight, cli.BootstrapPeriod)
			Expect(err).To(BeNil())
		})

		// TODO: uncomment when whole sequence of upgrade tests is ready
		// It("should match the expected module version map", func() {
		// 	By("loading the expected module version map")
		// 	var expected upgradetypes.QueryModuleVersionsResponse
		// 	_, err := Loader(filepath.Join(GeneratedJSONDir, "pre", "query - module-version-map", "v069.json"), &expected)
		// 	Expect(err).To(BeNil())

		// 	By("matching the expected module version map")
		// 	actual, err := cli.QueryModuleVersionMap(cli.Validator0)
		// 	Expect(err).To(BeNil())

		// 	Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		// })

		It("should load and run existing diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocCreatePayloads, err := RelGlob(GeneratedJSONDir, "pre", "create - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExistingDidDocCreatePayloads {
				var DidDocCreatePayload didtypesv1.MsgCreateDidPayload
				var DidDocCreateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				By("reading ")
				DidDocCreateSignInput, err = Loader(payload, &DidDocCreatePayload)
				Expect(err).To(BeNil())

				res, err := cli.CreateDidLegacy(DidDocCreatePayload, DidDocCreateSignInput, cli.Validator0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocUpdatePayloads, err := RelGlob(GeneratedJSONDir, "pre", "update - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExistingDidDocUpdatePayloads {
				var DidDocUpdatePayload didtypesv1.MsgUpdateDidPayload
				var DidDocUpdateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				By("reading ")
				DidDocUpdateSignInput, err = Loader(payload, &DidDocUpdatePayload)
				Expect(err).To(BeNil())

				resp, err := cli.QueryDidLegacy(DidDocUpdatePayload.Id, cli.Validator0)
				Expect(err).To(BeNil())
				Expect(resp.Did.Id).To(BeEquivalentTo(DidDocUpdatePayload.Id))

				DidDocUpdatePayload.VersionId = resp.Metadata.VersionId

				res, err := cli.UpdateDidLegacy(DidDocUpdatePayload, DidDocUpdateSignInput, cli.Validator0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExistingResourceCreatePayloads, err := RelGlob(GeneratedJSONDir, "pre", "create - resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExistingResourceCreatePayloads {
				var ResourceCreatePayload resourcetypesv1.MsgCreateResourcePayload
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

				// TODO: Add resource file. Right now, it is not possible to create a resource without a file. So we need to copy a file to the container home directory.
				res, err := cli.CreateResourceLegacy(
					ResourceCreatePayload.CollectionId,
					ResourceCreatePayload.Id,
					ResourceCreatePayload.Name,
					ResourceCreatePayload.ResourceType,
					filepath.Base(ResourceFile),
					ResourceCreateSignInput,
					cli.Validator0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		var UPGRADE_HEIGHT int64
		var VOTING_END_HEIGHT int64

		It("should calculate the upgrade height", func() {
			By("getting the current block height and calculating the voting end height")
			var err error
			UPGRADE_HEIGHT, VOTING_END_HEIGHT, err = cli.CalculateUpgradeHeight(cli.Validator0, cli.CLIBinaryName)
			Expect(err).To(BeNil())
			fmt.Printf("Upgrade height: %d\n", UPGRADE_HEIGHT)
			fmt.Printf("Voting end height: %d\n", VOTING_END_HEIGHT)
		})

		It("should submit a software upgrade proposal", func() {
			By("sending a SubmitUpgradeProposal transaction from `validator0` container")
			res, err := cli.SubmitUpgradeProposal(UPGRADE_HEIGHT, cli.Validator0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should deposit tokens for the software upgrade proposal", func() {
			By("sending a DepositGov transaction from `validator0` container")
			res, err := cli.DepositGov(cli.Validator0)
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator0` container", func() {
			By("sending a VoteProposal transaction from `validator0` container")
			res, err := cli.VoteProposal(cli.Validator0, "1", "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator1` container", func() {
			By("sending a VoteProposal transaction from `validator1` container")
			res, err := cli.VoteProposal(cli.Validator1, "1", "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator2` container", func() {
			By("sending a VoteProposal transaction from `validator2` container")
			res, err := cli.VoteProposal(cli.Validator2, "1", "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should vote for the software upgrade proposal from `validator3` container", func() {
			By("sending a VoteProposal transaction from `validator3` container")
			res, err := cli.VoteProposal(cli.Validator3, "1", "yes")
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))
		})

		It("should wait for the voting end height to be reached", func() {
			By("pinging the node status until the voting end height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CLIBinaryName, VOTING_END_HEIGHT, cli.VotingPeriod)
			Expect(err).To(BeNil())
		})

		It("should query the proposal status to ensure it has passed", func() {
			By("sending a QueryProposal Msg from `validator0` container")
			proposal, err := cli.QueryProposalLegacy(cli.Validator0, "1")
			Expect(err).To(BeNil())
			Expect(proposal.Status).To(BeEquivalentTo(govtypesv1beta1.StatusPassed))
		})

		It("should wait for the upgrade height to be reached", func() {
			By("pinging the node status until the upgrade height is reached")
			err := cli.WaitForChainHeight(cli.Validator0, cli.CLIBinaryName, UPGRADE_HEIGHT, cli.VotingPeriod)
			Expect(err).To(BeNil())
		})
	})
})
