package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidVersionID(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateDidVersionID: Starting migration")

	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		sctx.Logger().Debug("MigrateDidVersionID: DID: " + didDocWithMetadata.DidDoc.Id + " OldVersionId: " + didDocWithMetadata.Metadata.VersionId)
		versionID := uuid.NewSHA1(uuid.Nil, []byte(didDocWithMetadata.DidDoc.Id))
		didDocWithMetadata.Metadata.VersionId = versionID.String()
		sctx.Logger().Debug("MigrateDidVersionID: DID: " + didDocWithMetadata.DidDoc.Id + " NewVersionId: " + didDocWithMetadata.Metadata.VersionId)
	})
}
