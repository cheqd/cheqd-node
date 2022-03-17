package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *App) Migration05(ctx sdk.Context) {
	oldKey := "testnettestnet"
	namespase := app.cheqdKeeper.GetFromState(ctx, oldKey)
	app.cheqdKeeper.DeteteFromState(ctx, oldKey)
	app.cheqdKeeper.SetDidNamespace(ctx, namespase)
}
