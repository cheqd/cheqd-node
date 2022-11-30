package migrations

import (
	"sort"
	"time"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceVersionLinks(sctx sdk.Context, mctx MigrationContext) error {
	store := sctx.KVStore(mctx.resourceStoreKey)

	// Read all resources. Yes, this is memory intensive, but it's the simplest way to do it.
	// Resource size is limited to 200KB, so this should be fine.
	resources, err := mctx.resourceKeeperNew.GetAllResources(&sctx)
	if err != nil {
		return err
	}

	// Clean store
	keys := helpers.ReadAllKeys(store, nil)
	for _, key := range keys {
		store.Delete(key)
	}

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

	// Add resources to store in the same order as they were created. This will create proper links.
	for _, resource := range resources {
		err = mctx.resourceKeeperNew.AddNewResourceVersion(&sctx, &resource)
		if err != nil {
			return err
		}
	}

	return nil
}
