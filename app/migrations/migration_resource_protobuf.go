package migrations

import (
	"encoding/hex"
	"strings"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceProtobuf(sctx sdk.Context, mctx MigrationContext) error {
	sctx.Logger().Debug("MigrateResourceProtobuf: Starting migration")
	// Storage for old headers and data
	store := sctx.KVStore(mctx.resourceStoreKey)

	sctx.Logger().Debug("MigrateResourceProtobuf: Resetting counter")
	// Reset counter
	mctx.resourceKeeperNew.SetResourceCount(&sctx, 0)

	sctx.Logger().Debug("MigrateResourceProtobuf: Reading all keys")
	headerKeys := helpers.ReadAllKeys(store, didutils.StrBytes(resourcetypesv1.ResourceHeaderKey))

	for _, headerKey := range headerKeys {
		sctx.Logger().Debug("MigrateResourceProtobuf: Starting migration for resource with header key: " + string(headerKey))
		dataKey := ResourceV1HeaderkeyToDataKey(headerKey)

		var oldHeader resourcetypesv1.ResourceHeader
		mctx.codec.MustUnmarshal(store.Get(headerKey), &oldHeader)
		oldData := store.Get(dataKey)

		sctx.Logger().Debug("MigrateResourceProtobuf: Collecting new resource metadata")
		newMetadata := resourcetypes.Metadata{
			CollectionId:      oldHeader.CollectionId,
			Id:                oldHeader.Id,
			Name:              oldHeader.Name,
			Version:           "",
			ResourceType:      oldHeader.ResourceType,
			AlsoKnownAs:       []*resourcetypes.AlternativeUri{},
			MediaType:         oldHeader.MediaType,
			Created:           oldHeader.Created,
			Checksum:          hex.EncodeToString(oldHeader.Checksum),
			PreviousVersionId: oldHeader.PreviousVersionId,
			NextVersionId:     oldHeader.NextVersionId,
		}

		resourceWithMetadata := resourcetypes.ResourceWithMetadata{
			Metadata: &newMetadata,
			Resource: &resourcetypes.Resource{
				Data: oldData,
			},
		}

		sctx.Logger().Debug("MigrateResourceProtobuf: Remove old values")
		// Remove old resource data and header
		store.Delete(headerKey)
		store.Delete(dataKey)

		sctx.Logger().Debug("MigrateResourceProtobuf: Write new resource with metadata")
		// Write new resource
		err := mctx.resourceKeeperNew.SetResource(&sctx, &resourceWithMetadata)
		if err != nil {
			return err
		}
		sctx.Logger().Debug("MigrateResourceProtobuf: Migration finished for resource with header key: " + string(headerKey))
	}
	sctx.Logger().Debug("MigrateResourceProtobuf: Migration finished")
	return nil
}

func ResourceV1HeaderkeyToDataKey(headerKey []byte) []byte {
	return []byte(
		strings.Replace(
			string(headerKey),
			string(resourcetypesv1.ResourceHeaderKey),
			string(resourcetypesv1.ResourceDataKey),
			1))
}
