package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"

)

func BurnFee(bankKeeper BankKeeper, ctx sdk.Context, fee sdk.Coins) error {
	if fee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee to be burnt is zero")
	}

	return bankKeeper.BurnCoins(ctx, cheqdtypes.ModuleName, fee)
}
