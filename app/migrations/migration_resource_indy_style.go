package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		resource.Metadata.CollectionId = helpers.MigrateIndyStyleId(resource.Metadata.CollectionId)
	})
}
