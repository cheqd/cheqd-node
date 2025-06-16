package types

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	util "github.com/cheqd/cheqd-node/util"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		Image: []didtypes.FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(20000000000),
				MaxAmount: util.PtrInt(30000000000),
			},
		},
		Json: []didtypes.FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(3500000000),
				MaxAmount: util.PtrInt(60000000000),
			},
		},
		Default: []didtypes.FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(6000000000),
				MaxAmount: util.PtrInt(20000000000),
			},
		},
		BurnFactor: sdkmath.LegacyMustNewDecFromStr(DefaultBurnFactor),
	}
}

func DefaultLegacyFeeParams() *LegacyFeeParams {
	return &LegacyFeeParams{
		Image:      sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultCreateResourceImageFee)),
		Json:       sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultCreateResourceJSONFee)),
		Default:    sdk.NewCoin(BaseMinimalDenom, sdkmath.NewInt(DefaultCreateResourceDefaultFee)),
		BurnFactor: sdkmath.LegacyMustNewDecFromStr(DefaultBurnFactor),
	}
}

// USDRange represents a fee range converted to USD for comparison
type USDRange struct {
	MinUSD sdkmath.LegacyDec
	MaxUSD *sdkmath.LegacyDec // nil means no upper bound
}

// convertFeeRangeToUSD converts a FeeRange to USD using oracle price
func convertFeeRangeToUSD(ctx context.Context, oracleKeeper OracleKeeper, fr didtypes.FeeRange) (USDRange, error) {
	normalizeUSD := func(amount sdkmath.Int) sdkmath.LegacyDec {
		// Heuristic: If value >= 1e12, likely fixed-point with 18 decimals
		if amount.GTE(sdkmath.NewInt(1_000_000_000_000)) {
			return sdkmath.LegacyNewDecFromInt(amount).Quo(sdkmath.LegacyNewDec(1_000_000_000_000_000_000)) // 1e18
		}
		return sdkmath.LegacyNewDecFromInt(amount)
	}

	if fr.Denom == oracletypes.UsdDenom {
		minUSD := normalizeUSD(fr.MinAmount)

		var maxUSD *sdkmath.LegacyDec
		if fr.MaxAmount != nil {
			val := normalizeUSD(*fr.MaxAmount)
			maxUSD = &val
		}
		// If fr.MaxAmount is nil, maxUSD remains nil (no upper bound)

		return USDRange{
			MinUSD: minUSD,
			MaxUSD: maxUSD,
		}, nil
	}

	if fr.Denom == BaseMinimalDenom {
		// Get CHEQ/USD price from oracle
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		price, found := oracleKeeper.GetWMA(sdkCtx, oracletypes.CheqdSymbol, "BALANCED")
		if !found {
			return USDRange{}, fmt.Errorf("failed to get CHEQ/USD price")
		}

		// Validate price is positive
		if !price.IsPositive() {
			return USDRange{}, fmt.Errorf("invalid CHEQ/USD price: %s", price)
		}

		// Convert from base units (ncheq) to CHEQ, then to USD
		// 1 CHEQ = 1e9 ncheq
		cheqDecimals := sdkmath.LegacyNewDec(1_000_000_000) // 9 decimals

		minCheq := sdkmath.LegacyNewDecFromInt(fr.MinAmount).Quo(cheqDecimals)
		minUSD := minCheq.Mul(price)

		var maxUSD *sdkmath.LegacyDec
		if fr.MaxAmount != nil {
			maxCheq := sdkmath.LegacyNewDecFromInt(*fr.MaxAmount).Quo(cheqDecimals)
			val := maxCheq.Mul(price)
			maxUSD = &val
		}
		// If fr.MaxAmount is nil, maxUSD remains nil (no upper bound)

		return USDRange{
			MinUSD: minUSD,
			MaxUSD: maxUSD,
		}, nil
	}

	return USDRange{}, fmt.Errorf("unsupported denomination: %s", fr.Denom)
}

// validateFeeRangeOverlap validates that fee ranges for a message type have overlapping USD ranges
func validateFeeRangeOverlap(ctx context.Context, oracleKeeper OracleKeeper, msgType string, feeRanges []didtypes.FeeRange) error {
	if len(feeRanges) <= 1 {
		return nil // No overlap validation needed for single or no ranges
	}

	// Convert all ranges to USD
	usdRanges := make([]USDRange, len(feeRanges))
	for i, fr := range feeRanges {
		usdRange, err := convertFeeRangeToUSD(ctx, oracleKeeper, fr)
		if err != nil {
			return fmt.Errorf("failed to convert %s fee range %d to USD: %w", msgType, i, err)
		}
		usdRanges[i] = usdRange
	}

	// Check if all ranges have overlapping region
	// Start with first range as the overlap candidate
	overlapMin := usdRanges[0].MinUSD
	overlapMax := usdRanges[0].MaxUSD

	// For each subsequent range, find intersection
	for i := 1; i < len(usdRanges); i++ {
		// Update minimum to the higher of the two minimums
		if usdRanges[i].MinUSD.GT(overlapMin) {
			overlapMin = usdRanges[i].MinUSD
		}

		// Update maximum to the lower of the two maximums
		// Handle nil (unbounded) cases
		if usdRanges[i].MaxUSD != nil {
			if overlapMax == nil || usdRanges[i].MaxUSD.LT(*overlapMax) {
				val := *usdRanges[i].MaxUSD
				overlapMax = &val
			}
		}
		// If usdRanges[i].MaxUSD is nil but overlapMax is not, keep overlapMax
		// If both are nil, overlapMax remains nil (unbounded)
	}

	// Check if intersection is valid
	if overlapMax != nil && overlapMin.GT(*overlapMax) {
		return fmt.Errorf("no overlapping fee range found for %s: ranges do not have common USD value range", msgType)
	}

	return nil
}

// validateFeeRangeList is a generic validator for []didtypes.FeeRange
func validateFeeRangeList(name string, frs []didtypes.FeeRange) error {
	for i, f := range frs {
		if f.Denom != BaseMinimalDenom && f.Denom != oracletypes.UsdDenom {
			return fmt.Errorf("invalid denom in %s[%d]: got %s", name, i, f.Denom)
		}
		if f.MinAmount.IsNegative() || f.MinAmount.IsZero() {
			return fmt.Errorf("min_amount must be positive in %s[%d]: got %s", name, i, f.MinAmount.String())
		}
		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in %s[%d]: got max=%s, min=%s", name, i, f.MaxAmount, f.MinAmount)
		}
	}
	return nil
}

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	if err := validateFeeRangeList("image", tfp.Image); err != nil {
		return err
	}
	if err := validateFeeRangeList("json", tfp.Json); err != nil {
		return err
	}
	if err := validateFeeRangeList("default", tfp.Default); err != nil {
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
	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "create_resource_image", tfp.Image); err != nil {
		return err
	}

	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "create_resource_json", tfp.Json); err != nil {
		return err
	}

	if err := validateFeeRangeOverlap(ctx, oracleKeeper, "create_resource_default", tfp.Default); err != nil {
		return err
	}

	return nil
}

func validateImage(i interface{}) error {
	switch v := i.(type) {
	case []didtypes.FeeRange:
		return validateFeeRangeList("create_resource_image", v)
	case sdk.Coin:
		return validateCoin("create_resource_image", v)
	default:
		return fmt.Errorf("invalid type for create_resource_image: %T", i)
	}
}

func validateJSON(i interface{}) error {
	switch v := i.(type) {
	case []didtypes.FeeRange:
		return validateFeeRangeList("create_resource_json", v)
	case sdk.Coin:
		return validateCoin("create_resource_json", v)
	default:
		return fmt.Errorf("invalid type for create_resource_json: %T", i)
	}
}

func validateDefault(i interface{}) error {
	switch v := i.(type) {
	case []didtypes.FeeRange:
		return validateFeeRangeList("create_resource_default", v)
	case sdk.Coin:
		return validateCoin("create_resource_default", v)
	default:
		return fmt.Errorf("invalid type for create_resource_default: %T", i)
	}
}

// validateCoin is a helper to validate sdk.Coin
func validateCoin(name string, c sdk.Coin) error {
	if c.IsNil() || !c.IsPositive() {
		return fmt.Errorf("%s fee must be a positive coin: %s", name, c.String())
	}
	return nil
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

// ValidateBasic validates LegacyFeeParams structure using individual validators
func (tfp *LegacyFeeParams) ValidateBasic() error {
	if err := validateImage(tfp.Image); err != nil {
		return fmt.Errorf("invalid image fee: %w", err)
	}
	if err := validateJSON(tfp.Json); err != nil {
		return fmt.Errorf("invalid json fee: %w", err)
	}
	if err := validateDefault(tfp.Default); err != nil {
		return fmt.Errorf("invalid default fee: %w", err)
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
