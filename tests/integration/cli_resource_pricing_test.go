//go:build integration

package integration

import (
	"crypto/ed25519"

	"cosmossdk.io/math"
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
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli - positive resource pricing", func() {
	var tmpDir string
	var didFeeParams types.FeeParams
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
		useMin := false
		tax, err := cli.ResolveFeeFromParams(didFeeParams.CreateDid, useMin)
		Expect(err).To(BeNil())
		resp, err := cli.CreateDidDoc(tmpDir, didPayload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
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

		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Json, useMin)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		cheqp := cheqPrice.Price
		Expect(err).To(BeNil())
		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Json, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFeesToCheq, err := posthandler.ConvertToCheq(convertedFees, cheqp)
		Expect(err).To(BeNil())

		burnPotionInUsd := helpers.GetBurnFeePortion(resourceFeeParams.BurnFactor, convertedFees)
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

		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFeesToCheq.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

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

		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Image, useMin)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		cheqp := cheqPrice.Price
		Expect(err).To(BeNil())

		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Image, cheqp, nil)
		Expect(err).To(BeNil())

		convertedFeesIncheq, err := posthandler.ConvertToCheq(convertedFees, cheqp)
		Expect(err).To(BeNil())

		// Calculate fee portions in USD

		burnPotionInUsd := helpers.GetBurnFeePortion(resourceFeeParams.BurnFactor, convertedFees)
		rewardPortionInUsd := helpers.GetRewardPortion(convertedFees, burnPotionInUsd)

		// Convert USD portions back to ncheq tokens
		burnPotionInUsdToCheq, err := posthandler.ConvertToCheq(burnPotionInUsd, cheqp)
		Expect(err).To(BeNil())

		rewardPortionInUsdToCheq, err := posthandler.ConvertToCheq(rewardPortionInUsd, cheqp)
		Expect(err).To(BeNil())

		// Calculate oracle and fee collector rewards
		coin := rewardPortionInUsdToCheq.AmountOf(types.BaseMinimalDenom)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%

		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		// Total tax amount (burn + reward)
		taxIncheqd := burnPotionInUsdToCheq.Add(rewardPortionInUsdToCheq...)

		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFeesIncheq.AmountOf(types.BaseMinimalDenom).Int64(), 200_000_000))

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
		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Default, useMin)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price

		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		convertedFees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Default, cheqp, nil)
		Expect(err).To(BeNil())
		burnPortionUsd := helpers.GetBurnFeePortion(resourceFeeParams.BurnFactor, convertedFees)
		rewardPortionUsd := helpers.GetRewardPortion(convertedFees, burnPortionUsd)

		burnPortionCheq, err := posthandler.ConvertToCheq(burnPortionUsd, cheqp)
		Expect(err).To(BeNil())
		rewardPortionCheq, err := posthandler.ConvertToCheq(rewardPortionUsd, cheqp)
		Expect(err).To(BeNil())
		finalPrice, err := posthandler.ConvertToCheq(convertedFees, cheqp)
		Expect(err).To(BeNil())

		coin := rewardPortionCheq.AmountOf(types.BaseMinimalDenom)

		oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.5%

		oracleReward := oracleShareRate.MulInt(coin).TruncateInt()
		oracleRewardCoin := sdk.NewCoin(types.BaseMinimalDenom, oracleReward)
		feeCollectorReward := sdk.NewCoin(types.BaseMinimalDenom, coin.Sub(oracleReward))

		taxIncheqd := burnPortionCheq.Add(rewardPortionCheq...)

		res, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         resourceName,
			Version:      resourceVersion,
			ResourceType: resourceType,
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_4, helpers.GenerateFees(tax.Amount.String()+tax.Denom))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the altered account balance")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(balanceAfter.Denom).To(BeEquivalentTo(resourcetypes.BaseMinimalDenom))

		By("checking the balance difference")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", finalPrice.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

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
					{Key: "amount", Value: burnPortionCheq.String(), Index: true},
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
	It("should tax create resource json message with feegrant - case: fixed fee", func() {
		resourceID := uuid.NewString()
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Json, useMin)
		Expect(err).To(BeNil())

		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())

		cheqp := cheqPrice.Price

		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		fees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Json, cheqp, nil)

		Expect(err).To(BeNil())
		convertedFee, err := posthandler.ConvertToCheq(fees, cheqp)
		Expect(err).To(BeNil())

		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         "TestResource",
			Version:      "1.0",
			ResourceType: "TestType",
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.Amount.String()+tax.Denom)))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFee.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should tax create resource image with feegrant - case: fixed fee", func() {
		resourceID := uuid.NewString()
		resourceFile, err := testdata.CreateTestImage(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Image, useMin)
		Expect(err).To(BeNil())

		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price

		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		fees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Image, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFee, err := posthandler.ConvertToCheq(fees, cheqp)
		Expect(err).To(BeNil())

		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         "TestResource",
			Version:      "1.0",
			ResourceType: "TestType",
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.Amount.String()+tax.Denom)))

		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFee.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should tax create resource default with feegrant - case: fixed fee", func() {
		resourceID := uuid.NewString()
		resourceFile, err := testdata.CreateTestDefault(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		res, err := cli.GrantFees(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		granterBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		useMin := false
		tax, err := cli.ResolveFeeFromParams(resourceFeeParams.Default, useMin)
		Expect(err).To(BeNil())

		cheqPrice, err := cli.QueryWMA(types.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		cheqp := cheqPrice.Price

		userFee := sdk.NewCoins(sdk.NewCoin(tax.Denom, tax.Amount))

		fees, err := ante.GetFeeForMsg(userFee, resourceFeeParams.Default, cheqp, nil)
		Expect(err).To(BeNil())
		convertedFee, err := posthandler.ConvertToCheq(fees, cheqp)
		Expect(err).To(BeNil())

		resp, err := cli.CreateResource(tmpDir, resourcetypes.MsgCreateResourcePayload{
			CollectionId: collectionID,
			Id:           resourceID,
			Name:         "TestResource",
			Version:      "1.0",
			ResourceType: "TestType",
		}, signInputs, resourceFile, testdata.BASE_ACCOUNT_1, helpers.GenerateFeeGranter(testdata.BASE_ACCOUNT_4_ADDR, helpers.GenerateFees(tax.Amount.String()+tax.Denom)))
		Expect(err).To(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(0))

		granterBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_4_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		granteeBalanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, resourcetypes.BaseMinimalDenom)
		Expect(err).To(BeNil())

		diff := granterBalanceBefore.Amount.Sub(granterBalanceAfter.Amount)
		Expect(diff.Int64()).To(BeNumerically("~", convertedFee.AmountOf(types.BaseMinimalDenom).Int64(), 2_000_000))

		diff = granteeBalanceAfter.Amount.Sub(granteeBalanceBefore.Amount)
		Expect(diff.IsZero()).To(BeTrue())

		res, err = cli.RevokeFeeGrant(testdata.BASE_ACCOUNT_4_ADDR, testdata.BASE_ACCOUNT_1_ADDR, cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})
})
