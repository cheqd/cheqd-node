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

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	for i, f := range tfp.CreateDid {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in create_did[%d]: got %s, expected %s", i, f.Denom, BaseMinimalDenom)
		}

		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_did[%d]: got %s", i, f.MinAmount.String())
		}

		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in create_did[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}

	}

	for i, f := range tfp.UpdateDid {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in update_did[%d]: got %s, expected %s", i, f.Denom, BaseMinimalDenom)
		}
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in update_did[%d]: got %s", i, f.MinAmount.String())
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in update_did[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}
	}

	for i, f := range tfp.DeactivateDid {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in deactivate_did[%d]: got %s, expected %s", i, f.Denom, BaseMinimalDenom)
		}
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in deactivate_did[%d]: got %s", i, f.MinAmount.String())
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in deactivate_did[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}
	}

	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}

	return nil
}
func validateCreateDid(i interface{}) error {
	v, ok := i.([]*FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for create_did: %T", i)
	}

	for idx, f := range v {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in create_did[%d]: got %s, expected %s", idx, f.Denom, BaseMinimalDenom)
		}
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_did[%d]", idx)
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in create_did[%d]", idx)
		}
	}
	return nil
}

func validateUpdateDid(i interface{}) error {
	v, ok := i.([]*FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for update_did: %T", i)
	}

	for idx, f := range v {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in update_did[%d]: got %s, expected %s", idx, f.Denom, BaseMinimalDenom)
		}
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in update_did[%d]", idx)
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in update_did[%d]", idx)
		}
	}
	return nil
}

func validateDeactivateDid(i interface{}) error {
	v, ok := i.([]*FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for deactivate_did: %T", i)
	}

	for idx, f := range v {
		if f.Denom != BaseMinimalDenom {
			return fmt.Errorf("invalid denom in deactivate_did[%d]: got %s, expected %s", idx, f.Denom, BaseMinimalDenom)
		}
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in deactivate_did[%d]", idx)
		}
		if !f.MaxAmount.IsZero() && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in deactivate_did[%d]", idx)
		}
	}
	return nil
}

func validateBurnFactor(i interface{}) error {
	v, ok := i.(sdk.DecCoin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNil() {
		return fmt.Errorf("burn factor must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("burn factor must not be negative: %s", v)
	}

	if v.Amount.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("burn factor must be less than 1: %s", v)
	}

	return nil
}

func validateFeeParams(i interface{}) error {
	v, ok := i.(FeeParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

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
}
