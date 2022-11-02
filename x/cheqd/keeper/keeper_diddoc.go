package keeper

import (
	"strconv"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidDocCount(ctx *sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidCountKey))

	key := StrBytes(types.DidCountKey)
	valueBytes := store.Get(key)

	// Count doesn't exist: no element
	if valueBytes == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(valueBytes), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetDidCount set the total number of did
func (k Keeper) SetDidDocCount(ctx *sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidCountKey))

	key := StrBytes(types.DidCountKey)
	valueBytes := []byte(strconv.FormatUint(count, 10))

	store.Set(key, valueBytes)
}

// SetDid set a specific did in the store. Updates DID counter if the DID is new.
func (k Keeper) SetDidDoc(ctx *sdk.Context, value *types.DidDocWithMetadata) error {
	// Update counter
	if !k.HasDidDoc(ctx, value.DidDoc.Id) {
		count := k.GetDidDocCount(ctx)
		k.SetDidDocCount(ctx, count+1)
	}

	// Create the did
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidKey))

	key := StrBytes(value.DidDoc.Id)
	valueBytes := k.cdc.MustMarshal(value)
	store.Set(key, valueBytes)

	return nil
}

// GetDid returns a did from its id
func (k Keeper) GetDidDoc(ctx *sdk.Context, id string) (types.DidDocWithMetadata, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidKey))

	if !k.HasDidDoc(ctx, id) {
		return types.DidDocWithMetadata{}, sdkerrors.ErrNotFound.Wrap(id)
	}

	var value types.DidDocWithMetadata
	key := store.Get(StrBytes(id))
	k.cdc.MustUnmarshal(key, &value)

	return value, nil
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDidDoc(ctx *sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidKey))
	return store.Has(StrBytes(id))
}

// GetAllDidDocs returns all did
// Loads all DIDs in memory. Use only for genesis export.
func (k Keeper) GetAllDidDocs(ctx *sdk.Context) (list []types.DidDocWithMetadata) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err.Error())
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var val types.DidDocWithMetadata
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
