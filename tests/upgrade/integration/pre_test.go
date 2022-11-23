//go:build upgrade

package upgrade

import (
	"fmt"
	"path/filepath"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	Context("Before a softare upgrade execution is initiated", func() {
		It("should wait for chain to bootstrap", func() {
			By("pinging the node status until the dvoting end height is reached")
			err := cli.WaitForChainHeight(cli.VALIDATOR0, cli.CLI_BINARY_NAME, cli.BOOTSTRAP_HEIGHT, cli.BOOTSTRAP_PERIOD)
			Expect(err).To(BeNil())
		})

		It("should load and run existing diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocCreatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "diddoc", "create", "*.json")
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing diddoc sign input")
			ExistingSignInputCreatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "diddoc", "create", "signinput", "*.json")
			Expect(err).To(BeNil())

			for i, payload := range ExistingDidDocCreatePayloads {
				var DidDocCreatePayload didtypesv1.MsgCreateDidPayload
				var DidDocCreateSignInput []cli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocCreatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputCreatePayloads[i], &DidDocCreateSignInput)
				Expect(err).To(BeNil())

				res, err := cli.CreateDidLegacy(DidDocCreatePayload, DidDocCreateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExistingResourceCreatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "resource", "create", "*.json")
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing resource sign input")
			ExistingSignInputCreatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "resource", "create", "signinput", "*.json")
			Expect(err).To(BeNil())

			By("copying the existing resource file to the container")
			ResourceFile, err := integrationtestdata.CreateTestJson(GinkgoT().TempDir())
			Expect(err).To(BeNil())
			_, err = cli.LocalnetExecCopyAbsoluteWithPermissions(ResourceFile, cli.DOCKER_HOME, cli.VALIDATOR0)
			Expect(err).To(BeNil())

			for i, payload := range ExistingResourceCreatePayloads {
				var ResourceCreatePayload resourcetypesv1.MsgCreateResourcePayload
				var ResourceCreateSignInput []cli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				err = Loader(payload, &ResourceCreatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputCreatePayloads[i], &ResourceCreateSignInput)
				Expect(err).To(BeNil())

				// TODO: Add resource file. Right now, it is not possible to create a resource without a file. So we need to copy a file to the container home directory.
				res, err := cli.CreateResourceLegacy(ResourceCreatePayload.CollectionId, ResourceCreatePayload.Id, ResourceCreatePayload.Name, ResourceCreatePayload.ResourceType, filepath.Base(ResourceFile), ResourceCreateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocUpdatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "diddoc", "update", "*.json")
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing diddoc sign input")
			ExistingSignInputUpdatePayloads, err := RelGlob(GENERATED_JSON_DIR, "existing", "diddoc", "update", "signinput", "*.json")
			Expect(err).To(BeNil())

			for i, payload := range ExistingDidDocUpdatePayloads {
				var DidDocUpdatePayload didtypesv1.MsgUpdateDidPayload
				var DidDocUpdateSignInput []cli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocUpdatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputUpdatePayloads[i], &DidDocUpdateSignInput)
				Expect(err).To(BeNil())

				q, err := cli.QueryDidLegacy(DidDocUpdatePayload.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(q.Did.Id).To(BeEquivalentTo(DidDocUpdatePayload.Id))

				DidDocUpdatePayload.VersionId = q.Metadata.VersionId

				res, err := cli.UpdateDidLegacy(DidDocUpdatePayload, DidDocUpdateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		var UPGRADE_HEIGHT int64
		var VOTING_END_HEIGHT int64

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
			proposal, err := cli.QueryUpgradeProposalLegacy(cli.VALIDATOR0)
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
