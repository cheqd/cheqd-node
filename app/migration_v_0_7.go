package app

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *App) Migration07(ctx sdk.Context) {
	// TODO: Loading everything into memory is not the best approach.
	// Resources can be large. I would suggest to use iterator instead and load resources one by one.
	all_resources := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range all_resources {
		checksum := sha256.Sum256([]byte(resource.Data))
		resource.Header.Checksum = checksum[:]
		err := app.resourceKeeper.SetResource(&ctx, &resource)
		if err != nil {
			return
		}
	}
}
