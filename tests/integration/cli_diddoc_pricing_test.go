package integration

import (
	"crypto/ed25519"
	"fmt"

	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/ante"
	"github.com/cheqd/cheqd-node/tests/integration/cli"
	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"

	posthandler "github.com/cheqd/cheqd-node/post"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive diddoc pricing", func() {
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

	It("should tax create diddoc message - case: fixed fee", func() {
		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("submitting a create diddoc message")
		tax := feeParams.CreateDid[0]
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, *tax.MaxAmount))
		convertedFees, err := ante.GetFeeForMsg(userFee, feeParams.CreateDid, cheqp, nil)
		Expect(err).To(BeNil())
		burnPotionInUsd := helpers.GetBurnFeePortion(feeParams.BurnFactor, convertedFees)
		rewardPortionInUsd := helpers.GetRewardPortion(convertedFees, burnPotionInUsd)

		burnPotionInUsdToCheq, err := posthandler.ConvertToCheq(burnPotionInUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionInUsdToCheq, err := posthandler.ConvertToCheq(rewardPortionInUsd, cheqp)
		Expect(err).To(BeNil())

		coin := rewardPortionInUsdToCheq.AmountOf(types.BaseMinimalDenom)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%

		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		taxIncheqd := burnPotionInUsdToCheq.Add(rewardPortionInUsdToCheq...)

		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.MaxAmount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", taxIncheqd.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: taxIncheqd.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnPotionInUsdToCheq.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: feeCollectorReward.String(), Index: true},
				},
			},
		))
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.ORACLE_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: oracleRewardCoin.String(), Index: true},
				},
			},
		))
	})
	It("should tax update diddoc message - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

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

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("fetching CHEQ price and calculating expected tax")
		tax := feeParams.UpdateDid[0]
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.MinAmount.Mul(math.NewInt(2))))
		convertedFees, err := ante.GetFeeForMsg(userFee, feeParams.UpdateDid, cheqp, nil)
		Expect(err).To(BeNil())
		burnPotionInUsd := helpers.GetBurnFeePortion(feeParams.BurnFactor, convertedFees)
		rewardPortionInUsd := helpers.GetRewardPortion(convertedFees, burnPotionInUsd)

		burnPotionInUsdToCheq, err := posthandler.ConvertToCheq(burnPotionInUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionInUsdToCheq, err := posthandler.ConvertToCheq(rewardPortionInUsd, cheqp)
		Expect(err).To(BeNil())

		coin := rewardPortionInUsdToCheq.AmountOf(types.BaseMinimalDenom)
		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%
		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		taxIncheqd := burnPotionInUsdToCheq.Add(rewardPortionInUsdToCheq...)

		By("submitting an update diddoc message")
		res, err := cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.MinAmount.Mul(math.NewInt(2)).String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", taxIncheqd.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: taxIncheqd.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnPotionInUsdToCheq.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: feeCollectorReward.String(), Index: true},
				},
			},
		))
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.ORACLE_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: oracleRewardCoin.String(), Index: true},
				},
			},
		))
	})
	It("should tax deactivate diddoc message - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		fmt.Println("balanceBefore--------", balanceBefore)
		Expect(balanceBefore.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("fetching cheq EMA price and computing fees")
		tax := feeParams.DeactivateDid[0]
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, *tax.MaxAmount))
		convertedFees, err := ante.GetFeeForMsg(userFee, feeParams.DeactivateDid, cheqp, nil)
		Expect(err).To(BeNil())
		burnPotionInUsd := helpers.GetBurnFeePortion(feeParams.BurnFactor, convertedFees)
		rewardPortionInUsd := helpers.GetRewardPortion(convertedFees, burnPotionInUsd)

		fmt.Println("convertedFees------------", convertedFees)
		burnPotionInUsdToCheq, err := posthandler.ConvertToCheq(burnPotionInUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionInUsdToCheq, err := posthandler.ConvertToCheq(rewardPortionInUsd, cheqp)
		Expect(err).To(BeNil())

		coin := rewardPortionInUsdToCheq.AmountOf(types.BaseMinimalDenom)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%
		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		taxIncheqd := burnPotionInUsdToCheq.Add(rewardPortionInUsdToCheq...)
		fmt.Println("taxInCheqd------------", taxIncheqd)
		fmt.Println("tax------------", tax.MaxAmount)
		By("submitting a deactivate diddoc message")
		res, err := cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.MaxAmount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())
		fmt.Println("balanceAfter--------", balanceAfter)

		Expect(balanceAfter.Denom).To(BeEquivalentTo(types.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		fmt.Println("the diff is------------000", diff)
		Expect(diff.Int64()).To(BeNumerically("~", taxIncheqd.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: taxIncheqd.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnPotionInUsdToCheq.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution events")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: feeCollectorReward.String(), Index: true},
				},
			},
		))
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.ORACLE_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: oracleRewardCoin.String(), Index: true},
				},
			},
		))
	})

	It("should tax create diddoc message with feegrant - case: fixed fee", func() {
		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)

		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom)))
		Expect(err).To(BeNil())

		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.CreateDid[0].MaxAmount
		userFee := sdk.NewCoins(sdk.NewCoin(types.BaseMinimalDenom, *tax))
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		convertedFees, err := ante.GetFeeForMsg(userFee, feeParams.CreateDid, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFeesInCheq, err := posthandler.ConvertToCheq(convertedFees, cheqPrice.Price)
		Expect(err).To(BeNil())

		burnPortionCheq := helpers.GetBurnFeePortion(feeParams.BurnFactor, convertedFeesInCheq)
		rewardPortionCheq := helpers.GetRewardPortion(convertedFeesInCheq, burnPortionCheq)

		totalTax := burnPortionCheq.Add(rewardPortionCheq...)
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)

		// Allowing a small tolerance of 1,000,000 ncheq (0.001 CHEQ)
		Expect(diff.Int64()).To(BeNumerically("~", totalTax.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))
		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should tax update diddoc message with feegrant - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

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

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		fees := feeParams.UpdateDid[0].MinAmount.Mul(math.NewInt(2))

		By("submitting an update diddoc message")
		resp, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(fees.String()+feeParams.UpdateDid[0].Denom)))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		tax := feeParams.UpdateDid[0].MinAmount.Mul(math.NewInt(2))
		userFee := sdk.NewCoins(sdk.NewCoin(types.BaseMinimalDenom, tax))
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		convertedFees, err := ante.GetFeeForMsg(userFee, feeParams.UpdateDid, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFeesinCheq, err := posthandler.ConvertToCheq(convertedFees, cheqPrice.Price)
		Expect(err).To(BeNil())
		Expect(tax).ToNot(BeNil())
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFeesinCheq.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))
	})

	It("should tax deactivate diddoc message with feegrant - case: fixed fee", func() {
		By("submitting a create diddoc message")
		resp, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid[0].MaxAmount.String()+feeParams.CreateDid[0].Denom))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("creating a feegrant")
		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)

		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
		By("querying the fee granter account balance before the transaction")
		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance before the transaction")
		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("fetching cheq EMA price and computing fees")
		tax := feeParams.DeactivateDid[0]
		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price
		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, *tax.MaxAmount))
		FeeforMsg, err := ante.GetFeeForMsg(userFee, feeParams.DeactivateDid, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFees, err := posthandler.ConvertToCheq(FeeforMsg, cheqp)
		Expect(err).To(BeNil())

		burnPotionInUsd := helpers.GetBurnFeePortion(feeParams.BurnFactor, FeeforMsg)
		rewardPortionInUsd := helpers.GetRewardPortion(FeeforMsg, burnPotionInUsd)

		burnPotionInUsdToCheq, err := posthandler.ConvertToCheq(burnPotionInUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionInUsdToCheq, err := posthandler.ConvertToCheq(rewardPortionInUsd, cheqp)
		Expect(err).To(BeNil())

		coin := rewardPortionInUsdToCheq.AmountOf(types.BaseMinimalDenom)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%
		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		taxIncheqd := burnPotionInUsdToCheq.Add(rewardPortionInUsdToCheq...)

		By("submitting a deactivate diddoc message")
		resp, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.MaxAmount.String()+tax.Denom)))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		By("querying the fee granter account balance after the transaction")
		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("querying the fee grantee account balance after the transaction")
		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking the granter balance difference")
		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFees.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		By("checking the grantee balance difference")
		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		By("exporting a readable tx event log")
		txResp, err := cli.QueryTxn(resp.TxHash)
		Expect(err).To(BeNil())
		events := helpers.ReadableEvents(txResp.Events)

		By("ensuring the events contain the expected tax event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "tx",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "fee", Value: taxIncheqd.String(), Index: true},
					{Key: "fee_payer", Value: testdata.BASE_ACCOUNT_4_ADDR, Index: true},
				},
			},
		))

		By("ensuring the events contain the expected supply deflation event")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "burn",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "burner", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: burnPotionInUsdToCheq.String(), Index: true},
				},
			},
		))

		By("ensuring the events contain the expected reward distribution events")
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.FEE_COLLECTOR_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: feeCollectorReward.String(), Index: true},
				},
			},
		))
		Expect(events).To(ContainElement(
			helpers.HumanReadableEvent{
				Type: "transfer",
				Attributes: []helpers.HumanReadableEventAttribute{
					{Key: "recipient", Value: testdata.ORACLE_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "sender", Value: testdata.DID_MODULE_ACCOUNT_ADDR, Index: true},
					{Key: "amount", Value: oracleRewardCoin.String(), Index: true},
				},
			},
		))

		By("revoking the feegrant")
		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)

		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
})
