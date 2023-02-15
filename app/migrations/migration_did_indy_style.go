package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateDidIndyStyle: Starting migration")

	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		sctx.Logger().Debug("MigrateDidIndyStyle: OldDID: " + didDocWithMetadata.DidDoc.Id)
		// Migrate from old Indy style DIDs to new Indy DIDs
		newDid := helpers.MigrateIndyStyleDid(didDocWithMetadata.DidDoc.Id)
		didDocWithMetadata.ReplaceDids(didDocWithMetadata.DidDoc.Id, newDid)
		sctx.Logger().Debug("MigrateDidIndyStyle: NewDID: " + didDocWithMetadata.DidDoc.Id)
	})
}
