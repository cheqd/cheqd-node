package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	util "github.com/cheqd/cheqd-node/util"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
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

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	for i, f := range tfp.Image {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_resource_image[%d]: got %s", i, f.MinAmount.String())
		}

		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in create_resource_image[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}
	}

	for i, f := range tfp.Json {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_resource_json[%d]: got %s", i, f.MinAmount.String())
		}

		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in create_resource_json[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}
	}

	for i, f := range tfp.Default {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in default_fee[%d]: got %s", i, f.MinAmount.String())
		}

		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be greater than or equal to min_amount in default_fee[%d]: got max=%s, min=%s", i, f.MaxAmount.String(), f.MinAmount.String())
		}
	}

	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}

	return nil
}

func validateImage(i interface{}) error {
	v, ok := i.([]didtypes.FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for create_resource_image: %T", i)
	}

	for idx, f := range v {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_resource_image[%d]", idx)
		}
		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in create_resource_image[%d]", idx)
		}
	}

	return nil
}

func validateJSON(i interface{}) error {
	v, ok := i.([]didtypes.FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for create_resource_json: %T", i)
	}

	for idx, f := range v {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in create_resource_json[%d]", idx)
		}
		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in create_resource_json[%d]", idx)
		}
	}

	return nil
}

func validateDefault(i interface{}) error {
	v, ok := i.([]didtypes.FeeRange)
	if !ok {
		return fmt.Errorf("invalid parameter type for default_fee: %T", i)
	}

	for idx, f := range v {
		if f.MinAmount.IsNegative() {
			return fmt.Errorf("min_amount must be non-negative in default_fee[%d]", idx)
		}
		if f.MaxAmount != nil && f.MaxAmount.LT(f.MinAmount) {
			return fmt.Errorf("max_amount must be >= min_amount in default_fee[%d]", idx)
		}
	}

	return nil
}

func validateBurnFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("burn factor must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("burn factor must not be negative: %s", v)
	}

	if v.GTE(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("burn factor must be less than 1: %s", v)
	}

	return nil
}

func validateFeeParams(i interface{}) error {
	v, ok := i.(FeeParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := validateImage(v.Image); err != nil {
		return err
	}

	if err := validateJSON(v.Json); err != nil {
		return err
	}

	if err := validateDefault(v.Default); err != nil {
		return err
	}

	if err := validateBurnFactor(v.BurnFactor); err != nil {
		return err
	}

	return v.ValidateBasic()
}
