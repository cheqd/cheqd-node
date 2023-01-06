package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var ParamStoreKeyFeeParams = []byte("feeparams")

// ParamKeyTable returns the key declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable(
		paramstypes.NewParamSetPair(ParamStoreKeyFeeParams, FeeParams{}, validateFeeParams),
	)
}

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		Image:      sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceImageFee)),
		Json:       sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceJSONFee)),
		Default:    sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceDefaultFee)),
		BurnFactor: sdk.MustNewDecFromStr(DefaultBurnFactor),
	}
}

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	if !tfp.Image.IsPositive() || tfp.Image.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource image tx fee: %s", tfp.Image)
	}

	if !tfp.Json.IsPositive() || tfp.Json.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource json tx fee: %s", tfp.Json)
	}

	if !tfp.Json.IsPositive() || tfp.Json.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource default tx fee: %s", tfp.Json)
	}

	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdk.OneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}

	return nil
}

func validateImage(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("create resource image msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("create resource image msg fee param must be positive coin: %s", v)
	}

	return nil
}

func validateJSON(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("create resource json msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("create resource json msg fee param must be positive coin: %s", v)
	}

	return nil
}

func validateDefault(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("create resource default msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("create resource default msg fee param must be positive coin: %s", v)
	}

	return nil
}

func validateBurnFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("burn factor must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("burn factor must not be negative: %s", v)
	}

	if v.GTE(sdk.OneDec()) {
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
