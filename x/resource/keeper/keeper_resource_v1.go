package keeper

import (
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllResourcesV1(ctx *sdk.Context) (list []resourcetypesv1.Resource) {
	headerIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), resourcetypesv2.KeyPrefix(resourcetypesv2.ResourceMetadataKey))
	dataIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), resourcetypesv2.KeyPrefix(resourcetypesv2.ResourceDataKey))

	defer closeIteratorOrPanic(headerIterator)
	defer closeIteratorOrPanic(dataIterator)

	for headerIterator.Valid() {
		if !dataIterator.Valid() {
			panic("number of headers and data don't match")
		}

		var val resourcetypesv1.ResourceHeader
		k.cdc.MustUnmarshal(headerIterator.Value(), &val)

		list = append(list, resourcetypesv1.Resource{
			Header: &val,
			Data:   dataIterator.Value(),
		})

		headerIterator.Next()
		dataIterator.Next()
	}

	return
}
