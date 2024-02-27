//go:build integration

package integration

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource pricing", func() {
	var tmpDir string
	var didFeeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams
	var collectionID string
	var signInputs []didcli.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query did fee params
		res, err := cli.QueryParams(didtypes.ModuleName, string(didtypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &didFeeParams)
		Expect(err).To(BeNil())

		// Query resource fee params
		res, err = cli.QueryParams(resourcetypes.ModuleName, string(resourcetypes.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &resourceFeeParams)
		Expect(err).To(BeNil())

		// Create a new DID Doc
		collectionID = uuid.NewString()
		did := "did:cheqd:" + network.DidNamespace + ":" + collectionID
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(pubKey)

		didPayload := didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
		}

		signInputs = []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privKey,
			},
		}

		// Submit the DID Doc
		resp, err := cli.CreateDidDoc(tmpDir, didPayload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(didFeeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))
	})

	It("should tax json resource message - case: fixed fee", func() {
		By("preparing the json resource message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the json resource message")
		tax := resourceFeeParams.Json
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, resourceFeeParams.BurnFactor)
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
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the image resource message")
		tax := resourceFeeParams.Image
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, resourceFeeParams.BurnFactor)
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
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("submitting the default resource message")
		tax := resourceFeeParams.Default
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: tax.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		burnt := helpers.GetBurntPortion(tax, resourceFeeParams.BurnFactor)
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
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource json message")
		tax := resourceFeeParams.Json
		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.String())))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
	})

	It("should tax create resource image with feegrant - case: fixed fee", func() {
		By("preparing the create resource image message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource image message")
		tax := resourceFeeParams.Image
		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.String())))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
	})

	It("should tax create resource default with feegrant - case: fixed fee", func() {
		By("preparing the create resource default message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create resource default message")
		tax := resourceFeeParams.Default
		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.String())))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
	})
})
