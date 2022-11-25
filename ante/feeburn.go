package ante

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GenericFeeParams struct {
	BurnFactor sdk.Dec
}

type ModuleKeeper interface {
	GetParams(ctx sdk.Context) didtypes.FeeParams
}

func BurnFee(bankKeeper BankKeeper, ctx sdk.Context, fee sdk.Coins) error {
	if fee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee to be burnt is zero")
	}

	return bankKeeper.BurnCoins(ctx, didtypes.ModuleName, fee)
}

func GetBurnFeePortion(burnFactor sdk.Dec, fee sdk.Coins) sdk.Coins {
	feeDecCoins := sdk.NewDecCoinsFromCoins(fee...)

	burnFeePortion, _ := feeDecCoins.MulDec(burnFactor).TruncateDecimal()

	return burnFeePortion
}
