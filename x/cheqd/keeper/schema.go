package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	schema types.Schema,
	metadata *types.Metadata,
) (*string, error) {
	// Create the schema
	count := k.GetSchemaCount(ctx)

	err := k.SetSchema(ctx, schema, metadata)
	if err != nil {
		return nil, err
	}

	// Update schema count
	k.SetSchemaCount(ctx, count+1)

	return &schema.Id, nil
}

// SetSchema set a specific cred def in the store
func (k Keeper) SetSchema(ctx sdk.Context, schema types.Schema, metadata *types.Metadata) error {
	stateValue, err := types.NewStateValue(&schema, metadata)
	if err != nil {
		return types.ErrSetToState.Wrap(err.Error())
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	b := k.cdc.MustMarshal(stateValue)
	store.Set(GetSchemaIDBytes(schema.Id), b)
	return nil
}

// GetSchema returns a schema from its id
func (k Keeper) GetSchema(ctx sdk.Context, id string) (*types.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))

	if !k.HasSchema(ctx, id) {
		return nil, sdkerrors.ErrNotFound
	}

	var value types.StateValue
	var bytes = store.Get(GetSchemaIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return &value, nil
}

// HasSchema checks if the schema exists in the store
func (k Keeper) HasSchema(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	return store.Has(GetSchemaIDBytes(id))
}

// GetSchemaIDBytes returns the byte representation of the ID
func GetSchemaIDBytes(id string) []byte {
	return []byte(id)
}
