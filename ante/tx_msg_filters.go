package ante

import (
	"errors"
	"fmt"
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

func GetTaxableMsgFeeWithBurnPortion(ctx sdk.Context, msg interface{}, ncheqPrice sdkmath.LegacyDec) (sdk.Coins, sdk.Coins, bool, error) {
	switch msg := msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		fee, err := GetFeeForMsg(TaxableMsgFees[MsgCreateDidDoc], ncheqPrice)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *didtypes.MsgUpdateDidDoc:
		fee, err := GetFeeForMsg(TaxableMsgFees[MsgUpdateDidDoc], ncheqPrice)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *didtypes.MsgDeactivateDidDoc:
		fee, err := GetFeeForMsg(TaxableMsgFees[MsgDeactivateDidDoc], ncheqPrice)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *resourcetypes.MsgCreateResource:
		return GetResourceTaxableMsgFee(ctx, msg, ncheqPrice)
	default:
		return nil, nil, false, nil
	}
}

func GetRewardPortion(total sdk.Coins, burnPortion sdk.Coins) sdk.Coins {
	if burnPortion.IsZero() {
		return total
	}
	return total.Sub(burnPortion...)
}

func GetResourceTaxableMsgFee(ctx sdk.Context, msg *resourcetypes.MsgCreateResource, ncheqPrice sdkmath.LegacyDec) (sdk.Coins, sdk.Coins, bool, error) {
	mediaType := resourceutils.DetectMediaType(msg.GetPayload().ToResource().Resource.Data)

	// Mime type image
	if strings.HasPrefix(mediaType, "image/") {
		fee, err := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceImage], ncheqPrice)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	}

	// Mime type json
	if strings.HasPrefix(mediaType, "application/json") {
		fee, err := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceJSON], ncheqPrice)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	}

	fee, err := GetFeeForMsg(TaxableMsgFees[MsgCreateResourceDefault], ncheqPrice)
	if err != nil {
		return nil, nil, true, err
	}

	// Default mime type
	burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
	return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
}

func checkFeeParamsFromSubspace(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper) bool {
	didParams, err := didKeeper.GetParams(ctx)
	if err != nil {
		return false
	}
	TaxableMsgFees[MsgCreateDidDoc] = didParams.CreateDid
	TaxableMsgFees[MsgUpdateDidDoc] = didParams.UpdateDid
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

func IsTaxableTx(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx, oracleKeeper OracleKeeper) (bool, sdk.Coins, sdk.Coins, error) {
	ncheqPrice, exist := oracleKeeper.GetEMA(ctx, oracletypes.CheqdSymbol)
	if !exist {
		// fallback to fixed fee range in ncheq if defined
		ncheqPrice = sdkmath.LegacyZeroDec() // zero value, GetFeeForMsg will handle fallback
	}

	_ = checkFeeParamsFromSubspace(ctx, didKeeper, resourceKeeper)
	reward := (sdk.Coins)(nil)
	burn := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {

		rewardPortion, burnPortion, isIdentityMsg, err := GetTaxableMsgFeeWithBurnPortion(ctx, msg, ncheqPrice)
		if err != nil {
			return true, nil, nil, err
		}
		if !isIdentityMsg {
			continue
		}
		if rewardPortion != nil {
			reward = reward.Add(rewardPortion...)
			burn = burn.Add(burnPortion...)
		}
	}

	if !reward.IsZero() {
		return true, reward, burn, nil
	}

	return false, nil, nil, nil
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

func GetFeeForMsg(feeRanges []didtypes.FeeRange, cheqEmaPrice sdkmath.LegacyDec) (sdk.Coins, error) {
	type usdRange struct {
		denom   string
		minUSD  sdkmath.LegacyDec
		maxUSD  *sdkmath.LegacyDec
		minCoin sdkmath.Int
		maxCoin *sdkmath.Int
	}

	if len(feeRanges) == 0 {
		return nil, errors.New("fee ranges empty")
	}

	// Fallback: If CHEQ price is not available, try to return ncheq fixed fee
	if cheqEmaPrice.IsZero() {
		for _, fr := range feeRanges {
			if fr.Denom == oracletypes.CheqdDenom {
				// Use MinAmount directly as fallback fixed fee
				return sdk.NewCoins(sdk.NewCoin(oracletypes.CheqdDenom, fr.MinAmount)), nil
			}
		}
		return nil, errors.New("cheq price not available and no ncheq fallback fee defined")
	}

	var ranges []usdRange

	// Step 1: Convert fee ranges into USD values
	for _, fr := range feeRanges {
		var minUSD sdkmath.LegacyDec
		var maxUSD *sdkmath.LegacyDec

		switch fr.Denom {
		case oracletypes.CheqdDenom:
			minCHEQ := sdkmath.LegacyNewDecFromInt(fr.MinAmount).QuoInt64(1e9)
			minUSD = cheqEmaPrice.Mul(minCHEQ)

			if fr.MaxAmount != nil {
				maxCHEQ := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).QuoInt64(1e9)
				val := cheqEmaPrice.Mul(maxCHEQ)
				maxUSD = &val
			}
		case "usd":
			if fr.MinAmount.LT(sdkmath.NewInt(1e6)) {
				minUSD = sdkmath.LegacyNewDecFromInt(fr.MinAmount)
				if fr.MaxAmount != nil {
					val := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount)
					maxUSD = &val
				}
			} else {
				minUSD = sdkmath.LegacyNewDecFromInt(fr.MinAmount).QuoInt64(1e18)
				if fr.MaxAmount != nil {
					val := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).QuoInt64(1e18)
					maxUSD = &val
				}
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
		return nil, errors.New("no valid fee ranges could be converted")
	}

	// Step 2: Find overlapping USD range
	overlapMin := ranges[0].minUSD
	overlapMax := ranges[0].maxUSD

	for _, r := range ranges[1:] {
		if r.minUSD.GT(overlapMin) {
			overlapMin = r.minUSD
		}
		if r.maxUSD != nil {
			if overlapMax == nil || r.maxUSD.LT(*overlapMax) {
				val := *r.maxUSD
				overlapMax = &val
			}
		}
	}

	if overlapMax != nil && overlapMin.GT(*overlapMax) {
		return nil, errors.New("no valid overlapping range")
	}

	// Step 3: Prefer USD denom
	var chosen usdRange
	for _, r := range ranges {
		if r.denom == "usd" {
			chosen = r
			break
		}
	}
	if chosen.denom == "" {
		chosen = ranges[0]
	}

	// Step 4: Compute final amount
	var finalAmount sdkmath.Int
	switch chosen.denom {
	case oracletypes.CheqdDenom:
		finalAmount = overlapMin.Quo(cheqEmaPrice).TruncateInt().MulRaw(1e9)
	case "usd":
		if overlapMin.LT(sdkmath.LegacyNewDec(1e5)) {
			finalAmount = overlapMin.MulInt64(1e18).TruncateInt()
		} else {
			finalAmount = overlapMin.TruncateInt()
		}
	default:
		return nil, fmt.Errorf("unsupported denom selected: %s", chosen.denom)
	}

	// Step 5: Clamp
	if finalAmount.LT(chosen.minCoin) {
		finalAmount = chosen.minCoin
	} else if chosen.maxCoin != nil && finalAmount.GT(*chosen.maxCoin) {
		finalAmount = *chosen.maxCoin
	}

	return sdk.NewCoins(sdk.NewCoin(chosen.denom, finalAmount)), nil
}
