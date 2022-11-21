package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default FeeParams map keys
const (
	DefaultKeyCreateDid     = "create_did"
	DefaultKeyUpdateDid     = "update_did"
	DefaultKeyDeactivateDid = "deactivate_did"
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
		TxTypes: map[string]sdk.Coin{
			DefaultKeyCreateDid:     sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultCreateDidTxFee)),
			DefaultKeyUpdateDid:     sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultUpdateDidTxFee)),
			DefaultKeyDeactivateDid: sdk.NewCoin(BaseMinimalDenom, sdk.NewInt(DefaultDeactivateDidTxFee)),
		},
		BurnFactor: sdk.MustNewDecFromStr(DefaultBurnFactor),
	}
}

// ValidateBasic performs basic validation of cheqd module tx fee parameters
func (tfp *FeeParams) ValidateBasic() error {
	if !tfp.TxTypes[DefaultKeyCreateDid].IsPositive() || tfp.TxTypes[DefaultKeyCreateDid].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid create did tx fee: %s", tfp.TxTypes[DefaultKeyCreateDid])
	}

	if !tfp.TxTypes[DefaultKeyUpdateDid].IsPositive() || tfp.TxTypes[DefaultKeyUpdateDid].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid update did tx fee: %s", tfp.TxTypes[DefaultKeyUpdateDid])
	}

	if !tfp.TxTypes[DefaultKeyDeactivateDid].IsPositive() || tfp.TxTypes[DefaultKeyDeactivateDid].Denom != BaseMinimalDenom {
		return fmt.Errorf("invalid deactivate did tx fee: %s", tfp.TxTypes[DefaultKeyDeactivateDid])
	}

	if !tfp.TxTypes[DefaultKeyCreateDid].IsGTE(tfp.TxTypes[DefaultKeyUpdateDid]) {
		return fmt.Errorf("create did tx fee must be greater than or equal to update did tx fee: %s >= %s", tfp.TxTypes[DefaultKeyCreateDid], tfp.TxTypes[DefaultKeyUpdateDid])
	}

	if tfp.TxTypes[DefaultKeyUpdateDid].IsLTE(tfp.TxTypes[DefaultKeyDeactivateDid]) {
		return fmt.Errorf("update did tx fee must be greater than deactivate did tx fee: %s > %s", tfp.TxTypes[DefaultKeyUpdateDid], tfp.TxTypes[DefaultKeyDeactivateDid])
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

	if v.TxTypes == nil {
		return fmt.Errorf("tx types must not be nil")
	}

	if err := validateCreateDid(v.TxTypes[DefaultKeyCreateDid]); err != nil {
		return err
	}

	if err := validateUpdateDid(v.TxTypes[DefaultKeyUpdateDid]); err != nil {
		return err
	}

	if err := validateDeactivateDid(v.TxTypes[DefaultKeyDeactivateDid]); err != nil {
		return err
	}

	if err := validateBurnFactor(v.BurnFactor); err != nil {
		return err
	}

	return v.ValidateBasic()
}
