package migrations

import (
	"fmt"
	"sort"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceVersionLinks(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateResourceVersionLinks: Starting migration")
	store := sctx.KVStore(mctx.resourceStoreKey)

	sctx.Logger().Debug("MigrateResourceVersionLinks: Reading all resources")
	// Read all resources. Yes, this is memory intensive, but it's the simplest way to do it.
	// Resource size is limited to 200KB, so this should be fine.
	resources, err := mctx.resourceKeeperNew.GetAllResources(&sctx)
	if err != nil {
		return err
	}

	sctx.Logger().Debug("MigrateResourceVersionLinks: Reading all keys and Clean store")
	// Clean store
	keys := helpers.ReadAllKeys(store, nil)
	for _, key := range keys {
		store.Delete(key)
	}

	// Reset version links
	for _, resource := range resources {
		sctx.Logger().Debug(fmt.Sprintf(
			"MigrateResourceVersionLinks: Id: %s CollectionId: %s OldPreviousVersionId: %s OldNextVersionId: %s",
			resource.Metadata.Id,
			resource.Metadata.CollectionId,
			resource.Metadata.PreviousVersionId,
			resource.Metadata.NextVersionId))
		resource.Metadata.PreviousVersionId = ""
		resource.Metadata.NextVersionId = ""
	}

	sctx.Logger().Debug("MigrateResourceVersionLinks: Sorting resources by date created")
	// Sort resources by date created
	sort.Slice(resources, func(i, j int) bool {
		iCreated := resources[i].Metadata.Created
		jCreated := resources[j].Metadata.Created
		return iCreated.Before(jCreated)
	})

	sctx.Logger().Debug("MigrateResourceVersionLinks: Setting version links")
	// Add resources to store in the same order as they were created. This will create proper links.
	for _, resource := range resources {
		err = mctx.resourceKeeperNew.AddNewResourceVersion(&sctx, resource)
		sctx.Logger().Debug(fmt.Sprintf(
			"MigrateResourceVersionLinks: Id: %s CollectionId: %s NewPreviousVersionId: %s NewNextVersionId: %s",
			resource.Metadata.Id,
			resource.Metadata.CollectionId,
			resource.Metadata.PreviousVersionId,
			resource.Metadata.NextVersionId))
		if err != nil {
			return err
		}
	}
	sctx.Logger().Debug("MigrateResourceVersionLinks: Migration finished")

	return nil
}
