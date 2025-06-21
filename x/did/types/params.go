package types

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/util"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		CreateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: util.PtrInt(50000000000),
				MaxAmount: util.PtrInt(100000000000),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(1200000000000000000),
				MaxAmount: util.PtrInt(2000000000000000000),
			},
		},
		UpdateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: util.PtrInt(25000000000),
				MaxAmount: nil,
			},
		},

		DeactivateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: util.PtrInt(10000000000),
				MaxAmount: util.PtrInt(20000000000),
			},
		},
		BurnFactor: sdkmath.LegacyMustNewDecFromStr(DefaultBurnFactor),
	}
}

// DefaultLegacyFeeParams returns default fee params using sdk.Coins
func DefaultLegacyFeeParams() *LegacyFeeParams {
	return &LegacyFeeParams{
		CreateDid:     sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultCreateDidTxFee)),
		UpdateDid:     sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultUpdateDidTxFee)),
		DeactivateDid: sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultDeactivateDidTxFee)),
		BurnFactor:    sdkmath.LegacyMustNewDecFromStr(DefaultBurnFactor),
	}
}

// USDRange represents a fee range converted to USD for comparison
type USDRange struct {
	MinUSD *sdkmath.Int
	MaxUSD *sdkmath.Int
}

// convertFeeRangeToUSD converts a FeeRange to USD using oracle price
func convertFeeRangeToUSD(ctx context.Context, oracleKeeper OracleKeeper, fr FeeRange) (USDRange, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var minUSD *sdkmath.Int
	var maxUSD *sdkmath.Int

	switch fr.Denom {
	case oracletypes.UsdDenom:
		// Assume input is in 18 decimals (e.g., 1.2 USD = 1_200_000_000_000_000_000)
		// Convert to 6 decimals (1.2 USD = 1_200_000)
		normalizeUSD := func(amount sdkmath.Int) sdkmath.Int {
			return amount.Quo(sdkmath.NewInt(1_000_000_000_000)) // 1e12
		}

		if fr.MinAmount != nil {
			val := normalizeUSD(*fr.MinAmount)
			minUSD = &val
		}
		if fr.MaxAmount != nil {
			val := normalizeUSD(*fr.MaxAmount)
			maxUSD = &val
		}

	case BaseMinimalDenom:
		price, found := oracleKeeper.GetWMA(sdkCtx, oracletypes.CheqdSymbol, "BALANCED")
		if !found || !price.IsPositive() {
			return USDRange{}, fmt.Errorf("invalid or missing CHEQ/USD price")
		}

		cheqDecimals := sdkmath.NewInt(1_000_000_000) // 9 decimals for ncheq
		usdScaleFactor := sdkmath.NewInt(1_000_000)   // 6 decimals for USD

		if fr.MinAmount != nil {
			cheqDec := sdkmath.LegacyNewDecFromInt(*fr.MinAmount).QuoInt(cheqDecimals)
			usdValue := cheqDec.Mul(price).MulInt(usdScaleFactor)
			val := usdValue.TruncateInt()
			minUSD = &val
		}
		if fr.MaxAmount != nil {
			cheqDec := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).QuoInt(cheqDecimals)
			usdValue := cheqDec.Mul(price).MulInt(usdScaleFactor)
			val := usdValue.TruncateInt()
			maxUSD = &val
		}

	default:
		return USDRange{}, fmt.Errorf("unsupported denomination: %s", fr.Denom)
	}

	if minUSD == nil && maxUSD == nil {
		return USDRange{}, fmt.Errorf("both MinAmount and MaxAmount cannot be nil for denom: %s", fr.Denom)
	}

	return USDRange{
		MinUSD: minUSD,
		MaxUSD: maxUSD,
	}, nil
}

// validateFeeRangeOverlap ensures all fee ranges overlap in USD for a given message type
func validateFeeRangeOverlap(ctx context.Context, oracleKeeper OracleKeeper, msgType string, feeRanges []FeeRange) error {
	if len(feeRanges) <= 1 {
		return nil // No overlap check needed
	}
	// Convert all to USD
	usdRanges := make([]USDRange, len(feeRanges))
	for i, fr := range feeRanges {
		usdRange, err := convertFeeRangeToUSD(ctx, oracleKeeper, fr)
		if err != nil {
			return fmt.Errorf("failed to convert %s fee range %d to USD: %w", msgType, i, err)
		}
		usdRanges[i] = usdRange
	}

	// Initialize overlap range from first range
	overlapMin := usdRanges[0].MinUSD
	overlapMax := usdRanges[0].MaxUSD

	for i := 1; i < len(usdRanges); i++ {
		r := usdRanges[i]

		// Update overlapMin: max(current, r.MinUSD)
		if r.MinUSD != nil {
			if overlapMin == nil || r.MinUSD.GT(*overlapMin) {
				val := *r.MinUSD
				overlapMin = &val
			}
		}

		// Update overlapMax: min(current, r.MaxUSD), accounting for nils
		if r.MaxUSD != nil {
			if overlapMax == nil || r.MaxUSD.LT(*overlapMax) {
				val := *r.MaxUSD
				overlapMax = &val
			}
		}
	}

	// Final overlap check
	if overlapMax != nil && overlapMin != nil && overlapMin.GT(*overlapMax) {
		return fmt.Errorf("no overlapping fee range found for %s: USD ranges do not intersect", msgType)
	}

	return nil
}

// validateFeeRangeList is a generic validator for []FeeRange
func validateFeeRangeList(name string, frs []FeeRange) error {
	for i, f := range frs {
		if f.Denom != BaseMinimalDenom && f.Denom != oracletypes.UsdDenom {
			return fmt.Errorf("invalid denom in %s[%d]: got %s", name, i, f.Denom)
		}

		if f.MinAmount == nil && f.MaxAmount == nil {
			return fmt.Errorf("at least one of min_amount or max_amount must be set in %s[%d]", name, i)
		}

		if f.MinAmount != nil {
			if f.MinAmount.IsNegative() || f.MinAmount.IsZero() {
				return fmt.Errorf("min_amount must be positive if set in %s[%d]: got %s", name, i, f.MinAmount.String())
			}
		}

		if f.MaxAmount != nil {
			if f.MaxAmount.IsNegative() || f.MaxAmount.IsZero() {
				return fmt.Errorf("max_amount must be positive if set in %s[%d]: got %s", name, i, f.MaxAmount.String())
			}
		}

		if f.MinAmount != nil && f.MaxAmount != nil && f.MaxAmount.LT(*f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in %s[%d]: got max=%s, min=%s", name, i, f.MaxAmount, f.MinAmount)
		}
	}
	return nil
}

// validateCoin is a helper to validate sdk.Coin
func validateCoin(name string, c sdk.Coin) error {
	if c.IsNil() || !c.IsPositive() {
		return fmt.Errorf("%s fee must be a positive coin: %s", name, c.String())
	}
	return nil
}

func validateCreateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		return validateFeeRangeList("create_did", v)
	case sdk.Coin:
		return validateCoin("create_did", v)
	default:
		return fmt.Errorf("invalid type for create_did: %T", i)
	}
}

func validateUpdateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		return validateFeeRangeList("update_did", v)
	case sdk.Coin:
		return validateCoin("update_did", v)
	default:
		return fmt.Errorf("invalid type for update_did: %T", i)
	}
}

func validateDeactivateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		return validateFeeRangeList("deactivate_did", v)
	case sdk.Coin:
		return validateCoin("deactivate_did", v)
	default:
		return fmt.Errorf("invalid type for deactivate_did: %T", i)
	}
}

func validateBurnFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid type for burn_factor: %T", i)
	}
	if !v.IsPositive() || v.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("burn factor must be positive and < 1: %s", v)
	}
	return nil
}

// ValidateBasic validates FeeParams structure using individual validators
func (tfp *FeeParams) ValidateBasic() error {
	if err := validateCreateDid(tfp.CreateDid); err != nil {
		return err
	}
	if err := validateUpdateDid(tfp.UpdateDid); err != nil {
		return err
	}
	if err := validateDeactivateDid(tfp.DeactivateDid); err != nil {
		return err
	}
	if err := validateBurnFactor(tfp.BurnFactor); err != nil {
		return err
	}
	return nil
}

// ValidateWithOracle validates FeeParams with oracle price overlap checking
func (tfp *FeeParams) ValidateWithOracle(ctx context.Context, oracleKeeper OracleKeeper) error {
	// First do basic validation
	if err := tfp.ValidateBasic(); err != nil {
		return err
	}
	// Then validate overlaps for each message type
	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "create_did", tfp.CreateDid); err != nil {
		return err
	}

	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "update_did", tfp.UpdateDid); err != nil {
		return err
	}

	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "deactivate_did", tfp.DeactivateDid); err != nil {
		return err
	}

	return nil
}

// ValidateBasic validates LegacyFeeParams structure using individual validators
func (tfp *LegacyFeeParams) ValidateBasic() error {
	if err := validateCreateDid(tfp.CreateDid); err != nil {
		return fmt.Errorf("invalid create did tx fee: %w", err)
	}
	if err := validateUpdateDid(tfp.UpdateDid); err != nil {
		return fmt.Errorf("invalid update did tx fee: %w", err)
	}
	if err := validateDeactivateDid(tfp.DeactivateDid); err != nil {
		return fmt.Errorf("invalid deactivate did tx fee: %w", err)
	}
	if err := validateBurnFactor(tfp.BurnFactor); err != nil {
		return fmt.Errorf("invalid burn factor: %w", err)
	}
	return nil
}

// validateFeeParams validates either FeeParams or LegacyFeeParams
func validateFeeParams(i interface{}) error {
	switch v := i.(type) {
	case FeeParams:
		return v.ValidateBasic()
	case *FeeParams:
		return v.ValidateBasic()
	case LegacyFeeParams:
		return v.ValidateBasic()
	case *LegacyFeeParams:
		return v.ValidateBasic()
	default:
		return fmt.Errorf("invalid parameter type: %T", i)
	}
}
