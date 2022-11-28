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

	It("should tax create diddoc message - case: fixed fee", func() {
		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting a create diddoc message")
		tax := feeParams.CreateDid
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

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

	It("should tax create diddoc message - case: gas auto", func() {
		By("querying the fee payer account balance before the transaction")
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
		tax := feeParams.CreateDid
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

	It("should tax update diddoc message - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

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
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting an update diddoc message")
		tax := feeParams.UpdateDid
		res, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

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

	It("should tax update diddoc message - case: gas auto", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

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
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting an update diddoc message")
		tax := feeParams.UpdateDid
		res, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

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

	It("should tax deactivate diddoc message - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting an deactivate diddoc message")
		tax := feeParams.DeactivateDid
		res, err := cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

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

	It("should tax deactivate diddoc message - case: gas auto", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting an deactivate diddoc message")
		res, err := cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_2, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("checking the balance difference")
		tax := feeParams.DeactivateDid
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

	It("should tax create diddoc message with feegrant - case: fixed fee", func() {
		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.CreateDid
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff).To(BeEquivalentTo(0))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
	})

	It("should tax update diddoc message with feegrant - case: fixed fee", func() {
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

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting an update diddoc message")
		resp, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.UpdateDid
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff).To(BeEquivalentTo(0))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
	})

	It("should tax deactivate diddoc message with feegrant - case: fixed fee", func() {
		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.Id,
			VersionId: uuid.NewString(),
		}

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_2_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a deactivate diddoc message")
		resp, err := cli.DeactivateDidDoc(tmpDir, payload2, signInputs, testdata.BASE_ACCOUNT_1, cli.CLI_GAS_PARAMS)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_2_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.DeactivateDid
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
