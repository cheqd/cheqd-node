package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceUUID(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateResourceUUID: Starting migration")

	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		sctx.Logger().Debug("MigrateResourceUUID: OldId: " + resource.Metadata.Id)
		sctx.Logger().Debug("MigrateResourceUUID: OldCollectionId: " + resource.Metadata.CollectionId)
		resource.Metadata.CollectionId = helpers.MigrateUUIDId(resource.Metadata.CollectionId)
		resource.Metadata.Id = helpers.MigrateUUIDId(resource.Metadata.Id)
		sctx.Logger().Debug("MigrateResourceUUID: NewId: " + resource.Metadata.Id)
		sctx.Logger().Debug("MigrateResourceUUID: NewCollectionId: " + resource.Metadata.CollectionId)
	})
}
