package migrations

import (
	"fmt"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateDidIndyStyle: Starting migration")

	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		// Migrate from old Indy style DIDs to new Indy DIDs
		newDid := helpers.MigrateIndyStyleDid(didDocWithMetadata.DidDoc.Id)
		sctx.Logger().Debug(fmt.Sprintf("New version of DIDDoc: %s", newDid))
		didDocWithMetadata.ReplaceDids(didDocWithMetadata.DidDoc.Id, newDid)
	})
}
