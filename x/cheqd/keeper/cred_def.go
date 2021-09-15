package keeper

import (
	"encoding/binary"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetCredDefCount get the total number of credDef
func (k Keeper) GetCredDefCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefCountKey))
	byteKey := types.KeyPrefix(types.CredDefCountKey)
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

// SetCredDefCount set the total number of credDef
func (k Keeper) SetCredDefCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefCountKey))
	byteKey := types.KeyPrefix(types.CredDefCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendCredDef appends a credDef in the store with a new id and update the count
func (k Keeper) AppendCredDef(
	ctx sdk.Context,
	creator string,
	schema_id string,
	tag string,
	signature_type string,
	value string,
) uint64 {
	// Create the credDef
	count := k.GetCredDefCount(ctx)
	var credDef = types.CredDef{
		Creator:        creator,
		Id:             count,
		Schema_id:      schema_id,
		Tag:            tag,
		Signature_type: signature_type,
		Value:          value,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	value := k.cdc.MustMarshalBinaryBare(&credDef)
	store.Set(GetCredDefIDBytes(credDef.Id), value)

	// Update credDef count
	k.SetCredDefCount(ctx, count+1)

	return count
}

// SetCredDef set a specific credDef in the store
func (k Keeper) SetCredDef(ctx sdk.Context, credDef types.CredDef) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	b := k.cdc.MustMarshalBinaryBare(&credDef)
	store.Set(GetCredDefIDBytes(credDef.Id), b)
}

// GetCredDef returns a credDef from its id
func (k Keeper) GetCredDef(ctx sdk.Context, id uint64) types.CredDef {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	var credDef types.CredDef
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCredDefIDBytes(id)), &credDef)
	return credDef
}

// HasCredDef checks if the credDef exists in the store
func (k Keeper) HasCredDef(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	return store.Has(GetCredDefIDBytes(id))
}

// GetCredDefOwner returns the creator of the credDef
func (k Keeper) GetCredDefOwner(ctx sdk.Context, id uint64) string {
	return k.GetCredDef(ctx, id).Creator
}

// RemoveCredDef removes a credDef from the store
func (k Keeper) RemoveCredDef(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	store.Delete(GetCredDefIDBytes(id))
}

// GetAllCredDef returns all credDef
func (k Keeper) GetAllCredDef(ctx sdk.Context) (list []types.CredDef) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CredDef
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCredDefIDBytes returns the byte representation of the ID
func GetCredDefIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCredDefIDFromBytes returns ID in uint64 format from a byte array
func GetCredDefIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
