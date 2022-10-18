package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cheqdapp "github.com/cheqd/cheqd-node/app"
)


func (app *cheqdapp.App) Migration05(ctx sdk.Context) error {
	oldKey := "testnettestnet"
	namespase := app.cheqdKeeper.GetFromState(&ctx, oldKey)
	app.cheqdKeeper.DeleteFromState(&ctx, oldKey)
	app.cheqdKeeper.SetDidNamespace(&ctx, namespase)
	return nil
}
