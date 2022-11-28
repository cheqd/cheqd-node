package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateFees(amount string) []string {
	return []string{
		"--fees", amount,
	}
}

func GenerateFeeGranter(granter string, feeParams []string) []string {
	return append([]string{
		"--fee-granter", granter,
	}, feeParams...)
}

func GetBurntPortion(tax sdk.Coin, burnFactor sdk.Dec) sdk.Coin {
	taxDec := sdk.NewDecCoinFromCoin(tax)
	taxDec.Amount = taxDec.Amount.Mul(burnFactor)
	burnt, _ := taxDec.TruncateDecimal()
	return burnt
}

func GetRewardPortion(tax sdk.Coin, burnt sdk.Coin) sdk.Coin {
	return tax.Sub(burnt)
}
