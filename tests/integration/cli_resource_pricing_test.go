//go:build integration

package integration

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	clitypes "github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcecli "github.com/cheqd/cheqd-node/x/resource/client/cli"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource pricing", func() {
	var tmpDir string
	var feeParams resourcetypes.FeeParams
	var collectionId string
	var signInputs []clitypes.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryParams(resourcetypes.ModuleName, string(resourcetypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())

		// Create a new DID Doc
		collectionId = uuid.NewString()
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + collectionId
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		didPayload := didtypes.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*didtypes.VerificationMethod{
				{
					Id:                   keyId,
					Type:                 "Ed25519VerificationKey2020",
					Controller:           did,
					VerificationMaterial: "{\"publicKeyMultibase\": \"" + string(pubKeyMultibase58) + "\"}",
				},
			},
			Authentication: []string{keyId},
			VersionId:      uuid.NewString(),
		}

		signInputs = []clitypes.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		// Submit the DID Doc
		resp, err := cli.CreateDidDoc(tmpDir, didPayload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))
	})

	It("should tax json resource message - case: fixed fee", func() {
		By("preparing the json resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the json resource message")
		tax := feeParams.Json
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax json resource message - case: gas auto", func() {
		By("preparing the json resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the json resource message")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		tax := feeParams.Json
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax image resource message - case: fixed fee", func() {
		By("preparing the image resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the image resource message")
		tax := feeParams.Image
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax image resource message - case: gas auto", func() {
		By("preparing the image resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the image resource message")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		tax := feeParams.Image
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax default resource message - case: fixed fee", func() {
		By("preparing the default resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the default resource message")
		tax := feeParams.Default
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax default resource message - case: gas auto", func() {
		By("preparing the default resource message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the default resource message")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		tax := feeParams.Default
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(BeEquivalentTo(tax.Amount))

		By("exporting a readable tx event log")
		events := helpers.ReadableEvents(res.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: diff.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_2_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnt.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		reward := helpers.GetRewardPortion(tax, burnt)
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: reward.String(), Index: true},
				},
			},
		))
	})

	It("should tax create resource json message with feegrant - case: fixed fee", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource json message")
		resp, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.Json
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff).To(BeEquivalentTo(0))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
	})

	It("should tax create resource image with feegrant - case: fixed fee", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource image message")
		resp, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.Image
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff).To(BeEquivalentTo(0))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
	})

	It("should tax create resource default with feegrant - case: fixed fee", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource default message")
		resp, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.Default
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff).To(BeEquivalentTo(0))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
	})
})
