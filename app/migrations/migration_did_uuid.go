package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidUUID(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateDidUUID: Starting migration")

	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		sctx.Logger().Debug("MigrateDidUUID: OldDID: " + didDocWithMetadata.DidDoc.Id)
		// Migrate uuid dids, make them normalized
		newDid := helpers.MigrateUUIDDid(didDocWithMetadata.DidDoc.Id)
		didDocWithMetadata.ReplaceDids(didDocWithMetadata.DidDoc.Id, newDid)
		sctx.Logger().Debug("MigrateDidUUID: NewDID: " + didDocWithMetadata.DidDoc.Id)
	})
}
