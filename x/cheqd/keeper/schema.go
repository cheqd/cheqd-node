package keeper

import (
	"encoding/binary"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetSchemaCount get the total number of schema
func (k Keeper) GetSchemaCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaCountKey))
	byteKey := types.KeyPrefix(types.SchemaCountKey)
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

// SetSchemaCount set the total number of schema
func (k Keeper) SetSchemaCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaCountKey))
	byteKey := types.KeyPrefix(types.SchemaCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendSchema appends a schema in the store with a new id and update the count
func (k Keeper) AppendSchema(
	ctx sdk.Context,
	creator string,
	name string,
	version string,
	attr_names string,
) uint64 {
	// Create the schema
	count := k.GetSchemaCount(ctx)
	var schema = types.Schema{
		Creator:    creator,
		Id:         count,
		Name:       name,
		Version:    version,
		Attr_names: attr_names,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	value := k.cdc.MustMarshalBinaryBare(&schema)
	store.Set(GetSchemaIDBytes(schema.Id), value)

	// Update schema count
	k.SetSchemaCount(ctx, count+1)

	return count
}

// SetSchema set a specific schema in the store
func (k Keeper) SetSchema(ctx sdk.Context, schema types.Schema) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	b := k.cdc.MustMarshalBinaryBare(&schema)
	store.Set(GetSchemaIDBytes(schema.Id), b)
}

// GetSchema returns a schema from its id
func (k Keeper) GetSchema(ctx sdk.Context, id uint64) types.Schema {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	var schema types.Schema
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSchemaIDBytes(id)), &schema)
	return schema
}

// HasSchema checks if the schema exists in the store
func (k Keeper) HasSchema(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	return store.Has(GetSchemaIDBytes(id))
}

// GetSchemaOwner returns the creator of the schema
func (k Keeper) GetSchemaOwner(ctx sdk.Context, id uint64) string {
	return k.GetSchema(ctx, id).Creator
}

// RemoveSchema removes a schema from the store
func (k Keeper) RemoveSchema(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	store.Delete(GetSchemaIDBytes(id))
}

// GetAllSchema returns all schema
func (k Keeper) GetAllSchema(ctx sdk.Context) (list []types.Schema) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Schema
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetSchemaIDBytes returns the byte representation of the ID
func GetSchemaIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetSchemaIDFromBytes returns ID in uint64 format from a byte array
func GetSchemaIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
