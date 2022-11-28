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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// cases:
//   1. fixed fee - invalid fee denom image
//   2. fixed fee - invalid fee denom json
//   3. fixed fee - invalid fee denom default
//   4. fixed fee - invalid fee amount image
//   5. fixed fee - invalid fee amount json
//   6. fixed fee - invalid fee amount default
//   7. fixed fee - insufficient funds image
//   8. fixed fee - insufficient funds json
//   9. fixed fee - insufficient funds default
//   10. fixed fee - charge only tax if fee is more than tax image
//   11. fixed fee - charge only tax if fee is more than tax json
//   12. fixed fee - charge only tax if fee is more than tax default
//   13. gas auto - insufficient funds image
//   14. gas auto - insufficient funds json
//   15. gas auto - insufficient funds default

var _ = Describe("cheqd cli - negative resource pricing", func() {
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

	It("should not succeed in create resource json message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.Json.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in create resource image message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.Image.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in create resource default message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.Default.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not fail in create resource json message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.Json.Denom, sdk.NewInt(feeParams.Json.Amount.Int64()-1))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
	})

	It("should not fail in create resource image message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.Image.Denom, sdk.NewInt(feeParams.Image.Amount.Int64()-1))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
	})

	It("should not fail in create resource default message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.Default.Denom, sdk.NewInt(feeParams.Default.Amount.Int64()-1))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
	})

	It("should not succeed in create resource json message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with insufficient funds")
		tax := feeParams.Json
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource image message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with insufficient funds")
		tax := feeParams.Image
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource default message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with insufficient funds")
		tax := feeParams.Default
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource json - case: gas auto, insufficient funds", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with insufficient funds")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource image - case: gas auto, insufficient funds", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with insufficient funds")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource default - case: gas auto, insufficient funds", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with insufficient funds")
		res, err := cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_3, cli.CLI_GAS_PARAMS)
		Expect(err).ToNot(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not charge more than tax in create resource json message - case: fixed fee", func() {
		By("preparing the create resource json message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the json resource message with double the tax")
		tax := feeParams.Json
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax in create resource image message - case: fixed fee", func() {
		By("preparing the create resource image message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the image resource message with double the tax")
		tax := feeParams.Image
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax in create resource default message - case: fixed fee", func() {
		By("preparing the create resource default message")
		resourceId := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the default resource message with double the tax")
		tax := feeParams.Default
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcecli.CreateResourceOptions{
			CollectionId:    collectionId,
			ResourceId:      resourceId,
			ResourceName:    resourceName,
			ResourceVersion: resourceVersion,
			ResourceType:    resourceType,
			ResourceFile:    resourceFile,
		}, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2, resourcetypes.BaseMinimalDenom)

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})
})
