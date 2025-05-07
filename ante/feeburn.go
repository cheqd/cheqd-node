package ante

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetBurnFeePortion(burnFactor sdkmath.LegacyDec, fee sdk.Coins) sdk.Coins {
	feeDecCoins := sdk.NewDecCoinsFromCoins(fee...)

	burnFeePortion, _ := feeDecCoins.MulDec(burnFactor).TruncateDecimal()

	return burnFeePortion
}
