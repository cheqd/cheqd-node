package migrations

import (
	"crypto/sha256"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateDidV1(ctx sdk.Context, cheqdKeeper didkeeper.Keeper) error {
	// TODO: implement for cheqd module
	return nil
}

func MigrateResourceV1(ctx sdk.Context, cheqdKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// Resource Checksum migration
	err := MigrateResourceChecksumV1(ctx, cheqdKeeper, resourceKeeper)
	if err != nil {
		return err
	}
	// TODO: Add more migrations for resource module
	return nil
}

func MigrateResourceChecksumV1(ctx sdk.Context, cheqdKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	resources := resourceKeeper.GetAllResourcesV1(&ctx)
	for _, resource := range resources {
		checksum := sha256.Sum256([]byte(resource.Header.Checksum))
		resource.Header.Checksum = checksum[:]
		var migratedResource resourcetypesv2.ResourceWithMetadata
		migratedResource.Resource = &resourcetypesv2.Resource{Data: resource.Data}
		migratedResource.Metadata = &resourcetypesv2.Metadata{
			CollectionId:      resource.Header.CollectionId,
			Id:                resource.Header.Id,
			Name:              resource.Header.Name,
			ResourceType:      resource.Header.ResourceType,
			MediaType:         resource.Header.MediaType,
			Checksum:          resource.Header.Checksum,
			Created:           resource.Header.Created,
			PreviousVersionId: resource.Header.PreviousVersionId,
			NextVersionId:     resource.Header.NextVersionId,
		}
		err := resourceKeeper.SetResource(&ctx, &migratedResource)
		if err != nil {
			return err
		}
	}
	return nil
}
