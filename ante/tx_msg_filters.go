package ante

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgCreateDidDoc int = iota
	MsgUpdateDidDoc
	MsgDeactivateDidDoc
	MsgCreateResourceDefault
	MsgCreateResourceImage
	MsgCreateResourceJSON

	TaxableMsgFeeCount
)

const (
	BurnFactorDid int = iota
	BurnFactorResource

	BurnFactorCount
)

type TaxableMsgFee = [TaxableMsgFeeCount][]didtypes.FeeRange

type BurnFactor = [BurnFactorCount]sdkmath.LegacyDec

var TaxableMsgFees = TaxableMsgFee{
	MsgCreateDidDoc:          []didtypes.FeeRange{},
	MsgUpdateDidDoc:          []didtypes.FeeRange{},
	MsgDeactivateDidDoc:      []didtypes.FeeRange{},
	MsgCreateResourceDefault: []didtypes.FeeRange{},
	MsgCreateResourceImage:   []didtypes.FeeRange{},
	MsgCreateResourceJSON:    []didtypes.FeeRange{},
}

var BurnFactors = BurnFactor{
	BurnFactorDid:      sdkmath.LegacyNewDec(0),
	BurnFactorResource: sdkmath.LegacyNewDec(0),
}

func GetTaxableMsg(msg interface{}) bool {
	switch msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		return true
	case *didtypes.MsgUpdateDidDoc:
		return true
	case *didtypes.MsgDeactivateDidDoc:
		return true
	case *resourcetypes.MsgCreateResource:
		return true
	default:
		return false
	}
}

func GetTaxableMsgFeeWithBurnPortion(ctx sdk.Context, msg interface{}, ncheqPrice sdkmath.LegacyDec) (sdk.Coins, sdk.Coins, bool) {
	switch msg := msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		fee := GetFeeForMsg(TaxableMsgFees[MsgCreateDidDoc], ncheqPrice)
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true
	case *didtypes.MsgUpdateDidDoc:
		fee := GetFeeForMsg(TaxableMsgFees[MsgUpdateDidDoc], ncheqPrice)
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true
	case *didtypes.MsgDeactivateDidDoc:
		fee := GetFeeForMsg(TaxableMsgFees[MsgDeactivateDidDoc], ncheqPrice)
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true
	case *resourcetypes.MsgCreateResource:
		return GetResourceTaxableMsgFee(ctx, msg, ncheqPrice)
	default:
		return nil, nil, false
	}
}

func GetRewardPortion(total sdk.Coins, burnPortion sdk.Coins) sdk.Coins {
	if burnPortion.IsZero() {
		return total
	}
	return total.Sub(burnPortion...)
}

func GetResourceTaxableMsgFee(ctx sdk.Context, msg *resourcetypes.MsgCreateResource, ncheqPrice sdkmath.LegacyDec) (sdk.Coins, sdk.Coins, bool) {
	mediaType := resourceutils.DetectMediaType(msg.GetPayload().ToResource().Resource.Data)

	// Mime type image
	if strings.HasPrefix(mediaType, "image/") {
		fee := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceImage], ncheqPrice)
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true
	}

	// Mime type json
	if strings.HasPrefix(mediaType, "application/json") {
		fee := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceJSON], ncheqPrice)
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true
	}

	fee := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceDefault], ncheqPrice)

	// Default mime type
	burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
	return GetRewardPortion(fee, burnPortion), burnPortion, true
}

func checkFeeParamsFromSubspace(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper) bool {
	didParams, err := didKeeper.GetParams(ctx)
	if err != nil {
		return false
	}
	TaxableMsgFees[MsgCreateDidDoc] = didParams.CreateDid
	TaxableMsgFees[MsgUpdateDidDoc] = didParams.DeactivateDid
	TaxableMsgFees[MsgDeactivateDidDoc] = didParams.DeactivateDid

	resourceParams, err := resourceKeeper.GetParams(ctx)
	if err != nil {
		return false
	}
	TaxableMsgFees[MsgCreateResourceImage] = resourceParams.Image
	TaxableMsgFees[MsgCreateResourceJSON] = resourceParams.Json
	TaxableMsgFees[MsgCreateResourceDefault] = resourceParams.Default

	BurnFactors[BurnFactorDid] = didParams.BurnFactor
	BurnFactors[BurnFactorResource] = resourceParams.BurnFactor

	return true
}

func IsTaxableTx(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx, oracleKeeper OracleKeeper) (bool, sdk.Coins, sdk.Coins) {
	ncheqPrice, exist := oracleKeeper.GetEMA(ctx, oracletypes.CheqdSymbol)
	if !exist {
		return false, nil, nil
	}

	_ = checkFeeParamsFromSubspace(ctx, didKeeper, resourceKeeper)
	reward := (sdk.Coins)(nil)
	burn := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		rewardPortion, burnPortion, isIdentityMsg := GetTaxableMsgFeeWithBurnPortion(ctx, msg, ncheqPrice)
		if !isIdentityMsg {
			continue
		}
		if rewardPortion != nil {
			reward = reward.Add(rewardPortion...)
			burn = burn.Add(burnPortion...)
		}
	}

	if !reward.IsZero() {
		return true, reward, burn
	}

	return false, nil, nil
}

func IsTaxableTxLite(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		if GetTaxableMsg(msg) {
			return true
		}
	}

	return false
}

func GetFeeForMsg(feeRanges []didtypes.FeeRange, cheqEmaPrice sdkmath.LegacyDec) sdk.Coins {
	if len(feeRanges) == 0 || cheqEmaPrice.IsZero() {
		return nil
	}

	type usdRange struct {
		denom   string
		minUSD  sdkmath.LegacyDec
		maxUSD  sdkmath.LegacyDec
		minCoin sdkmath.Int
		maxCoin sdkmath.Int
	}

	var ranges []usdRange

	// Step 1: Convert each denom’s fee range into USD values
	for _, fr := range feeRanges {
		var minUSD, maxUSD sdkmath.LegacyDec

		switch fr.Denom {
		case "ncheq":
			// Convert native CHEQ to USD using EMA price
			minCHEQ := sdkmath.LegacyNewDecFromInt(fr.MinAmount).QuoInt64(1e9)
			maxCHEQ := sdkmath.LegacyNewDecFromInt(fr.MaxAmount).QuoInt64(1e9)

			minUSD = cheqEmaPrice.MulInt(sdkmath.NewInt(minCHEQ.TruncateInt64()))
			maxUSD = cheqEmaPrice.MulInt(sdkmath.NewInt(maxCHEQ.TruncateInt64()))
		case "usd":
			// Treat USD as 18-decimal fixed point
			if fr.MinAmount.LT(sdkmath.NewInt(1e6)) {
				minUSD = sdkmath.LegacyNewDecFromInt(fr.MinAmount)
				maxUSD = sdkmath.LegacyNewDecFromInt(fr.MaxAmount)
			} else {
				// Assume it's 18-dec format
				minUSD = sdkmath.LegacyNewDecFromInt(fr.MinAmount).QuoInt64(1e18)
				maxUSD = sdkmath.LegacyNewDecFromInt(fr.MaxAmount).QuoInt64(1e18)
			}
		default:
			continue
		}

		ranges = append(ranges, usdRange{
			denom:   fr.Denom,
			minUSD:  minUSD,
			maxUSD:  maxUSD,
			minCoin: fr.MinAmount,
			maxCoin: fr.MaxAmount,
		})
	}

	if len(ranges) == 0 {
		return nil
	}

	// Step 2: Find overlap: [max(minA, minB...), min(maxA, maxB...)]
	overlapMin := ranges[0].minUSD
	overlapMax := ranges[0].maxUSD

	for _, r := range ranges[1:] {
		if r.minUSD.GT(overlapMin) {
			overlapMin = r.minUSD
		}
		if r.maxUSD.LT(overlapMax) {
			overlapMax = r.maxUSD
		}
	}

	if overlapMin.GT(overlapMax) {
		return nil // No valid overlapping range
	}

	// Step 3: Choose a denom to return the fee in — using the first one
	var chosen usdRange
	for _, r := range ranges {
		if r.denom == "usd" {
			chosen = r
			break
		}
	}
	var finalAmount sdkmath.Int
	// fallback to first if usd is not available
	if chosen.denom == "" {
		chosen = ranges[0]
	}

	switch chosen.denom {
	case "ncheq":
		finalAmount = overlapMin.Quo(cheqEmaPrice).TruncateInt()

	case "usd":
		if overlapMin.LT(sdkmath.LegacyNewDec(1e5)) {
			finalAmount = overlapMin.MulInt64(1e18).TruncateInt()
		} else {
			finalAmount = overlapMin.TruncateInt()
		}

	default:
		return nil
	}

	// Clamp to denom-specific coin range
	if finalAmount.LT(chosen.minCoin) {
		finalAmount = chosen.minCoin
	} else if finalAmount.GT(chosen.maxCoin) {
		finalAmount = chosen.maxCoin
	}

	return sdk.NewCoins(sdk.NewCoin(chosen.denom, finalAmount))
}
