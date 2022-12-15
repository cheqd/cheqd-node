package migrations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksum(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateResourceChecksum: Starting migration")

	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		sctx.Logger().Debug(fmt.Sprintf(
			"MigrateResourceChecksum: Id: %s CollectionId: %s OldChecksum: %s",
			resource.Metadata.Id,
			resource.Metadata.CollectionId,
			resource.Metadata.Checksum))
		checksum := sha256.Sum256(resource.Resource.Data)
		resource.Metadata.Checksum = hex.EncodeToString(checksum[:])
		sctx.Logger().Debug(fmt.Sprintf(
			"MigrateResourceChecksum: Id: %s CollectionId: %s NewChecksum: %s",
			resource.Metadata.Id,
			resource.Metadata.CollectionId,
			resource.Metadata.Checksum))
	})
}
