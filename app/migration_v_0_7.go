package app

import (
	"crypto/sha256"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (app *App) Migration07(ctx sdk.Context) {
	app.FixResourceChecksums(ctx)
	app.FixResourceVersionLinks(ctx)
}

func (app *App) FixResourceChecksums(ctx sdk.Context) {
	// TODO: Loading everything into memory is not the best approach.
	// Resources can be large. I would suggest to use iterator instead and load resources one by one.
	all_resources := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range all_resources {
		checksum := sha256.Sum256([]byte(resource.Resource.Data))
		resource.Metadata.Checksum = hex.EncodeToString(checksum[:])
		err := app.resourceKeeper.SetResource(&ctx, &resource)
		if err != nil {
			return
		}
	}
}

func (app *App) FixResourceVersionLinks(ctx sdk.Context) {
	// TODO: We have to reset links first. Then we can use GetLastResourceVersionHeader
	// but only because resources in state are corted by creation time.
	// Also, we need to avoid loading all resources in memory.
	resourceList := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range resourceList {
		previousResourceVersionHeader, found := app.resourceKeeper.GetLastResourceVersionMetadata(&ctx, resource.Metadata.CollectionId, resource.Metadata.Name, resource.Metadata.ResourceType)
		if found {
			// Set links
			previousResourceVersionHeader.NextVersionId = resource.Metadata.Id
			resource.Metadata.PreviousVersionId = previousResourceVersionHeader.Id

			// Update previous version
			err := app.resourceKeeper.UpdateResourceMetadata(&ctx, &previousResourceVersionHeader)
			if err != nil {
				return
			}
		}
	}
}
