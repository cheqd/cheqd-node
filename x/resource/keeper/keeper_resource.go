package keeper

import (
	"strconv"

	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

// SetResource set a specific resource in the store
func (k Keeper) SetResource(ctx *sdk.Context, resource *types.Resource) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))
	key := GetResourceKeyBytes(resource.CollectionId, resource.Id)
	bytes := k.cdc.MustMarshal(resource)

	if !store.Has(key) {
		count := k.GetResourceCount(ctx)
		k.SetResourceCount(ctx, count+1)
	}

	store.Set(key, bytes)
	return nil
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx *sdk.Context, collectionId string, id string) (types.Resource, error) {
	if !k.HasResource(ctx, collectionId, id) {
		return types.Resource{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionId + ":" + id)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))
	bytes := store.Get(GetResourceKeyBytes(collectionId, id))

	var value types.Resource
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return types.Resource{}, sdkerrors.ErrInvalidType.Wrap(err.Error())
	}

	return value, nil
}

// HasResource checks if the resource exists in the store
func (k Keeper) HasResource(ctx *sdk.Context, collectionId string, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))
	return store.Has(GetResourceKeyBytes(collectionId, id))
}

// GetResourceKeyBytes returns the byte representation of resource key
func GetResourceKeyBytes(collectionId string, id string) []byte {
	return []byte(collectionId + ":" + id)
}

func (k Keeper) GetLastResourceVersion(collectionId, id, name, resourceType, mimeType string) (uint64, error) {
	store := prefix.NewStore(k.ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))


// GetAllResources returns all resources as a list
// Loads everything in memory. Use only for genesis export!
func (k Keeper) GetAllResources(ctx *sdk.Context) (list []types.Resource) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err.Error())
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.Resource
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
