package keeper

import (
	"strconv"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidCount(ctx *sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidCountKey))
	byteKey := types.KeyPrefix(types.DidCountKey)
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

// SetDidCount set the total number of did
func (k Keeper) SetDidCount(ctx *sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidCountKey))
	byteKey := types.KeyPrefix(types.DidCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendDid appends a did in the store with a new id and updates the count
func (k Keeper) AppendDid(ctx *sdk.Context, did *types.Did, metadata *types.Metadata) error {
	// Check that did doesn't exist
	if k.HasDid(ctx, did.Id) {
		return types.ErrDidDocExists.Wrapf(did.Id)
	}

	// Create the did
	count := k.GetDidCount(ctx)
	err := k.SetDid(ctx, did, metadata)
	if err != nil {
		return err
	}

	// Update did count
	k.SetDidCount(ctx, count+1)
	return nil
}

// SetDid set a specific did in the store
func (k Keeper) SetDid(ctx *sdk.Context, did *types.Did, metadata *types.Metadata) error {
	stateValue, err := types.NewStateValue(did, metadata)
	if err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	b := k.cdc.MustMarshal(&stateValue)
	store.Set(GetDidIDBytes(did.Id), b)
	return nil
}

// GetDid returns a did from its id
func (k Keeper) GetDid(ctx *sdk.Context, id string) (types.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))

	if !k.HasDid(ctx, id) {
		return types.StateValue{}, sdkerrors.ErrNotFound.Wrap(id)
	}

	var value types.StateValue
	bytes := store.Get(GetDidIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return types.StateValue{}, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return value, nil
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDid(ctx *sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	return store.Has(GetDidIDBytes(id))
}

// GetDidIDBytes returns the byte representation of the ID
func GetDidIDBytes(id string) []byte {
	return []byte(id)
}

// GetAllDid returns all did
// Loads all DIDs in memory. Use only for genesis export.
func (k Keeper) GetAllDid(ctx *sdk.Context) (list []types.StateValue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err.Error())
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.StateValue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
