//go:build integration

package integration

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	clitypes "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
//   7. fixed fee - insufficient funds create
//   8. fixed fee - insufficient funds update
//   9. fixed fee - insufficient funds deactivate
//   10. fixed fee - charge only tax if fee is more than tax create
//   11. fixed fee - charge only tax if fee is more than tax update
//   12. fixed fee - charge only tax if fee is more than tax deactivate
//   13. gas auto - insufficient funds create
//   14. gas auto - insufficient funds update
//   15. gas auto - insufficient funds deactivate

var _ = Describe("cheqd cli - negative diddoc pricing", func() {
	var tmpDir string
	var feeParams types.FeeParams
	var payload types.MsgCreateDidDocPayload
	var signInputs []clitypes.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryParams(types.ModuleName, string(types.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())

		// Create a new DID Doc
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload = types.MsgCreateDidDocPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
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
	})

	It("should not succeed in create diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting create diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.GetCreateDid().Amount.Int64()))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in update diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := types.MsgUpdateDidDocPayload{
			Id: payload.Id,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   payload.VerificationMethod[0].Id,
					Controller:           payload.VerificationMethod[0].Controller,
					Type:                 payload.VerificationMethod[0].Type,
					VerificationMaterial: payload.VerificationMethod[0].VerificationMaterial,
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0].Id}, // <-- changed
			VersionId:       uuid.NewString(),                           // <-- changed
		}

		By("submitting update diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.GetUpdateDid().Amount.Int64()))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("submitting deactivate diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.GetDeactivateDid().Amount.Int64()))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not fail in create diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting create diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.CreateDid.Denom, sdk.NewInt(feeParams.CreateDid.Amount.Int64()-1))
		_, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
	})

	It("should not fail in update diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := types.MsgUpdateDidDocPayload{
			Id: payload.Id,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   payload.VerificationMethod[0].Id,
					Controller:           payload.VerificationMethod[0].Controller,
					Type:                 payload.VerificationMethod[0].Type,
					VerificationMaterial: payload.VerificationMethod[0].VerificationMaterial,
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0].Id}, // <-- changed
			VersionId:       uuid.NewString(),                           // <-- changed
		}

		By("submitting update diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.UpdateDid.Denom, sdk.NewInt(feeParams.UpdateDid.Amount.Int64()-1))
		_, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
	})

	It("should not fail in deactivate diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("submitting deactivate diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.DeactivateDid.Denom, sdk.NewInt(feeParams.DeactivateDid.Amount.Int64()-1))
		_, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(lowerTax.String()))
		Expect(err).NotTo(BeNil())
	})

	It("should not succeed in create diddoc create message - case: fixed fee, insufficient funds", func() {
		By("submitting create diddoc message with insufficient funds")
		tax := feeParams.CreateDid
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in update diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := types.MsgUpdateDidDocPayload{
			Id: payload.Id,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   payload.VerificationMethod[0].Id,
					Controller:           payload.VerificationMethod[0].Controller,
					Type:                 payload.VerificationMethod[0].Type,
					VerificationMaterial: payload.VerificationMethod[0].VerificationMaterial,
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0].Id}, // <-- changed
			VersionId:       uuid.NewString(),                           // <-- changed
		}

		By("submitting update diddoc message with insufficient funds")
		tax := feeParams.UpdateDid
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(feeParams.GetCreateDid().String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("submitting deactivate diddoc message with insufficient funds")
		tax := feeParams.DeactivateDid
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees(tax.String()))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in create diddoc message - case: gas auto, insufficient funds", func() {
		By("submitting create diddoc message with insufficient funds")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees("0.0001stake"))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in update diddoc message - case: gas auto, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := types.MsgUpdateDidDocPayload{
			Id: payload.Id,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   payload.VerificationMethod[0].Id,
					Controller:           payload.VerificationMethod[0].Controller,
					Type:                 payload.VerificationMethod[0].Type,
					VerificationMaterial: payload.VerificationMethod[0].VerificationMaterial,
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0].Id}, // <-- changed
			VersionId:       uuid.NewString(),                           // <-- changed
		}

		By("submitting update diddoc message with insufficient funds")
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees("0.0001stake"))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in deactivate diddoc message - case: gas auto, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("submitting deactivate diddoc message with insufficient funds")
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_3, helpers.GenerateFees("0.0001stake"))
		Expect(err).NotTo(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not charge more than tax for create diddoc message - case: fixed fee", func() {
		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("submitting the create diddoc message with double the tax")
		tax := feeParams.CreateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax for update diddoc message - case: fixed fee", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := types.MsgUpdateDidDocPayload{
			Id: payload.Id,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                   payload.VerificationMethod[0].Id,
					Controller:           payload.VerificationMethod[0].Controller,
					Type:                 payload.VerificationMethod[0].Type,
					VerificationMaterial: payload.VerificationMethod[0].VerificationMaterial,
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0].Id}, // <-- changed
			VersionId:       uuid.NewString(),                           // <-- changed
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("submitting the update diddoc message with double the tax")
		tax := feeParams.UpdateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax for deactivate diddoc message - case: fixed fee", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("submitting the deactivate diddoc message with double the tax")
		tax := feeParams.DeactivateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(tmpDir, testdata.BASE_ACCOUNT_2)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})
})
