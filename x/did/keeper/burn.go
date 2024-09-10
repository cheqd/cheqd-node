package keeper

import (
	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) burnFrom(ctx sdk.Context, amount sdk.Coins, burnFrom string) error {
	addr, err := sdk.AccAddressFromBech32(burnFrom)
	if err != nil {
		return err
	}

	err = k.bankkeeper.SendCoinsFromAccountToModule(ctx,
		addr,
		types.ModuleName,
		amount)
	if err != nil {
		return err
	}

	return k.bankkeeper.BurnCoins(ctx, types.ModuleName, amount)
}
