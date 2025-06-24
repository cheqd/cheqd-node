//go:build integration

package integration

import (
	"crypto/ed25519"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/ante"
	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		didRes, err := cli.QueryDidParams()
		Expect(err).To(BeNil())

		didFeeParams = didRes.Params

		// Query resource fee params
		resourceRes, err := cli.QueryResourceParams()
		Expect(err).To(BeNil())

		resourceFeeParams = resourceRes.Params

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
		resp, err := cli.CreateDidDoc(tmpDir, didPayload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFees(didFeeParams.CreateDid[0].MaxAmount.String()+didFeeParams.CreateDid[0].Denom))
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
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Json[0].MinAmount.Int64()))
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
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Image[0].MinAmount.Int64()))
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
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(resourceFeeParams.Default[0].MinAmount.Int64()))
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

	// It("should not fail in create resource json message - case: fixed fee, lower amount than required", func() {
	// 	By("preparing the create resource json message")
	// 	resourceID := uuid.NewString()
	// 	resourceName := "TestResource"
	// 	resourceVersion := "1.0"
	// 	resourceType := "TestType"
	// 	resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
	// 	Expect(err).To(BeNil())

	// 	By("submitting the json resource message with lower amount than required")
	// 	lowerTax := sdk.NewCoin(resourceFeeParams.Json[0].Denom, sdkmath.NewInt(resourceFeeParams.Json[0].MinAmount.Int64()-1000000))
	// 	res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
	// 		CollectionId: collectionID,
	// 		Id:           resourceID,
	// 		Name:         resourceName,
	// 		Version:      resourceVersion,
	// 		ResourceType: resourceType,
	// 	}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
	// 	Expect(err).To(BeNil())
	// 	Expect(res.Code).To(BeEquivalentTo(1))
	// })

	// It("should fail in create resource image message - case: fee ranging between two values so lower value than lower bound fails", func() {
	// 	By("preparing the create resource image message")
	// 	resourceID := uuid.NewString()
	// 	resourceName := "TestResource"
	// 	resourceVersion := "1.0"
	// 	resourceType := "TestType"
	// 	resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
	// 	Expect(err).To(BeNil())

	// 	By("submitting the image resource message with lower amount than required")

	// 	cheqPrice, err := cli.QueryWMA(types.BaseDenom)
	// 	Expect(err).To(BeNil())
	// 	cheqp := cheqPrice.Price

	// 	minNcheq := resourceFeeParams.Image[0].MinAmount

	// 	// Constants
	// 	cheqScale := sdkmath.NewInt(1_000_000_000) // 1e9
	// 	usdScale := sdkmath.NewInt(1_000_000)      // 1e6

	// 	// Convert minNcheq to USD (6 decimals)
	// 	minNcheqDec := sdkmath.LegacyNewDecFromInt(*minNcheq).QuoInt(cheqScale)
	// 	minUsdDec := minNcheqDec.Mul(cheqp).MulInt(usdScale) // Final in 6-decimal USD

	// 	// Go slightly below the lower bound (1 µUSD = 1)
	// 	usdBelowMin := minUsdDec.TruncateInt().SubRaw(1)

	// 	// Convert back to ncheq
	// 	usdBelowMinDec := sdkmath.LegacyNewDecFromInt(usdBelowMin)
	// 	cheqBelowMin := usdBelowMinDec.
	// 		Quo(cheqp).
	// 		MulInt(cheqScale).
	// 		QuoInt(usdScale).
	// 		TruncateInt()

	// 	// Final lower-than-min fee
	// 	lowerTax := sdk.NewCoin(types.BaseMinimalDenom, cheqBelowMin)

	// 	res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
	// 		CollectionId: collectionID,
	// 		Id:           resourceID,
	// 		Name:         resourceName,
	// 		Version:      resourceVersion,
	// 		ResourceType: resourceType,
	// 	}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
	// 	Expect(res.Code).To(BeEquivalentTo(13))
	// 	// insufficient fee
	// })

	// It("should not fail in create resource default message - case: fixed fee, lower amount than required", func() {
	// 	By("preparing the create resource default message")
	// 	resourceID := uuid.NewString()
	// 	resourceName := "TestResource"
	// 	resourceVersion := "1.0"
	// 	resourceType := "TestType"
	// 	resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
	// 	Expect(err).To(BeNil())

	// 	By("submitting the default resource message with lower amount than required")
	// 	lowerTax := sdk.NewCoin(resourceFeeParams.Default[0].Denom, sdkmath.NewInt(resourceFeeParams.Default[0].MinAmount.Int64()-1000000))
	// 	res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
	// 		CollectionId: collectionID,
	// 		Id:           resourceID,
	// 		Name:         resourceName,
	// 		Version:      resourceVersion,
	// 		ResourceType: resourceType,
	// 	}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
	// 	Expect(err).To(BeNil())
	// 	Expect(res.Code).To(BeEquivalentTo(1))
	// })

	It("should not succeed in create resource json message - case: fixed fee, insufficient funds", func() {
		By("preparing the create resource json message")
		resourceID := uuid.NewString()
		resourceName := "TestResource"
		resourceVersion := "1.0"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		By("submitting the json resource message with insufficient funds")
		tax := resourceFeeParams.Json[0]
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.MaxAmount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
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
		tax := resourceFeeParams.Image[0]
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.MaxAmount.BigInt().String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
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
		tax := resourceFeeParams.Default[0]
		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.MaxAmount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
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
		tax := resourceFeeParams.Json[0]
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.MaxAmount.Mul(sdkmath.NewInt(2)))

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(doubleTax)
		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Json, cheqp, nil)

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
		Expect(diff).To(Equal(convertedFees.AmountOf(types.BaseMinimalDenom)))
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
		tax := resourceFeeParams.Image[0]
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.MaxAmount.Mul(sdkmath.NewInt(2)))

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(doubleTax)
		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Image, cheqp, nil)

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
		Expect(diff).To(Equal(convertedFees.AmountOf(types.BaseMinimalDenom)))
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
		tax := resourceFeeParams.Default[0]
		doubleTax := sdk.NewCoin(resourcetypes.BaseMinimalDenom, tax.MaxAmount.Mul(sdkmath.NewInt(2)))

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(doubleTax)
		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Default, cheqp, nil)
		Expect(err).To(BeNil())
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
		Expect(diff).To(Equal(convertedFees.AmountOf(types.BaseMinimalDenom)))
	})
})
