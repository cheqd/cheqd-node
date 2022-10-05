package ante

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GenericFeeParams struct {
	BurnFactor sdk.Dec
}

type ModuleKeeper interface {
	GetParams(ctx sdk.Context) cheqdtypes.FeeParams
}

func BurnFee(bankKeeper BankKeeper, ctx sdk.Context, fee sdk.Coins) error {
	if fee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee to be burnt is zero")
	}

	return bankKeeper.BurnCoins(ctx, cheqdtypes.ModuleName, fee)
}

func GetBurnFeePortion(ctx sdk.Context, burnFactor sdk.Dec, fee sdk.Coins) sdk.Coins {
	feeDecCoins := sdk.NewDecCoinsFromCoins(fee...)

	burnFeePortion, _ := feeDecCoins.MulDec(burnFactor).TruncateDecimal();

	return burnFeePortion
}