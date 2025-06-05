package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		CreateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(50000000000),
				MaxAmount: sdkmath.NewInt(100000000000),
			},
			{
				Denom:     "usd",
				MinAmount: sdkmath.NewInt(1200000000000000000),
				MaxAmount: sdkmath.NewInt(2000000000000000000),
			},
		},
		UpdateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(25000000000),
				MaxAmount: sdkmath.NewInt(0),
			},
		},

		DeactivateDid: []FeeRange{
			{
				Denom:     BaseMinimalDenom,
				MinAmount: sdkmath.NewInt(10000000000),
				MaxAmount: sdkmath.NewInt(20000000000),
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

// ValidateBasic validates FeeParams structure
func (tfp *FeeParams) ValidateBasic() error {
	for i, f := range tfp.CreateDid {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_did[%d]: got %s", i, f.MinAmount.String())
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in create_did[%d]: got max=%s, min=%s", i, f.MaxAmount, f.MinAmount)
		}
	}
	for i, f := range tfp.UpdateDid {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in update_did[%d]: got %s", i, f.MinAmount.String())
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in update_did[%d]: got max=%s, min=%s", i, f.MaxAmount, f.MinAmount)
		}
	}
	for i, f := range tfp.DeactivateDid {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in deactivate_did[%d]: got %s", i, f.MinAmount.String())
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in deactivate_did[%d]: got max=%s, min=%s", i, f.MaxAmount, f.MinAmount)
		}
	}
	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}
	return nil
}

// ValidateBasic validates LegacyFeeParams structure
func (tfp *LegacyFeeParams) ValidateBasic() error {
	if !tfp.CreateDid.IsPositive() || tfp.CreateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create did tx fee: %s", tfp.CreateDid)
	}
	if !tfp.UpdateDid.IsPositive() || tfp.UpdateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid update did tx fee: %s", tfp.UpdateDid)
	}
	if !tfp.DeactivateDid.IsPositive() || tfp.DeactivateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid deactivate did tx fee: %s", tfp.DeactivateDid)
	}
	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}
	return nil
}

// Validators

func validateCreateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		for idx, f := range v {
			if f.MinAmount.IsNegative() {
				return fmt.Errorf("min_amount must be non-negative in create_did[%d]", idx)
			}
			if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
				return fmt.Errorf("max_amount must be >= min_amount in create_did[%d]", idx)
			}
		}
	case sdk.Coin:
		if v.IsNil() || !v.IsPositive() {
			return fmt.Errorf("create did fee must be a positive coin: %s", v)
		}
	default:
		return fmt.Errorf("invalid type for create_did: %T", i)
	}
	return nil
}

func validateUpdateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		for idx, f := range v {
			if f.MinAmount.IsNegative() {
				return fmt.Errorf("min_amount must be non-negative in update_did[%d]", idx)
			}
			if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
				return fmt.Errorf("max_amount must be >= min_amount in update_did[%d]", idx)
			}
		}
	case sdk.Coin:
		if v.IsNil() || !v.IsPositive() {
			return fmt.Errorf("update did fee must be a positive coin: %s", v)
		}
	default:
		return fmt.Errorf("invalid type for update_did: %T", i)
	}
	return nil
}

func validateDeactivateDid(i interface{}) error {
	switch v := i.(type) {
	case []FeeRange:
		for idx, f := range v {
			if f.MinAmount.IsNegative() {
				return fmt.Errorf("min_amount must be non-negative in deactivate_did[%d]", idx)
			}
			if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
				return fmt.Errorf("max_amount must be >= min_amount in deactivate_did[%d]", idx)
			}
		}
	case sdk.Coin:
		if v.IsNil() || !v.IsPositive() {
			return fmt.Errorf("deactivate did fee must be a positive coin: %s", v)
		}
	default:
		return fmt.Errorf("invalid type for deactivate_did: %T", i)
	}
	return nil
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

func validateFeeParams(i interface{}) error {
	switch v := i.(type) {
	case FeeParams:
		if err := validateCreateDid(v.CreateDid); err != nil {
			return err
		}
		if err := validateUpdateDid(v.UpdateDid); err != nil {
			return err
		}
		if err := validateDeactivateDid(v.DeactivateDid); err != nil {
			return err
		}
		if err := validateBurnFactor(v.BurnFactor); err != nil {
			return err
		}
		return v.ValidateBasic()

	case LegacyFeeParams:
		if err := validateCreateDid(v.CreateDid); err != nil {
			return fmt.Errorf("invalid create_did: %w", err)
		}
		if err := validateUpdateDid(v.UpdateDid); err != nil {
			return fmt.Errorf("invalid update_did: %w", err)
		}
		if err := validateDeactivateDid(v.DeactivateDid); err != nil {
			return fmt.Errorf("invalid deactivate_did: %w", err)
		}
		if err := validateBurnFactor(v.BurnFactor); err != nil {
			return fmt.Errorf("invalid burn_factor: %w", err)
		}
		return v.ValidateBasic()

	default:
		return fmt.Errorf("invalid parameter type: %T", i)
	}
}
