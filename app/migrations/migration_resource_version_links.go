package migrations

import (
	"sort"
	"time"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceVersionLinks(sctx sdk.Context, mctx MigrationContext) error {
	println("Resource version links migration. Start")
	store := sctx.KVStore(mctx.resourceStoreKey)

	println("Resource version links migration. Read all resources")
	// Read all resources. Yes, this is memory intensive, but it's the simplest way to do it.
	// Resource size is limited to 200KB, so this should be fine.
	resources, err := mctx.resourceKeeperNew.GetAllResources(&sctx)
	if err != nil {
		return err
	}

	println("Resource version links migration. Read all keys and Clean store")
	// Clean store
	keys := helpers.ReadAllKeys(store, nil)
	for _, key := range keys {
		store.Delete(key)
	}

	// Reset version links
	for _, resource := range resources {
		resource.Metadata.PreviousVersionId = ""
		resource.Metadata.NextVersionId = ""
	}

	println("Resource version links migration. Sort resources by date created")
	// Sort resources by date created
	sort.Slice(resources[:], func(i, j int) bool {
		iCreated, err := time.Parse(time.RFC3339, resources[i].Metadata.Created)
		if err != nil {
			panic(err)
		}

		jCreated, err := time.Parse(time.RFC3339, resources[j].Metadata.Created)
		if err != nil {
			panic(err)
		}

		return iCreated.Before(jCreated)
	})

	println("Resource version links migration. Set version links")
	// Add resources to store in the same order as they were created. This will create proper links.
	for _, resource := range resources {
		err = mctx.resourceKeeperNew.AddNewResourceVersion(&sctx, &resource)
		if err != nil {
			return err
		}
	}
	println("Resource version links migration. End")

	return nil
}
