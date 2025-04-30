//go:build integration

package integration

import (
	"crypto/ed25519"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
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

var _ = Describe("cheqd cli - negative resource pricing", func() {
	var tmpDir string
	var didFeeParams didtypes.FeeParams
	var resourceFeeParams resourcetypes.FeeParams
	var collectionID string
	var signInputs []didcli.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query did fee params
		_, err := cli.QueryDidParams()
		Expect(err).To(BeNil())

		// Query resource fee params
		_, err = cli.QueryResourceParams()
		Expect(err).To(BeNil())

		// Create a new DID Doc
		collectionID = uuid.NewString()
		did := "did:cheqd:" + network.DidNamespace + ":" + collectionID
		keyId := did + "#key1"

		publicKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

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
		resp, err := cli.CreateDidDoc(tmpDir, didPayload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(didFeeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))
	})

	It("should not succeed in create resource json message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource json message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Json.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in create resource image message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource image message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Image.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in create resource default message - case: fixed fee, invalid denom", func() {
		By("preparing the create resource default message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Default.Amount.Int64()))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not fail in create resource json message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource json message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with lower amount than required")
		lowerTax := sdk.NewCoin(resourceFeeParams.Json.Denom, sdkmath.NewInt(resourceFeeParams.Json.Amount.Int64()-1))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not fail in create resource image message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource image message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with lower amount than required")
		lowerTax := sdk.NewCoin(resourceFeeParams.Image.Denom, sdkmath.NewInt(resourceFeeParams.Image.Amount.Int64()-1))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not fail in create resource default message - case: fixed fee, lower amount than required", func() {
		By("preparing the create resource default message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with lower amount than required")
		lowerTax := sdk.NewCoin(resourceFeeParams.Default.Denom, sdkmath.NewInt(resourceFeeParams.Default.Amount.Int64()-1))
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not succeed in create resource json message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource json message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with insufficient funds")
		tax := resourceFeeParams.Json
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource image message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource image message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the image resource message with insufficient funds")
		tax := resourceFeeParams.Image
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create resource default message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource default message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the default resource message with insufficient funds")
		tax := resourceFeeParams.Default
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not charge more than tax in create resource json message - case: fixed fee", func() {
		By("preparing the create resource json message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the json resource message with double the tax")
		tax := resourceFeeParams.Json
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdkmath.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax in create resource image message - case: fixed fee", func() {
		By("preparing the create resource image message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the image resource message with double the tax")
		tax := resourceFeeParams.Image
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdkmath.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax in create resource default message - case: fixed fee", func() {
		By("preparing the create resource default message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the default resource message with double the tax")
		tax := resourceFeeParams.Default
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.Amount.Mul(sdkmath.NewInt(2)))
		_, err = cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})
})
