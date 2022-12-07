package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetBurnFeePortion(burnFactor sdk.Dec, fee sdk.Coins) sdk.Coins {
	feeDecCoins := sdk.NewDecCoinsFromCoins(fee...)

	burnFeePortion, _ := feeDecCoins.MulDec(burnFactor).TruncateDecimal()

	return burnFeePortion
}
