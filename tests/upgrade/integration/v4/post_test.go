//go:build upgrade_integration

package integration

import (
	"fmt"
	"path/filepath"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/v4/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Post", func() {
	var feeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams

	BeforeEach(func() {
		// query fee params - case: did
		didRes, err := cli.QueryDidFeeParams(cli.Validator0)
		Expect(err).To(BeNil())
		feeParams = didRes.Params
		Expect(feeParams).NotTo(BeNil())

		// query fee params - case: resource
		resourceRes, err := cli.QueryResourceFeeParams(cli.Validator0)
		Expect(err).To(BeNil())
		resourceFeeParams = resourceRes.Params
		Expect(resourceFeeParams).NotTo(BeNil())
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
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*2)
			Expect(err).To(BeNil())
		})

		It("should match the expected module version map", func() {
			By("loading the expected module version map")
			var expected upgradetypes.QueryModuleVersionsResponse
			_, err := Loader(filepath.Join(GeneratedJSONDir, "post", "query - module-version-map", "v4.json"), &expected)
			Expect(err).To(BeNil())

			By("matching the expected module version map")
			actual, err := cli.QueryModuleVersionMap(cli.Validator0)
			Expect(err).To(BeNil())

			Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		})

		It("should export command work", func() {
			By("running export command at latest height")
			_, err := cli.RunExportCommand(cli.Validator2, cli.CliBinaryName)
			Expect(err).To(BeNil())
		})

		It("should restart stopped container", func() {
			By("restart container which was stopped while export command")
			_, err := cli.LocalnetStartContainer(cli.Validator2)
			Expect(err).To(BeNil())

			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top")
			err = cli.WaitForChainHeight(cli.Validator2, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*2)
			Expect(err).To(BeNil())
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

		It("should burn the coins from the given address (here container/validator)", func() {
			By("querying the account balance")
			Operator0_address, err := integrationcli.QueryKeys(integrationcli.Operator0)
			Expect(err).To(BeNil())
			balanceBefore, err := integrationcli.QueryBalance(Operator0_address, didtypes.BaseMinimalDenom)
			Expect(err).To(BeNil())
			Expect(balanceBefore.Denom).To(Equal(didtypes.BaseMinimalDenom))

			By("burning the coins")
			coins := sdk.NewCoins(sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(1000)})
			res, err := cli.BurnMsg(cli.Validator0, coins.String())
			Expect(err).To(BeNil())
			Expect(res.Code).To(BeEquivalentTo(0))

			By("querying the account balance again")
			balanceAfter, err := integrationcli.QueryBalance(Operator0_address, didtypes.BaseMinimalDenom)
			Expect(err).To(BeNil())
			Expect(balanceAfter.Denom).To(Equal(didtypes.BaseMinimalDenom))

			By("checking if the balance has been reduced by the expected amount")
			expectedBalance := balanceBefore.Sub(sdk.NewCoin(didtypes.BaseMinimalDenom, sdkmath.NewInt(1000)))
			diff := balanceAfter.Amount.Sub(expectedBalance.Amount)
			Expect(diff).To(Equal(sdkmath.NewInt(1000)))
		})

		It("should ensure whitelisted IBC messages bypass global fees", func() {
			By("querying the globalfee bypass messages")
			bypassMessages, err := integrationcli.QueryBypassMessages()
			Expect(err).To(BeNil())

			expectedBypassMessages := []string{
				"/ibc.core.channel.v1.MsgAcknowledgement",
				"/ibc.core.client.v1.MsgUpdateClient",
				"/ibc.core.channel.v1.MsgRecvPacket",
				"/ibc.core.channel.v1.MsgTimeout",
			}

			for _, typeURL := range expectedBypassMessages {
				Expect(bypassMessages).To(ContainElement(typeURL))
			}

			tmpDir := GinkgoT().TempDir()

			By("broadcasting IBC MsgAcknowledgement with zero fees")
			res, err := integrationcli.IBCAcknowledgementTx(tmpDir, testdata.RELAYER_ACCOUNT, testdata.IBCAcknowledgementMsg, testdata.IBCAcknowledgementGasLimit)
			Expect(err).To(BeNil())
			Expect(res.Code).NotTo(BeEquivalentTo(13))

			By("broadcasting IBC MsgRecvPacket with zero fees")
			res, err = integrationcli.IBCRecvPacketTx(tmpDir, testdata.RELAYER_ACCOUNT, testdata.IBCRecvPacketMsg, testdata.IBCRecvPacketGasLimit)
			Expect(err).To(BeNil())
			Expect(res.Code).NotTo(BeEquivalentTo(13))

			By("broadcasting IBC MsgTimeout with zero fees")
			res, err = integrationcli.IBCTimeoutTx(tmpDir, testdata.RELAYER_ACCOUNT, testdata.IBCTimeoutMsg, testdata.IBCTimeoutGasLimit)
			Expect(err).To(BeNil())
			Expect(res.Code).NotTo(BeEquivalentTo(13))
		})

		It("should ensure governance expedited deposit params are restored", func() {
			By("querying the governance module params")
			govParams, err := cli.QueryGovParams(cli.Validator0)
			Expect(err).To(BeNil())
			Expect(govParams.Params).NotTo(BeNil())

			By("verifying the expedited minimum deposit matches the expected coin")
			expeditedDeposit := govParams.Params.ExpeditedMinDeposit
			Expect(expeditedDeposit).To(HaveLen(1))

			expeditedCoin := expeditedDeposit[0]
			Expect(expeditedCoin.Denom).To(Equal(didtypes.BaseMinimalDenom))
			Expect(expeditedCoin.Amount.String()).To(Equal("8000000000000"))
		})
	})
})
