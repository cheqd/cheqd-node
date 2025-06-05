package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		CreateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(50000000000),
				MaxAmount: util.PtrInt(100000000000),
			},
			{
				Denom:     "usd",
				MinAmount: sdkmath.NewInt(1200000000000000000),
				MaxAmount: util.PtrInt(2000000000000000000),
			},
		},
		UpdateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(25000000000),
				MaxAmount: nil,
			},
		},

		DeactivateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(10000000000),
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

// validateFeeRangeList is a generic validator for []FeeRange
func validateFeeRangeList(name string, frs []FeeRange) error {
	for i, f := range frs {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in %s[%d]: got %s", name, i, f.MinAmount.String())
		}
		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
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
	if v.IsNil() || v.IsNegative() || v.GTE(sdkmath.LegacyOneDec()) {
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
