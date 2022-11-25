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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive diddoc pricing", func() {
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

	// TODO: Fix race condition
	/* It("should tax create diddoc message - case: fixed fee", func() {
		feeParams := helpers.GenerateFees("5000000000ncheq")

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, feeParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		fmt.Println("res", res)

		// WIP:
		// 	 1. Check the balance of fee payer account
		//   2. Check supply deflation
		//   3. Check events
		//      - fee: {amount: 5000000000ncheq, payer: <fee payer>, granter: <granter>}
	}) */

	It("should tax create diddoc message - case: gas auto", func() {
		By("querying the total supply")
		supplyBeforeDeflation, err := cli.QuerySupplyOf(types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the current account balance")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting a create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("checking the balance difference")
		minusTax := balanceBefore.Amount.Sub(balanceAfter.Amount)
		tax := feeParams.TxTypes[types.DefaultKeyCreateDid]
		Expect(minusTax).To(Equal(tax.Amount))

		By("querying the deflated total supply")
		supplyAfterDeflation, err := cli.QuerySupplyOf(types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the deflation")
		burnt := helpers.GetBurntPortion(tax, feeParams.BurnFactor)
		Expect(supplyBeforeDeflation.Amount.Sub(supplyAfterDeflation.Amount)).To(Equal(burnt.Amount))

		// WIP:
		// [x] 1. Check the balance of fee payer account
		// [x] 2. Check supply deflation
		// [_] 3. Check events
		//      - fee: {amount: 5000000000ncheq, payer: <fee payer>, granter: <granter>}
	})
})
