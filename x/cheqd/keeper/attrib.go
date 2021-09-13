package keeper

import (
	"encoding/binary"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetAttribCount get the total number of attrib
func (k Keeper) GetAttribCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribCountKey))
	byteKey := types.KeyPrefix(types.AttribCountKey)
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

// SetAttribCount set the total number of attrib
func (k Keeper) SetAttribCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribCountKey))
	byteKey := types.KeyPrefix(types.AttribCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendAttrib appends a attrib in the store with a new id and update the count
func (k Keeper) AppendAttrib(
	ctx sdk.Context,
	creator string,
	did string,
	raw string,
) uint64 {
	// Create the attrib
	count := k.GetAttribCount(ctx)
	var attrib = types.Attrib{
		Creator: creator,
		Id:      count,
		Did:     did,
		Raw:     raw,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	value := k.cdc.MustMarshalBinaryBare(&attrib)
	store.Set(GetAttribIDBytes(attrib.Id), value)

	// Update attrib count
	k.SetAttribCount(ctx, count+1)

	return count
}

// SetAttrib set a specific attrib in the store
func (k Keeper) SetAttrib(ctx sdk.Context, attrib types.Attrib) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	b := k.cdc.MustMarshalBinaryBare(&attrib)
	store.Set(GetAttribIDBytes(attrib.Id), b)
}

// GetAttrib returns a attrib from its id
func (k Keeper) GetAttrib(ctx sdk.Context, id uint64) types.Attrib {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	var attrib types.Attrib
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetAttribIDBytes(id)), &attrib)
	return attrib
}

// HasAttrib checks if the attrib exists in the store
func (k Keeper) HasAttrib(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	return store.Has(GetAttribIDBytes(id))
}

// GetAttribOwner returns the creator of the attrib
func (k Keeper) GetAttribOwner(ctx sdk.Context, id uint64) string {
	return k.GetAttrib(ctx, id).Creator
}

// RemoveAttrib removes a attrib from the store
func (k Keeper) RemoveAttrib(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	store.Delete(GetAttribIDBytes(id))
}

// GetAllAttrib returns all attrib
func (k Keeper) GetAllAttrib(ctx sdk.Context) (list []types.Attrib) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Attrib
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAttribIDBytes returns the byte representation of the ID
func GetAttribIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetAttribIDFromBytes returns ID in uint64 format from a byte array
func GetAttribIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
