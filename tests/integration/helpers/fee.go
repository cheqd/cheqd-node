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

func GetBurntPortion(tax sdk.Coin, burnFactor sdkmath.LegacyDec) sdk.Coin {
	taxDec := sdk.NewDecCoinFromCoin(tax)
	taxDec.Amount = taxDec.Amount.Mul(burnFactor)
	burnt, _ := taxDec.TruncateDecimal()
	return burnt
}

func GetRewardPortion(tax sdk.Coin, burnt sdk.Coin) sdk.Coin {
	return tax.Sub(burnt)
}
