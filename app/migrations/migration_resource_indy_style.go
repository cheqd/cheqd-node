package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateResourceIndyStyle: Starting migration")

	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		sctx.Logger().Debug("MigrateResourceIndyStyle: OldCollectionId: " + resource.Metadata.CollectionId)
		resource.Metadata.CollectionId = helpers.MigrateIndyStyleId(resource.Metadata.CollectionId)
		sctx.Logger().Debug("MigrateResourceIndyStyle: NewCollectionId: " + resource.Metadata.CollectionId)
	})
}
