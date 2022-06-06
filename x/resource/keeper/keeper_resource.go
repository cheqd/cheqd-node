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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceCountKey))
	byteKey := types.KeyPrefix(types.ResourceCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetResourceCount set the total number of resource
func (k Keeper) SetResourceCount(ctx *sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceCountKey))
	byteKey := types.KeyPrefix(types.ResourceCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendResource appends a resource in the store with a new id and updates the count
func (k Keeper) AppendResource(ctx *sdk.Context, resource *types.Resource) error {
	// Check that resource doesn't exist
	if k.HasResource(ctx, resource.CollectionId, resource.Id) {
		return types.ErrResourceExists.Wrapf(resource.Id)
	}

	// Create the resource
	count := k.GetResourceCount(ctx)
	err := k.SetResource(ctx, resource)
	if err != nil {
		return err
	}

	// Update resource count
	k.SetResourceCount(ctx, count+1)
	return nil
}

// SetResource set a specific resource in the store
func (k Keeper) SetResource(ctx *sdk.Context, resource *types.Resource) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))
	b := k.cdc.MustMarshal(&resource)
	store.Set(GetResourceKeyBytes(resource.CollectionId, resource.Id), b)
	return nil
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx *sdk.Context, collectionId string, id string) (types.Resource, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))

	if !k.HasResource(ctx, collectionId, id) {
		return types.Resource{}, sdkerrors.ErrNotFound.Wrap(id)
	}

	var value types.Resource
	bytes := store.Get(GetResourceKeyBytes(collectionId, id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return types.Resource{}, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
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

// GetAllResource returns all resource
func (k Keeper) GetAllResources(ctx *sdk.Context) (list []types.Resource) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ResourceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

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
