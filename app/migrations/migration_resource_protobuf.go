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
	// Storage for old headers and data
	store := sctx.KVStore(mctx.resourceStoreKey)

	// Reset counter
	mctx.resourceKeeperNew.SetResourceCount(&sctx, 0)

	headerKeys := helpers.ReadAllKeys(store, didutils.StrBytes(resourcetypesv1.ResourceHeaderKey))

	for _, headerKey := range headerKeys {
		dataKey := ResourceV1HeaderkeyToDataKey(headerKey)

		var oldHeader resourcetypesv1.ResourceHeader
		mctx.codec.MustUnmarshal(store.Get(headerKey), &oldHeader)
		oldData := store.Get(dataKey)

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

		// Remove old resource data and header
		store.Delete(headerKey)
		store.Delete(dataKey)

		// Write new resource
		err := mctx.resourceKeeperNew.SetResource(&sctx, &resourceWithMetadata)
		if err != nil {
			return err
		}
	}
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
