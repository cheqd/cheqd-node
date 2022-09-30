package app

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *App) Migration07(ctx sdk.Context) {
	all_resources := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range all_resources {
		h := sha256.New()
		h.Write(resource.Data)
		resource.Header.Checksum = h.Sum(nil)
		err := app.resourceKeeper.SetResource(&ctx, &resource)
		if err != nil {
			return
		}
	}
}
