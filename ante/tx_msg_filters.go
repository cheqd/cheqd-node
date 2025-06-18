package ante

import (
	"errors"
	"fmt"
	"math"
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
		minUSD  *sdkmath.Int
		maxUSD  *sdkmath.Int
		minCoin *sdkmath.Int
		maxCoin *sdkmath.Int
	}

	if len(feeRanges) == 0 {
		return nil, errors.New("fee ranges empty")
	}
	// Fallback: If CHEQ price is not available, return ncheq fixed fee
	if cheqEmaPrice.IsZero() {
		for _, fr := range feeRanges {
			if fr.Denom != oracletypes.CheqdDenom {
				continue
			}
			switch {
			case fr.MinAmount != nil:
				return sdk.NewCoins(sdk.NewCoin(oracletypes.CheqdDenom, *fr.MinAmount)), nil
			case fr.MaxAmount != nil:
				return sdk.NewCoins(sdk.NewCoin(oracletypes.CheqdDenom, *fr.MaxAmount)), nil
			default:
				return nil, errors.New("cheq price not available and no valid fee fallback")
			}
		}
		return nil, errors.New("cheq price not available and no ncheq fallback fee defined")
	}

	const cheqExponent = 9
	const usdExponent = 6

	cheqScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(cheqExponent)))
	usdScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(usdExponent)))
	usdFrom18To6 := sdkmath.NewInt(1_000_000_000_000) // 1e12

	var ranges []usdRange

	// Step 1: Convert all fee ranges into scaled USD ints
	for _, fr := range feeRanges {
		if fr.MinAmount == nil && fr.MaxAmount == nil {
			continue
		}

		var minUSD *sdkmath.Int
		var maxUSD *sdkmath.Int

		switch fr.Denom {
		case oracletypes.CheqdDenom:
			if fr.MinAmount != nil {
				cheqDec := sdkmath.LegacyNewDecFromInt(*fr.MinAmount).QuoInt(cheqScale)
				usdVal := cheqDec.Mul(cheqEmaPrice).MulInt(usdScale).TruncateInt()
				minUSD = &usdVal
			}
			if fr.MaxAmount != nil {
				cheqDec := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).QuoInt(cheqScale)
				usdVal := cheqDec.Mul(cheqEmaPrice).MulInt(usdScale).TruncateInt()
				maxUSD = &usdVal
			}

		case oracletypes.UsdDenom:
			if fr.MinAmount != nil {
				val := fr.MinAmount.Quo(usdFrom18To6)
				minUSD = &val
			}
			if fr.MaxAmount != nil {
				val := fr.MaxAmount.Quo(usdFrom18To6)
				maxUSD = &val
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

	// Step 2: Find overlapping USD range (support free-ranging)
	var overlapMin *sdkmath.Int
	var overlapMax *sdkmath.Int

	for i, r := range ranges {
		if i == 0 {
			overlapMin = r.minUSD
			overlapMax = r.maxUSD
		} else {
			if r.minUSD != nil {
				if overlapMin == nil || r.minUSD.GT(*overlapMin) {
					overlapMin = r.minUSD
				}
			}
			if r.maxUSD != nil {
				if overlapMax == nil || r.maxUSD.LT(*overlapMax) {
					overlapMax = r.maxUSD
				}
			}
		}
	}

	// Reject only if defined overlap is invalid
	if overlapMin != nil && overlapMax != nil && overlapMin.GT(*overlapMax) {
		return nil, errors.New("no valid overlapping USD range")
	}

	// Step 3: Pick preferred range (prefer USD if available)
	var chosen usdRange
	for _, r := range ranges {
		if r.denom == oracletypes.UsdDenom {
			chosen = r
			break
		}
	}
	if chosen.denom == "" {
		chosen = ranges[0]
	}

	// Step 4: Pick target USD value to convert
	var usdToUse *sdkmath.Int
	if overlapMin != nil {
		usdToUse = overlapMin
	} else if overlapMax != nil {
		usdToUse = overlapMax
	} else {
		return nil, errors.New("cannot determine fee: no min or max USD bound")
	}

	var finalAmount sdkmath.Int
	switch chosen.denom {
	case oracletypes.CheqdDenom:
		overlapDec := sdkmath.LegacyNewDecFromInt(*usdToUse)
		cheqAmount := overlapDec.Quo(cheqEmaPrice).MulInt(cheqScale).QuoInt(usdScale).TruncateInt()
		finalAmount = cheqAmount

	case oracletypes.UsdDenom:
		finalAmount = usdToUse.Mul(usdFrom18To6)

	default:
		return nil, fmt.Errorf("unsupported denom selected: %s", chosen.denom)
	}

	// Step 5: Clamp to original coin range
	if chosen.minCoin != nil && finalAmount.LT(*chosen.minCoin) {
		finalAmount = *chosen.minCoin
	} else if chosen.maxCoin != nil && finalAmount.GT(*chosen.maxCoin) {
		finalAmount = *chosen.maxCoin
	}

	return sdk.NewCoins(sdk.NewCoin(chosen.denom, finalAmount)), nil
}
