package ante

import (
	"errors"
	"fmt"
	"math"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"
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

func GetTaxableMsgFeeWithBurnPortion(ctx sdk.Context, msg interface{}, ncheqPrice sdkmath.LegacyDec, userFee sdk.Coins, fak feeabskeeper.Keeper) (sdk.Coins, sdk.Coins, bool, error) {
	denom := userFee[0].Denom
	var nativeFees sdk.Coins
	if hostChainConfig, found := fak.GetHostZoneConfig(ctx, denom); found {
		var err error
		nativeFees, err = fak.CalculateNativeFromIBCCoins(ctx, userFee, hostChainConfig)
		if err != nil {
			return nil, nil, true, fmt.Errorf("failed to convert IBC fee to native denom: %w", err)
		}
	}
	switch msg := msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgCreateDidDoc], ncheqPrice, nativeFees)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *didtypes.MsgUpdateDidDoc:
		fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgUpdateDidDoc], ncheqPrice, nativeFees)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *didtypes.MsgDeactivateDidDoc:
		fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgDeactivateDidDoc], ncheqPrice, nativeFees)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	case *resourcetypes.MsgCreateResource:
		return GetResourceTaxableMsgFee(ctx, msg, ncheqPrice, userFee, nativeFees)
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

func GetResourceTaxableMsgFee(ctx sdk.Context, msg *resourcetypes.MsgCreateResource, ncheqPrice sdkmath.LegacyDec, userFee sdk.Coins, nativeFee sdk.Coins) (sdk.Coins, sdk.Coins, bool, error) {
	mediaType := resourceutils.DetectMediaType(msg.GetPayload().ToResource().Resource.Data)

	// Mime type image
	if strings.HasPrefix(mediaType, "image/") {
		fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgCreateResourceImage], ncheqPrice, nativeFee)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	}

	// Mime type json
	if strings.HasPrefix(mediaType, "application/json") {
		fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgCreateResourceJSON], ncheqPrice, nativeFee)
		if err != nil {
			return nil, nil, true, err
		}
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], fee)
		return GetRewardPortion(fee, burnPortion), burnPortion, true, nil
	}

	fee, err := GetFeeForMsg(userFee, TaxableMsgFees[MsgCreateResourceDefault], ncheqPrice, nativeFee)
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

func IsTaxableTx(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx, oracleKeeper OracleKeeper, feeabsKeeper feeabskeeper.Keeper) (bool, sdk.Coins, sdk.Coins, error) {
	ncheqPrice, exist := oracleKeeper.GetEMA(ctx, oracletypes.CheqdSymbol)
	if !exist {
		// fallback to fixed fee range in ncheq if defined
		ncheqPrice = sdkmath.LegacyZeroDec() // zero value, GetFeeForMsg will handle fallback
	}
	_ = checkFeeParamsFromSubspace(ctx, didKeeper, resourceKeeper)

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return false, sdk.Coins{}, sdk.Coins{}, errorsmod.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", tx)
	}
	reward := (sdk.Coins)(nil)
	burn := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {

		rewardPortion, burnPortion, isIdentityMsg, err := GetTaxableMsgFeeWithBurnPortion(ctx, msg, ncheqPrice, feeTx.GetFee(), feeabsKeeper)
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

func GetFeeForMsg(txFee sdk.Coins, feeRanges []didtypes.FeeRange, cheqEmaPrice sdkmath.LegacyDec, nativeFees sdk.Coins) (sdk.Coins, error) {
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

	// for fixed fees
	for _, fr := range feeRanges {
		if fr.MinAmount != nil && fr.MaxAmount != nil && fr.MinAmount.Equal(*fr.MaxAmount) {
			expectedFee := sdk.NewCoin(fr.Denom, *fr.MinAmount)

			if txFee.AmountOf(fr.Denom).IsPositive() {
				if !txFee.AmountOf(fr.Denom).GTE(expectedFee.Amount) {
					return nil, fmt.Errorf("invalid fixed fee: expected at least %s, got %s", expectedFee, txFee)
				}
				return sdk.NewCoins(expectedFee), nil
			}

			return validateCrossDenomFixedFee(txFee, expectedFee, cheqEmaPrice, nativeFees, fr.Denom)
		}
	}

	// Fallback: If CHEQ price is not available
	if cheqEmaPrice.IsZero() {
		return getFallbackFee(txFee, feeRanges)
	}

	// Setup
	const cheqExponent = 9
	const usdExponent = 6

	cheqScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(cheqExponent)))
	usdScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(usdExponent)))
	usdFrom18To6 := sdkmath.NewInt(1_000_000_000_000) // 1e12

	var ranges []usdRange

	// Convert all ranges to scaled USD values
	for _, fr := range feeRanges {
		if fr.MinAmount == nil && fr.MaxAmount == nil {
			continue
		}

		var minUSD, maxUSD *sdkmath.Int

		switch fr.Denom {
		case oracletypes.CheqdDenom:
			if fr.MinAmount != nil {
				usd := sdkmath.LegacyNewDecFromInt(*fr.MinAmount).QuoInt(cheqScale).Mul(cheqEmaPrice).MulInt(usdScale).TruncateInt()
				minUSD = &usd
			}
			if fr.MaxAmount != nil {
				usd := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).QuoInt(cheqScale).Mul(cheqEmaPrice).MulInt(usdScale).TruncateInt()
				maxUSD = &usd
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

	// Find overlapping range
	var overlapMin, overlapMax *sdkmath.Int
	for _, r := range ranges {
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
	if overlapMin != nil && overlapMax != nil && overlapMin.GT(*overlapMax) {
		return nil, errors.New("no valid overlapping USD range")
	}

	// Select denom to compute with (prefer USD)
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
	// Convert user fee into USD and validate against overlap
	var userUsdAmount *sdkmath.Int
	var handled bool

	for _, coin := range txFee {
		switch coin.Denom {
		case oracletypes.CheqdDenom:
			cheqDec := sdkmath.LegacyNewDecFromInt(coin.Amount).QuoInt(cheqScale)
			usd := cheqDec.Mul(cheqEmaPrice).MulInt(usdScale)
			usd1 := usd.TruncateInt()
			userUsdAmount = &usd1
			handled = true

		case oracletypes.UsdDenom:
			val := coin.Amount.Quo(usdFrom18To6)
			userUsdAmount = &val
			handled = true
		}
	}

	// IBC fallback: use converted nativeFees (always in ncheq)
	if !handled && len(nativeFees) > 0 {
		nativeCoin := nativeFees[0]
		if nativeCoin.Denom != oracletypes.CheqdDenom {
			return nil, fmt.Errorf("unexpected native denom: %s", nativeCoin.Denom)
		}

		cheqDec := sdkmath.LegacyNewDecFromInt(nativeCoin.Amount).QuoInt(cheqScale)
		usd := cheqDec.Mul(cheqEmaPrice).MulInt(usdScale)
		usd1 := usd.TruncateInt()
		userUsdAmount = &usd1
		handled = true
	}

	if !handled || userUsdAmount == nil {
		return nil, fmt.Errorf("user fee denom not supported and no nativeFees fallback")
	}

	// Reject if below minimum
	if overlapMin != nil && userUsdAmount.LT(*overlapMin) {
		return nil, fmt.Errorf("fee too low: expected ≥ %s USD, got %s USD", overlapMin, userUsdAmount)
	}

	// Use userUsdAmount directly, capped to overlapMax if needed
	effectiveUsd := *userUsdAmount
	if overlapMax != nil && effectiveUsd.GT(*overlapMax) {
		effectiveUsd = *overlapMin // ← Or cap to overlapMax if desired instead
	}

	// Convert effectiveUsd to target denom
	var finalAmount sdkmath.Int
	var finalDenom string

	// Use same denom as user for simplicity
	switch chosen.denom {
	case oracletypes.CheqdDenom:
		cheqAmt := sdkmath.LegacyNewDecFromInt(effectiveUsd).
			Quo(cheqEmaPrice).
			MulInt(cheqScale).
			QuoInt(usdScale).
			TruncateInt()
		finalAmount = cheqAmt
		finalDenom = oracletypes.CheqdDenom

	case oracletypes.UsdDenom:
		finalAmount = effectiveUsd.Mul(usdFrom18To6)
		finalDenom = oracletypes.UsdDenom

	default:
		return nil, fmt.Errorf("unsupported user fee denom: %s", txFee[0].Denom)
	}

	// Return final fee to be deducted (not necessarily full user amount)
	return sdk.NewCoins(sdk.NewCoin(finalDenom, finalAmount)), nil
}

func getFallbackFee(txFee sdk.Coins, feeRanges []didtypes.FeeRange) (sdk.Coins, error) {
	for _, fr := range feeRanges {
		if fr.Denom != oracletypes.CheqdDenom {
			continue
		}

		feeAmt := txFee.AmountOf(fr.Denom)

		if fr.MinAmount != nil {
			expected := sdk.NewCoin(fr.Denom, *fr.MinAmount)
			if !feeAmt.GTE(expected.Amount) {
				return nil, fmt.Errorf("cheq price unavailable; expected fee ≥ %s", expected)
			}
			return sdk.NewCoins(expected), nil
		}

		if fr.MaxAmount != nil {
			expected := sdk.NewCoin(fr.Denom, *fr.MaxAmount)
			if !feeAmt.GTE(expected.Amount) {
				return nil, fmt.Errorf("cheq price unavailable; expected fee ≥ %s", expected)
			}
			return sdk.NewCoins(expected), nil
		}
	}

	return nil, errors.New("cheq price not available and no ncheq fallback fee defined")
}

func validateCrossDenomFixedFee(
	txFee sdk.Coins,
	expectedFee sdk.Coin,
	cheqEmaPrice sdkmath.LegacyDec,
	nativeFees sdk.Coins,
	feeDenom string,
) (sdk.Coins, error) {
	const cheqExponent = 9
	const usdExponent = 6
	cheqScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(cheqExponent)))
	usdScale := sdkmath.NewIntFromUint64(uint64(math.Pow10(usdExponent)))

	for _, coin := range txFee {
		switch coin.Denom {
		case oracletypes.CheqdDenom:
			if cheqEmaPrice.IsZero() {
				return nil, errors.New("cannot verify cross-denom fixed fee: cheq price not available")
			}

			usdAmount := sdkmath.LegacyNewDecFromInt(expectedFee.Amount).QuoInt(sdkmath.NewInt(1_000_000_000_000)) // 1e18 → 1e6
			requiredCheq := usdAmount.Quo(cheqEmaPrice).MulInt(cheqScale).QuoInt(usdScale).TruncateInt()

			if coin.Amount.LT(requiredCheq) {
				return nil, fmt.Errorf("insufficient ncheq: need at least %s, got %s", requiredCheq, coin.Amount)
			}
			return sdk.NewCoins(sdk.NewCoin(coin.Denom, requiredCheq)), nil

		case oracletypes.UsdDenom:
			if cheqEmaPrice.IsZero() {
				return nil, errors.New("cannot verify cross-denom fixed fee: cheq price not available")
			}

			cheqAmountDec := sdkmath.LegacyNewDecFromInt(expectedFee.Amount).QuoInt(cheqScale)
			requiredUsd := cheqAmountDec.Mul(cheqEmaPrice).MulInt(usdScale).TruncateInt()

			userUsd := coin.Amount.Quo(sdkmath.NewInt(1_000_000_000_000)) // 1e18 → 1e6

			if userUsd.LT(requiredUsd) {
				return nil, fmt.Errorf("insufficient usd: need at least %s, got %s", requiredUsd, userUsd)
			}
			return sdk.NewCoins(sdk.NewCoin(coin.Denom, requiredUsd.Mul(sdkmath.NewInt(1_000_000_000_000)))), nil

		default:
			// Handle IBC-denom equivalents via nativeFee
			nativeCoin := nativeFees.AmountOf(oracletypes.CheqdDenom)
			if nativeCoin.IsZero() {
				return nil, fmt.Errorf("unsupported cross-denom fixed fee: user paid %s", coin.Denom)
			}
			if cheqEmaPrice.IsZero() {
				return nil, errors.New("cannot verify fixed fee for IBC denom: cheq price not available")
			}

			switch feeDenom {
			case oracletypes.UsdDenom:
				usdAmount := sdkmath.LegacyNewDecFromInt(expectedFee.Amount).QuoInt(sdkmath.NewInt(1_000_000_000_000)) // 1e18 → 1e6
				requiredCheq := usdAmount.Quo(cheqEmaPrice).MulInt(cheqScale).QuoInt(usdScale).TruncateInt()

				if nativeCoin.LT(requiredCheq) {
					return nil, fmt.Errorf("insufficient IBC-equivalent ncheq: required %s, got %s", requiredCheq, nativeCoin)
				}
				return sdk.NewCoins(sdk.NewCoin(oracletypes.CheqdDenom, requiredCheq)), nil

			case oracletypes.CheqdDenom:
				if nativeCoin.LT(expectedFee.Amount) {
					return nil, fmt.Errorf("insufficient IBC-equivalent ncheq: expected %s, got %s", expectedFee.Amount, nativeCoin)
				}
				return sdk.NewCoins(expectedFee), nil
			}

			return nil, fmt.Errorf("unsupported fixed fee denom %s for IBC conversion", feeDenom)
		}
	}

	return nil, fmt.Errorf("no valid cross-denom fee found")
}
