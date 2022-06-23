package keeper

import (
	"strconv"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetResourceCount get the total number of resource
func (k Keeper) GetResourceCount(ctx *sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.ResourceCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetResourceCount set the total number of resource
func (k Keeper) SetResourceCount(ctx *sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.ResourceCountKey)

	// Set bytes
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// SetResource create or update a specific resource in the store
func (k Keeper) SetResource(ctx *sdk.Context, resource *types.Resource) error {
	if !k.HasResource(ctx, resource.Header.CollectionId, resource.Header.Id) {
		count := k.GetResourceCount(ctx)
		k.SetResourceCount(ctx, count+1)
	}

	store := ctx.KVStore(k.storeKey)

	// Set header
	headerKey := GetResourceHeaderKeyBytes(resource.Header.CollectionId, resource.Header.Id)
	headerBytes := k.cdc.MustMarshal(resource.Header)
	store.Set(headerKey, headerBytes)

	// Set data
	dataKey := GetResourceDataKeyBytes(resource.Header.CollectionId, resource.Header.Id)
	store.Set(dataKey, resource.Data)

	return nil
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx *sdk.Context, collectionId string, id string) (types.Resource, error) {
	if !k.HasResource(ctx, collectionId, id) {
		return types.Resource{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionId + ":" + id)
	}

	store := ctx.KVStore(k.storeKey)

	headerBytes := store.Get(GetResourceHeaderKeyBytes(collectionId, id))
	var header types.ResourceHeader
	if err := k.cdc.Unmarshal(headerBytes, &header); err != nil {
		return types.Resource{}, sdkerrors.ErrInvalidType.Wrap(err.Error())
	}

	dataBytes := store.Get(GetResourceDataKeyBytes(collectionId, id))

	return types.Resource{
		Header: &header,
		Data:   dataBytes,
	}, nil
}

// HasResource checks if the resource exists in the store
func (k Keeper) HasResource(ctx *sdk.Context, collectionId string, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetResourceHeaderKeyBytes(collectionId, id))
}

func (k Keeper) GetAllResourceVersions(ctx *sdk.Context, collectionId, name, resourceType, mimeType string) []*types.ResourceHeader {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetResourceHeaderCollectionPrefixBytes(collectionId))

	defer closeIteratorOrPanic(iterator)

	var result []*types.ResourceHeader

	for ; iterator.Valid(); iterator.Next() {
		var val types.ResourceHeader
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if val.Name == name &&
			val.ResourceType == resourceType &&
			val.MimeType == mimeType {
			result = append(result, &val)
		}
	}

	return result
}

func (k Keeper) GetResourceCollection(ctx *sdk.Context, collectionId string) []*types.ResourceHeader {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetResourceHeaderCollectionPrefixBytes(collectionId))

	var resources []*types.ResourceHeader

	defer closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.ResourceHeader
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		resources = append(resources, &val)

	}

	return resources
}

func (k Keeper) GetLastResourceVersionHeader(ctx *sdk.Context, collectionId, name, resourceType, mimeType string) (types.ResourceHeader, bool) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), GetResourceHeaderCollectionPrefixBytes(collectionId))

	defer closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.ResourceHeader
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if val.Name == name &&
			val.ResourceType == resourceType &&
			val.MimeType == mimeType &&
			val.NextVersionId == "" {
			return val, true
		}
	}

	return types.ResourceHeader{}, false
}

// UpdateResourceHeader update the header of a resource. Returns an error if the resource doesn't exist
func (k Keeper) UpdateResourceHeader(ctx *sdk.Context, header *types.ResourceHeader) error {
	if !k.HasResource(ctx, header.CollectionId, header.Id) {
		return sdkerrors.ErrNotFound.Wrap("resource " + header.CollectionId + ":" + header.Id)
	}

	store := ctx.KVStore(k.storeKey)

	// Set header
	headerKey := GetResourceHeaderKeyBytes(header.CollectionId, header.Id)
	headerBytes := k.cdc.MustMarshal(header)
	store.Set(headerKey, headerBytes)

	return nil
}

// GetAllResources returns all resources as a list
// Loads everything in memory. Use only for genesis export!
func (k Keeper) GetAllResources(ctx *sdk.Context) (list []types.Resource) {
	headerIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceHeaderKey))
	dataIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceDataKey))

	defer closeIteratorOrPanic(headerIterator)
	defer closeIteratorOrPanic(dataIterator)


	for headerIterator.Valid() {
		if !dataIterator.Valid() {
			panic("number of headers and data don't match")
		}

		var val types.ResourceHeader
		k.cdc.MustUnmarshal(headerIterator.Value(), &val)

		list = append(list, types.Resource{
			Header: &val,
			Data:   dataIterator.Value(),
		})

		headerIterator.Next()
		dataIterator.Next()
	}

	return
}

// GetResourceHeaderKeyBytes returns the byte representation of resource key
func GetResourceHeaderKeyBytes(collectionId string, id string) []byte {
	return []byte(types.ResourceHeaderKey + collectionId + ":" + id)
}

// GetResourceHeaderCollectionPrefixBytes used to iterate over all resource headers in a collection
func GetResourceHeaderCollectionPrefixBytes(collectionId string) []byte {
	return []byte(types.ResourceHeaderKey + collectionId + ":")
}

// GetResourceDataKeyBytes returns the byte representation of resource key
func GetResourceDataKeyBytes(collectionId string, id string) []byte {
	return []byte(types.ResourceDataKey + collectionId + ":" + id)
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}