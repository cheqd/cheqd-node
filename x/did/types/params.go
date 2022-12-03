package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	ParamStoreKeyFeeParams = []byte("feeparams")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyFeeParams, FeeParams{}, validateFeeParams),
	)
}

// DefaultFeeParams returns default cheqd module tx fee parameters
func DefaultFeeParams() *FeeParams {
	return &FeeParams{
		CreateDid:     sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateDidTxFee)),
		UpdateDid:     sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultUpdateDidTxFee)),
		DeactivateDid: sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultDeactivateDidTxFee)),
		BurnFactor:    sdk.MustNewDecFromStr(DefaultBurnFactor),
	}
}

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	if !tfp.CreateDid.IsPositive() || tfp.CreateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create did tx fee: %s", tfp.CreateDid)
	}

	if !tfp.UpdateDid.IsPositive() || tfp.UpdateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid update did tx fee: %s", tfp.UpdateDid)
	}

	if !tfp.DeactivateDid.IsPositive() || tfp.DeactivateDid.Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid deactivate did tx fee: %s", tfp.DeactivateDid)
	}

	if !tfp.BurnFactor.IsPositive() || tfp.BurnFactor.GTE(sdk.OneDec()) {
		return fmt.Errorf("invalid burn factor: %s", tfp.BurnFactor)
	}

	return nil
}

func validateCreateDid(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("create did msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("create did msg fee param must be positive coin: %s", v)
	}

	return nil
}

func validateUpdateDid(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("update did msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("update did msg fee param must be positive coin: %s", v)
	}

	return nil
}

func validateDeactivateDid(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("deactivate did msg fee param must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("deactivate did msg fee param must be positive coin: %s", v)
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
