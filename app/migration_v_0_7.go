package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Here we need to make all the links the 
func (app *App) Migration07(ctx sdk.Context) {
	resourceList := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range resourceList {
		previousResourceVersionHeader, found := app.resourceKeeper.GetLastResourceVersionHeader(&ctx, resource.Header.CollectionId, resource.Header.Name, resource.Header.ResourceType)
		if found {
			// Set links
			previousResourceVersionHeader.NextVersionId = resource.Header.Id
			resource.Header.PreviousVersionId = previousResourceVersionHeader.Id

			// Update previous version
			app.resourceKeeper.UpdateResourceHeader(&ctx, &previousResourceVersionHeader)
		}
	}
}
