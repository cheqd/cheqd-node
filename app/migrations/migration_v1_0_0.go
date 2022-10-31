package migrations

import (
	"crypto/sha256"

	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateCheqdV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper) error {
	// TODO: implement for cheqd module
	return nil
}

func MigrateResourceV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// Resource Checksum migration
	err := MigrateResourceChecksumV1(ctx, cheqdKeeper, resourceKeeper)
	if err != nil {
		return err
	}
	// TODO: Add more migrations for resource module
	return nil
}

func MigrateResourceChecksumV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	resources := resourceKeeper.GetAllResources(&ctx)
	for _, resource := range resources {
		checksum := sha256.Sum256([]byte(resource.Data))
		resource.Header.Checksum = checksum[:]
		err := resourceKeeper.SetResource(&ctx, &resource)
		if err != nil {
			return err
		}
	}
	return nil
}
