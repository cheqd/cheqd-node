//go:build integration

package integration

import (
	"crypto/ed25519"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/ante"
	posthandler "github.com/cheqd/cheqd-node/post"
	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	"github.com/google/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// cases:
//   1. fixed fee - invalid fee denom create
//   2. fixed fee - invalid fee denom update
//   3. fixed fee - invalid fee denom deactivate
//   4. fixed fee - invalid fee amount create
//   5. fixed fee - invalid fee amount update
//   6. fixed fee - invalid fee amount deactivate
//   10. fixed fee - charge only tax if fee is more than tax create
//   11. fixed fee - charge only tax if fee is more than tax update
//   12. fixed fee - charge only tax if fee is more than tax deactivate
//   13. fixed fee - insufficient funds create
//   14. fixed fee - insufficient funds update
//   15. fixed fee - insufficient funds deactivate

var _ = Describe("cheqd cli - negative diddoc pricing", func() {
	var tmpDir string
	var feeParams types.FeeParams
	var payload didcli.DIDDocument
	var signInputs []didcli.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryDidParams()
		Expect(err).To(BeNil())

		feeParams = res.Params

		// Create a new DID Doc
		did := "did:cheqd:" + network.DidNamespace + ":" + uuid.NewString()
		keyId := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

		payload = didcli.DIDDocument{
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
				PrivKey:              privateKey,
			},
		}
	})

	It("should not succeed in create diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting create diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(feeParams.CreateDid[0].MaxAmount.Int64()))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in update diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(feeParams.GetUpdateDid()[0].MinAmount.Int64()))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdkmath.NewInt(feeParams.GetDeactivateDid()[0].MinAmount.Int64()))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should fail in create diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting create diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.CreateDid[0].Denom, sdkmath.NewInt(feeParams.CreateDid[0].MinAmount.Int64()-1000000))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(1))
	})

	It("should fail in update diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.UpdateDid[0].Denom, sdkmath.NewInt(feeParams.UpdateDid[0].MinAmount.Int64()-1000000))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(1))
	})

	It("should fail in deactivate diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.DeactivateDid[0].Denom, sdkmath.NewInt(feeParams.DeactivateDid[0].MinAmount.Int64()-1000000))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(1))
	})

	It("should not charge more than tax for update diddoc message - case: fixed fee", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)},
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("fetching cheq EMA price and computing fees")
		tax := feeParams.UpdateDid[0]
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.MinAmount.Mul(sdkmath.NewInt(2)))

		convertedFees, err := ante.GetFeeForMsg(sdk.NewCoins(doubleTax), feeParams.UpdateDid, cheqp, nil)
		Expect(err).To(BeNil())
		burnPortionUsd := helpers.GetBurnFeePortion(feeParams.BurnFactor, convertedFees)
		rewardPortionUsd := helpers.GetRewardPortion(convertedFees, burnPortionUsd)

		burnPortionCheq, err := posthandler.ConvertToCheq(burnPortionUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionCheq, err := posthandler.ConvertToCheq(rewardPortionUsd, cheqp)
		Expect(err).To(BeNil())

		taxInCheqd := burnPortionCheq.Add(rewardPortionCheq...)

		By("submitting the update diddoc message with double the tax")
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased only by the actual tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(BeEquivalentTo(taxInCheqd.AmountOf(types.BaseMinimalDenom)))
	})

	It("should charge more than tax for deactivate diddoc message - case: fee range between min and max", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.ID,
			VersionId: uuid.NewString(),
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the deactivate diddoc message with double the tax")
		tax := feeParams.DeactivateDid[0].MinAmount
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Mul(sdkmath.NewInt(2)))
		price, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)

		Expect(err).To(BeNil())
		userFee := sdk.NewCoins(doubleTax)
		fees, err := ante.GetFeeForMsg(userFee, feeParams.DeactivateDid, price.Price, nil)
		Expect(err).To(BeNil())

		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(fees.AmountOf(types.BaseMinimalDenom)))
	})

	It("should not succeed in create diddoc create message - case: fixed fee, insufficient funds", func() {
		By("submitting create diddoc message with insufficient funds")
		tax := feeParams.CreateDid[0].MaxAmount
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
	})

	It("should not succeed in update diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               payload.VerificationMethod[0]["type"],
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with insufficient funds")
		tax := feeParams.UpdateDid[0].MinAmount
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()+feeParams.UpdateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with insufficient funds")
		tax := feeParams.DeactivateDid[0].MinAmount
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()+feeParams.DeactivateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(res.RawLog).To(ContainSubstring(sdkerrors.ErrInsufficientFunds.Error()))
	})
})
