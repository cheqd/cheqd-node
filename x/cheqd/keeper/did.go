package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidCountKey))
	byteKey := v1.KeyPrefix(v1.DidCountKey)
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
func (k Keeper) SetDidCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidCountKey))
	byteKey := v1.KeyPrefix(v1.DidCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendDid appends a did in the store with a new id and updates the count
func (k Keeper) AppendDid(ctx sdk.Context, did v1.Did, metadata *v1.Metadata) (*string, error) {
	// Create the did
	count := k.GetDidCount(ctx)
	err := k.SetDid(ctx, did, metadata)
	if err != nil {
		return nil, err
	}

	// Update did count
	k.SetDidCount(ctx, count+1)
	return &did.Id, nil
}

// SetDid set a specific did in the store
func (k Keeper) SetDid(ctx sdk.Context, did v1.Did, metadata *v1.Metadata) error {
	stateValue, err := v1.NewStateValue(&did, metadata)
	if err != nil {
		return v1.ErrSetToState.Wrap(err.Error())
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidKey))
	b := k.cdc.MustMarshal(stateValue)
	store.Set(GetDidIDBytes(did.Id), b)
	return nil
}

// GetDid returns a did from its id
func (k Keeper) GetDid(ctx *sdk.Context, id string) (*v1.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidKey))

	if !k.HasDid(*ctx, id) {
		return nil, sdkerrors.ErrNotFound
	}

	var value v1.StateValue
	var bytes = store.Get(GetDidIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return &value, nil
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDid(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidKey))
	return store.Has(GetDidIDBytes(id))
}

// GetDidIDBytes returns the byte representation of the ID
func GetDidIDBytes(id string) []byte {
	return []byte(id)
}

// GetAllDid returns all did
func (k Keeper) GetAllDid(ctx sdk.Context) (list []v1.StateValue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), v1.KeyPrefix(v1.DidKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val v1.StateValue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
