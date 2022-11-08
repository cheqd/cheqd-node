package keeper

import (
	"strconv"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetResourceCount get the total number of resource
func (k Keeper) GetResourceCount(ctx *sdk.Context) uint64 {
	store := ctx.KVStore(k.StoreKey)
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
	store := ctx.KVStore(k.StoreKey)
	byteKey := types.KeyPrefix(types.ResourceCountKey)

	// Set bytes
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// SetResource create or update a specific resource in the store
func (k Keeper) SetResource(ctx *sdk.Context, resource *types.ResourceWithMetadata) error {
	if !k.HasResource(ctx, resource.Metadata.CollectionId, resource.Metadata.Id) {
		count := k.GetResourceCount(ctx)
		k.SetResourceCount(ctx, count+1)
	}

	store := ctx.KVStore(k.StoreKey)

	// Set metadata
	metadataKey := GetResourceMetadataKeyBytes(resource.Metadata.CollectionId, resource.Metadata.Id)
	metadataBytes := k.Cdc.MustMarshal(resource.Metadata)
	store.Set(metadataKey, metadataBytes)

	// Set data
	dataKey := GetResourceDataKeyBytes(resource.Metadata.CollectionId, resource.Metadata.Id)
	store.Set(dataKey, resource.Resource.Data)

	return nil
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx *sdk.Context, collectionId string, id string) (types.ResourceWithMetadata, error) {
	if !k.HasResource(ctx, collectionId, id) {
		return types.ResourceWithMetadata{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionId + ":" + id)
	}

	store := ctx.KVStore(k.StoreKey)

	metadataBytes := store.Get(GetResourceMetadataKeyBytes(collectionId, id))
	var metadata types.Metadata
	if err := k.Cdc.Unmarshal(metadataBytes, &metadata); err != nil {
		return types.ResourceWithMetadata{}, sdkerrors.ErrInvalidType.Wrap(err.Error())
	}

	dataBytes := store.Get(GetResourceDataKeyBytes(collectionId, id))
	data := types.Resource{Data: dataBytes}

	return types.ResourceWithMetadata{
		Metadata: &metadata,
		Resource: &data,
	}, nil
}

func (k Keeper) GetResourceMetadata(ctx *sdk.Context, collectionId string, id string) (types.Metadata, error) {
	if !k.HasResource(ctx, collectionId, id) {
		return types.Metadata{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionId + ":" + id)
	}

func (k Keeper) GetAllResourceVersions(ctx *sdk.Context, collectionId, name, resourceType string) []*types.ResourceHeader {
	store := ctx.KVStore(k.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, GetResourceHeaderCollectionPrefixBytes(collectionId))

	metadataBytes := store.Get(GetResourceMetadataKeyBytes(collectionId, id))
	var metadata types.Metadata
	if err := k.Cdc.Unmarshal(metadataBytes, &metadata); err != nil {
		return types.Metadata{}, sdkerrors.ErrInvalidType.Wrap(err.Error())
	}

	return metadata, nil
}

// HasResource checks if the resource exists in the store
func (k Keeper) HasResource(ctx *sdk.Context, collectionId string, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetResourceMetadataKeyBytes(collectionId, id))
}

func (k Keeper) GetResourceCollection(ctx *sdk.Context, collectionId string) []*types.Metadata {
	store := ctx.KVStore(k.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, GetResourceMetadataCollectionPrefixBytes(collectionId))

	var resources []*types.Metadata

	defer CloseIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.Metadata
		k.Cdc.MustUnmarshal(iterator.Value(), &val)
		resources = append(resources, &val)

	}

	return resources
}

func (k Keeper) GetLastResourceVersionMetadata(ctx *sdk.Context, collectionId, name, resourceType string) (types.Metadata, bool) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), GetResourceMetadataCollectionPrefixBytes(collectionId))

	defer CloseIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var metadata types.Metadata
		k.Cdc.MustUnmarshal(iterator.Value(), &metadata)

		if metadata.Name == name && metadata.ResourceType == resourceType && metadata.NextVersionId == "" {
			return metadata, true
		}
	}

	return types.Metadata{}, false
}

// UpdateResourceMetadata update the metadata of a resource. Returns an error if the resource doesn't exist
func (k Keeper) UpdateResourceMetadata(ctx *sdk.Context, metadata *types.Metadata) error {
	if !k.HasResource(ctx, metadata.CollectionId, metadata.Id) {
		return sdkerrors.ErrNotFound.Wrap("resource " + metadata.CollectionId + ":" + metadata.Id)
	}

	store := ctx.KVStore(k.StoreKey)

	// Set metadata
	metadataKey := GetResourceMetadataKeyBytes(metadata.CollectionId, metadata.Id)
	metadataBytes := k.Cdc.MustMarshal(metadata)
	store.Set(metadataKey, metadataBytes)

	return nil
}

// GetAllResources returns all resources as a list
// Loads everything in memory. Use only for genesis export!
func (k Keeper) GetAllResources(ctx *sdk.Context) (list []types.ResourceWithMetadata) {
	metadataIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.KeyPrefix(types.ResourceMetadataKey))
	dataIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.KeyPrefix(types.ResourceDataKey))

	defer closeIteratorOrPanic(metadataIterator)
	defer closeIteratorOrPanic(dataIterator)

	for metadataIterator.Valid() {
		if !dataIterator.Valid() {
			panic("number of headers and data don't match")
		}

		var metadata types.Metadata
		k.Cdc.MustUnmarshal(metadataIterator.Value(), &metadata)

		data := types.Resource{Data: dataIterator.Value()}

		list = append(list, types.ResourceWithMetadata{
			Metadata: &metadata,
			Resource: &data,
		})

		metadataIterator.Next()
		dataIterator.Next()
	}

	return
}

func (k Keeper) GetHeaderIterator(ctx *sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.KeyPrefix(types.ResourceHeaderKey))
}

func (k Keeper) GetDataIterator(ctx *sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.KeyPrefix(types.ResourceDataKey))
}

// GetResourceMetadataKeyBytes returns the byte representation of resource key
func GetResourceMetadataKeyBytes(collectionId string, id string) []byte {
	return []byte(types.ResourceMetadataKey + collectionId + ":" + id)
}

// GetResourceMetadataCollectionPrefixBytes used to iterate over all resource metadatas in a collection
func GetResourceMetadataCollectionPrefixBytes(collectionId string) []byte {
	return []byte(types.ResourceMetadataKey + collectionId + ":")
}

// GetResourceDataKeyBytes returns the byte representation of resource key
func GetResourceDataKeyBytes(collectionId string, id string) []byte {
	return []byte(types.ResourceDataKey + collectionId + ":" + id)
}

func CloseIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
