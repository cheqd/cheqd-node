package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cheqd-id/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetNymCount get the total number of nym
func (k Keeper) GetNymCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymCountKey))
	byteKey := types.KeyPrefix(types.NymCountKey)
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

// SetNymCount set the total number of nym
func (k Keeper) SetNymCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymCountKey))
	byteKey := types.KeyPrefix(types.NymCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendNym appends a nym in the store with a new id and update the count
func (k Keeper) AppendNym(
	ctx sdk.Context,
	creator string,
	alias string,
	verkey string,
	did string,
	role string,
) uint64 {
	// Create the nym
	count := k.GetNymCount(ctx)
	var nym = types.Nym{
		Creator: creator,
		Id:      count,
		Alias:   alias,
		Verkey:  verkey,
		Did:     did,
		Role:    role,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	value := k.cdc.MustMarshalBinaryBare(&nym)
	store.Set(GetNymIDBytes(nym.Id), value)

	// Update nym count
	k.SetNymCount(ctx, count+1)

	return count
}

// SetNym set a specific nym in the store
func (k Keeper) SetNym(ctx sdk.Context, nym types.Nym) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	b := k.cdc.MustMarshalBinaryBare(&nym)
	store.Set(GetNymIDBytes(nym.Id), b)
}

// GetNym returns a nym from its id
func (k Keeper) GetNym(ctx sdk.Context, id uint64) types.Nym {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	var nym types.Nym
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetNymIDBytes(id)), &nym)
	return nym
}

// HasNym checks if the nym exists in the store
func (k Keeper) HasNym(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	return store.Has(GetNymIDBytes(id))
}

// GetNymOwner returns the creator of the nym
func (k Keeper) GetNymOwner(ctx sdk.Context, id uint64) string {
	return k.GetNym(ctx, id).Creator
}

// RemoveNym removes a nym from the store
func (k Keeper) RemoveNym(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	store.Delete(GetNymIDBytes(id))
}

// GetAllNym returns all nym
func (k Keeper) GetAllNym(ctx sdk.Context) (list []types.Nym) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nym
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetNymIDBytes returns the byte representation of the ID
func GetNymIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetNymIDFromBytes returns ID in uint64 format from a byte array
func GetNymIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
