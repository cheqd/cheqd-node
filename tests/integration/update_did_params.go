package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"cosmossdk.io/math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// var _ = Describe("Integration - update the did params", func() {
// 	var (
// 		oracleParams *oracletypes.QueryParamsResponse
// 		wmaPrice     math.LegacyDec
// 		Proposal_id  string
// 	)

// 	It("Wait for CHEQ price", func() {
// 		By("querying oracle params")
// 		var err error
// 		oracleParamsResp, err := cli.QueryOracleParams()
// 		Expect(err).To(BeNil())
// 		oracleParams = &oracleParamsResp

// 		// Wait for ComputeAllAverages to naturally trigger
// 		historicStampPeriod := oracleParams.Params.HistoricStampPeriod
// 		averagingWindow := oraclekeeper.AveragingWindow
// 		targetHeight := int64(historicStampPeriod) * int64(averagingWindow)

// 		fmt.Printf("⏳ Waiting until block height ≥ %d to trigger ComputeAllAverages...\n", targetHeight)

// 		for {
// 			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
// 			Expect(err).To(BeNil())

// 			if currentHeight >= targetHeight {
// 				fmt.Printf("✅ Reached block height %d — proceeding...\n", currentHeight)
// 				break
// 			}

// 			fmt.Printf("  → Current height: %d (waiting for %d)\n", currentHeight, targetHeight)
// 			time.Sleep(2 * time.Second)
// 		}

// 		By("querying updated WMA oracle price")
// 		wmaRes, err := cli.QueryWMA(didtypes.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
// 		Expect(err).To(BeNil())
// 		wmaPrice = wmaRes.Price
// 	})

// 	It("should generate update-did-fee-params.json with usd.min = ncheq.min * price", func() {
// 		By("querying current DID fee params")
// 		didParams, err := cli.QueryDidParams()
// 		Expect(err).To(BeNil())

// 		convertNcheqToUSD := func(ncheqAmt math.Int) math.Int {
// 			return math.LegacyNewDecFromInt(ncheqAmt). // ncheq as LegacyDec
// 									Quo(math.LegacyNewDec(1_000_000_000)).             // normalize 9 decimals
// 									Mul(wmaPrice).                                     // LegacyDec from oracle
// 									Mul(math.LegacyNewDec(1_000_000_000_000_000_000)). // scale to 18 decimals
// 									TruncateInt()                                      // back to Int
// 		}

// 		double := func(i math.Int) *math.Int {
// 			d := i.MulRaw(2)
// 			return &d
// 		}

// 		transform := func(fees []didtypes.FeeRange) []didtypes.FeeRange {
// 			var ncheqMin math.Int
// 			foundNcheq := false

// 			for _, fee := range fees {
// 				if fee.Denom == "ncheq" && fee.MinAmount != nil {
// 					ncheqMin = *fee.MinAmount
// 					foundNcheq = true
// 					break
// 				}
// 			}

// 			if !foundNcheq {
// 				ncheqMin = math.NewInt(50_000_000_000) // fallback: 50 CHEQ
// 			}

// 			ncheqMax := double(ncheqMin)
// 			usdMin := convertNcheqToUSD(ncheqMin)

// 			return []didtypes.FeeRange{
// 				{
// 					Denom:     "ncheq",
// 					MinAmount: &ncheqMin,
// 					MaxAmount: ncheqMax,
// 				},
// 				{
// 					Denom:     "usd",
// 					MinAmount: &usdMin,
// 					MaxAmount: &usdMin,
// 				},
// 			}
// 		}

// 		newParams := didtypes.FeeParams{
// 			CreateDid:     transform(didParams.Params.CreateDid),
// 			UpdateDid:     transform(didParams.Params.UpdateDid),
// 			DeactivateDid: transform(didParams.Params.DeactivateDid),
// 			BurnFactor:    didParams.Params.BurnFactor,
// 		}
// 		proposal := map[string]interface{}{
// 			"messages": []interface{}{
// 				map[string]interface{}{
// 					"@type":     "/cheqd.did.v2.MsgUpdateParams",
// 					"authority": "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
// 					"params":    newParams,
// 				},
// 			},
// 			"deposit": "10000000000ncheq",
// 			"title":   "params update",
// 			"summary": "update params using gov authority address",
// 		}

// 		outPath := "update-did-fee-params.json"
// 		jsonBytes, err := json.MarshalIndent(proposal, "", "  ")
// 		Expect(err).To(BeNil())
// 		err = os.WriteFile(outPath, jsonBytes, 0644)
// 		Expect(err).To(BeNil())
// 		fmt.Printf("✅ Fee params proposal written to: %s\n", outPath)
// 	})

// 	It("should submit a mint proposal ", func() {
// 		By("passing the proposal file to the container")
// 		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join("update-did-fee-params.json"), cli.DockerHome, cli.Validator0)
// 		Expect(err).To(BeNil())

// 		By("sending a SubmitProposal transaction from `validator0` container")
// 		res, err := cli.SubmitProposalTx(cli.Operator0, "update-did-fee-params.json", cli.CliGasParams)
// 		Expect(err).To(BeNil())
// 		res, err = cli.QueryTxn(res.TxHash)
// 		Expect(err).To(BeNil())

// 		proposal_id, err := cli.GetProposalID(res.Events)
// 		Proposal_id = proposal_id
// 		Expect(err).To(BeNil())

// 		Expect(res.Code).To(BeEquivalentTo(0))
// 	})

// 	It("should wait for the proposal submission to be included in a block", func() {
// 		By("getting the current block height")
// 		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
// 		Expect(err).To(BeNil())

// 		By("waiting for the proposal to be included in a block")
// 		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, 2)
// 		Expect(err).To(BeNil())
// 	})

// 	It("should vote for the mint proposal from `validator1` container", func() {
// 		By("sending a VoteProposal transaction from `validator1` container")
// 		res, err := cli.VoteProposalTx(cli.Operator1, Proposal_id, "yes", cli.CliGasParams)
// 		Expect(err).To(BeNil())
// 		Expect(res.Code).To(BeEquivalentTo(0))
// 	})

// 	It("should vote for the mint proposal from `validator2` container", func() {
// 		By("sending a VoteProposal transaction from `validator2` container")
// 		res, err := cli.VoteProposalTx(cli.Operator2, Proposal_id, "yes", cli.CliGasParams)
// 		Expect(err).To(BeNil())
// 		Expect(res.Code).To(BeEquivalentTo(0))
// 	})

// 	It("should vote for the mint proposal from `validator0` container", func() {
// 		By("sending a VoteProposal transaction from `validator0` container")
// 		res, err := cli.VoteProposalTx(cli.Operator0, Proposal_id, "yes", cli.CliGasParams)
// 		Expect(err).To(BeNil())
// 		Expect(res.Code).To(BeEquivalentTo(0))
// 	})

// 	It("should vote for the mint proposal from `validator0` container", func() {
// 		By("sending a VoteProposal transaction from `validator0` container")
// 		res, err := cli.VoteProposalTx(cli.Operator3, Proposal_id, "yes", cli.CliGasParams)
// 		Expect(err).To(BeNil())
// 		res, err = cli.QueryTxn(res.TxHash)
// 		Expect(err).To(BeNil())
// 		Expect(res.Code).To(BeEquivalentTo(0))
// 	})

// 	It("should wait for the proposal to pass", func() {
// 		By("getting the current block height")
// 		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
// 		Expect(err).To(BeNil())

// 		By("waiting for the proposal to pass")
// 		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, 25)
// 		Expect(err).To(BeNil())
// 	})

// 	It("should check the proposal status to ensure it has passed", func() {
// 		By("sending a QueryProposal query from `validator0` container")
// 		proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
// 		Expect(err).To(BeNil())
// 		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
// 	})
// })

var _ = Describe("Integration - update resource params", func() {
	var (
		oracleParams *oracletypes.QueryParamsResponse
		wmaPrice     math.LegacyDec
		Proposal_id  string
	)

	It("Wait for CHEQ price", func() {
		By("querying oracle params")
		var err error
		oracleParamsResp, err := cli.QueryOracleParams()
		Expect(err).To(BeNil())
		oracleParams = &oracleParamsResp

		// Wait for ComputeAllAverages to naturally trigger
		historicStampPeriod := oracleParams.Params.HistoricStampPeriod
		averagingWindow := oraclekeeper.AveragingWindow
		targetHeight := int64(historicStampPeriod) * int64(averagingWindow)

		fmt.Printf("⏳ Waiting until block height ≥ %d to trigger ComputeAllAverages...\n", targetHeight)

		for {
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			if currentHeight >= targetHeight {
				fmt.Printf("✅ Reached block height %d — proceeding...\n", currentHeight)
				break
			}

			fmt.Printf("  → Current height: %d (waiting for %d)\n", currentHeight, targetHeight)
			time.Sleep(2 * time.Second)
		}

		By("querying updated WMA oracle price")
		wmaRes, err := cli.QueryWMA(didtypes.BaseDenom, string(oraclekeeper.WmaStrategyBalanced), nil)
		Expect(err).To(BeNil())
		wmaPrice = wmaRes.Price
	})

	It("should generate update-did-fee-params.json with usd.min = ncheq.min * price", func() {
		By("querying current DID fee params")
		resourceParams, err := cli.QueryResourceParams()
		Expect(err).To(BeNil())

		convertNcheqToUSD := func(ncheqAmt math.Int) math.Int {
			return math.LegacyNewDecFromInt(ncheqAmt). // ncheq as LegacyDec
									Quo(math.LegacyNewDec(1_000_000_000)).             // normalize 9 decimals
									Mul(wmaPrice).                                     // LegacyDec from oracle
									Mul(math.LegacyNewDec(1_000_000_000_000_000_000)). // scale to 18 decimals
									TruncateInt()                                      // back to Int
		}

		double := func(i math.Int) *math.Int {
			d := i.MulRaw(2)
			return &d
		}

		transform := func(fees []didtypes.FeeRange) []didtypes.FeeRange {
			var ncheqMin math.Int
			foundNcheq := false

			for _, fee := range fees {
				if fee.Denom == "ncheq" && fee.MinAmount != nil {
					ncheqMin = *fee.MinAmount
					foundNcheq = true
					break
				}
			}

			if !foundNcheq {
				ncheqMin = math.NewInt(50_000_000_000) // fallback: 50 CHEQ
			}

			ncheqMax := double(ncheqMin)
			usdMin := convertNcheqToUSD(ncheqMin)

			return []didtypes.FeeRange{
				{
					Denom:     "ncheq",
					MinAmount: &ncheqMin,
					MaxAmount: ncheqMax,
				},
				{
					Denom:     "usd",
					MinAmount: &usdMin,
					MaxAmount: &usdMin,
				},
			}
		}

		newParams := resourcetypes.FeeParams{
			Default:    transform(resourceParams.Params.Default),
			Json:       transform(resourceParams.Params.Json),
			Image:      transform(resourceParams.Params.Image),
			BurnFactor: resourceParams.Params.BurnFactor,
		}
		proposal := map[string]interface{}{
			"messages": []interface{}{
				map[string]interface{}{
					"@type":     "/cheqd.resource.v2.MsgUpdateParams",
					"authority": "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
					"params":    newParams,
				},
			},
			"deposit": "10000000000ncheq",
			"title":   "params update",
			"summary": "update params using gov authority address",
		}

		outPath := "update-resource-fee-params.json"
		jsonBytes, err := json.MarshalIndent(proposal, "", "  ")
		Expect(err).To(BeNil())
		err = os.WriteFile(outPath, jsonBytes, 0644)
		Expect(err).To(BeNil())
		fmt.Printf("✅ Fee params proposal written to: %s\n", outPath)
	})

	It("should submit a mint proposal ", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join("update-resource-fee-params.json"), cli.DockerHome, cli.Validator0)
		Expect(err).To(BeNil())

		By("sending a SubmitProposal transaction from `validator0` container")
		res, err := cli.SubmitProposalTx(cli.Operator0, "update-resource-fee-params.json", cli.CliGasParams)
		Expect(err).To(BeNil())
		res, err = cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		proposal_id, err := cli.GetProposalID(res.Events)
		Proposal_id = proposal_id
		Expect(err).To(BeNil())

		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal submission to be included in a block", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to be included in a block")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, 2)
		Expect(err).To(BeNil())
	})

	It("should vote for the mint proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposalTx(cli.Operator1, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposalTx(cli.Operator2, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposalTx(cli.Operator0, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposalTx(cli.Operator3, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		res, err = cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, 25)
		Expect(err).To(BeNil())
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})
})
