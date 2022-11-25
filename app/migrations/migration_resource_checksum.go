package migrations

import (
	"crypto/sha256"
	"errors"

	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksum(sctx sdk.Context, mctx MigrationContext) error {
	metadataStore := sctx.KVStore(mctx.resourceStoreKey)
	dataStore := sctx.KVStore(mctx.resourceStoreKey)
	metadataIterator := sdk.KVStorePrefixIterator(
		metadataStore,
		didutils.StrBytes(resourcetypes.ResourceMetadataKey))
	dataIterator := sdk.KVStorePrefixIterator(
		dataStore,
		didutils.StrBytes(resourcetypes.ResourceDataKey))

	defer closeIteratorOrPanic(metadataIterator)
	defer closeIteratorOrPanic(dataIterator)

	for metadataIterator.Valid() {
		if !dataIterator.Valid() {
			return errors.New("number of headers and data don't match")
		}

		var metadata resourcetypes.Metadata
		var data []byte

		// Get metadata and data from storage
		mctx.codec.MustUnmarshal(metadataIterator.Value(), &metadata)
		data = dataIterator.Value()

		// Fix checksum
		checksum := sha256.Sum256(data)
		metadata.Checksum = checksum[:]

		// Update HeaderInfo
		err := mctx.resourceKeeperNew.UpdateResourceMetadata(&sctx, &metadata)
		if err != nil {
			return err
		}

		// Iterate next
		metadataIterator.Next()
		dataIterator.Next()
	}
	return nil
}
