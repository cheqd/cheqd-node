package keeper

import (
	"encoding/binary"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetCred_defCount get the total number of cred_def
func (k Keeper) GetCred_defCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defCountKey))
	byteKey := types.KeyPrefix(types.Cred_defCountKey)
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

// SetCred_defCount set the total number of cred_def
func (k Keeper) SetCred_defCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defCountKey))
	byteKey := types.KeyPrefix(types.Cred_defCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendCred_def appends a cred_def in the store with a new id and update the count
func (k Keeper) AppendCred_def(
	ctx sdk.Context,
	creator string,
	schema_id string,
	tag string,
	signature_type string,
	value string,
) uint64 {
	// Create the cred_def
	count := k.GetCred_defCount(ctx)
	var cred_def = types.Cred_def{
		Creator:        creator,
		Id:             count,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	value := k.cdc.MustMarshalBinaryBare(&cred_def)
	store.Set(GetCred_defIDBytes(cred_def.Id), value)

	// Update cred_def count
	k.SetCred_defCount(ctx, count+1)

	return count
}

// SetCred_def set a specific cred_def in the store
func (k Keeper) SetCred_def(ctx sdk.Context, cred_def types.Cred_def) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	b := k.cdc.MustMarshalBinaryBare(&cred_def)
	store.Set(GetCred_defIDBytes(cred_def.Id), b)
}

// GetCred_def returns a cred_def from its id
func (k Keeper) GetCred_def(ctx sdk.Context, id uint64) types.Cred_def {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	var cred_def types.Cred_def
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCred_defIDBytes(id)), &cred_def)
	return cred_def
}

// HasCred_def checks if the cred_def exists in the store
func (k Keeper) HasCred_def(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	return store.Has(GetCred_defIDBytes(id))
}

// GetCred_defOwner returns the creator of the cred_def
func (k Keeper) GetCred_defOwner(ctx sdk.Context, id uint64) string {
	return k.GetCred_def(ctx, id).Creator
}

// RemoveCred_def removes a cred_def from the store
func (k Keeper) RemoveCred_def(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	store.Delete(GetCred_defIDBytes(id))
}

// GetAllCred_def returns all cred_def
func (k Keeper) GetAllCred_def(ctx sdk.Context) (list []types.Cred_def) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Cred_def
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCred_defIDBytes returns the byte representation of the ID
func GetCred_defIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCred_defIDFromBytes returns ID in uint64 format from a byte array
func GetCred_defIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
