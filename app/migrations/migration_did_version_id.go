package migrations

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

// This migration should be run after protobuf that's why we use new DidDocWithMetadata
func MigrateDidVersionId(sctx sdk.Context, mctx MigrationContext) error {
	return MigrateDidSimple(sctx, mctx, func(didDocWithMetadata *didtypes.DidDocWithMetadata) {
		versionId := uuid.NewSHA1(uuid.Nil, []byte(didDocWithMetadata.DidDoc.Id))
		didDocWithMetadata.Metadata.VersionId = versionId.String()
	})
}
