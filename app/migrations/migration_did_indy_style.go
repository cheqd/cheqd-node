package migrations

import (
	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		// Migrate all dids, make them indy style
		newDid := helpers.MigrateIndyStyleDid(didDocWithMetadata.DidDoc.Id)
		didDocWithMetadata.ReplaceDids(didDocWithMetadata.DidDoc.Id, newDid)
	})
}
