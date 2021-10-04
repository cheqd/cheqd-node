package keeper

import (
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

// AppendCredDef TODO add Value
// AppendCredDef appends a credDef in the store with a new id and update the count
func (k Keeper) AppendCredDef(
	ctx sdk.Context,
	id string,
	schemaId string,
	tag string,
	signatureType string,
	clValue *types.CredDef_ClType,
) string {
	// Create the credDef
	count := k.GetCredDefCount(ctx)

	var credDef = types.CredDef{
		Id:            id,
		SchemaId:      schemaId,
		Tag:           tag,
		SignatureType: signatureType,
		Value:         clValue,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	value := k.cdc.MustMarshal(&credDef)
	store.Set(GetCredDefIDBytes(credDef.Id), value)

	// Update credDef count
	k.SetCredDefCount(ctx, count+1)

	return id
}

// SetCredDef set a specific credDef in the store
func (k Keeper) SetCredDef(ctx sdk.Context, credDef types.CredDef) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	b := k.cdc.MustMarshal(&credDef)
	store.Set(GetCredDefIDBytes(credDef.Id), b)
}

// GetCredDef returns a credDef from its id
func (k Keeper) GetCredDef(ctx sdk.Context, id string) types.CredDef {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	var credDef types.CredDef
	k.cdc.MustUnmarshal(store.Get(GetCredDefIDBytes(id)), &credDef)
	return credDef
}

// HasCredDef checks if the credDef exists in the store
func (k Keeper) HasCredDef(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	return store.Has(GetCredDefIDBytes(id))
}

// RemoveCredDef removes a credDef from the store
func (k Keeper) RemoveCredDef(ctx sdk.Context, id string) {
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
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCredDefIDBytes returns the byte representation of the ID
func GetCredDefIDBytes(id string) []byte {
	return []byte(id)
}

// GetCredDefIDFromBytes returns ID in string format from a byte array
func GetCredDefIDFromBytes(bz []byte) string {
	return string(bz)
}
