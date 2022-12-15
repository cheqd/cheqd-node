package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidVersionId(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateDidVersionId: Starting migration")

	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		sctx.Logger().Debug("MigrateDidVersionId: DID: " + didDocWithMetadata.DidDoc.Id + " OldVersionId: " + didDocWithMetadata.Metadata.VersionId)
		versionId := uuid.NewSHA1(uuid.Nil, []byte(didDocWithMetadata.DidDoc.Id))
		didDocWithMetadata.Metadata.VersionId = versionId.String()
		sctx.Logger().Debug("MigrateDidVersionId: DID: " + didDocWithMetadata.DidDoc.Id + " NewVersionId: " + didDocWithMetadata.Metadata.VersionId)
	})
}
