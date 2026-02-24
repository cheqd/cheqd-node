package helpers

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateFees(amount string) []string {
	return []string{
		"--fees", amount,
		"--gas", "auto",
		"--gas-adjustment", "3.5",
	}
}

func GenerateFeeGranter(granter string, feeParams []string) []string {
	return append([]string{
		"--fee-granter", granter,
	}, feeParams...)
}

func GetBurnFeePortion(burnFactor sdkmath.LegacyDec, fee sdk.Coins) sdk.Coins {
	feeDecCoins := sdk.NewDecCoinsFromCoins(fee...)

	burnFeePortion, _ := feeDecCoins.MulDec(burnFactor).TruncateDecimal()

	return burnFeePortion
}

func GetRewardPortion(total sdk.Coins, burnPortion sdk.Coins) sdk.Coins {
	if burnPortion.IsZero() {
		return total
	}
	return total.Sub(burnPortion...)
}
