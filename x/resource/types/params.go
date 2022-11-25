package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default FeeParams map keys
const (
	DefaultKeyCreateResourceImage = "image"
	DefaultKeyCreateResourceJson  = "json"
	DefaultKeyCreateResource      = "default"
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
		MediaTypes: map[string]sdk.Coin{
			DefaultKeyCreateResourceImage: sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceImageFee)),
			DefaultKeyCreateResourceJson:  sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceJsonFee)),
			DefaultKeyCreateResource:      sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateResourceDefaultFee)),
		},
		BurnFactor: sdk.MustNewDecFromStr(DefaultBurnFactor),
	}
}

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	if !tfp.MediaTypes[DefaultKeyCreateResourceImage].IsPositive() || tfp.MediaTypes[DefaultKeyCreateResourceImage].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource image tx fee: %s", tfp.MediaTypes[DefaultKeyCreateResourceImage])
	}

	if !tfp.MediaTypes[DefaultKeyCreateResourceJson].IsPositive() || tfp.MediaTypes[DefaultKeyCreateResourceJson].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource json tx fee: %s", tfp.MediaTypes[DefaultKeyCreateResourceJson])
	}

	if !tfp.MediaTypes[DefaultKeyCreateResourceJson].IsPositive() || tfp.MediaTypes[DefaultKeyCreateResourceJson].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create resource default tx fee: %s", tfp.MediaTypes[DefaultKeyCreateResourceJson])
	}

	if !tfp.MediaTypes[DefaultKeyCreateResourceImage].IsGTE(tfp.MediaTypes[DefaultKeyCreateResourceJson]) {
		return fmt.Errorf("create resource image tx fee must be greater than or equal to create resource json tx fee: %s >= %s", tfp.MediaTypes[DefaultKeyCreateResourceImage], tfp.MediaTypes[DefaultKeyCreateResourceJson])
	}

	if tfp.MediaTypes[DefaultKeyCreateResourceJson].IsLTE(tfp.MediaTypes[DefaultKeyCreateResource]) {
		return fmt.Errorf("create resource json tx fee must be greater than create resource default tx fee: %s > %s", tfp.MediaTypes[DefaultKeyCreateResourceJson], tfp.MediaTypes[DefaultKeyCreateResource])
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

func validateJson(i interface{}) error {
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

	if err := validateImage(v.MediaTypes[DefaultKeyCreateResourceImage]); err != nil {
		return err
	}

	if err := validateJson(v.MediaTypes[DefaultKeyCreateResourceJson]); err != nil {
		return err
	}

	if err := validateDefault(v.MediaTypes[DefaultKeyCreateResource]); err != nil {
		return err
	}

	if err := validateBurnFactor(v.BurnFactor); err != nil {
		return err
	}

	return v.ValidateBasic()
}
