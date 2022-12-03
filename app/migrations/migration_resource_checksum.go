package migrations

import (
	"crypto/sha256"
	"encoding/hex"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksum(sctx sdk.Context, mctx MigrationContext) error {
	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		checksum := sha256.Sum256(resource.Resource.Data)
		resource.Metadata.Checksum = hex.EncodeToString(checksum[:])
	})
}
